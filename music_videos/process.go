package music_videos

import (
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"os"
)

func Process(mf *media_file.MediaFile) error {
	video := &MusicVideo{MediaFile: mf}
	if mf == nil || (video.nfoExist() && video.thumbExist()) {
		return nil
	}

	if err := video.prepare(); err != nil {
		return err
	}

	if probe, err := video.getProbe(); err != nil {
		return err
	} else {
		video.VideoStream = probe.FirstVideoStream()
		video.AudioStream = probe.FirstAudioStream()
		if video.VideoStream == nil || video.AudioStream == nil {
			return errors.New("video or audio stream is nil")
		}
	}

	if err := video.drawThumb(); err != nil {
		return err
	}

	if err := video.saveToNfo(); err != nil {
		return err
	}

	return nil
}

func (mv *MusicVideo) prepare() error {
	dir := mv.cacheDir()
	if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0755); err != nil {
			return err
		}
	}

	return nil
}
