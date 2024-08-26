package ffmpeg

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
)

func InitFfmpeg() {
	SetFfmpeg()
	SetFfprobe()
}

func SetFfmpeg() {
	if !utils.FileExist(config.Ffmpeg.FfmpegPath) {
		return
	}

	ffmpeg = config.Ffmpeg.FfmpegPath
}

func SetFfprobe() {
	if !utils.FileExist(config.Ffmpeg.FfprobePath) {
		return
	}

	ffprobe = config.Ffmpeg.FfprobePath
}
