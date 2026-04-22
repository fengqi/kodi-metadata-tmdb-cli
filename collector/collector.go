package collector

import (
	"fengqi/kodi-metadata-tmdb-cli/common/watcher"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"sync"
)

type scanTask struct {
	file *media_file.MediaFile
	done *sync.WaitGroup
}

type collector struct {
	channel   chan *scanTask
	watcher   *watcher.Watcher
	closeOnce sync.Once
	scanMu    sync.Mutex
}

var ins *collector

// Run 运行扫描
func Run() {
	ins = &collector{
		channel: make(chan *scanTask, 100),
		watcher: watcher.InitWatcher("collector"),
	}

	if config.Collector.RunMode == config.CollectorRunModeOnce || config.Collector.RunMode == config.CollectorRunModeSpec {
		go ins.runScan()
	} else {
		go ins.watcher.Run(ins.watcherCallback)
		go ins.runScan()
		go ins.runCronScan()
	}

	ins.runProcess()
}
