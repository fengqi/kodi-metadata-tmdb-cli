package collector

import (
	"fengqi/kodi-metadata-tmdb-cli/common/watcher"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"sync"
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

	if config.Collector.RunMode == config.CollectorRunModeOnce || config.Collector.RunMode == config.CollectorRunModeSpec {
		go ins.runScan()
		ins.runProcess()
		return
	}

	go ins.watcher.Run(ins.watcherCallback)
	go ins.runScan()
	go ins.runCronScan()

	ins.runProcess()
}
