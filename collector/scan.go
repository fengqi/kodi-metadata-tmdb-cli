package collector

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/kodi"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

// collector 运行扫描
func (c *collector) runScan() {
	wg := sync.WaitGroup{}
	wg.Add(3)
	go c.scanDir(config.Collector.MoviesDir, media_file.Movies)
	go c.scanDir(config.Collector.ShowsDir, media_file.TvShows)
	go c.scanDir(config.Collector.MusicVideosDir, media_file.MusicVideo)
	wg.Wait()

	// 扫描完成，通知kodi刷新媒体库
	if config.Collector.CronScanKodi {
		kodi.Rpc.VideoLibrary.Scan("", false)
	}

	// 扫描完成后，通知kodi清理媒体库
	if config.Kodi.CleanLibrary {
		kodi.Rpc.AddCleanTask("")
	}
}

// scanDir 扫描目录
func (c *collector) scanDir(roots []string, videoType media_file.VideoType) {
	for _, root := range roots {
		if f, err := os.Stat(root); err != nil || !f.IsDir() {
			utils.Logger.WarningF("%s is not a directory", root)
			return
		}

		c.watcher.Add(root)

		filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.Name()[0:1] == "." {
				return fs.SkipDir
			}

			if d.IsDir() {
				if c.skipFolders(path, d.Name()) {
					utils.Logger.DebugF("skip folder by config: %s", d.Name())
					return fs.SkipDir
				}
			}

			mf := media_file.NewMediaFile(path, d.Name(), videoType)
			if mf.IsBluRay() {
				c.channel <- mf
				return fs.SkipDir
			}
			if mf.IsVideo() {
				c.channel <- mf
			}

			return nil
		})
	}
}

// skipFolders 检查是否跳过目录
func (c *collector) skipFolders(path, filename string) bool {
	base := filepath.Base(path)
	return utils.InArray(config.Collector.SkipFolders, base) ||
		utils.InArray(config.Collector.SkipFolders, filename)
}
