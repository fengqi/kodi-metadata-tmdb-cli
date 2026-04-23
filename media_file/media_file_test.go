package media_file

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"path/filepath"
	"testing"
)

func setCollectorForTest(t *testing.T, runMode int) {
	t.Helper()

	oldCollector := config.Collector
	config.Collector = &config.CollectorConfig{RunMode: runMode}
	t.Cleanup(func() {
		config.Collector = oldCollector
	})
}

func TestNewMediaFile(t *testing.T) {
	t.Run("hidden file should return nil", func(t *testing.T) {
		setCollectorForTest(t, config.CollectorRunModeDaemon)

		got := NewMediaFile(" /library/movies/.hidden ", ".hidden", Movies)
		if got != nil {
			t.Fatalf("NewMediaFile should return nil for hidden file")
		}
	})

	tests := []struct {
		name       string
		runMode    int
		wantTask   TaskType
		path       string
		wantPath   string
		filename   string
		videoType  VideoType
		wantMedia  MediaType
		wantSuffix string
	}{
		{
			name:       "spec mode should map to task spec",
			runMode:    config.CollectorRunModeSpec,
			wantTask:   TaskSpec,
			path:       " C:\\library\\movies\\Movie.Name.2024.mkv ",
			wantPath:   "C:/library/movies/Movie.Name.2024.mkv",
			filename:   "Movie.Name.2024.mkv",
			videoType:  Movies,
			wantMedia:  VIDEO,
			wantSuffix: ".mkv",
		},
		{
			name:       "daemon mode should map to task scan",
			runMode:    config.CollectorRunModeDaemon,
			wantTask:   TaskScan,
			path:       "/library/movies/movie.nfo",
			wantPath:   "/library/movies/movie.nfo",
			filename:   "movie.nfo",
			videoType:  Movies,
			wantMedia:  NFO,
			wantSuffix: ".nfo",
		},
		{
			name:       "once mode should map to task scan",
			runMode:    config.CollectorRunModeOnce,
			wantTask:   TaskScan,
			path:       "/library/shows/trailers",
			wantPath:   "/library/shows/trailers",
			filename:   "episode.mp4",
			videoType:  TvShows,
			wantMedia:  TRAILER,
			wantSuffix: ".mp4",
		},
		{
			name:       "unknown mode should fallback to watcher task",
			runMode:    99,
			wantTask:   TaskWatcher,
			path:       "/library/music/track.mp3",
			wantPath:   "/library/music/track.mp3",
			filename:   "track.mp3",
			videoType:  MusicVideo,
			wantMedia:  AUDIO,
			wantSuffix: ".mp3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setCollectorForTest(t, tt.runMode)

			got := NewMediaFile(tt.path, tt.filename, tt.videoType)
			if got == nil {
				t.Fatalf("NewMediaFile should not return nil")
			}

			wantPath := filepath.ToSlash(filepath.Clean(tt.wantPath))
			if got.Path != wantPath {
				t.Fatalf("Path = %q; want %q", got.Path, wantPath)
			}
			if got.Dir != filepath.Dir(got.Path) {
				t.Fatalf("Dir = %q; want %q", got.Dir, filepath.Dir(got.Path))
			}
			if got.Filename != tt.filename {
				t.Fatalf("Filename = %q; want %q", got.Filename, tt.filename)
			}
			if got.VideoType != tt.videoType {
				t.Fatalf("VideoType = %v; want %v", got.VideoType, tt.videoType)
			}
			if got.MediaType != tt.wantMedia {
				t.Fatalf("MediaType = %v; want %v", got.MediaType, tt.wantMedia)
			}
			if got.Suffix != tt.wantSuffix {
				t.Fatalf("Suffix = %q; want %q", got.Suffix, tt.wantSuffix)
			}
			if got.TaskType != tt.wantTask {
				t.Fatalf("TaskType = %v; want %v", got.TaskType, tt.wantTask)
			}
		})
	}
}

