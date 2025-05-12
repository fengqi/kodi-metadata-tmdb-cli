package collector

import (
	"fengqi/kodi-metadata-tmdb-cli/common/watcher"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"time"
)

type collector struct {
	channel chan *media_file.MediaFile
	watcher *watcher.Watcher
}

var ins *collector

// Run 运行扫描
func Run() {
	ins = &collector{
		channel: make(chan *media_file.MediaFile, 100),
		watcher: watcher.InitWatcher("collector"),
	}

	go ins.watcher.Run(ins.watcherCallback)
	go ins.runScan()
	go ins.runProcess()

	ticker := time.NewTicker(time.Second * time.Duration(config.Collector.CronSeconds))
	for {
		select {
		case <-ticker.C:
			ins.runScan()
		}
	}
}
