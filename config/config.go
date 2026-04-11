package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

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

	Collector.RunMode = lrace.Ternary(Collector.RunMode == 0, CollectorRunModeDaemon, Collector.RunMode)
	Collector.RunMode = lrace.Ternary(runMode > 0, runMode, Collector.RunMode)
	validateConfigEnums()
	Collector.ShowsDir = clearPath(Collector.ShowsDir)
	Collector.MoviesDir = clearPath(Collector.MoviesDir)
}

// validateConfigEnums validates and normalizes configuration enum values for Collector and Log to ensure compatibility.
func validateConfigEnums() {
	if Collector != nil {
		if !inIntSet(Collector.RunMode, CollectorRunModeDaemon, CollectorRunModeOnce, CollectorRunModeSpec) {
			log.Printf("invalid collector.run_mode=%d, fallback to %d", Collector.RunMode, CollectorRunModeDaemon)
			Collector.RunMode = CollectorRunModeDaemon
		}

		if !inIntSet(Collector.MoviesNfoMode, CollectorMoviesNfoModeMovieNfo, CollectorMoviesNfoModeVideoNfo) {
			log.Printf("invalid collector.movies_nfo_mode=%d, fallback to %d", Collector.MoviesNfoMode, CollectorMoviesNfoModeVideoNfo)
			Collector.MoviesNfoMode = CollectorMoviesNfoModeVideoNfo
		}
	}

	if Log != nil {
		if !inIntSet(Log.Mode, LogModeStdout, LogModeLogfile, LogModeBoth) {
			log.Printf("invalid log.mode=%d, fallback to %d", Log.Mode, LogModeStdout)
			Log.Mode = LogModeStdout
		}

		if !inIntSet(Log.Level, LogLevelDebug, LogLevelInfo, LogLevelWarning, LogLevelError, LogLevelFatal) {
			log.Printf("invalid log.level=%d, fallback to %d", Log.Level, LogLevelInfo)
			Log.Level = LogLevelInfo
		}
	}
}

func inIntSet(val int, set ...int) bool {
	for _, item := range set {
		if val == item {
			return true
		}
	}
	return false
}

func clearPath(name []string) []string {
	for i, item := range name {
		name[i] = filepath.Clean(item)
	}
	return name
}