func TestParseMediaType(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		filename string
		want     MediaType
	}{
		{name: "extra folder", path: "/library/extras", filename: "extras", want: EXTRA},
		{name: "extra folder single", path: "/library/extra", filename: "extra", want: EXTRA},
		{name: "nfo extension should be nfo", path: "/library/movie.nfo", filename: "movie.NFO", want: NFO},
		{name: "graphic extension should be graphic", path: "/library/poster.jpg", filename: "poster.jpg", want: GRAPHIC},
		{name: "audio extension should be audio", path: "/library/theme.mp3", filename: "theme.mp3", want: AUDIO},
		{name: "subtitle extension should be subtitle", path: "/library/movie.srt", filename: "movie.srt", want: SUBTITLE},
		{name: "dvd file should be disc", path: "/library/vts_01_1.vob", filename: "VTS_01_1.VOB", want: DISC},
		{name: "bluray file should be disc", path: "/library/00001.m2ts", filename: "00001.M2TS", want: DISC},
		{name: "hddvd folder should be disc", path: "/library/hdvd_ts", filename: "movie.idx", want: DISC},
		{name: "movie trailer should be trailer", path: "/library/movie-trailer.mp4", filename: "movie-trailer.mp4", want: TRAILER},
		{name: "trailer folder should be trailer", path: "/library/trailers", filename: "clip.mp4", want: TRAILER},
		{name: "trailer regex should be trailer", path: "/library/My.Movie-trailer-2.mp4", filename: "My.Movie-trailer-2.mp4", want: TRAILER},
		{name: "sample basename should be sample", path: "/library/sample.MP4", filename: "sample.mp4", want: SAMPLE},
		{name: "sample folder should be sample", path: "/library/sample", filename: "clip.mp4", want: SAMPLE},
		{name: "sample regex should be sample", path: "/library/My.Movie.sample-.mp4", filename: "My.Movie.sample-.mp4", want: SAMPLE},
		{name: "normal video should be video", path: "/library/movie.mp4", filename: "movie.mp4", want: VIDEO},
		{name: "disc folder fallback should be video", path: "/library/my_video_ts/", filename: "note.xyz", want: VIDEO},
		{name: "unknown extension should be unknown", path: "/library/movie.binlog", filename: "movie.binlog", want: UNKNOWN},
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

func TestDiscHelpers(t *testing.T) {
	if !isDiscFile("video_ts", "/library/dvd/video_ts") {
		t.Fatalf("isDiscFile should identify dvd")
	}
	if !isDiscFile("index.bdmv", "/library/BDMV/index.bdmv") {
		t.Fatalf("isDiscFile should identify bluray")
	}
	if !isDiscFile("movie.idx", "/library/hdvd_ts") {
		t.Fatalf("isDiscFile should identify hddvd")
	}
	if isDiscFile("movie.mp4", "/library/movie.mp4") {
		t.Fatalf("isDiscFile should reject non-disc file")
	}

	if !isDVDFile("video_ts", "/library/video_ts") {
		t.Fatalf("isDVDFile should identify video_ts")
	}
	if isDVDFile("movie.mp4", "/library/movie.mp4") {
		t.Fatalf("isDVDFile should reject non-dvd file")
	}

	if !isBluRayFile("bdmv", "/library/bdmv") {
		t.Fatalf("isBluRayFile should identify bdmv")
	}
	if isBluRayFile("movie.mp4", "/library/movie.mp4") {
		t.Fatalf("isBluRayFile should reject non-bluray file")
	}

	if !isHDDVDFile("hdvd_ts", "/library/hdvd_ts") {
		t.Fatalf("isHDDVDFile should identify hddvd_ts")
	}
	if isHDDVDFile("movie.mp4", "/library/movie.mp4") {
		t.Fatalf("isHDDVDFile should reject non-hddvd file")
	}
}

func TestMediaFileMethods(t *testing.T) {
	cases := []struct {
		name          string
		mf            *MediaFile
		wantIsNFO     bool
		wantIsVideo   bool
		wantIsBluRay  bool
		wantIsDvd     bool
		wantIsDisc    bool
		wantPathNoExt string
	}{
		{
			name: "nfo file",
			mf: &MediaFile{
				Path:      "/library/movie.nfo",
				Filename:  "movie.nfo",
				Suffix:    ".nfo",
				MediaType: NFO,
			},
			wantIsNFO:     true,
			wantPathNoExt: "/library/movie",
		},
		{
			name: "video file",
			mf: &MediaFile{
				Path:      "/library/movie.mp4",
				Filename:  "movie.mp4",
				Suffix:    ".mp4",
				MediaType: VIDEO,
			},
			wantIsVideo:   true,
			wantPathNoExt: "/library/movie",
		},
		{
			name: "bluray disc by bdmv",
			mf: &MediaFile{
				Path:      "/library/BDMV",
				Filename:  "BDMV",
				Suffix:    "",
				MediaType: DISC,
			},
			wantIsBluRay:  true,
			wantIsDisc:    true,
			wantPathNoExt: "/library/BDMV",
		},
		{
			name: "bluray disc by hddvd_ts",
			mf: &MediaFile{
				Path:      "/library/HDVD_TS",
				Filename:  "HDVD_TS",
				Suffix:    "",
				MediaType: DISC,
			},
			wantIsBluRay:  true,
			wantIsDisc:    true,
			wantPathNoExt: "/library/HDVD_TS",
		},
		{
			name: "dvd disc by video_ts",
			mf: &MediaFile{
				Path:      "/library/VIDEO_TS",
				Filename:  "VIDEO_TS",
				Suffix:    "",
				MediaType: DISC,
			},
			wantIsDvd:     true,
			wantIsDisc:    true,
			wantPathNoExt: "/library/VIDEO_TS",
		},
		{
			name: "dvd disc by dvd",
			mf: &MediaFile{
				Path:      "/library/DVD",
				Filename:  "DVD",
				Suffix:    "",
				MediaType: DISC,
			},
			wantIsDvd:     true,
			wantIsDisc:    true,
			wantPathNoExt: "/library/DVD",
		},
		{
			name: "non-disc should not match disc helpers",
			mf: &MediaFile{
				Path:      "/library/poster.jpg",
				Filename:  "poster.jpg",
				Suffix:    ".jpg",
				MediaType: GRAPHIC,
			},
			wantPathNoExt: "/library/poster",
		},
	}

	for _, tt := range cases {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.mf.IsNFO(); got != tt.wantIsNFO {
				t.Fatalf("IsNFO = %v; want %v", got, tt.wantIsNFO)
			}
			if got := tt.mf.IsVideo(); got != tt.wantIsVideo {
				t.Fatalf("IsVideo = %v; want %v", got, tt.wantIsVideo)
			}
			if got := tt.mf.IsBluRay(); got != tt.wantIsBluRay {
				t.Fatalf("IsBluRay = %v; want %v", got, tt.wantIsBluRay)
			}
			if got := tt.mf.IsDvd(); got != tt.wantIsDvd {
				t.Fatalf("IsDvd = %v; want %v", got, tt.wantIsDvd)
			}
			if got := tt.mf.IsDisc(); got != tt.wantIsDisc {
				t.Fatalf("IsDisc = %v; want %v", got, tt.wantIsDisc)
			}
			if got := tt.mf.PathWithoutSuffix(); got != tt.wantPathNoExt {
				t.Fatalf("PathWithoutSuffix = %q; want %q", got, tt.wantPathNoExt)
			}
		})
	}
}
