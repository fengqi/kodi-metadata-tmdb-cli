package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"os"
	"path/filepath"
)

func (c *Collector) watcherCallback(filename string, fileInfo os.FileInfo) {
	// 根目录电视剧不允许以单文件的形式存在
	if !fileInfo.IsDir() && utils.InArray(config.Collector.ShowsDir, filepath.Dir(filename)) {
		utils.Logger.WarningF("shows file not allow root: %s", filename)
		return
	}

	// 新增文件夹
	if fileInfo.IsDir() {
		utils.Logger.InfoF("created dir: %s", filename)

		showsDir := c.parseShowsDir(filepath.Dir(filename), fileInfo)
		if showsDir != nil {
			c.dirChan <- showsDir
		}

		c.watcher.Add(filename)
		return
	}

	// 新增剧集文件
	if utils.IsVideo(filename) != "" {
		utils.Logger.InfoF("created file: %s", filename)

		filePath := filepath.Dir(filename)
		dirInfo, _ := os.Stat(filePath)
		dir := c.parseShowsDir(filepath.Dir(filePath), dirInfo)
		if dir != nil {
			c.dirChan <- dir
		}
	}
}
