package music_videos

import (
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"os"
	"path/filepath"
)

func (c *Collector) watcherCallback(filename string, fileInfo os.FileInfo) {
	//  新增目录
	if fileInfo.IsDir() {
		c.watcher.Add(filename)

		videos, err := c.scanDir(filename)
		if err != nil || len(videos) == 0 {
			utils.Logger.WarningF("new dir %s scan err: %v or no videos", filename, err)
			return
		}

		for _, video := range videos {
			c.channel <- video
		}

		return
	}

	// 单个文件
	if utils.IsVideo(filename) != "" {
		video := c.parseVideoFile(filepath.Dir(filename), fileInfo)
		c.channel <- video
	}
}
