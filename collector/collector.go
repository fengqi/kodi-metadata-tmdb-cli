package collector

import (
	"fengqi/kodi-metadata-tmdb-cli/common/watcher"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"sync"
	"time"
)

type collector struct {
	channel chan *media_file.MediaFile
	watcher *watcher.Watcher
	wg      *sync.WaitGroup
}

var ins *collector

// Run 运行扫描
func Run() {
	ins = &collector{
		channel: make(chan *media_file.MediaFile, 100),
		watcher: watcher.InitWatcher("collector"),
		wg:      &sync.WaitGroup{},
	}

	go ins.watcher.Run(ins.watcherCallback)
	go ins.runScan()
	go ins.runProcess()

	ticker := time.NewTicker(time.Second * time.Duration(config.Collector.CronSeconds))
	for range ticker.C {
		ins.runScan()
	}
}
