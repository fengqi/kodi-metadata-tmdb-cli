package shows

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"github.com/agiledragon/gomonkey/v2"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseShowFile(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()
	patches.ApplyGlobalVar(&config.Collector, &config.CollectorConfig{ShowsDir: []string{"/data/tmp/shows"}})
	patches.ApplyFuncReturn(os.Mkdir, func(...any) error { return nil })
	patches.ApplyPrivateMethod(&Show{}, "checkCacheDir", func() {})
	patches.ApplyPrivateMethod(&Show{}, "checkTvCacheDir", func() {})

	tests := []struct {
		name      string
		mediaFile *media_file.MediaFile
		show      *Show
		want      *Show
	}{
		{
			name: "S01E01.mp4",
			mediaFile: &media_file.MediaFile{
				Path:      "/data/tmp/shows/庆余年/庆余年S02/S02E03.mp4",
				Filename:  "S02E03.mp4",
				Suffix:    ".mp4",
				MediaType: media_file.VIDEO,
				VideoType: media_file.TvShows,
			},
			want: &Show{
				Title:   "庆余年",
				Season:  2,
				Episode: 3,
			},
		},
		{
			name: "Season Dir",
			mediaFile: &media_file.MediaFile{
				Path:      "/data/tmp/shows/庆余年/庆余年/S02/E03.mp4",
				Filename:  "E03.mp4",
				Suffix:    ".mp4",
				MediaType: media_file.VIDEO,
				VideoType: media_file.TvShows,
			},
			want: &Show{
				Title:   "庆余年",
				Season:  2,
				Episode: 3,
			},
		},
		{
			name: "Season Dir",
			mediaFile: &media_file.MediaFile{
				Path:      "/data/tmp/shows/9-1-1.Lone.Star.S01.1080p.DSNP.WEB-DL.DDP5.1.H264-HHWEB/9-1-1.Lone.Star.S01E08.Monster.Inside.1080p.DSNP.WEB-DL.DDP5.1.H264-HHWEB.mkv",
				Filename:  "9-1-1.Lone.Star.S01E08.Monster.Inside.1080p.DSNP.WEB-DL.DDP5.1.H264-HHWEB.mkv",
				Suffix:    ".mkv",
				MediaType: media_file.VIDEO,
				VideoType: media_file.TvShows,
			},
			want: &Show{
				Title:   "9 1 1 Lone Star",
				Season:  1,
				Episode: 8,
			},
		},
		{
			name: "Season Dir",
			mediaFile: &media_file.MediaFile{
				Path:      "/data/tmp/shows/Gannibal.2022.Disney+.WEB-DL.4K.HEVC.HDR.DDP-HDCTV/Gannibal.E01.2022.Disney+.WEB-DL.4K.HEVC.HDR.DDP-HDCTV.mkv",
				Filename:  "Gannibal.E01.2022.Disney+.WEB-DL.4K.HEVC.HDR.DDP-HDCTV.mkv",
				Suffix:    ".mkv",
				MediaType: media_file.VIDEO,
				VideoType: media_file.TvShows,
			},
			want: &Show{
				Title:   "Gannibal",
				Season:  1,
				Episode: 1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			show := &Show{MediaFile: tt.mediaFile}
			ParseShowFile(show, tt.mediaFile.Path)

			assert.Equal(t, tt.want.Title, show.Title)
			assert.Equal(t, tt.want.Season, show.Season)
			assert.Equal(t, tt.want.Episode, show.Episode)
		})
	}
}
