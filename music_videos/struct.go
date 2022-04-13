package music_videos

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"github.com/fsnotify/fsnotify"
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
