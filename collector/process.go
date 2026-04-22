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

	for task := range c.channel {
		if task == nil || task.file == nil {
			continue
		}

		utils.Logger.DebugF("receive task: %s", task.file.Filename)

		switch task.file.VideoType {
		case media_file.Movies:
			if err := movies.Process(task.file); err != nil {
				log.Printf("prcess movies error: %v\n", err)
			}
		case media_file.TvShows:
			if err := shows.Process(task.file); err != nil {
				log.Printf("pricess shows error: %v\n", err)
			}
		case media_file.MusicVideo:
			if err := music_videos.Process(task.file); err != nil {
				log.Printf("pricess music videos error: %v\n", err)
			}
		}

		if task.done != nil {
			task.done.Done()
		}
	}

	log.Println("run process done")
}
