package kodi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var Rpc *JsonRpc

func InitKodi(c config.KodiConfig) {
	Rpc = &JsonRpc{
		config: c,
		queue:  make(map[string]*JsonRpcRequest, 0),
		lock:   &sync.RWMutex{},
		VideoLibrary: &VideoLibrary{
			scanLimiter:   NewLimiter(300),
			refreshMovie:  NewLimiter(300),
			refreshTVShow: NewLimiter(300),
		},
	}
}

func (r *JsonRpc) AddTask(name string, req *JsonRpcRequest) bool {
	if !r.config.Enable {
		return false
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.queue[name]; !ok {
		r.queue[name] = req
	}

	return true
}

func (r *JsonRpc) RunNotify() {
	if !r.config.Enable {
		return
	}

	task := func() {
		if len(r.queue) == 0 {
			return
		}

		if !r.Ping() {
			return
		}

		r.lock.RLock()
		defer r.lock.RUnlock()

		utils.Logger.DebugF("kodi request queue size: %d", len(r.queue))
		for k, req := range r.queue {
			_, err := r.request(req)
			if err != nil {
				utils.Logger.ErrorF("request kodi %s err: %v", req.Method, err)
				continue
			}

			delete(r.queue, k)

			time.Sleep(time.Second * 30)
		}
	}

	ticker := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-ticker.C:
			task()
		}
	}
}

func (r *JsonRpc) Ping() bool {
	_, err := r.request(&JsonRpcRequest{Method: "JSONRPC.Ping"})
	if err != nil {
		utils.Logger.WarningF("ping kodi err: %v", err)
	}
	return err == nil
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
		r.VideoLibrary.RefreshMovie(&RefreshMovieRequest{MovieId: item.MovieId, IgnoreNfo: false})
	}

	return true
}

func (r *JsonRpc) RefreshShows(name string) bool {
	kodiTvShowsReq := &GetTVShowsRequest{
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

	kodiShowsResp := r.VideoLibrary.GetTVShows(kodiTvShowsReq)
	if kodiShowsResp == nil || kodiShowsResp.Limits.Total == 0 {
		return false
	}

	for _, item := range kodiShowsResp.TvShows {
		utils.Logger.DebugF("find tv shows by name :%s, refresh detail", item.Title)
		kodiRefreshReq := &RefreshTVShowRequest{TvShowId: item.TvShowId, IgnoreNfo: false, RefreshEpisodes: true}
		r.VideoLibrary.RefreshTVShow(kodiRefreshReq)
	}

	return true
}

// 发送json rpc请求
func (r *JsonRpc) request(rpcReq *JsonRpcRequest) ([]byte, error) {
	utils.Logger.InfoF("request kodi: %s", rpcReq.Method)

	if rpcReq.JsonRpc == "" {
		rpcReq.JsonRpc = "2.0"
	}

	if rpcReq.Id == "" {
		rpcReq.Id = time.Now().String()
	}

	jsonBytes, err := json.Marshal(rpcReq)
	if err != nil {
		utils.Logger.WarningF("request kodi marshal err: %v", err)
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, r.config.JsonRpc, bytes.NewReader(jsonBytes))
	if err != nil {
		utils.Logger.WarningF("request kodi NewRequest err: %v", err)
		return nil, err
	}

	req.SetBasicAuth(r.config.Username, r.config.Password)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout:   time.Duration(r.config.Timeout) * time.Second,
		Transport: http.DefaultTransport,
	}

	resp, err := client.Do(req)
	if err != nil {
		utils.Logger.WarningF("request kodi Do err: %v", err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.Logger.WarningF("request kodi closeBody err: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}
