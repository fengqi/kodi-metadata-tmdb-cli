package collector

import (
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"fengqi/kodi-metadata-tmdb-cli/movies"
	"fengqi/kodi-metadata-tmdb-cli/music_videos"
	"fengqi/kodi-metadata-tmdb-cli/shows"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"log"
)

// runProcess 处理扫描到的文件
func (c *collector) runProcess() {
	utils.Logger.Debug("run process")

	for file := range c.channel {
		utils.Logger.DebugF("receive task: %s", file.Filename)

		switch file.VideoType {
		case media_file.Movies:
			if err := movies.Process(file); err != nil {
				log.Printf("prcess movies error: %v\n", err)
			}
		case media_file.TvShows:
			if err := shows.Process(file); err != nil {
				log.Printf("pricess shows error: %v\n", err)
			}
		case media_file.MusicVideo:
			if err := music_videos.Process(file); err != nil {
				log.Printf("pricess music videos error: %v\n", err)
			}
		}

		c.wg.Done()
	}
}
