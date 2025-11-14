package config

import (
	"encoding/json"
	"log"
	"os"

	"github.com/fengqi/lrace"
)

var (
	Log       *LogConfig
	Ffmpeg    *FfmpegConfig
	Tmdb      *TmdbConfig
	Kodi      *KodiConfig
	Collector *CollectorConfig
)

func LoadConfig(file string, runMode int) {
	bytes, err := os.ReadFile(file)
	if err != nil {
		log.Fatalf("load config err: %v", err)
	}

	c := &Config{}
	err = json.Unmarshal(bytes, c)
	if err != nil {
		log.Fatalf("parse config err: %v", err)
	}

	Log = c.Log
	Ffmpeg = c.Ffmpeg
	Tmdb = c.Tmdb
	Kodi = c.Kodi
	Collector = c.Collector

	Collector.RunMode = lrace.Ternary(Collector.RunMode == 0, 1, Collector.RunMode)
	Collector.RunMode = lrace.Ternary(runMode > 0, runMode, Collector.RunMode)
}
