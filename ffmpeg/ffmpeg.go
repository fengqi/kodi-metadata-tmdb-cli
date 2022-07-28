package ffmpeg

import "fengqi/kodi-metadata-tmdb-cli/config"

func InitFfmpeg(config *config.FfmpegConfig) {
	SetFfmpeg(config.FfmpegPath)
	SetFfprobe(config.FfprobePath)
}

func SetFfmpeg(path string) {
	ffmpeg = path
}

func SetFfprobe(path string) {
	ffprobe = path
}
