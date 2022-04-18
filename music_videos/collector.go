package music_videos

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"github.com/fsnotify/fsnotify"
	"io/ioutil"
	"path/filepath"
	"runtime"
	"time"
)

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
	// todo
}

// 运行处理器
func (c *Collector) runProcessor() {
	utils.Logger.Debug("run music videos processor")

	limiter := make(chan struct{}, runtime.NumCPU())
	for {
		select {
		case video := <-c.channel:
			utils.Logger.DebugF("receive music video task: %v", video)

			limiter <- struct{}{}
			go func() {
				err := video.drawThumb()
				if err == nil {
					_ = video.saveToNfo()
				}
				<-limiter
			}()
		}
	}
}

// 运行扫描器
func (c *Collector) runScanner() {
	utils.Logger.DebugF("run music video scanner cron_seconds: %d", c.config.CronSeconds)

	task := func() {
		for _, item := range c.config.MusicVideosDir {
			videos, err := c.scanDir(item)
			if err != nil {
				continue
			}

			for _, video := range videos {
				c.channel <- video
			}
		}
	}

	task()
	ticker := time.NewTicker(time.Second * time.Duration(c.config.CronSeconds))
	for {
		select {
		case <-ticker.C:
			task()

			utils.Logger.Debug("run music video scanner finished")
		}
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
	for _, item := range c.config.MusicVideosSkipFolders {
		if item == base || item == filename {
			return true
		}
	}
	return false
}
