package music_videos

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/ffmpeg"
	"github.com/fsnotify/fsnotify"
	"os"
)

type Collector struct {
	config  *config.Config
	watcher *fsnotify.Watcher
	channel chan *MusicVideo
}

type MusicVideo struct {
	Dir         string
	Title       string
	OriginTitle string
	DateAdded   string
	VideoStream *ffmpeg.Stream
	AudioStream *ffmpeg.Stream
}

func (m *MusicVideo) getFullPath() string {
	return m.Dir + "/" + m.OriginTitle
}

func (m *MusicVideo) getNfoThumb() string {
	return m.Dir + "/" + m.Title + "-thumb.jpg"
}

func (m *MusicVideo) getNfoFile() string {
	return m.Dir + "/" + m.Title + ".nfo"
}

func (m *MusicVideo) NfoExist() bool {
	nfo := m.getNfoFile()

	if info, err := os.Stat(nfo); err == nil && info.Size() > 0 {
		return true
	}

	return false
}

func (m *MusicVideo) ThumbExist() bool {
	thumb := m.getNfoThumb()
	if info, err := os.Stat(thumb); err == nil && info.Size() > 0 {
		return true
	}

	return false
}
