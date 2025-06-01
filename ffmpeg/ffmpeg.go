package ffmpeg

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"github.com/fengqi/lrace"
)

func InitFfmpeg() {
	SetFfmpeg()
	SetFfprobe()
}

func SetFfmpeg() {
	if !lrace.FileExist(config.Ffmpeg.FfmpegPath) {
		return
	}

	ffmpeg = config.Ffmpeg.FfmpegPath
}

func SetFfprobe() {
	if !lrace.FileExist(config.Ffmpeg.FfprobePath) {
		return
	}

	ffprobe = config.Ffmpeg.FfprobePath
}
