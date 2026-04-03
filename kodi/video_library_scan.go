package kodi

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"strings"
	"time"
)

// AddScanTaskByName 通过名称添加扫描任务，Type 1 为电视剧，2 为电影
func (r *JsonRpc) AddScanTaskByName(Type int, name string) {
	if !config.Kodi.Enable {
		return
	}

	if Type == 1 {
		kodiShowsResp := r.VideoLibrary.GetTVShowsByField("title", "contains", name)
		if kodiShowsResp == nil || kodiShowsResp.Limits.Total == 0 {
			return
		}
		r.AddScanTask(kodiShowsResp.TvShows[0].File)
	}
}

func (r *JsonRpc) AddScanTask(directory string) {
	if !config.Kodi.Enable {
		return
	}

	utils.Logger.DebugF("AddScanTask %s", directory)
	if directory != "" {
		sources := r.Files.GetSources("video")
		if sources != nil {
			for _, item := range sources {
				if strings.Contains(item.File, directory) || strings.Contains(directory, item.File) {
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
	if !config.Kodi.Enable {
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

		for directory := range r.scanQueue {
			r.scanLock.Lock()

			r.VideoLibrary.Scan(directory, true)

			delete(r.scanQueue, directory)
			r.scanLock.Unlock()
		}
	}
}
