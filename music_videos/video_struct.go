package music_videos

import (
	"fengqi/kodi-metadata-tmdb-cli/ffmpeg"
)

type MusicVideo struct {
	Dir         string
	Title       string
	OriginTitle string
	DateAdded   string
	VideoStream *ffmpeg.Stream
	AudioStream *ffmpeg.Stream
}
