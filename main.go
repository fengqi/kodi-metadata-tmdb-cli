package main

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/movies"
	"fengqi/kodi-metadata-tmdb-cli/music"
	"fengqi/kodi-metadata-tmdb-cli/shows"
	"fengqi/kodi-metadata-tmdb-cli/stock"
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
	utils.InitLogger(c.LogLevel, c.LogFile)
	tmdb.InitTmdb(c)

	wg := &sync.WaitGroup{}
	wg.Add(3)
	go shows.RunCollector(c)
	go movies.RunCollector(c)
	go stock.RunCollector(c)
	go music.RunCollector(c)
	wg.Wait()
}
