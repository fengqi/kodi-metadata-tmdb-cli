package shows

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseShowFile_UseRuleWhenTvIdCached(t *testing.T) {
	config.Log = &config.LogConfig{Mode: config.LogModeStdout, Level: config.LogLevelDebug}
	utils.InitLogger()
	config.Ai = &config.AiConfig{
		Enable:         true,
		BaseURL:        "http://127.0.0.1",
		ApiKey:         "test",
		Model:          "test",
		MatchMode:      config.AiMatchModeAiThenRule,
		SearchMode:     config.AiSearchModeAiDecision,
		TimeoutSeconds: 1,
	}
	showsDir := filepath.Join(t.TempDir(), "shows")
	config.Collector = &config.CollectorConfig{ShowsDir: []string{showsDir}}
	t.Cleanup(func() {
		config.Ai = nil
		config.Collector = nil
	})

	showRoot := filepath.Join(showsDir, "Foundation")
	seasonRoot := filepath.Join(showRoot, "Season 02")
	require.NoError(t, os.MkdirAll(filepath.Join(showRoot, "tmdb"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(showRoot, "tmdb", "id.txt"), []byte("93740"), 0644))

	mf := &media_file.MediaFile{
		Path:      filepath.Join(seasonRoot, "Foundation.S02E03.mkv"),
		Dir:       seasonRoot,
		Filename:  "Foundation.S02E03.mkv",
		Suffix:    ".mkv",
		MediaType: media_file.VIDEO,
		VideoType: media_file.TvShows,
	}

	got, err := parseShowFile(mf)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Foundation", got.Title)
	assert.Equal(t, 2, got.Season)
	assert.Equal(t, 3, got.Episode)
	assert.Equal(t, 93740, got.TvId)
}

func TestLoadShowCache_UseEpisodeCacheFirst(t *testing.T) {
	config.Log = &config.LogConfig{Mode: config.LogModeStdout, Level: config.LogLevelDebug}
	utils.InitLogger()
	showsDir := filepath.Join(t.TempDir(), "shows")
	config.Collector = &config.CollectorConfig{ShowsDir: []string{showsDir}}
	t.Cleanup(func() {
		config.Collector = nil
	})

	showRoot := filepath.Join(showsDir, "Foundation")
	seasonRoot := filepath.Join(showRoot, "Season 02")
	require.NoError(t, os.MkdirAll(filepath.Join(showRoot, "tmdb"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(seasonRoot, "tmdb"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(showRoot, "tmdb", "id.txt"), []byte("93740"), 0644))

	tvBytes, err := json.Marshal(&tmdb.TvDetail{Id: 93740, Name: "Foundation", LastAirDate: "2024-01-01", FirstAirDate: "2023-01-01"})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(showRoot, "tmdb", "tv.json"), tvBytes, 0644))

	mf := &media_file.MediaFile{
		Path:      filepath.Join(seasonRoot, "episode-a.mkv"),
		Dir:       seasonRoot,
		Filename:  "episode-a.mkv",
		Suffix:    ".mkv",
		MediaType: media_file.VIDEO,
		VideoType: media_file.TvShows,
	}

	show := &Show{MediaFile: mf}
	fillShowPathMeta(show)
	episodeBytes, err := json.Marshal(&tmdb.TvEpisodeDetail{Id: 1001, Name: "Episode A", AirDate: "2024-01-01", SeasonNumber: 2, EpisodeNumber: 3})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(show.EpisodeCacheFile(), episodeBytes, 0644))

	gotShow, gotTv, gotEpisode, err := loadShowCache(mf)
	require.NoError(t, err)
	require.NotNil(t, gotShow)
	require.NotNil(t, gotTv)
	require.NotNil(t, gotEpisode)
	assert.Equal(t, 93740, gotShow.TvId)
	assert.Equal(t, 2, gotShow.Season)
	assert.Equal(t, 3, gotShow.Episode)
	assert.True(t, gotTv.FromCache)
	assert.True(t, gotEpisode.FromCache)
}

