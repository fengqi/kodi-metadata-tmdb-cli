package shows

import (
	"encoding/json"
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/common/memcache"
	"fengqi/kodi-metadata-tmdb-cli/tmdb"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"

	"os"
	"path/filepath"
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
	detail, err = s.loadTvDetailFromCache()
	if err != nil {
		utils.Logger.WarningF("load tv detail cache err: %v", err)
	}
	cacheExpire := detail == nil
	if detail == nil {
		detail = new(tmdb.TvDetail)
	}

	// 缓存失效，重新搜索
	if detail.Id == 0 || cacheExpire {
		detail.FromCache = false
		if s.TvId == 0 {
			searchResults, err := tmdb.Api.SearchShows(s.ChsTitle, s.EngTitle, s.Year)
			if err != nil || searchResults == nil {
				utils.Logger.ErrorF("search title: %s year: %d failed", s.Title, s.Year)
				return detail, err
			}
			s.TvId = searchResults.Id
		}

		// 获取详情
		detail, err = tmdb.Api.GetTvDetail(s.TvId)
		if err != nil || detail == nil || detail.Id == 0 || detail.Name == "" {
			utils.Logger.ErrorF("get tv: %d detail err: %v", s.TvId, err)
			return nil, err
		}

		// 保存到缓存
		detail.SaveToCache(s.GetTvCacheDir() + "/tv.json")
	}

	// 剧集分组详情的附加在Process中处理，避免分组信息混入memcache

	if s.TvId > 0 {
		s.CacheTvId()
		cacheKey = fmt.Sprintf("show:%d", s.TvId)
		memcache.Cache.SetDefault(cacheKey, detail)
	}

	return detail, nil
}

