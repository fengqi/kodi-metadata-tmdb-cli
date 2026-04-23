package media_file

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"testing"
)

func TestParseMediaType_CaseInsensitiveByExtension(t *testing.T) {
	oldCollector := config.Collector
	config.Collector = &config.CollectorConfig{
		TmpSuffix: []string{".part", ".!qB", ".!qb", ".!ut"},
	}
	defer func() { config.Collector = oldCollector }()

	tests := []struct {
		name     string
		path     string
		filename string
		want     MediaType
	}{
		{
			name:     "image should be image",
			path:     "/shows/Invincible/Season 1",
			filename: "Invincible.S01E08.2021.1080p.AMZN.WEB-DL.H264.DDP5.1-ADWeb-poster.jpg",
			want:     GRAPHIC,
		},
		{
			name:     "upper nfo should be nfo",
			path:     "/shows/Invincible/Season 1",
			filename: "Invincible.S01E08.2021.1080p.AMZN.WEB-DL.H264.DDP5.1-ADWeb.mp4.NFO",
			want:     NFO,
		},
		{
			name:     "lower nfo should be nfo",
			path:     "/shows/Invincible/Season 1",
			filename: "Invincible.S01E08.2021.1080p.AMZN.WEB-DL.H264.DDP5.1-ADWeb.mp4.nfo",
			want:     NFO,
		},
		{
			name:     "upper mp4 should be video",
			path:     "/shows/Invincible/Season 1",
			filename: "Invincible.S01E08.2021.1080p.AMZN.WEB-DL.H264.DDP5.1-ADWeb.MP4",
			want:     VIDEO,
		},
		{
			name:     "tmp suffix file should not be video",
			path:     "/shows/Invincible/Season 1",
			filename: "Invincible.S01E08.2021.1080p.AMZN.WEB-DL.H264.DDP5.1-ADWeb.mp4.!qB",
			want:     UNKNOWN,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseMediaType(tt.path, tt.filename)
			if got != tt.want {
				t.Fatalf("parseMediaType(%q, %q) = %v; want %v", tt.path, tt.filename, got, tt.want)
			}
		})
	}
}
