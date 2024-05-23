package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

var watcher *fsnotify.Watcher

// todo 封装+回调，可以电视剧、电影、普通视频复用代码
func (c *Collector) initWatcher() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		utils.Logger.FatalF("new shows watcher err: %v", err)
	}
}

// 目录监听，新增的增加到队列，删除的移除监听
func (c *Collector) runWatcher() {
	if !c.config.Collector.Watcher {
		return
	}

	utils.Logger.Debug("run shows watcher")

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				continue
			}

			fileInfo, err := os.Stat(event.Name)
			if fileInfo == nil || err != nil {
				utils.Logger.WarningF("get shows stat err: %v", err)
				continue
			}

			// 根目录电视剧不允许以单文件的形式存在
			if !fileInfo.IsDir() && utils.InArray(c.config.Collector.ShowsDir, filepath.Dir(event.Name)) {
				utils.Logger.WarningF("shows file not allow root: %s", event.Name)
				continue
			}

			// 删除文件夹
			if event.Has(fsnotify.Remove) && fileInfo.IsDir() {
				utils.Logger.InfoF("removed dir: %s", event.Name)

				err := watcher.Remove(event.Name)
				if err != nil {
					utils.Logger.WarningF("remove shows watcher: %s error: %v", event.Name, err)
				}
				continue
			}

			// 新增文件夹
			if event.Has(fsnotify.Create) && fileInfo.IsDir() {
				utils.Logger.InfoF("created dir: %s", event.Name)

				showsDir := c.parseShowsDir(filepath.Dir(event.Name), fileInfo)
				if showsDir != nil {
					c.dirChan <- showsDir
				}

				c.watchDir(event.Name)
				continue
			}

			// 新增剧集文件
			if event.Has(fsnotify.Create) && utils.IsVideo(event.Name) != "" {
				utils.Logger.InfoF("created file: %s", event.Name)

				filePath := filepath.Dir(event.Name)
				dirInfo, _ := os.Stat(filePath)
				dir := c.parseShowsDir(filepath.Dir(filePath), dirInfo)
				if dir != nil {
					c.dirChan <- dir
				}
			}

		case err, ok := <-watcher.Errors:
			utils.Logger.ErrorF("shows watcher error: %v", err)

			if !ok {
				return
			}
		}
	}
}

// 监听目录
// todo 判断是否是目录
func (c *Collector) watchDir(name string) {
	if !c.config.Collector.Watcher {
		return
	}

	utils.Logger.DebugF("add shows dir: %s to watcher", name)

	err := watcher.Add(name)
	if err != nil {
		utils.Logger.FatalF("add shows dir: %s to watcher err: %v", name, err)
	}
}
