package movies

import (
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"os"
	"path/filepath"
)

func (c *Collector) watcherCallback(filename string, fileInfo os.FileInfo) {
	if fileInfo == nil || (!fileInfo.IsDir() && utils.IsVideo(filename) == "") {
		return
	}

	moviesDir := parseMoviesDir(filepath.Dir(filename), fileInfo)
	if moviesDir != nil {
		c.channel <- moviesDir
	}
}
