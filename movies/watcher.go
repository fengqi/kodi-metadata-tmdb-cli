package movies

import (
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
		utils.Logger.FatalF("new movies watcher err: %v", err)
	}
}

// 开启文件夹监听
func (c *Collector) runWatcher() {
	if !c.config.Collector.Watcher {
		return
	}

	utils.Logger.Debug("run movies watcher")

	for {
		select {
		// 接受事件，增删改查都会收到，需要过滤，部分情况下可能收不到create而是chmod
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}

			if !event.Has(fsnotify.Create) {
				continue
			}

			fileInfo, _ := os.Stat(event.Name)
			if fileInfo == nil || (!fileInfo.IsDir() && utils.IsVideo(event.Name) == "") {
				continue
			}

			utils.Logger.InfoF("created file: %s", event.Name)

			moviesDir := parseMoviesDir(filepath.Dir(event.Name), fileInfo)
			if moviesDir != nil {
				c.channel <- moviesDir
			}

		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}

			utils.Logger.ErrorF("movies watcher error: %v", err)
		}
	}
}

// 监听目录
func (c *Collector) watchDir(name string) {
	if !c.config.Collector.Watcher {
		return
	}

	utils.Logger.DebugF("add movies dir: %s to watcher", name)

	err := watcher.Add(name)
	if err != nil {
		utils.Logger.FatalF("add movies dir: %s to watcher err: %v", name, err)
	}
}
