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
		go ins.runDaemon()
	}

	ins.runProcess()
}

// runDaemon 守护进程模式启动流程
func (c *collector) runDaemon() {
	c.watcher.Run(c.watcherCallback)

	// watcher 目录注册不依赖扫描，定时扫描和启动扫描都关闭时 watcher 依然生效
	c.registerWatcherDirs()

	// 启动后立即执行一次扫描
	if config.Collector.CronScanBoot {
		c.runScan()
	}

	c.runCronScan()
}
