package music_videos

import "os"

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
