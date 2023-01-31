package kodi

import (
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"strconv"
	"strings"
	"time"
)

// AddRefreshTask 添加刷新数据任务
func (r *JsonRpc) AddRefreshTask(task TaskRefresh, value string) {
	if !r.config.Enable {
		return
	}

	utils.Logger.DebugF("AddRefreshTask %d %s", task, value)

	r.refreshLock.Lock()
	defer r.refreshLock.Unlock()

	taskName := fmt.Sprintf("%.02d|-|%s", task, value)
	if _, ok := r.refreshQueue[taskName]; !ok {
		r.refreshQueue[taskName] = struct{}{}
	}

	return
}

// ConsumerRefreshTask 消费刷新数据任务
func (r *JsonRpc) ConsumerRefreshTask() {
	if !r.config.Enable {
		return
	}

	for {
		if len(r.refreshQueue) == 0 || !r.Ping() || r.VideoLibrary.IsScanning() {
			time.Sleep(time.Second * 30)
			continue
		}

		for queue, _ := range r.refreshQueue {
			_task, _ := strconv.Atoi(queue[0:2])
			task := TaskRefresh(_task)

			r.refreshLock.Lock()

			switch task {
			case TaskRefreshTVShow:
				r.RefreshShows(queue[5:])
				break
			case TaskRefreshEpisode:
				r.RefreshEpisode(queue[5:])
				break
			case TaskRefreshMovie:
				r.RefreshMovie(queue[5:])
				break
			}

			delete(r.refreshQueue, queue)
			r.refreshLock.Unlock()
		}
	}
}

func (r *JsonRpc) RefreshMovie(name string) bool {
	kodiMoviesReq := &GetMoviesRequest{
		Filter: &Filter{
			Field:    "originaltitle",
			Operator: "is",
			Value:    name,
		},
		Limit: &Limits{
			Start: 0,
			End:   5,
		},
		Properties: MovieFields,
	}

	kodiMoviesResp := r.VideoLibrary.GetMovies(kodiMoviesReq)
	if kodiMoviesResp == nil || kodiMoviesResp.Limits.Total == 0 {
		r.VideoLibrary.Scan("", true) // 同剧集，新电影，刷新变扫描库
		return false
	}

	for _, item := range kodiMoviesResp.Movies {
		utils.Logger.DebugF("find movie by name: %s, refresh detail", item.Title)
		refresh := r.VideoLibrary.RefreshMovie(item.MovieId)
		time.Sleep(time.Second * 2) // 不sleep就返回-32602，原因未知

		// 补回丢失的观看记录
		if refresh && item.PlayCount > 0 && item.LastPlayed != "" {
			newKodiMoviesResp := r.VideoLibrary.GetMovies(kodiMoviesReq)
			if newKodiMoviesResp == nil || newKodiMoviesResp.Limits.Total == 0 {
				continue
			}
			for _, item2 := range newKodiMoviesResp.Movies {
				// 名称搜索有概率重复，所以用TMDB匹配，期望将来可以用TMDB过滤
				// 极端情况下，媒体库中有同一个电影的两个版本，也会错乱，需要辅助其他信息确定唯一
				if item2.UniqueId.Tmdb != "" && item2.UniqueId.Tmdb == item.UniqueId.Tmdb {
					played := map[string]interface{}{
						"playcount":  1,
						"lastplayed": item.LastPlayed,
					}
					r.VideoLibrary.SetMovieDetails(item2.MovieId, played)
				}
			}
		}
	}

	return true
}

func (r *JsonRpc) RefreshShows(name string) bool {
	kodiShowsResp := r.VideoLibrary.GetTVShowsByField("originaltitle", "contains", name)
	if kodiShowsResp == nil || kodiShowsResp.Limits.Total == 0 {
		r.VideoLibrary.Scan("", true) // 新剧集，刷新变扫描库，不知道在Kodi的路径所以路径为空
		return false
	}

	for _, item := range kodiShowsResp.TvShows {
		utils.Logger.DebugF("refresh tv shows %s", item.Title)
		r.VideoLibrary.RefreshTVShow(item.TvShowId)
	}

	return true
}

func (r *JsonRpc) RefreshEpisode(taskVal string) bool {
	taskInfo := strings.Split(taskVal, "|-|")
	if len(taskInfo) != 3 {
		return false
	}

	kodiShowsResp := r.VideoLibrary.GetTVShowsByField("originaltitle", "contains", taskInfo[0])
	if kodiShowsResp == nil || kodiShowsResp.Limits.Total == 0 {
		return false
	}

	for _, item := range kodiShowsResp.TvShows {
		filter := &Filter{
			Field:    "episode",
			Operator: "is",
			Value:    taskInfo[2],
		}
		season, err := strconv.Atoi(taskInfo[1])
		if err != nil || season == 0 {
			continue
		}

		episodes, err := r.VideoLibrary.GetEpisodes(item.TvShowId, season, filter)
		if err != nil {
			continue
		}

		// 新增的剧集，需要扫描库
		if episodes == nil || len(episodes) == 0 {
			r.AddScanTask(item.File)
			continue
		}

		episode := episodes[0]
		refresh := r.VideoLibrary.RefreshEpisode(episode.EpisodeId)
		time.Sleep(time.Second * 2)

		// 刷新可能会导致episodeId变化，丢失观看记录，需要补回来
		if refresh && episode.PlayCount > 0 && episode.LastPlayed != "" {
			newFilter := &Filter{
				Field:    "episode",
				Operator: "is",
				Value:    strconv.Itoa(episode.Episode),
			}
			newEpisodes, err := r.VideoLibrary.GetEpisodes(episode.TvShowId, season, newFilter)
			if err != nil || len(newEpisodes) == 0 {
				continue
			}
			played := map[string]interface{}{
				"playcount":  1,
				"lastplayed": episode.LastPlayed,
			}
			r.VideoLibrary.SetEpisodeDetails(newEpisodes[0].EpisodeId, played)
		}
	}

	return true
}
