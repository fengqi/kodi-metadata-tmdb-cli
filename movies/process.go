package movies

import (
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/common/ai"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/kodi"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fengqi/lrace"
)

// Process 处理扫描到的电影文件
func Process(mf *media_file.MediaFile) error {
	movie, detail, err := loadMovieCache(mf)
	if err != nil {
		return err
	}

	if movie == nil {
		movie, err = parseMoviesFile(mf)
		if err != nil {
			return err
		}
	}
	if movie == nil {
		return errors.New("movie file empty")
	}

	utils.Logger.DebugF("receive movies task: %v", movie)

	movie.checkCacheDir()
	if detail == nil {
		detail, err = movie.getMovieDetail()
		if err != nil {
			return err
		}
	}
	if detail == nil {
		return errors.New("get movie detail empty")
	}

	if !detail.FromCache || !lrace.FileExist(movie.NfoFile) {
		if err := movie.saveToNfo(detail); err != nil {
			return err
		}
		kodi.Rpc.AddRefreshTask(kodi.TaskRefreshMovie, detail.Title)
	}

	return movie.downloadImage(detail)
}

// 解析文件, 返回详情：年份、中文名称、英文名称等
func parseMoviesFile(mf *media_file.MediaFile) (*Movie, error) {
	// 未启用ai
	if !ai.Enabled() {
		return parseMoviesFileByRule(mf)
	}

	// 规则优先，匹配不到再 AI 介入
	if config.Ai.MatchMode == config.AiMatchModeRuleThenAi {
		ruleMovie, ruleErr := parseMoviesFileByRule(mf)
		if ruleErr == nil {
			return ruleMovie, nil
		}
		aiMovie, aiErr := parseMoviesFileByAI(mf, ruleMovie)
		if aiErr == nil {
			return aiMovie, nil
		}
		return nil, ruleErr
	}

	// AI 优先，匹配不到再规则介入
	if config.Ai.MatchMode == config.AiMatchModeAiThenRule {
		aiMovie, aiErr := parseMoviesFileByAI(mf, nil)
		if aiErr == nil {
			return aiMovie, nil
		}
		return parseMoviesFileByRule(mf)
	}

	// 规则优先，结果给 AI 参考，最终使用 AI 结果
	ruleMovie, ruleErr := parseMoviesFileByRule(mf)
	aiMovie, aiErr := parseMoviesFileByAI(mf, ruleMovie)
	if aiErr == nil {
		return aiMovie, nil
	}
	return ruleMovie, ruleErr
}

// parseMoviesFileByRule 规则解析：年份、中英文名、别名
func parseMoviesFileByRule(mf *media_file.MediaFile) (*Movie, error) {
	movieName := utils.FilterTmpSuffix(mf.Filename)
	if mf.IsDisc() { // 光盘类文件使用目录名刮削
		movieName = filepath.Base(filepath.Dir(mf.Path))
	}

	// 过滤无用文件
	if movieName[0:1] == "." || lrace.InArray(config.Collector.SkipFolders, movieName) {
		return nil, errors.New("invalid movie name")
	}

	// 过滤可选字符
	movieName = utils.FilterOptionals(movieName)

	// 使用自定义方法切割
	split := utils.Split(strings.TrimRight(movieName, mf.Suffix))

	// 文件名识别
	nameStart := false
	nameStop := false
	movie := newMovieWithPaths(mf)

	for _, item := range split {
		if resolution := utils.IsResolution(item); resolution != "" {
			nameStop = true
			continue
		}

		if year := utils.IsYear(item); year > 0 {
			movie.Year = year
			nameStop = true
			continue
		}

		if format := utils.IsFormat(item); len(format) > 0 {
			nameStop = true
			continue
		}

		if source := utils.IsSource(item); len(source) > 0 {
			nameStop = true
			continue
		}

		if studio := utils.IsStudio(item); len(studio) > 0 {
			nameStop = true
			continue
		}

		if channel := utils.IsChannel(item); len(channel) > 0 {
			nameStop = true
			continue
		}

		if !nameStart {
			nameStart = true
			nameStop = false
		}

		if !nameStop {
			movie.Title += item + " "
		}
	}

	movie.Title, movie.AliasTitle = utils.SplitTitleAlias(movie.Title)
	movie.ChsTitle, movie.EngTitle = utils.SplitChsEngTitle(movie.Title)
	movie.Title = strings.TrimSpace(movie.Title)
	if movie.Title == "" {
		return nil, fmt.Errorf("file: %s parse title empty: %v", mf.Filename, movie)
	}

	return movie, nil
}

// 调用AI模型解析
func parseMoviesFileByAI(mf *media_file.MediaFile, ruleMovie *Movie) (*Movie, error) {
	if !ai.Enabled() {
		return nil, errors.New("ai disabled")
	}

	rule := map[string]any{}
	if ruleMovie != nil {
		rule = map[string]any{
			"title":       ruleMovie.Title,
			"alias_title": ruleMovie.AliasTitle,
			"chs_title":   ruleMovie.ChsTitle,
			"eng_title":   ruleMovie.EngTitle,
			"year":        ruleMovie.Year,
		}
	}

	result, err := ai.ParseMedia(&ai.ParseInput{
		MediaType: "movie",
		Path:      mf.Path,
		Filename:  mf.Filename,
		Rule:      rule,
	})
	if err != nil {
		return nil, err
	}
	if !ai.ParseUsable(result) {
		return nil, errors.New("ai parse unusable")
	}

	movie := newMovieWithPaths(mf)
	movie.Title = result.Title
	movie.AliasTitle = result.AliasTitle
	movie.ChsTitle = result.ChsTitle
	movie.EngTitle = result.EngTitle
	movie.Year = result.Year

	if movie.Title == "" {
		return nil, errors.New("ai parse title empty")
	}

	utils.Logger.DebugF("ai parse movie success: %s (%d) confidence %.2f", movie.Title, movie.Year, result.Confidence)
	return movie, nil
}

func newMovieWithPaths(mf *media_file.MediaFile) *Movie {
	movie := &Movie{MediaFile: mf}
	if mf.IsDisc() {
		movie.PosterFile = mf.Dir + "/poster.jpg"
		movie.FanArtFile = mf.Dir + "/fanart.jpg"
		movie.ClearLogoFile = mf.Dir + "/clearlogo.png"
		if mf.IsBluRay() {
			movie.NfoFile = mf.Path + "/index.nfo"
		} else if mf.IsDvd() {
			movie.NfoFile = mf.Path + "/VIDEO_TS/VIDEO_TS.nfo"
		}
		return movie
	}

	prefix := mf.Dir + "/" + strings.Replace(mf.Filename, mf.Suffix, "", 1)
	movie.PosterFile = prefix + "-poster.jpg"
	movie.FanArtFile = prefix + "-fanart.jpg"
	movie.ClearLogoFile = prefix + "-clearlogo.png"
	movie.NfoFile = prefix + ".nfo"
	return movie
}

// 加载缓存
func loadMovieCache(mf *media_file.MediaFile) (*Movie, *tmdb.MovieDetail, error) {
	if mf == nil {
		return nil, nil, errors.New("movie file empty")
	}

	movie := newMovieWithPaths(mf)
	movie.checkCacheDir()
	detail, err := movie.loadMovieDetailFromCache()
	if err != nil {
		return nil, nil, err
	}

	if detail != nil {
		movie.fillByDetail(detail)
		return movie, detail, nil
	}

	if movie.hasIdCache() {
		return movie, nil, nil
	}

	return nil, nil, nil
}