func TestGetEpisodeDetail_UseLegacyCacheWhenNewCacheMissing(t *testing.T) {
	config.Log = &config.LogConfig{Mode: config.LogModeStdout, Level: config.LogLevelDebug}
	utils.InitLogger()

	showsDir := filepath.Join(t.TempDir(), "shows")
	config.Collector = &config.CollectorConfig{ShowsDir: []string{showsDir}}
	t.Cleanup(func() {
		config.Collector = nil
	})

	showRoot := filepath.Join(showsDir, "Foundation")
	seasonRoot := filepath.Join(showRoot, "Season 02")
	require.NoError(t, os.MkdirAll(filepath.Join(seasonRoot, "tmdb"), 0755))

	mf := &media_file.MediaFile{
		Path:      filepath.Join(seasonRoot, "episode-a.mkv"),
		Dir:       seasonRoot,
		Filename:  "episode-a.mkv",
		Suffix:    ".mkv",
		MediaType: media_file.VIDEO,
		VideoType: media_file.TvShows,
	}

	show := &Show{MediaFile: mf, TvRoot: showRoot, SeasonRoot: seasonRoot, Season: 2, Episode: 3, TvId: 93740}
	legacyFile := filepath.Join(seasonRoot, "tmdb", fmt.Sprintf("s%02de%02d.json", show.Season, show.Episode))
	episodeBytes, err := json.Marshal(&tmdb.TvEpisodeDetail{Id: 1001, Name: "Legacy Episode", AirDate: "2024-01-01", SeasonNumber: 2, EpisodeNumber: 3})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(legacyFile, episodeBytes, 0644))

	got, err := show.getEpisodeDetail()
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.FromCache)
	assert.Equal(t, "Legacy Episode", got.Name)
	assert.Equal(t, 2, show.Season)
	assert.Equal(t, 3, show.Episode)
	_, err = os.Stat(show.EpisodeCacheFile())
	assert.NoError(t, err)
	_, err = os.Stat(legacyFile)
	assert.True(t, os.IsNotExist(err))
}

func TestGetEpisodeDetail_PreferNewCacheOverLegacy(t *testing.T) {
	config.Log = &config.LogConfig{Mode: config.LogModeStdout, Level: config.LogLevelDebug}
	utils.InitLogger()

	showsDir := filepath.Join(t.TempDir(), "shows")
	config.Collector = &config.CollectorConfig{ShowsDir: []string{showsDir}}
	t.Cleanup(func() {
		config.Collector = nil
	})

	showRoot := filepath.Join(showsDir, "Foundation")
	seasonRoot := filepath.Join(showRoot, "Season 02")
	require.NoError(t, os.MkdirAll(filepath.Join(seasonRoot, "tmdb"), 0755))

	mf := &media_file.MediaFile{
		Path:      filepath.Join(seasonRoot, "episode-a.mkv"),
		Dir:       seasonRoot,
		Filename:  "episode-a.mkv",
		Suffix:    ".mkv",
		MediaType: media_file.VIDEO,
		VideoType: media_file.TvShows,
	}

	show := &Show{MediaFile: mf, TvRoot: showRoot, SeasonRoot: seasonRoot, Season: 2, Episode: 3, TvId: 93740}
	newBytes, err := json.Marshal(&tmdb.TvEpisodeDetail{Id: 2001, Name: "New Episode", AirDate: "2024-01-01", SeasonNumber: 2, EpisodeNumber: 3})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(show.EpisodeCacheFile(), newBytes, 0644))

	legacyFile := filepath.Join(seasonRoot, "tmdb", fmt.Sprintf("s%02de%02d.json", show.Season, show.Episode))
	legacyBytes, err := json.Marshal(&tmdb.TvEpisodeDetail{Id: 1001, Name: "Legacy Episode", AirDate: "2024-01-01", SeasonNumber: 2, EpisodeNumber: 3})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(legacyFile, legacyBytes, 0644))

	got, err := show.getEpisodeDetail()
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.True(t, got.FromCache)
	assert.Equal(t, 2001, got.Id)
	assert.Equal(t, "New Episode", got.Name)
}

func TestGetEpisodeDetail_IgnoreExpiredLegacyCache(t *testing.T) {
	config.Log = &config.LogConfig{Mode: config.LogModeStdout, Level: config.LogLevelDebug}
	utils.InitLogger()

	showsDir := filepath.Join(t.TempDir(), "shows")
	config.Collector = &config.CollectorConfig{ShowsDir: []string{showsDir}}
	t.Cleanup(func() {
		config.Collector = nil
	})

	showRoot := filepath.Join(showsDir, "Foundation")
	seasonRoot := filepath.Join(showRoot, "Season 02")
	require.NoError(t, os.MkdirAll(filepath.Join(seasonRoot, "tmdb"), 0755))

	mf := &media_file.MediaFile{
		Path:      filepath.Join(seasonRoot, "episode-a.mkv"),
		Dir:       seasonRoot,
		Filename:  "episode-a.mkv",
		Suffix:    ".mkv",
		MediaType: media_file.VIDEO,
		VideoType: media_file.TvShows,
	}

	show := &Show{MediaFile: mf, TvRoot: showRoot, SeasonRoot: seasonRoot, Season: 2, Episode: 3}
	legacyFile := filepath.Join(seasonRoot, "tmdb", fmt.Sprintf("s%02de%02d.json", show.Season, show.Episode))
	episodeBytes, err := json.Marshal(&tmdb.TvEpisodeDetail{Id: 1001, Name: "Legacy Episode", AirDate: "2026-04-01", SeasonNumber: 2, EpisodeNumber: 3})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(legacyFile, episodeBytes, 0644))
	oldTime := time.Date(2026, time.April, 3, 0, 0, 0, 0, time.UTC)
	require.NoError(t, os.Chtimes(legacyFile, oldTime, oldTime))

	got, err := show.getEpisodeDetail()
	require.Error(t, err)
	assert.Nil(t, got)
}
