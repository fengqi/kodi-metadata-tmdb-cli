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

	MoviesNfoMode int      `json:"movies_nfo_mode"` // 电影NFO写入模式：1 movie.nfo， 2 <VideoFileName>.nfo
	ShowsDir      []string `json:"shows_dir"`
	MoviesDir     []string `json:"movies_dir"`
	CronSeconds   int      `json:"cron_seconds"` // todo、shows、movies 分别设置

	Kodi KodiConfig `json:"kodi"`
}

type KodiConfig struct {
	JsonRpc  string `json:"json_rpc"`
	Timeout  int    `json:"timeout"`
	Username string `json:"username"`
	Password string `json:"password"`
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
