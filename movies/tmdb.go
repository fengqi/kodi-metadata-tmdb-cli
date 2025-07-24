package movies

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"github.com/fengqi/lrace"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// getMovieDetail 获取电影详情
func (m *Movie) getMovieDetail() (*tmdb.MovieDetail, error) {
	var err error
	detail := &tmdb.MovieDetail{}
	cacheExpire := false

	// 缓存文件路径
	// todo `tmdb/movie.json` 这种格式后期删除掉
	oldCacheFile := m.GetCacheDir() + "/movie.json"
	cacheFile := m.GetCacheDir() + "/" + m.MediaFile.Filename + ".movie.json"
	if _, err := os.Stat(oldCacheFile); err == nil {
		utils.Logger.DebugF("rename old cache file %s to %s", oldCacheFile, cacheFile)
		if _, err = lrace.CopyFile(oldCacheFile, cacheFile); err == nil {
			_ = os.Remove(oldCacheFile)
		} else {
			utils.Logger.WarningF("rename old cache file %s to %s err: %v", oldCacheFile, cacheFile, err)
		}
	}

	// 从缓存读取
	if cf, err := os.Stat(cacheFile); err == nil {
		utils.Logger.DebugF("get movie detail from cache: %s", cacheFile)

		bytes, err := os.ReadFile(cacheFile)
		if err != nil {
			utils.Logger.WarningF("read movie.json cache: %s err: %v", cacheFile, err)
		}

		if err = json.Unmarshal(bytes, detail); err != nil {
			utils.Logger.WarningF("parse movie: %s file err: %v", cacheFile, err)
		}

		airTime, _ := time.Parse("2006-01-02", detail.ReleaseDate)
		cacheExpire = utils.CacheExpire(cf.ModTime(), airTime)
		detail.FromCache = true
	}

	// 缓存失效，重新搜索
	if detail.Id == 0 || cacheExpire {
		detail.FromCache = false
		movieId := 0

		// todo 兼容 tmdb/id.txt，后期删除
		oldIdFile := m.GetCacheDir() + "/id.txt"
		idFile := m.GetCacheDir() + "/" + m.MediaFile.Filename + ".id.txt"
		if _, err := os.Stat(oldIdFile); err == nil {
			utils.Logger.DebugF("rename old id file %s to %s", oldIdFile, idFile)
			if _, err = lrace.CopyFile(oldIdFile, idFile); err == nil {
				_ = os.Remove(oldIdFile)
			} else {
				utils.Logger.WarningF("rename old id file %s to %s err: %v", oldIdFile, idFile, err)
			}
		}

		if _, err = os.Stat(idFile); err == nil {
			bytes, err := os.ReadFile(idFile)
			if err != nil {
				utils.Logger.WarningF("id file: %s read err: %v", idFile, err)
			} else {
				movieId, _ = strconv.Atoi(strings.Trim(string(bytes), "\r\n "))
				utils.Logger.DebugF("read id from %s", idFile)
			}
		}

		if movieId == 0 {
			SearchResults, err := tmdb.Api.SearchMovie(m.ChsTitle, m.EngTitle, m.Year)
			if err != nil || SearchResults == nil {
				return detail, err
			}

			// 保存movieId
			movieId = SearchResults.Id
			err = os.WriteFile(idFile, []byte(strconv.Itoa(movieId)), 0664)
			if err != nil {
				utils.Logger.ErrorF("save movieId %d to %s err: %v", movieId, idFile, err)
			}
		}

		// 获取详情
		if detail, err = tmdb.Api.GetMovieDetail(movieId); err != nil {
			return nil, err
		}

		// 保存到缓存
		if err = detail.SaveToCache(cacheFile); err != nil {
			utils.Logger.ErrorF("save detail to: %s err: %v", cacheFile, err)
		}
	}

	if detail.Id == 0 || m.Title == "" {
		return nil, err
	}

	return detail, err
}

// downloadImage 下载图片
func (m *Movie) downloadImage(detail *tmdb.MovieDetail) error {
	utils.Logger.DebugF("download %s images", m.Title)

	if len(detail.PosterPath) > 0 {
		err := tmdb.DownloadFile(tmdb.Api.GetImageOriginal(detail.PosterPath), m.PosterFile)
		if err != nil {
			return err
		}
	}

	if len(detail.BackdropPath) > 0 {
		err := tmdb.DownloadFile(tmdb.Api.GetImageOriginal(detail.BackdropPath), m.FanArtFile)
		if err != nil {
			return err
		}
	}

	if detail.Images != nil && len(detail.Images.Logos) > 0 {
		sort.SliceStable(detail.Images.Logos, func(i, j int) bool {
			return detail.Images.Logos[i].VoteAverage > detail.Images.Logos[j].VoteAverage
		})

		image := detail.Images.Logos[0]
		for _, item := range detail.Images.Logos {
			if image.FilePath == "" && item.FilePath != "" {
				image = item
			}
			if item.Iso6391 == "zh" && image.Iso6391 != "zh" { // todo 语言可选
				image = item
				break
			}
		}
		if image.FilePath != "" {
			err := tmdb.DownloadFile(tmdb.Api.GetImageOriginal(image.FilePath), m.ClearLogoFile)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
