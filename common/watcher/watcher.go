package watcher

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"github.com/fsnotify/fsnotify"
	"os"
	"path/filepath"
)

type callback func(filename string, fileInfo os.FileInfo)

type Watcher struct {
	name     string
	watcher  *fsnotify.Watcher
	callback callback
}

func InitWatcher(taskName string) *Watcher {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		utils.Logger.FatalF("new %s watcher err: %v", taskName, err)
	}

	w := &Watcher{
		name:    taskName,
		watcher: watcher,
	}

	return w
}

func (w *Watcher) Run(callback callback) {
	w.callback = callback
	go w.runWatcher()
}

func (w *Watcher) runWatcher() {
	if !config.Collector.Watcher || w.callback == nil {
		return
	}

	utils.Logger.DebugF("run %s watcher", w.name)

	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				continue
			}

			if !event.Has(fsnotify.Create) || w.skipFolders(filepath.Dir(event.Name), event.Name) {
				continue
			}

			utils.Logger.InfoF("created file: %s", event.Name)

			fileInfo, err := os.Stat(event.Name)
			if fileInfo == nil || err != nil {
				utils.Logger.WarningF("get file: %s stat err: %v", event.Name, err)
				continue
			}

			w.callback(event.Name, fileInfo)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}

			utils.Logger.ErrorF("%s watcher error: %v", w.name, err)
		}
	}
}

func (w *Watcher) Add(name string) {
	if !config.Collector.Watcher {
		return
	}

	utils.Logger.DebugF("add dir: %s to %s watcher", name, w.name)

	err := w.watcher.Add(name)
	if err != nil {
		utils.Logger.FatalF("add dir: %s to %s watcher err: %v", name, w.name, err)
	}
}

// todo 代码复用
func (w *Watcher) skipFolders(path, filename string) bool {
	base := filepath.Base(path)
	for _, item := range config.Collector.SkipFolders {
		if filename[0:1] == "." || item == base || item == filename {
			return true
		}
	}
	return false
}
