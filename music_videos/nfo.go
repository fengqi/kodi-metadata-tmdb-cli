package music_videos

import (
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"os"
)

// 保存到nfo
func (mv *MusicVideo) saveToNfo() error {
	nfoPath := mv.nfoPath()

	utils.Logger.InfoF("save music video nfo to: %s", nfoPath)

	fileInfo := &FileInfo{
		StreamDetails: StreamDetails{
			Video: []Video{
				{
					Codec:             mv.VideoStream.CodecName,
					Aspect:            mv.VideoStream.DisplayAspectRatio,
					Width:             mv.VideoStream.Width,
					Height:            mv.VideoStream.Height,
					DurationInSeconds: mv.VideoStream.Duration,
					StereoMode:        "progressive",
				},
			},
			Audio: []Audio{
				{
					Language: "zho",
					Codec:    mv.VideoStream.CodecName,
					Channels: mv.VideoStream.Channels,
				},
			},
		},
	}

	fi, err := os.Stat(mv.MediaFile.Path)
	if err != nil {
		return err
	}

	top := MusicVideoNfo{
		Title:     mv.MediaFile.Filename,
		DateAdded: fi.ModTime().Format("2006-01-02 15:04:05"),
		FileInfo:  fileInfo,
		Thumb: []Thumb{
			{
				Aspect:  "thumb",
				Preview: mv.thumbPath(),
			},
		},
		Poster: mv.thumbPath(),
	}

	return utils.SaveNfo(nfoPath, top)
}
