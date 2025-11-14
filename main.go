package main

import (
	"fengqi/kodi-metadata-tmdb-cli/collector"
	"fengqi/kodi-metadata-tmdb-cli/common/memcache"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/ffmpeg"
	"fengqi/kodi-metadata-tmdb-cli/kodi"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"flag"
	"fmt"
	"runtime"
)

var (
	configFile   string
	version      bool
	runMode      int
	buildVersion = "dev-master"
)

func init() {
	flag.StringVar(&configFile, "config", "./config.json", "config file")
	flag.BoolVar(&version, "version", false, "display version")
	flag.IntVar(&runMode, "mode", 0, "run mode: 1: daemon, 2: once, 3: spec")
	flag.Parse()
}

func main() {
	if version {
		fmt.Printf("version: %s, build with: %s\n", buildVersion, runtime.Version())
		return
	}

	config.LoadConfig(configFile, runMode)

	memcache.InitCache()
	utils.InitLogger()
	tmdb.InitTmdb()
	kodi.InitKodi()
	ffmpeg.InitFfmpeg()

	collector.Run()
}
