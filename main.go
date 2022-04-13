package main

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/ffmpeg"
	"fengqi/kodi-metadata-tmdb-cli/kodi"
	"fengqi/kodi-metadata-tmdb-cli/movies"
	"fengqi/kodi-metadata-tmdb-cli/music_videos"
	"fengqi/kodi-metadata-tmdb-cli/shows"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"flag"
	"sync"
)

var (
	configFile string
)

func init() {
	flag.StringVar(&configFile, "config", "./config.json", "config file")
	flag.Parse()
}

func main() {
	c := config.LoadConfig(configFile)

	utils.InitLogger(c.LogMode, c.LogLevel, c.LogFile)
	tmdb.InitTmdb(c)
	kodi.InitKodi(c.Kodi)
	ffmpeg.InitFfmpeg(c)

	wg := &sync.WaitGroup{}
	wg.Add(4)
	go kodi.Rpc.RunNotify()
	go shows.RunCollector(c)
	go movies.RunCollector(c)
	go music_videos.RunCollector(c)
	wg.Wait()
}
