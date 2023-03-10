package shows

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

// GetTvDetail 获取详情 TODO err 只返回不输出，由调用方自行处理
func (d *Dir) getTvDetail() (*tmdb.TvDetail, error) {
	var err error
	var detail = new(tmdb.TvDetail)

	// 从缓存读取
	tvCacheFile := d.GetCacheDir() + "/tv.json"
	cacheExpire := false
	if cf, err := os.Stat(tvCacheFile); err == nil {
		utils.Logger.DebugF("get tv detail from cache: %s", tvCacheFile)

		bytes, err := ioutil.ReadFile(tvCacheFile)
		if err != nil {
			utils.Logger.WarningF("read tv.json cache: %s err: %v", tvCacheFile, err)
		}

		err = json.Unmarshal(bytes, detail)
		if err != nil {
			utils.Logger.WarningF("parse tv file: %s err: %v", tvCacheFile, err)
		}

		airTime, _ := time.Parse("2006-01-02", detail.LastAirDate)
		cacheExpire = utils.CacheExpire(cf.ModTime(), airTime)
		detail.FromCache = true
	}

	// 缓存失效，重新搜索
	if detail == nil || detail.Id == 0 || cacheExpire {
		detail.FromCache = false
		tvId := detail.Id
		idFile := d.GetCacheDir() + "/id.txt"
		if _, err = os.Stat(idFile); err == nil {
			bytes, err := ioutil.ReadFile(idFile)
			if err != nil {
				utils.Logger.WarningF("id file: %s read err: %v", idFile, err)
			} else {
				tvId, _ = strconv.Atoi(strings.Trim(string(bytes), "\r\n "))
			}
		}

		if tvId == 0 {
			SearchResults, err := tmdb.Api.SearchShows(d.ChsTitle, d.EngTitle, d.Year)
			if err != nil || SearchResults == nil {
				utils.Logger.ErrorF("search title: %s year: %d failed", d.Title, d.Year)
				return detail, err
			}

			tvId = SearchResults.Id

			// 保存tvId
			err = ioutil.WriteFile(idFile, []byte(strconv.Itoa(tvId)), 0664)
			if err != nil {
				utils.Logger.ErrorF("save tvId %d to %s err: %v", tvId, idFile, err)
			}
		}

		// 获取详情
		detail, err = tmdb.Api.GetTvDetail(tvId)
		if err != nil {
			utils.Logger.ErrorF("get tv: %d detail err: %v", tvId, err)
			return nil, err
		}

		// 保存到缓存
		d.checkCacheDir()
		detail.SaveToCache(tvCacheFile)
	}

	if detail.Id == 0 || detail.Name == "" {
		return nil, err
	}

	// 剧集分组：不同的季版本
	if d.GroupId != "" {
		groupDetail, err := d.getTvEpisodeGroupDetail()
		if err == nil {
			detail.TvEpisodeGroupDetail = groupDetail
		}
	}

	d.TvId = detail.Id

	return detail, err
}

func (f *File) getTvEpisodeDetail() (*tmdb.TvEpisodeDetail, error) {
	var err error
	var detail = new(tmdb.TvEpisodeDetail)

	cacheFile := f.getCacheDir() + "/" + f.SeasonEpisode + ".json"
	cacheExpire := false
	if cf, err := os.Stat(cacheFile); err == nil {
		utils.Logger.DebugF("get episode from cache: %s", cacheFile)

		bytes, err := ioutil.ReadFile(cacheFile)
		if err != nil {
			utils.Logger.WarningF("read episode cache: %s err: %v", cacheFile, err)
		}

		err = json.Unmarshal(bytes, &detail)
		if err != nil {
			utils.Logger.WarningF("parse episode cache: %s err: %v", cacheFile, err)
		}

		airTime, _ := time.Parse("2006-01-02", detail.AirDate)
		cacheExpire = utils.CacheExpire(cf.ModTime(), airTime)
		detail.FromCache = true
	}

	// 请求tmdb
	if detail == nil || detail.Id == 0 || cacheExpire {
		detail.FromCache = false
		detail, err = tmdb.Api.GetTvEpisodeDetail(f.TvId, f.Season, f.Episode)
		if err != nil {
			utils.Logger.ErrorF("get tv episode error %v", err)
			return nil, err
		}

		if detail == nil {
			utils.Logger.WarningF("get episode from tmdb: %d season: %d episode: %d failed", f.TvId, f.Season, f.Episode)
			return detail, err
		}

		// 保存到缓存
		detail.SaveToCache(cacheFile)
	}

	if detail.Id == 0 || detail.Name == "" {
		return nil, err
	}

	return detail, err
}

func (d *Dir) getTvEpisodeGroupDetail() (*tmdb.TvEpisodeGroupDetail, error) {
	if d.GroupId == "" {
		return nil, nil
	}

	var err error
	var detail = new(tmdb.TvEpisodeGroupDetail)

	// 从缓存读取
	cacheFile := d.GetCacheDir() + "/group.json"
	cacheExpire := false
	if cf, err := os.Stat(cacheFile); err == nil {
		utils.Logger.DebugF("get tv episode group detail from cache: %s", cacheFile)

		bytes, err := ioutil.ReadFile(cacheFile)
		if err != nil {
			utils.Logger.WarningF("read group.json cache: %s err: %v", cacheFile, err)
		}

		err = json.Unmarshal(bytes, detail)
		if err != nil {
			utils.Logger.WarningF("parse group.json file: %s err: %v", cacheFile, err)
		}

		airTime, _ := time.Parse("2006-01-02", detail.Groups[len(detail.Groups)-1].Episodes[0].AirDate)
		cacheExpire = utils.CacheExpire(cf.ModTime(), airTime)
		detail.FromCache = true
	}

	// 缓存失效，重新搜索
	if detail == nil || detail.Id == "" || cacheExpire {
		detail.FromCache = false
		detail, err = tmdb.Api.GetTvEpisodeGroupDetail(d.GroupId)
		if err != nil {
			utils.Logger.ErrorF("get tv episode group: %s detail err: %v", d.GroupId, err)
			return nil, err
		}

		// 保存到缓存
		d.checkCacheDir()
		detail.SaveToCache(cacheFile)
	}

	return detail, nil
}
