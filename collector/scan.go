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
	"time"

	"github.com/fengqi/lrace"
)

// collector 运行扫描
func (c *collector) runScan() {

	if config.Collector.RunMode == 3 {
		pwd, err := os.Getwd()
		if err != nil {
			utils.Logger.ErrorF("get pwd error: %s", err)
			return
		}

		for _, item := range config.Collector.MoviesDir {
			if strings.HasPrefix(pwd, item) {
				go c.scanDir([]string{pwd}, media_file.Movies)
				break
			}
		}
		for _, item := range config.Collector.ShowsDir {
			if strings.HasPrefix(pwd, item) {
				go c.scanDir([]string{pwd}, media_file.TvShows)
				break
			}
		}
		for _, item := range config.Collector.MusicVideosDir {
			if strings.HasPrefix(pwd, item) {
				go c.scanDir([]string{pwd}, media_file.MusicVideo)
				break
			}
		}

	} else {
		go c.scanDir(config.Collector.MoviesDir, media_file.Movies)
		go c.scanDir(config.Collector.ShowsDir, media_file.TvShows)
		go c.scanDir(config.Collector.MusicVideosDir, media_file.MusicVideo)
	}

	time.Sleep(time.Second * 3)
	c.wg.Wait()

	// 扫描完成，通知kodi刷新媒体库
	if config.Collector.CronScanKodi {
		log.Println("scan done, refresh kodi library")
		kodi.Rpc.VideoLibrary.Scan("", false)
	}

	// 扫描完成后，通知kodi清理媒体库
	if config.Kodi.CleanLibrary {
		log.Println("scan done, clean kodi library")
		kodi.Rpc.AddCleanTask("")
	}

	// 单次模式，关闭channel
	if config.Collector.RunMode == 2 || config.Collector.RunMode == 3 {
		close(c.channel)
	}
}

// scanDir 扫描目录
func (c *collector) scanDir(roots []string, videoType media_file.VideoType) {
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
				c.wg.Add(1)
				c.channel <- mf
				return fs.SkipDir
			}
			if mf.IsVideo() {
				c.wg.Add(1)
				c.channel <- mf
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
