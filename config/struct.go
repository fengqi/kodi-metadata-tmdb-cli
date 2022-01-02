package config

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io/ioutil"
)

type Config struct {
	LogLevel int    `json:"log_level"`
	LogFile  string `json:"log_file"`

	Rating   string `json:"rating"`   // 内容分级
	ApiKey   string `json:"api_key"`  // api key
	Language string `json:"language"` // 语言

	ShowsDir    []string `json:"shows_dir"`
	MoviesDir   []string `json:"movies_dir"`
	StockDir    []string `json:"stock_dir"`
	MusicDir    []string `json:"music_dir"`
	CronSeconds int      `json:"cron_seconds"` // todo、shows、movies 分别设置
}

func LoadConfig(file string) *Config {
	bytes, err := ioutil.ReadFile(file)
	if err != nil {
		utils.Logger.FatalF("load config err: %v", err)
	}

	c := &Config{}
	err = json.Unmarshal(bytes, c)
	if err != nil {
		utils.Logger.FatalF("parse config err: %v", err)
	}

	return c
}
