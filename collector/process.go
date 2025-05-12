package collector

import (
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"fengqi/kodi-metadata-tmdb-cli/movies"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"log"
)

// runProcess 处理扫描到的文件
func (c *collector) runProcess() {
	utils.Logger.Debug("run process")

	for {
		select {
		case file := <-c.channel:
			utils.Logger.DebugF("receive task: %s", file.Filename)

			switch file.VideoType {
			case media_file.Movies:
				if err := movies.Process(file); err != nil {
					log.Printf("prcess movies error: %s\n", err)
				}
			case media_file.TvShows:
				// todo
			case media_file.MusicVideo:
				// todo
			}
		}
	}
}
