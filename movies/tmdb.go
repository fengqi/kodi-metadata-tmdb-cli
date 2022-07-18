package movies

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"time"
)

func (d *Movie) getMovieDetail() (*tmdb.MovieDetail, error) {
	var err error
	var detail = new(tmdb.MovieDetail)

	// 从缓存读取
	cacheFile := d.GetCacheDir() + "/movie.json"
	if d.IsFile {
		cacheFile = d.GetCacheDir() + "/" + d.OriginTitle + ".movie.json"
	}
	cacheExpire := false
	if cf, err := os.Stat(cacheFile); err == nil {
		utils.Logger.DebugF("get movie detail from cache: %s", cacheFile)

		bytes, err := ioutil.ReadFile(cacheFile)
		if err != nil {
			utils.Logger.WarningF("read movie.json cache: %s err: %v", cacheFile, err)
		}

		err = json.Unmarshal(bytes, detail)
		if err != nil {
			utils.Logger.WarningF("parse movie: %s file err: %v", cacheFile, err)
		}

		airTime, _ := time.Parse("2006-01-02", detail.ReleaseDate)
		cacheExpire = utils.CacheExpire(cf.ModTime(), airTime)
		detail.FromCache = true
	}

	// 缓存失效，重新搜索
	if detail == nil || detail.Id == 0 || cacheExpire {
		detail.FromCache = false
		movieId := detail.Id
		idFile := d.GetCacheDir() + "/id.txt"
		if d.IsFile {
			idFile = d.Dir + "/tmdb/" + d.OriginTitle + ".id.txt"
		}
		if _, err = os.Stat(idFile); err == nil {
			bytes, err := ioutil.ReadFile(idFile)
			if err != nil {
				utils.Logger.WarningF("id file: %s read err: %v", idFile, err)
			} else {
				movieId, _ = strconv.Atoi(strings.Trim(string(bytes), "\r\n "))
			}
		}

		if movieId == 0 {
			SearchResults, err := tmdb.Api.SearchMovie(d.ChsTitle, d.EngTitle, d.Year)
			if err != nil || SearchResults == nil {
				utils.Logger.ErrorF("search title: %s or %s, year: %d failed", d.ChsTitle, d.EngTitle, d.Year)
				return detail, err
			}

			movieId = SearchResults.Id

			// 保存movieId
			err = ioutil.WriteFile(idFile, []byte(strconv.Itoa(movieId)), 0664)
			if err != nil {
				utils.Logger.ErrorF("save movieId %d to %s err: %v", movieId, idFile, err)
			}
		}

		// 获取详情
		detail, err = tmdb.Api.GetMovieDetail(movieId)
		if err != nil {
			utils.Logger.ErrorF("get movie: %d detail err: %v", movieId, err)
			return nil, err
		}

		// 保存到缓存
		d.checkCacheDir()
		detail.SaveToCache(cacheFile)
	}

	if detail.Id == 0 || d.Title == "" {
		return nil, err
	}

	return detail, err
}
