package config

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io/ioutil"
)

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
