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
				// todo
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
		Properties: []string{"title", "originaltitle", "year"},
	}

	kodiMoviesResp := r.VideoLibrary.GetMovies(kodiMoviesReq)
	if kodiMoviesResp == nil || kodiMoviesResp.Limits.Total == 0 {
		return false
	}

	for _, item := range kodiMoviesResp.Movies {
		utils.Logger.DebugF("find movie by name: %s, refresh detail", item.Title)
		r.VideoLibrary.RefreshMovie(item.MovieId)
	}

	return true
}

func (r *JsonRpc) RefreshShows(name string) bool {
	kodiShowsResp := r.VideoLibrary.GetTVShowsByField("originaltitle", "contains", name)
	if kodiShowsResp == nil || kodiShowsResp.Limits.Total == 0 {
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
		}

		for _, item := range episodes {
			r.VideoLibrary.RefreshEpisode(item.EpisodeId)
		}
	}

	return true
}
