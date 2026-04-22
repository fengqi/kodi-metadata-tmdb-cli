package movies

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseMoviesFile_ByRuleWhenIdCached(t *testing.T) {
	config.Log = &config.LogConfig{Mode: config.LogModeStdout, Level: config.LogLevelDebug}
	utils.InitLogger()
	config.Collector = &config.CollectorConfig{}
	config.Ai = &config.AiConfig{
		Enable:         true,
		BaseURL:        "http://127.0.0.1",
		ApiKey:         "test",
		Model:          "test",
		MatchMode:      config.AiMatchModeAiThenRule,
		SearchMode:     config.AiSearchModeAiDecision,
		TimeoutSeconds: 1,
	}
	t.Cleanup(func() {
		config.Ai = nil
		config.Collector = nil
	})

	dir := t.TempDir()
	mf := &media_file.MediaFile{
		Path:      filepath.Join(dir, "Inception.2010.mkv"),
		Dir:       dir,
		Filename:  "Inception.2010.mkv",
		Suffix:    ".mkv",
		MediaType: media_file.VIDEO,
		VideoType: media_file.Movies,
	}

	movie := newMovieWithPaths(mf)
	require.NoError(t, os.MkdirAll(movie.GetCacheDir(), 0755))
	require.NoError(t, os.WriteFile(movie.IdFile(), []byte("27205"), 0644))

	got, err := parseMoviesFile(mf)
	require.NoError(t, err)
	require.NotNil(t, got)
	assert.Equal(t, "Inception", got.Title)
	assert.Equal(t, 2010, got.Year)
}

func TestLoadMovieCache_UseMovieDetailCacheFirst(t *testing.T) {
	config.Log = &config.LogConfig{Mode: config.LogModeStdout, Level: config.LogLevelDebug}
	utils.InitLogger()

	dir := t.TempDir()
	mf := &media_file.MediaFile{
		Path:      filepath.Join(dir, "Inception.2010.mkv"),
		Dir:       dir,
		Filename:  "Inception.2010.mkv",
		Suffix:    ".mkv",
		MediaType: media_file.VIDEO,
		VideoType: media_file.Movies,
	}

	movie := newMovieWithPaths(mf)
	require.NoError(t, os.MkdirAll(movie.GetCacheDir(), 0755))
	detailBytes, err := json.Marshal(&tmdb.MovieDetail{Id: 27205, Title: "鐩楁ⅵ绌洪棿", OriginalTitle: "Inception", ReleaseDate: "2010-07-16"})
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(movie.DetailCacheFile(), detailBytes, 0644))
	require.NoError(t, os.WriteFile(movie.IdFile(), []byte("27205"), 0644))

	gotMovie, gotDetail, err := loadMovieCache(mf)
	require.NoError(t, err)
	require.NotNil(t, gotMovie)
	require.NotNil(t, gotDetail)
	assert.Equal(t, 27205, gotMovie.MovieId)
	assert.Equal(t, "鐩楁ⅵ绌洪棿", gotMovie.Title)
	assert.Equal(t, "Inception", gotMovie.EngTitle)
	assert.Equal(t, 2010, gotMovie.Year)
	assert.True(t, gotDetail.FromCache)
}

func TestLoadMovieCache_UseIdCacheOnly(t *testing.T) {
	config.Log = &config.LogConfig{Mode: config.LogModeStdout, Level: config.LogLevelDebug}
	utils.InitLogger()

	dir := t.TempDir()
	mf := &media_file.MediaFile{
		Path:      filepath.Join(dir, "Inception.2010.mkv"),
		Dir:       dir,
		Filename:  "Inception.2010.mkv",
		Suffix:    ".mkv",
		MediaType: media_file.VIDEO,
		VideoType: media_file.Movies,
	}

	movie := newMovieWithPaths(mf)
	require.NoError(t, os.MkdirAll(movie.GetCacheDir(), 0755))
	require.NoError(t, os.WriteFile(movie.IdFile(), []byte("27205"), 0644))

	gotMovie, gotDetail, err := loadMovieCache(mf)
	require.NoError(t, err)
	require.NotNil(t, gotMovie)
	assert.Nil(t, gotDetail)
	assert.Equal(t, mf, gotMovie.MediaFile)
	assert.Equal(t, movie.NfoFile, gotMovie.NfoFile)
}

func TestLoadMovieCache_NoCacheReturnNil(t *testing.T) {
	config.Log = &config.LogConfig{Mode: config.LogModeStdout, Level: config.LogLevelDebug}
	utils.InitLogger()

	dir := t.TempDir()
	mf := &media_file.MediaFile{
		Path:      filepath.Join(dir, "Inception.2010.mkv"),
		Dir:       dir,
		Filename:  "Inception.2010.mkv",
		Suffix:    ".mkv",
		MediaType: media_file.VIDEO,
		VideoType: media_file.Movies,
	}

	gotMovie, gotDetail, err := loadMovieCache(mf)
	require.NoError(t, err)
	assert.Nil(t, gotMovie)
	assert.Nil(t, gotDetail)
}

func TestLoadMovieCache_EmptyInput(t *testing.T) {
	gotMovie, gotDetail, err := loadMovieCache(nil)
	require.Error(t, err)
	assert.Nil(t, gotMovie)
	assert.Nil(t, gotDetail)
}

func TestNewMovieWithPaths_SingleFile(t *testing.T) {
	dir := t.TempDir()
	mf := &media_file.MediaFile{
		Path:      filepath.Join(dir, "Inception.2010.mkv"),
		Dir:       dir,
		Filename:  "Inception.2010.mkv",
		Suffix:    ".mkv",
		MediaType: media_file.VIDEO,
		VideoType: media_file.Movies,
	}

	movie := newMovieWithPaths(mf)
	require.NotNil(t, movie)
	assert.Equal(t, dir+"/Inception.2010-poster.jpg", movie.PosterFile)
	assert.Equal(t, dir+"/Inception.2010-fanart.jpg", movie.FanArtFile)
	assert.Equal(t, dir+"/Inception.2010-clearlogo.png", movie.ClearLogoFile)
	assert.Equal(t, dir+"/Inception.2010.nfo", movie.NfoFile)
}

func TestNewMovieWithPaths_BluRay(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "Movie")
	mf := &media_file.MediaFile{
		Path:      dir,
		Dir:       dir,
		Filename:  media_file.BDMV,
		MediaType: media_file.DISC,
		VideoType: media_file.Movies,
	}

	movie := newMovieWithPaths(mf)
	require.NotNil(t, movie)
	assert.Equal(t, dir+"/poster.jpg", movie.PosterFile)
	assert.Equal(t, dir+"/fanart.jpg", movie.FanArtFile)
	assert.Equal(t, dir+"/clearlogo.png", movie.ClearLogoFile)
	assert.Equal(t, dir+"/index.nfo", movie.NfoFile)
}

func TestNewMovieWithPaths_Dvd(t *testing.T) {
	dir := filepath.Join(t.TempDir(), "Movie")
	mf := &media_file.MediaFile{
		Path:      dir,
		Dir:       dir,
		Filename:  media_file.VideoTs,
		MediaType: media_file.DISC,
		VideoType: media_file.Movies,
	}

	movie := newMovieWithPaths(mf)
	require.NotNil(t, movie)
	assert.Equal(t, dir+"/poster.jpg", movie.PosterFile)
	assert.Equal(t, dir+"/fanart.jpg", movie.FanArtFile)
	assert.Equal(t, dir+"/clearlogo.png", movie.ClearLogoFile)
	assert.Equal(t, dir+"/VIDEO_TS/VIDEO_TS.nfo", movie.NfoFile)
}
