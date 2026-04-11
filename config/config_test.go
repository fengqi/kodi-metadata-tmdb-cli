package config

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func writeTestConfig(t *testing.T, runMode, moviesNfoMode, logMode, logLevel int) string {
	t.Helper()

	content := fmt.Sprintf(`{
  "log": {"mode": %d, "level": %d, "file": "./test.log"},
  "tmdb": {"api_host": "https://api.themoviedb.org", "api_key": "k", "image_host": "https://image.tmdb.org", "language": "zh-CN", "rating": "US", "proxy": ""},
  "collector": {
    "run_mode": %d,
    "watcher": true,
    "cron_seconds": 60,
    "cron_scan_kodi": false,
    "filter_tmp_suffix": true,
    "tmp_suffix": [".part"],
    "nfo_field": {"tag": true, "genre": true},
    "skip_folders": [],
    "skip_keywords": [],
    "movies_nfo_mode": %d,
    "movies_dir": ["./movies/../movies"],
    "shows_dir": ["./shows/../shows"],
    "music_videos_dir": []
  },
  "kodi": {"enable": false, "clean_library": false, "json_rpc": "", "timeout": 1, "username": "", "password": ""},
  "ffmpeg": {"max_worker": 1, "ffmpeg_path": "ffmpeg", "ffprobe_path": "ffprobe"}
}`, logMode, logLevel, runMode, moviesNfoMode)

	file := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(file, []byte(content), 0644); err != nil {
		t.Fatalf("write config file: %v", err)
	}
	return file
}

func TestLoadConfig_InvalidEnumFallback(t *testing.T) {
	configFile := writeTestConfig(t, 99, 99, 99, 99)
	LoadConfig(configFile, 0)

	if Collector.RunMode != CollectorRunModeDaemon {
		t.Fatalf("Collector.RunMode = %d; want %d", Collector.RunMode, CollectorRunModeDaemon)
	}
	if Collector.MoviesNfoMode != CollectorMoviesNfoModeVideoNfo {
		t.Fatalf("Collector.MoviesNfoMode = %d; want %d", Collector.MoviesNfoMode, CollectorMoviesNfoModeVideoNfo)
	}
	if Log.Mode != LogModeStdout {
		t.Fatalf("Log.Mode = %d; want %d", Log.Mode, LogModeStdout)
	}
	if Log.Level != LogLevelInfo {
		t.Fatalf("Log.Level = %d; want %d", Log.Level, LogLevelInfo)
	}
}

func TestLoadConfig_ValidEnumsKeepValue(t *testing.T) {
	configFile := writeTestConfig(t, CollectorRunModeSpec, CollectorMoviesNfoModeMovieNfo, LogModeBoth, LogLevelFatal)
	LoadConfig(configFile, 0)

	if Collector.RunMode != CollectorRunModeSpec {
		t.Fatalf("Collector.RunMode = %d; want %d", Collector.RunMode, CollectorRunModeSpec)
	}
	if Collector.MoviesNfoMode != CollectorMoviesNfoModeMovieNfo {
		t.Fatalf("Collector.MoviesNfoMode = %d; want %d", Collector.MoviesNfoMode, CollectorMoviesNfoModeMovieNfo)
	}
	if Log.Mode != LogModeBoth {
		t.Fatalf("Log.Mode = %d; want %d", Log.Mode, LogModeBoth)
	}
	if Log.Level != LogLevelFatal {
		t.Fatalf("Log.Level = %d; want %d", Log.Level, LogLevelFatal)
	}
}
