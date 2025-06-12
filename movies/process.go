package movies

import (
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/kodi"
	"fengqi/kodi-metadata-tmdb-cli/media_file"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"github.com/fengqi/lrace"
	"path/filepath"
	"strings"
)

// Process 处理扫描到的电影文件
func Process(mf *media_file.MediaFile) error {
	movie, err := parseMoviesFile(mf)
	if err != nil {
		return err
	}

	if movie == nil {
		return errors.New("movie file empty")
	}

	utils.Logger.DebugF("receive movies task: %v", movie)

	movie.checkCacheDir()
	detail, err := movie.getMovieDetail()
	if err != nil {
		return err
	}
	if detail == nil {
		return errors.New("get movie detail empty")
	}

	if !detail.FromCache || !lrace.FileExist(movie.NfoFile) {
		if err := movie.saveToNfo(detail); err != nil {
			return err
		}
		kodi.Rpc.AddRefreshTask(kodi.TaskRefreshMovie, detail.OriginalTitle)
	}

	return movie.downloadImage(detail)
}

// 解析文件, 返回详情：年份、中文名称、英文名称等
func parseMoviesFile(mf *media_file.MediaFile) (*Movie, error) {
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
	} else {
		prefix := mf.Dir + "/" + strings.Replace(mf.Filename, mf.Suffix, "", 1)
		movie.PosterFile = prefix + "-poster.jpg"
		movie.FanArtFile = prefix + "-fanart.jpg"
		movie.ClearLogoFile = prefix + "-clearlogo.png"
		movie.NfoFile = prefix + ".nfo"
	}
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
	if len(movie.Title) == 0 {
		return nil, errors.New(fmt.Sprintf("file: %s parse title empty: %v", mf.Filename, movie))
	}

	return movie, nil
}
