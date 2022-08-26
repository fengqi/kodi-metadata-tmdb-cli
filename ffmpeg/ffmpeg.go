package ffmpeg

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
)

func InitFfmpeg(config *config.FfmpegConfig) {
	SetFfmpeg(config.FfmpegPath)
	SetFfprobe(config.FfprobePath)
}

func SetFfmpeg(path string) {
	if !utils.FileExist(path) {
		return
	}

	ffmpeg = path
}

func SetFfprobe(path string) {
	if !utils.FileExist(path) {
		return
	}

	ffprobe = path
}
