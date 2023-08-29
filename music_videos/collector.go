package music_videos

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/kodi"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

type Collector struct {
	config  *config.Config
	watcher *fsnotify.Watcher
	channel chan *MusicVideo
}

var collector *Collector

func RunCollector(config *config.Config) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.Logger.FatalF("new movies watcher err: %v", err)
	}

	collector = &Collector{
		config:  config,
		watcher: watcher,
		channel: make(chan *MusicVideo, runtime.NumCPU()),
	}

	go collector.runWatcher()
	go collector.runProcessor()
	collector.runScanner()
}

// 运行文件变动监听
func (c *Collector) runWatcher() {
	if !c.config.Collector.Watcher {
		return
	}

	utils.Logger.Debug("run music videos watcher")

	for {
		select {
		case event, ok := <-c.watcher.Events:
			if !ok {
				continue
			}

			fileInfo, err := os.Stat(event.Name)
			if fileInfo == nil || err != nil {
				utils.Logger.WarningF("get music video stat err: %v", err)
				continue
			}

			// 删除文件夹
			if event.Has(fsnotify.Remove) && fileInfo.IsDir() {
				utils.Logger.InfoF("removed dir: %s", event.Name)

				err := c.watcher.Remove(event.Name)
				if err != nil {
					utils.Logger.WarningF("remove shows watcher: %s error: %v", event.Name, err)
				}
				continue
			}

			if !event.Has(fsnotify.Create) || c.skipFolders(filepath.Dir(event.Name), event.Name) {
				continue
			}

			//  新增目录
			if fileInfo.IsDir() {
				err = c.watcher.Add(event.Name)
				if err != nil {
					utils.Logger.WarningF("add music video dir: %s to watcher err: %v", event.Name, err)
				}

				videos, err := c.scanDir(event.Name)
				if err != nil || len(videos) == 0 {
					utils.Logger.WarningF("new dir %s scan err: %v", event.Name, err)
					continue
				}

				for _, video := range videos {
					c.channel <- video
				}

				continue
			}

			// 单个文件
			if utils.IsVideo(event.Name) != "" {
				video := c.parseVideoFile(filepath.Dir(event.Name), fileInfo)
				c.channel <- video
				continue
			}

		case err, ok := <-c.watcher.Errors:
			if !ok {
				return
			}

			utils.Logger.ErrorF("music videos watcher error: %v", err)
		}
	}
}

// 处理扫描队列
func (c *Collector) runProcessor() {
	utils.Logger.Debug("run music videos processor")

	limiter := make(chan struct{}, c.config.Ffmpeg.MaxWorker)
	for {
		select {
		case video := <-c.channel:
			utils.Logger.DebugF("receive music video task: %v", video)

			limiter <- struct{}{}
			go func() {
				c.videoProcessor(video)
				<-limiter
			}()
		}
	}
}

// 视频文件处理
func (c *Collector) videoProcessor(video *MusicVideo) {
	if video == nil || (video.NfoExist() && video.ThumbExist()) {
		return
	}

	probe, err := video.getProbe()
	if err != nil {
		utils.Logger.WarningF("parse video %s probe err: %v", video.Dir+"/"+video.OriginTitle, err)
		return
	}

	video.VideoStream = probe.FirstVideoStream()
	video.AudioStream = probe.FirstAudioStream()
	if video.VideoStream == nil || video.AudioStream == nil {
		return
	}

	err = video.drawThumb()
	if err != nil {
		utils.Logger.WarningF("draw thumb err: %v", err)
		return
	}

	err = video.saveToNfo()
	if err != nil {
		utils.Logger.WarningF("save to NFO err: %v", err)
		return
	}

	split := strings.Split(video.BaseDir, "/")
	kodi.Rpc.AddScanTask(split[len(split)-1])
}

// 运行扫描器
func (c *Collector) runScanner() {
	utils.Logger.DebugF("run music video scanner cron_seconds: %d", c.config.Collector.CronSeconds)

	task := func() {
		for _, item := range c.config.Collector.MusicVideosDir {
			_ = c.watcher.Add(item)
			utils.Logger.DebugF("runCronScan add music videos dir: %s to watcher", item)

			videos, err := c.scanDir(item)
			if len(videos) == 0 || err != nil {
				continue
			}

			// 刮削信息缓存目录
			cacheDir := item + "/tmdb"
			if _, err := os.Stat(cacheDir); err != nil && os.IsNotExist(err) {
				err := os.Mkdir(cacheDir, 0755)
				if err != nil {
					utils.Logger.ErrorF("create probe cache: %s dir err: %v", cacheDir, err)
					continue
				}
			}

			for _, video := range videos {
				c.channel <- video
			}
		}
	}

	task()
	ticker := time.NewTicker(time.Second * time.Duration(c.config.Collector.CronSeconds))
	for range ticker.C {
		task()
		utils.Logger.Debug("run music video scanner finished")
	}
}

func (c *Collector) scanDir(dir string) ([]*MusicVideo, error) {
	videos := make([]*MusicVideo, 0)
	dirInfo, err := ioutil.ReadDir(dir)
	if err != nil {
		utils.Logger.WarningF("scanDir %s err: %v", dir, err)
		return nil, err
	}

	for _, file := range dirInfo {
		if file.IsDir() {
			if c.skipFolders(dir, file.Name()) {
				utils.Logger.DebugF("passed in skip folders: %s", file.Name())
				continue
			}

			_ = c.watcher.Add(dir + "/" + file.Name())
			utils.Logger.DebugF("scanDir add music videos dir: %s to watcher", file.Name())

			subVideos, err := c.scanDir(dir + "/" + file.Name())
			if err != nil {
				continue
			}

			if len(subVideos) > 0 {
				videos = append(videos, subVideos...)
			}

			continue
		}

		video := c.parseVideoFile(dir, file)
		if video != nil {
			videos = append(videos, video)
		}
	}

	return videos, err
}

func (c *Collector) skipFolders(path, filename string) bool {
	base := filepath.Base(path)
	for _, item := range c.config.Collector.SkipFolders {
		if item == base || item == filename {
			return true
		}
	}
	return false
}
