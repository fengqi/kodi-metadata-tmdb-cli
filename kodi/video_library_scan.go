package kodi

import (
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"path/filepath"
	"strings"
	"time"
)

func (r *JsonRpc) AddScanTask(directory string) {
	if !r.config.Enable {
		return
	}

	utils.Logger.DebugF("AddScanTask %s", directory)
	if directory != "" {
		directory = filepath.Clean(directory)
		sources := r.Files.GetSources("video")
		if sources != nil {
			for _, item := range sources {
				if strings.Contains(item.File, directory) {
					directory = item.File
					break
				}
			}
		}
	}

	r.scanLock.Lock()
	defer r.scanLock.Unlock()

	if _, ok := r.scanQueue[directory]; !ok {
		r.scanQueue[directory] = struct{}{}
	}

	return
}

func (r *JsonRpc) ConsumerScanTask() {
	if !r.config.Enable {
		return
	}

	for {
		if len(r.scanQueue) == 0 || !r.Ping() || r.VideoLibrary.IsScanning() {
			time.Sleep(time.Second * 30)
			continue
		}

		if !r.VideoLibrary.scanLimiter.take() {
			time.Sleep(time.Second * 30)
			continue
		}

		for directory, _ := range r.scanQueue {
			r.scanLock.Lock()

			r.VideoLibrary.Scan(directory, true)

			delete(r.scanQueue, directory)
			r.scanLock.Unlock()
		}
	}
}
