package shows

import (
	"encoding/json"
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/common/memcache"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"github.com/fengqi/lrace"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"
)

func (s *Show) getTvDetail() (*tmdb.TvDetail, error) {
	var err error
	var detail = new(tmdb.TvDetail)

	cacheKey := fmt.Sprintf("show:%d", s.TvId)
	if val, ok := memcache.Cache.Get(cacheKey); ok {
		if detail, ok = val.(*tmdb.TvDetail); ok {
			utils.Logger.DebugF("get tv detail from memcache: %d", s.TvId)
			return detail, nil
		}
	}

	// 从缓存读取
	tvCacheFile := s.GetTvCacheDir() + "/tv.json"
	cacheExpire := false
	if cf, err := os.Stat(tvCacheFile); err == nil {
		utils.Logger.DebugF("get tv detail from cache: %s", tvCacheFile)

		bytes, err := os.ReadFile(tvCacheFile)
		if err != nil {
			utils.Logger.WarningF("read tv.json cache: %s err: %v", tvCacheFile, err)
			goto search
		}

		err = json.Unmarshal(bytes, detail)
		if err != nil {
			utils.Logger.WarningF("parse tv file: %s err: %v", tvCacheFile, err)
			_ = os.Remove(tvCacheFile)
			goto search
		}

		airTime, _ := time.Parse("2006-01-02", detail.LastAirDate)
		cacheExpire = utils.CacheExpire(cf.ModTime(), airTime)
		detail.FromCache = true
		s.TvId = detail.Id
	}

search:
	// 缓存失效，重新搜索
	if detail == nil || detail.Id == 0 || cacheExpire {
		detail.FromCache = false
		if s.TvId == 0 {
			SearchResults, err := tmdb.Api.SearchShows(s.ChsTitle, s.EngTitle, s.Year)
			if err != nil || SearchResults == nil {
				utils.Logger.ErrorF("search title: %s year: %d failed", s.Title, s.Year)
				return detail, err
			}

			s.TvId = SearchResults.Id
		}

		// 获取详情
		detail, err = tmdb.Api.GetTvDetail(s.TvId)
		if err != nil || detail == nil || detail.Id == 0 || detail.Name == "" {
			utils.Logger.ErrorF("get tv: %d detail err: %v", s.TvId, err)
			return nil, err
		}

		// 保存到缓存
		detail.SaveToCache(tvCacheFile)
	}

	// 剧集分组：不同的季版本
	if s.GroupId != "" {
		groupDetail, err := s.getTvEpisodeGroupDetail()
		if err == nil {
			detail.TvEpisodeGroupDetail = groupDetail
		}
	}

	if s.TvId > 0 {
		s.CacheTvId()
		cacheKey = fmt.Sprintf("show:%d", s.TvId)
		memcache.Cache.SetDefault(cacheKey, detail)
	}

	return detail, nil
}

func (s *Show) getEpisodeDetail() (*tmdb.TvEpisodeDetail, error) {
	var err error
	var detail = new(tmdb.TvEpisodeDetail)

	cacheFile := fmt.Sprintf("%s/tmdb/s%02de%02d.json", lrace.Ternary(s.SeasonRoot != "", s.SeasonRoot, s.TvRoot), s.Season, s.Episode)
	cacheExpire := false
	if cf, err := os.Stat(cacheFile); err == nil {
		utils.Logger.DebugF("get episode from cache: %s", cacheFile)

		bytes, err := os.ReadFile(cacheFile)
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
		detail, err = tmdb.Api.GetTvEpisodeDetail(s.TvId, s.Season, s.Episode)
		if err != nil {
			return nil, errors.Join(errors.New("get tv episode error"), err)
		}

		if detail == nil || detail.Id == 0 {
			return nil, errors.New(fmt.Sprintf("get episode from tmdb: %d season: %d episode: %d failed", s.TvId, s.Season, s.Episode))
		}

		// 保存到缓存
		detail.SaveToCache(cacheFile)
	}

	if detail.Id == 0 || detail.Name == "" {
		return nil, err
	}

	return detail, err
}

func (s *Show) getTvEpisodeGroupDetail() (*tmdb.TvEpisodeGroupDetail, error) {
	if s.GroupId == "" {
		return nil, nil
	}

	var err error
	var detail = new(tmdb.TvEpisodeGroupDetail)

	// 从缓存读取
	cacheFile := s.SeasonRoot + "/tmdb/group.json"
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
		detail, err = tmdb.Api.GetTvEpisodeGroupDetail(s.GroupId)
		if err != nil {
			utils.Logger.ErrorF("get tv episode group: %s detail err: %v", s.GroupId, err)
			return nil, err
		}

		// 保存到缓存
		detail.SaveToCache(cacheFile)
	}

	return detail, nil
}

// 下载电视剧的相关图片
// TODO 下载失败后，没有重复以及很长一段时间都不会再触发下载
func (s *Show) downloadTvImage(detail *tmdb.TvDetail) {
	if len(detail.PosterPath) > 0 {
		_ = tmdb.DownloadFile(tmdb.Api.GetImageOriginal(detail.PosterPath), s.TvRoot+"/poster.jpg")
	}

	if len(detail.BackdropPath) > 0 {
		_ = tmdb.DownloadFile(tmdb.Api.GetImageOriginal(detail.BackdropPath), s.TvRoot+"/fanart.jpg")
	}

	// TODO group的信息里可能 season poster不全
	if len(detail.Seasons) > 0 {
		for _, item := range detail.Seasons {
			if /*!s.IsCollection &&*/ item.SeasonNumber != s.Season || item.PosterPath == "" {
				continue
			}
			seasonPoster := fmt.Sprintf("season%02d-poster.jpg", item.SeasonNumber)
			_ = tmdb.DownloadFile(tmdb.Api.GetImageOriginal(item.PosterPath), s.TvRoot+"/"+seasonPoster)
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
			if item.Iso6391 == "zh" && image.Iso6391 != "zh" {
				image = item
				break
			}
		}
		if image.FilePath != "" {
			logoFile := s.TvRoot + "/clearlogo.png"
			_ = tmdb.DownloadFile(tmdb.Api.GetImageOriginal(image.FilePath), logoFile)
		}
	}
}

// 下载剧集的相关图片
func (s *Show) downloadEpisodeImage(d *tmdb.TvEpisodeDetail) {
	file := strings.Replace(s.MediaFile.Path, s.MediaFile.Suffix, "-thumb.jpg", 1)
	if len(d.StillPath) > 0 {
		_ = tmdb.DownloadFile(tmdb.Api.GetImageOriginal(d.StillPath), file)
	}
}
