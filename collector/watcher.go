package collector

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"os"
	"strings"
)

// watcherCallback 监听文件变化的回调函数
// todo 部分逻辑和scanDir重复，考虑复用
func (c *collector) watcherCallback(filename string, fi os.FileInfo) {
	if fi.Name()[0:1] == "." {
		return
	}

	if fi.IsDir() {
		if c.skipFolders(filename, fi.Name()) {
			utils.Logger.DebugF("skip folder by config: %s", fi.Name())
			return
		}

		c.watcher.Add(filename)
	}

	var videoType media_file.VideoType
	for _, item := range config.Collector.MoviesDir {
		if strings.HasPrefix(filename, item) {
			videoType = media_file.Movies
			break
		}
	}
	for _, item := range config.Collector.ShowsDir {
		if strings.HasPrefix(filename, item) {
			videoType = media_file.TvShows
			break
		}
	}
	for _, item := range config.Collector.MusicVideosDir {
		if strings.HasPrefix(filename, item) {
			videoType = media_file.MusicVideo
			break
		}
	}

	mf := media_file.NewMediaFile(filename, fi.Name(), videoType)
	if mf.IsBluRay() {
		c.wg.Add(1)
		c.channel <- mf
	}
	if mf.IsVideo() {
		c.wg.Add(1)
		c.channel <- mf
	}
}
