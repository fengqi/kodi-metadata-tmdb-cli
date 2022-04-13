package music_videos

import (
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io/fs"
	"strings"
	"syscall"
	"time"
)

func (c *Collector) parseVideoFile(dir string, file fs.FileInfo) *MusicVideo {
	ext := utils.IsVideo(file.Name())
	if ext == "" {
		utils.Logger.DebugF("not a video file: %s", file.Name())
		return nil
	}

	stat := file.Sys().(*syscall.Stat_t)
	mTime := time.Unix(stat.Mtim.Sec, stat.Mtim.Nsec)

	return &MusicVideo{
		Dir:         dir,
		OriginTitle: file.Name(),
		Title:       strings.Replace(file.Name(), "."+ext, "", 1),
		DateAdded:   mTime.Format("2006-01-02 15:04:05"),
	}
}
