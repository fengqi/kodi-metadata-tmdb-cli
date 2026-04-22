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
	"strings"
	"sync"
	"time"

	"github.com/fengqi/lrace"
)

// collector 运行扫描
func (c *collector) runScan() {
	c.scanMu.Lock()
	defer c.scanMu.Unlock()

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

// scanDir 扫描目录
func (c *collector) scanDir(roots []string, videoType media_file.VideoType, producerWG, scanTaskWG *sync.WaitGroup) {
	defer producerWG.Done()

	for _, root := range roots {
		if f, err := os.Stat(root); err != nil || !f.IsDir() {
			utils.Logger.WarningF("%s is not a directory", root)
			continue
		}

		err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.Name()[0:1] == "." {
				return nil
			}

			if d.IsDir() {
				if c.skipFolders(path, d.Name()) {
					utils.Logger.DebugF("skip folder by config: %s", d.Name())
					return fs.SkipDir
				}

				c.watcher.Add(path) // todo 定时执行，等于会重复Add，不确定有没有问题，后续确认
			}

			mf := media_file.NewMediaFile(path, d.Name(), videoType)
			if mf.IsBluRay() {
				scanTaskWG.Add(1)
				c.channel <- &scanTask{file: mf, done: scanTaskWG}
				return fs.SkipDir
			}
			if mf.IsVideo() {
				scanTaskWG.Add(1)
				c.channel <- &scanTask{file: mf, done: scanTaskWG}
			}

			return nil
		})

		if err != nil {
			utils.Logger.WarningF("scan dir %s error: %s", root, err)
		}
	}
}

// skipFolders 检查是否跳过目录
func (c *collector) skipFolders(path, filename string) bool {
	base := filepath.Base(path)
	return lrace.InArray(config.Collector.SkipFolders, base) ||
		lrace.InArray(config.Collector.SkipFolders, filename)
}
