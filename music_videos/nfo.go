package music_videos

import (
	"fengqi/kodi-metadata-tmdb-cli/ffmpeg"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"github.com/fengqi/lrace"
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
				Preview: m.Title + "-thumb.jpg",
			},
		},
		Poster: m.Title + "-thumb.jpg",
	}

	return utils.SaveNfo(nfo, top)
}

// 缩略图提取
// TODO 截取开始位置可配置
func (m *MusicVideo) drawThumb() error {
	if m.ThumbExist() {
		return nil
	}

	// 如果有视频文件同名后缀的图片，尝试直接使用
	filename := m.getFullPath()
	thumb := m.getNfoThumb()
	for _, i := range ThumbImagesFormat {
		check := m.Dir + "/" + m.Title + "." + i
		if lrace.FileExist(check) {
			n, err := lrace.CopyFile(check, thumb)
			if n > 0 && err == nil {
				return nil
			}
		}
	}

	// 对于大文件，尝试偏移30秒，防止读到的是黑屏白屏或者logo
	ss := "00:00:00"
	second, _ := strconv.ParseFloat(m.VideoStream.Duration, 10)
	if m.VideoStream != nil && second > 30 {
		ss = "00:00:30"
	}

	base := filepath.Base(filename)
	if (len(base) > 2 && base[0:2] == "03") || (len(base) > 5 && strings.ToLower(base[0:5]) == "heyzo") {
		ss = "00:01:10"
	}

	utils.Logger.InfoF("draw thumb start: %s, %s to %s", ss, m.OriginTitle, thumb)
	return ffmpeg.Frame(filename, thumb, "-ss", ss)
}