func (s *Show) getEpisodeDetail() (*tmdb.TvEpisodeDetail, error) {
	// 从缓存读取，优先新缓存，未命中再读旧版缓存（读取后会迁移到新缓存文件）
	detail, err := s.loadEpisodeDetailFromCache()
	if err != nil {
		utils.Logger.WarningF("load episode detail cache err: %v", err)
		return nil, err
	}

	if detail == nil {
		detail, err = s.loadLegacyEpisodeDetailFromCache()
		if err != nil {
			utils.Logger.WarningF("load legacy episode detail cache err: %v", err)
			return nil, err
		}
	}

	cacheExpire := detail == nil
	if detail == nil {
		detail = new(tmdb.TvEpisodeDetail)
	}

	// 请求tmdb
	if detail.Id == 0 || cacheExpire {
		detail.FromCache = false

		if s.GroupId != "" {
			// 指定了剧集分组，文件名的season/episode是分组内编号，需通过分组映射获取
			detail, err = s.getEpisodeDetailFromGroup()
		} else {
			detail, err = tmdb.Api.GetTvEpisodeDetail(s.TvId, s.Season, s.Episode)
		}
		if err != nil {
			return nil, errors.Join(errors.New("get tv episode error"), err)
		}

		if detail == nil || detail.Id == 0 {
			return nil, errors.New(fmt.Sprintf("get episode from tmdb: %d season: %d episode: %d failed", s.TvId, s.Season, s.Episode))
		}

		// 保存到缓存
		detail.SaveToCache(s.EpisodeCacheFile())
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

	detail := new(tmdb.TvEpisodeGroupDetail)
	cacheFile := s.SeasonRoot + "/tmdb/group.json"
	cacheExpire := false
	if cf, err := os.Stat(cacheFile); err == nil {
		utils.Logger.DebugF("get tv episode group detail from cache: %s", cacheFile)

		if bytes, err := os.ReadFile(cacheFile); err != nil {
			utils.Logger.WarningF("read group.json cache: %s err: %v", cacheFile, err)
		} else if err = json.Unmarshal(bytes, detail); err != nil {
			utils.Logger.WarningF("parse group.json file: %s err: %v", cacheFile, err)
		}

		if len(detail.Groups) > 0 {
			lastGroup := detail.Groups[len(detail.Groups)-1]
			if len(lastGroup.Episodes) > 0 {
				airTime, _ := time.Parse("2006-01-02", lastGroup.Episodes[len(lastGroup.Episodes)-1].AirDate)
				cacheExpire = utils.CacheExpire(cf.ModTime(), airTime)
			}
		}
		detail.FromCache = true
	}

	// 缓存失效，重新请求
	if detail.Id == "" || cacheExpire {
		detail.FromCache = false
		newDetail, err := tmdb.Api.GetTvEpisodeGroupDetail(s.GroupId)
		if err != nil || newDetail == nil || newDetail.Id == "" {
			utils.Logger.ErrorF("get tv episode group: %s detail err: %v", s.GroupId, err)
			return nil, err
		}
		detail = newDetail

		// 保存到缓存
		if err := os.MkdirAll(filepath.Dir(cacheFile), 0755); err != nil {
			utils.Logger.WarningF("create tv episode group cache dir: %s err: %v", filepath.Dir(cacheFile), err)
		}
		detail.SaveToCache(cacheFile)
	}

	return detail, nil
}

// getEpisodeDetailFromGroup 从剧集分组获取分集详情
// 文件名解析出的season/episode是分组内的编号：season对应分组的order，episode对应组内位置；
// 分组内剧集自带标题、简介、播出日期、剧照、评分等信息，直接取用，
// 再用分组编号覆盖season/episode，保证写入NFO的编号和文件名一致
func (s *Show) getEpisodeDetailFromGroup() (*tmdb.TvEpisodeDetail, error) {
	groupDetail, err := s.getTvEpisodeGroupDetail()
	if err != nil {
		return nil, err
	}
	if groupDetail == nil || len(groupDetail.Groups) == 0 {
		return nil, errors.New("get tv episode group detail empty")
	}

	var group *tmdb.TvEpisodeGroup
	for i := range groupDetail.Groups {
		if groupDetail.Groups[i].Order == s.Season {
			group = &groupDetail.Groups[i]
			break
		}
	}
	if group == nil {
		return nil, fmt.Errorf("season: %d not found in episode group: %s", s.Season, groupDetail.Name)
	}

	episodes := make([]tmdb.TvEpisodeGroupEpisode, len(group.Episodes))
	copy(episodes, group.Episodes)
	sort.SliceStable(episodes, func(i, j int) bool { return episodes[i].Order < episodes[j].Order })

	if s.Episode < 1 || s.Episode > len(episodes) {
		return nil, fmt.Errorf("episode: %d out of range in group: %s", s.Episode, group.Name)
	}
	groupEpisode := &episodes[s.Episode-1]

	return &tmdb.TvEpisodeDetail{
		AirDate:        groupEpisode.AirDate,
		Name:           groupEpisode.Name,
		Overview:       groupEpisode.Overview,
		Id:             groupEpisode.Id,
		ProductionCode: groupEpisode.ProductionCode,
		StillPath:      groupEpisode.StillPath,
		VoteAverage:    groupEpisode.VoteAverage,
		VoteCount:      groupEpisode.VoteCount,
		// 使用分组内的编号
		SeasonNumber:  group.Order,
		EpisodeNumber: s.Episode,
	}, nil
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

// loadTvDetailFromCache 从缓存中加载电视剧详情
func (s *Show) loadTvDetailFromCache() (*tmdb.TvDetail, error) {
	detail := new(tmdb.TvDetail)
	tvCacheFile := s.GetTvCacheDir() + "/tv.json"
	cf, err := os.Stat(tvCacheFile)
	if err != nil {
		return nil, nil
	}

	utils.Logger.DebugF("get tv detail from cache: %s", tvCacheFile)
	bytes, err := os.ReadFile(tvCacheFile)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, detail); err != nil {
		return nil, err
	}

	airTime, _ := time.Parse("2006-01-02", detail.LastAirDate)
	if detail.Id == 0 || utils.CacheExpire(cf.ModTime(), airTime) {
		return nil, nil
	}

	detail.FromCache = true
	s.TvId = detail.Id
	return detail, nil
}

// loadEpisodeDetailFromCache 从缓存中加载剧集详情
func (s *Show) loadEpisodeDetailFromCache() (*tmdb.TvEpisodeDetail, error) {
	detail := new(tmdb.TvEpisodeDetail)
	cacheFile := s.EpisodeCacheFile()
	cf, err := os.Stat(cacheFile)
	if err != nil {
		return nil, nil
	}

	utils.Logger.DebugF("get episode from cache: %s", cacheFile)
	bytes, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, &detail); err != nil {
		return nil, err
	}

	airTime, _ := time.Parse("2006-01-02", detail.AirDate)
	if detail.Id == 0 || utils.CacheExpire(cf.ModTime(), airTime) {
		return nil, nil
	}

	detail.FromCache = true
	s.Season = detail.SeasonNumber
	s.Episode = detail.EpisodeNumber
	return detail, nil
}

// TODO 从旧版缓存中加载剧集详情，加载后重命名到新缓存文件，后续删除该逻辑
func (s *Show) loadLegacyEpisodeDetailFromCache() (*tmdb.TvEpisodeDetail, error) {
	base := s.SeasonRoot
	if base == "" {
		base = s.TvRoot
	}

	if base == "" {
		return nil, errors.New("season root or tv root is empty")
	}

	if s.Season == 0 || s.Episode == 0 {
		return nil, errors.New("season or episode is zero")
	}

	detail := new(tmdb.TvEpisodeDetail)
	cacheFile := fmt.Sprintf("%s/tmdb/s%02de%02d.json", base, s.Season, s.Episode)
	cf, err := os.Stat(cacheFile)
	if err != nil {
		return nil, nil
	}

	utils.Logger.DebugF("get legacy episode from cache: %s", cacheFile)
	bytes, err := os.ReadFile(cacheFile)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(bytes, &detail); err != nil {
		return nil, err
	}

	airTime, _ := time.Parse("2006-01-02", detail.AirDate)
	if detail.Id == 0 || utils.CacheExpire(cf.ModTime(), airTime) {
		return nil, nil
	}

	detail.FromCache = true
	s.Season = detail.SeasonNumber
	s.Episode = detail.EpisodeNumber
	newCacheFile := s.EpisodeCacheFile()
	if newCacheFile != "" && newCacheFile != cacheFile {
		if err = os.Rename(cacheFile, newCacheFile); err != nil {
			utils.Logger.WarningF("rename legacy episode cache %s to %s err: %v", cacheFile, newCacheFile, err)
		}
	}

	return detail, nil
}
