package main

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/kodi"
	"fengqi/kodi-metadata-tmdb-cli/movies"
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
	utils.InitLogger(c.LogLevel, c.LogFile)
	tmdb.InitTmdb(c)
	kodi.InitKodi(c.Kodi)

	wg := &sync.WaitGroup{}
	wg.Add(2)
	go shows.RunCollector(c)
	go movies.RunCollector(c)
	wg.Wait()
}
