package music_videos

import (
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io/fs"
	"strings"
)

func (c *Collector) parseVideoFile(dir string, file fs.FileInfo) *MusicVideo {
	ext := utils.IsVideo(file.Name())
	if ext == "" {
		utils.Logger.DebugF("not a video file: %s", file.Name())
		return nil
	}

	baseDir := ""
	for _, base := range c.config.Collector.MusicVideosDir {
		if dir == baseDir || strings.Contains(dir, base) {
			baseDir = base
			break
		}
	}

	return &MusicVideo{
		Dir:         dir,
		BaseDir:     baseDir,
		OriginTitle: file.Name(),
		Title:       strings.Replace(file.Name(), "."+ext, "", 1),
		DateAdded:   file.ModTime().Format("2006-01-02 15:04:05"),
	}
}
