package music_videos

import (
	"fengqi/kodi-metadata-tmdb-cli/ffmpeg"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"path/filepath"
	"strconv"
	"strings"
)

func (m *MusicVideo) saveToNfo() error {
	if m.NfoExist() {
		return nil
	}

	nfo := m.getNfoFile()
	utils.Logger.InfoF("save %s to %s", m.OriginTitle, nfo)

	fileInfo := &FileInfo{
		StreamDetails: StreamDetails{
			Video: []Video{
				{
					Codec:             m.VideoStream.CodecName,
					Aspect:            m.VideoStream.DisplayAspectRatio,
					Width:             m.VideoStream.Width,
					Height:            m.VideoStream.Height,
					DurationInSeconds: m.VideoStream.Duration,
					StereoMode:        "progressive",
				},
			},
			Audio: []Audio{
				{
					Language: "zho",
					Codec:    m.VideoStream.CodecName,
					Channels: m.VideoStream.Channels,
				},
			},
		},
	}

	top := MusicVideoNfo{
		Title:     m.Title,
		DateAdded: m.DateAdded,
		FileInfo:  fileInfo,
		Thumb: []Thumb{
			{
				Aspect:  "thumb",
				Preview: m.Title + ".jpg",
			},
		},
		Poster: m.Title + ".jpg",
	}

	return utils.SaveNfo(nfo, top)
}

func (m *MusicVideo) drawThumb() error {
	if m.ThumbExist() {
		return nil
	}

	ss := "00:00:00"
	second, _ := strconv.ParseFloat(m.VideoStream.Duration, 10)
	if m.VideoStream != nil && second > 30 {
		ss = "00:00:30"
	}

	filename := m.getFullPath()
	base := filepath.Base(filename)
	if (len(base) > 2 && base[0:2] == "03") || (len(base) > 5 && strings.ToLower(base[0:5]) == "heyzo") {
		ss = "00:01:10"
	}

	thumb := m.getNfoThumb()
	utils.Logger.InfoF("draw thumb start: %s, %s to %s", ss, m.OriginTitle, thumb)

	err := ffmpeg.Frame(filename, thumb, "-ss", ss)
	if err != nil {
		utils.Logger.WarningF("draw thumb err: %v", err)
	}

	return err
}
