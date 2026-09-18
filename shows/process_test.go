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

func TestReadGroupId_SeasonRootFirst(t *testing.T) {
	showsDir := filepath.Join(t.TempDir(), "shows")
	showRoot := filepath.Join(showsDir, "Money.Heist")
	seasonRoot := filepath.Join(showRoot, "Season 02")
	require.NoError(t, os.MkdirAll(filepath.Join(seasonRoot, "tmdb"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(showRoot, "tmdb"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(seasonRoot, "tmdb", "group.txt"), []byte("season-group"), 0644))
	require.NoError(t, os.WriteFile(filepath.Join(showRoot, "tmdb", "group.txt"), []byte("tv-root-group"), 0644))

	mf := &media_file.MediaFile{
		Path:      filepath.Join(seasonRoot, "episode-a.mkv"),
		Dir:       seasonRoot,
		Filename:  "episode-a.mkv",
		Suffix:    ".mkv",
		MediaType: media_file.VIDEO,
		VideoType: media_file.TvShows,
	}

	show := &Show{MediaFile: mf, TvRoot: showRoot, SeasonRoot: seasonRoot}
	show.ReadGroupId()
	assert.Equal(t, "season-group", show.GroupId)
}

func TestReadGroupId_FallbackTvRoot(t *testing.T) {
	showsDir := filepath.Join(t.TempDir(), "shows")
	showRoot := filepath.Join(showsDir, "Money.Heist")
	seasonRoot := filepath.Join(showRoot, "Season 02")
	require.NoError(t, os.MkdirAll(filepath.Join(seasonRoot, "tmdb"), 0755))
	require.NoError(t, os.MkdirAll(filepath.Join(showRoot, "tmdb"), 0755))
	require.NoError(t, os.WriteFile(filepath.Join(showRoot, "tmdb", "group.txt"), []byte(" 5eb7353b0cb3350020ce402e \r\n"), 0644))

	mf := &media_file.MediaFile{
		Path:      filepath.Join(seasonRoot, "episode-a.mkv"),
		Dir:       seasonRoot,
		Filename:  "episode-a.mkv",
		Suffix:    ".mkv",
		MediaType: media_file.VIDEO,
		VideoType: media_file.TvShows,
	}

	show := &Show{MediaFile: mf, TvRoot: showRoot, SeasonRoot: seasonRoot}
	show.ReadGroupId()
	assert.Equal(t, "5eb7353b0cb3350020ce402e", show.GroupId)
}

// buildGroupDetail 构造测试用的剧集分组缓存：
// 组1(order=1)有13集，组2(order=2)有9集，组内剧集指向原始编号（模拟纸房子 Netflix 分组）
func buildGroupDetail() *tmdb.TvEpisodeGroupDetail {
	group1 := tmdb.TvEpisodeGroup{Id: "g1", Name: "Part 1", Order: 1}
	for i := range 13 {
		group1.Episodes = append(group1.Episodes, tmdb.TvEpisodeGroupEpisode{
			Id: 2000 + i, SeasonNumber: 0, EpisodeNumber: i + 1, Order: i,
			AirDate: "2017-12-20", Name: fmt.Sprintf("Part1-E%d", i+1), StillPath: "/p1.jpg",
		})
	}

	group2 := tmdb.TvEpisodeGroup{Id: "g2", Name: "Part 2", Order: 2}
	for i := range 9 {
		group2.Episodes = append(group2.Episodes, tmdb.TvEpisodeGroupEpisode{
			Id: 3000 + i, SeasonNumber: 0, EpisodeNumber: 14 + i, Order: i,
			AirDate: "2018-04-06", Name: fmt.Sprintf("Part2-E%d", i+1), StillPath: "/p2.jpg",
		})
	}

	return &tmdb.TvEpisodeGroupDetail{
		Id: "5eb730dfca7ec6001f7beb51", Name: "Parts", GroupCount: 2, EpisodeCount: 22,
		Groups: []tmdb.TvEpisodeGroup{group1, group2},
	}
}

func newGroupTestShow(t *testing.T) *Show {
	showsDir := filepath.Join(t.TempDir(), "shows")
	showRoot := filepath.Join(showsDir, "Money.Heist")
	seasonRoot := showRoot
	require.NoError(t, os.MkdirAll(filepath.Join(seasonRoot, "tmdb"), 0755))

	groupBytes, err := json.Marshal(buildGroupDetail())
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(seasonRoot, "tmdb", "group.json"), groupBytes, 0644))

	mf := &media_file.MediaFile{
		Path:      filepath.Join(seasonRoot, "Money.Heist.S02E01.1080p.NF.WEB-DL.mkv"),
		Dir:       seasonRoot,
		Filename:  "Money.Heist.S02E01.1080p.NF.WEB-DL.mkv",
		Suffix:    ".mkv",
		MediaType: media_file.VIDEO,
		VideoType: media_file.TvShows,
	}

	return &Show{
		MediaFile: mf, TvRoot: showRoot, SeasonRoot: seasonRoot,
		TvId: 71446, GroupId: "5eb730dfca7ec6001f7beb51", Season: 2, Episode: 1,
	}
}

func TestGetEpisodeDetail_GroupMapping(t *testing.T) {
	config.Log = &config.LogConfig{Mode: config.LogModeStdout, Level: config.LogLevelDebug}
	utils.InitLogger()

	// 纯分组信息，不调用任何TMDB接口（Api为nil，一旦被调用会panic）
	tmdb.Api = nil

	show := newGroupTestShow(t)

	got, err := show.getEpisodeDetail()
	require.NoError(t, err)
	require.NotNil(t, got)
	// S02E01(分组编号) -> 组2(order=2)第1个 -> 直接使用分组内的信息
	assert.Equal(t, 3000, got.Id)
	assert.Equal(t, "Part2-E1", got.Name)
	assert.Equal(t, 2, got.SeasonNumber)
	assert.Equal(t, 1, got.EpisodeNumber)
	assert.Equal(t, "/p2.jpg", got.StillPath)
	assert.Equal(t, "2018-04-06", got.AirDate)

	// 分集缓存应保存分组编号
	cached := new(tmdb.TvEpisodeDetail)
	bytes, err := os.ReadFile(show.EpisodeCacheFile())
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal(bytes, cached))
	assert.Equal(t, 2, cached.SeasonNumber)
	assert.Equal(t, 1, cached.EpisodeNumber)
}

func TestGetEpisodeDetail_GroupSeasonOrEpisodeOutOfRange(t *testing.T) {
	config.Log = &config.LogConfig{Mode: config.LogModeStdout, Level: config.LogLevelDebug}
	utils.InitLogger()

	tmdb.Api = nil

	show := newGroupTestShow(t)
	show.Season = 3
	_, err := show.getEpisodeDetail()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found in episode group")

	show2 := newGroupTestShow(t)
	show2.Episode = 10
	_, err = show2.getEpisodeDetail()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "out of range")
}
