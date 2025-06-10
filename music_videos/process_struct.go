package music_videos

import (
	"crypto/md5"
	"fengqi/kodi-metadata-tmdb-cli/ffmpeg"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"fmt"
	"github.com/fengqi/lrace"
	"io"
	"path/filepath"
	"strings"
)

type MusicVideo struct {
	MediaFile   *media_file.MediaFile `json:"media_file"` // 媒体文件
	NfoPath     string
	ThumbPath   string
	VideoStream *ffmpeg.Stream
	AudioStream *ffmpeg.Stream
}

func (mv *MusicVideo) nameMd5() string {
	h := md5.New()
	_, _ = io.WriteString(h, mv.MediaFile.Path)
	sum := fmt.Sprintf("%x", h.Sum(nil))
	return sum
}

func (mv *MusicVideo) cacheDir() string {
	return filepath.Dir(mv.MediaFile.Path) + "/tmdb"
}

func (mv *MusicVideo) thumbPath() string {
	if mv.ThumbPath == "" {
		mv.ThumbPath = strings.TrimRight(mv.MediaFile.Path, mv.MediaFile.Suffix) + "-thumb.jpg"
	}
	return mv.ThumbPath
}

func (mv *MusicVideo) thumbExist() bool {
	return lrace.FileExist(mv.thumbPath())
}

// 判断nfo文件是否存在
func (mv *MusicVideo) nfoExist() bool {
	return lrace.FileExist(mv.nfoPath())
}

// 获取nfo文件路径
func (mv *MusicVideo) nfoPath() string {
	if mv.NfoPath == "" {
		mv.NfoPath = strings.TrimRight(mv.MediaFile.Path, mv.MediaFile.Suffix) + ".nfo"
	}
	return mv.NfoPath
}
