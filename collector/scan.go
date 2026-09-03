package collector

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/kodi"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"time"
)

// collector 运行扫描
func (c *collector) runScan() {
	producerWG := &sync.WaitGroup{}
	scanTaskWG := &sync.WaitGroup{}

	if config.Collector.RunMode == config.CollectorRunModeSpec {
		pwd, err := os.Getwd()
		if err != nil {
			utils.Logger.ErrorF("get pwd error: %s", err)
			return
		}

		for _, item := range config.Collector.MoviesDir {
			if strings.HasPrefix(pwd, item) {
				producerWG.Add(1)
				go c.scanDir([]string{pwd}, media_file.Movies, producerWG, scanTaskWG)
				break
			}
		}
		for _, item := range config.Collector.ShowsDir {
			if strings.HasPrefix(pwd, item) {
				producerWG.Add(1)
				go c.scanDir([]string{pwd}, media_file.TvShows, producerWG, scanTaskWG)
				break
			}
		}
		for _, item := range config.Collector.MusicVideosDir {
			if strings.HasPrefix(pwd, item) {
				producerWG.Add(1)
				go c.scanDir([]string{pwd}, media_file.MusicVideo, producerWG, scanTaskWG)
				break
			}
		}
	} else {
		producerWG.Add(3)
		go c.scanDir(config.Collector.MoviesDir, media_file.Movies, producerWG, scanTaskWG)
		go c.scanDir(config.Collector.ShowsDir, media_file.TvShows, producerWG, scanTaskWG)
		go c.scanDir(config.Collector.MusicVideosDir, media_file.MusicVideo, producerWG, scanTaskWG)
	}

	producerWG.Wait()
	scanTaskWG.Wait()

	// 扫描完成，通知kodi刷新媒体库
	if config.Kodi.Enable && config.Collector.CronScanKodi {
		log.Println("scan done, refresh kodi library")
		kodi.Rpc.VideoLibrary.Scan("", false)
	}

	// 扫描完成后，通知kodi清理媒体库
	if config.Kodi.CleanLibrary {
		log.Println("scan done, clean kodi library")
		kodi.Rpc.AddCleanTask("")
	}

	// 单次模式，关闭channel
	if config.Collector.RunMode == config.CollectorRunModeOnce || config.Collector.RunMode == config.CollectorRunModeSpec {
		c.closeOnce.Do(func() { close(c.channel) })
	}
}

// runCronScan 运行定时扫描
func (c *collector) runCronScan() {
	if !config.Collector.CronScan || config.Collector.CronSeconds <= 0 {
		return
	}

	ticker := time.NewTicker(time.Second * time.Duration(config.Collector.CronSeconds))
	defer ticker.Stop()
	for range ticker.C {
		c.runScan()
	}
}

// registerWatcherDirs 仅注册 watcher 目录，不产生扫描任务
// watcher 目录原本只在扫描时注册，定时扫描和启动扫描都关闭时需在此注册，否则 watcher 失效
func (c *collector) registerWatcherDirs() {
	c.walkMediaDir(config.Collector.MoviesDir, media_file.Movies, nil)
	c.walkMediaDir(config.Collector.ShowsDir, media_file.TvShows, nil)
	c.walkMediaDir(config.Collector.MusicVideosDir, media_file.MusicVideo, nil)
}

// scanDir 扫描目录
func (c *collector) scanDir(roots []string, videoType media_file.VideoType, producerWG, scanTaskWG *sync.WaitGroup) {
	defer producerWG.Done()

	c.walkMediaDir(roots, videoType, func(mf *media_file.MediaFile) {
		scanTaskWG.Add(1)
		c.channel <- &scanTask{file: mf, done: scanTaskWG}
	})
}

// walkMediaDir 遍历媒体目录，统一处理隐藏目录、skip_folders、watcher注册、蓝光目录识别
// emit 接收发现的蓝光目录和视频文件，nil 表示只注册watcher不投递
func (c *collector) walkMediaDir(roots []string, videoType media_file.VideoType, emit func(mf *media_file.MediaFile)) {
	for _, root := range roots {
		if f, err := os.Stat(root); err != nil || !f.IsDir() {
			utils.Logger.WarningF("%s is not a directory", root)
			continue
		}

		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			// 隐藏文件跳过自身，隐藏目录跳过整个子树
			// 文件不能返回SkipDir，否则会跳过同目录剩余文件
			if d.Name()[0:1] == "." {
				if d.IsDir() {
					return fs.SkipDir
				}
				return nil
			}

			if d.IsDir() {
				if c.skipFolders(path, d.Name()) {
					utils.Logger.DebugF("skip folder by config: %s", d.Name())
					return fs.SkipDir
				}
				c.watcher.Add(path)
			}

			mf := media_file.NewMediaFile(path, d.Name(), videoType)

			// 蓝光目录只监听目录本身，不下探内部结构
			if mf.IsBluRay() {
				if emit != nil {
					emit(mf)
				}
				return fs.SkipDir
			}

			if emit != nil && mf.IsVideo() {
				emit(mf)
			}

			return nil
		})

		if err != nil {
			utils.Logger.WarningF("walk dir %s error: %s", root, err)
		}
	}
}

// skipFolders 检查是否跳过目录
func (c *collector) skipFolders(path, filename string) bool {
	base := filepath.Base(path)
	return slices.Contains(config.Collector.SkipFolders, base) ||
		slices.Contains(config.Collector.SkipFolders, filename)
}
