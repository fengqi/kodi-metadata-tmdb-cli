package music_videos

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

var watcher *fsnotify.Watcher

func (c *Collector) initWatcher() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		utils.Logger.FatalF("new videos watcher err: %v", err)
	}
}

// 运行文件变动监听
func (c *Collector) runWatcher() {
	if !config.Collector.Watcher {
		return
	}

	utils.Logger.Debug("run videos watcher")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				continue
			}

			fileInfo, err := os.Stat(event.Name)
			if fileInfo == nil || err != nil {
				utils.Logger.WarningF("get videos stat err: %v", err)
				continue
			}

			// 删除文件夹
			if event.Has(fsnotify.Remove) && fileInfo.IsDir() {
				utils.Logger.InfoF("removed dir: %s", event.Name)

				err := watcher.Remove(event.Name)
				if err != nil {
					utils.Logger.WarningF("remove videos watcher: %s error: %v", event.Name, err)
				}
				continue
			}

			if !event.Has(fsnotify.Create) || c.skipFolders(filepath.Dir(event.Name), event.Name) {
				continue
			}

			//  新增目录
			if fileInfo.IsDir() {
				err = watcher.Add(event.Name)
				if err != nil {
					utils.Logger.WarningF("add video dir: %s to watcher err: %v", event.Name, err)
				}

				videos, err := c.scanDir(event.Name)
				if err != nil || len(videos) == 0 {
					utils.Logger.WarningF("new dir %s scan err: %v or no videos", event.Name, err)
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

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			utils.Logger.ErrorF("videos watcher error: %v", err)
		}
	}
}

// 监听目录
func (c *Collector) watchDir(name string) {
	if !config.Collector.Watcher {
		return
	}

	utils.Logger.DebugF("add videos dir: %s to watcher", name)

	err := watcher.Add(name)
	if err != nil {
		utils.Logger.FatalF("add videos dir: %s to watcher err: %v", name, err)
	}
}
