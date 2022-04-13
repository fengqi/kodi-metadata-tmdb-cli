package music_videos

import (
	"fengqi/kodi-metadata-tmdb-cli/ffmpeg"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"os"
	"path/filepath"
	"strings"
)

func (m *MusicVideo) saveToNfo() error {
	nfo := m.getNfoFile()

	if info, err := os.Stat(nfo); err == nil && info.Size() > 0 {
		return nil
	}

	utils.Logger.InfoF("save %s to %s", m, nfo)

	top := MusicVideoNfo{
		Title:     m.Title,
		DateAdded: m.DateAdded,
		FileInfo:  m.coverNfoFileInfo(),
		Thumb: []Thumb{
			{
				Aspect:  "thumb",
				Preview: m.Title + ".jpg",
			},
		},
		Poster: m.Title + ".jpg",
	}

	return utils.SaveNfo(m.getNfoFile(), top)
}

func (m *MusicVideo) coverNfoFileInfo() *FileInfo {
	probe, err := ffmpeg.Probe(m.Dir + "/" + m.OriginTitle)
	if err != nil {
		utils.Logger.WarningF("parse %s probe err: %v", m.OriginTitle, err)
		return nil
	}

	if probe == nil {
		return nil
	}

	audio := probe.FirstAudioStream()
	video := probe.FirstVideoStream()

	return &FileInfo{
		StreamDetails: StreamDetails{
			Video: []Video{
				{
					Codec:             video.CodecName,
					Aspect:            video.DisplayAspectRatio,
					Width:             video.Width,
					Height:            video.Height,
					DurationInSeconds: video.Duration,
					StereoMode:        "progressive",
				},
			},
			Audio: []Audio{
				{
					Language: "zho",
					Codec:    audio.CodecName,
					Channels: audio.Channels,
				},
			},
		},
	}
}

func (m *MusicVideo) drawThumb() error {
	thumb := m.getNfoThumb()
	if info, err := os.Stat(thumb); err == nil && info.Size() > 0 {
		return nil
	}

	filename := m.getFullPath()
	ss := "00:00:30"
	base := filepath.Base(filename)
	if (len(base) > 2 && base[0:2] == "03") || (len(base) > 5 && strings.ToLower(base[0:5]) == "heyzo") {
		ss = "00:01:10"
	}

	utils.Logger.InfoF("draw thumb start: %s, %s to %s", ss, m, thumb)

	err := ffmpeg.Frame(filename, thumb, "-ss", ss)
	if err != nil {
		utils.Logger.WarningF("draw thumb err: %v", err)
		panic(err)
	}

	return nil
}
