package kodi

import (
	"encoding/json"
	"sync"
)

var (
	vlOnce sync.Once
	vl     *VideoLibrary
)

type VideoLibrary struct {
	scanLimiter   *Limiter
	refreshMovie  *Limiter
	refreshTVShow *Limiter
}

// RefreshMovieRequest 刷新电影请求参数
type RefreshMovieRequest struct {
	MovieId   int    `json:"movieid"`
	IgnoreNfo bool   `json:"ignorenfo"`
	Title     string `json:"title"`
}

// RefreshTVShowRequest 刷新电视剧请求参数
type RefreshTVShowRequest struct {
	TvShowId        int    `json:"tvshowid"`
	IgnoreNfo       bool   `json:"ignorenfo"`
	RefreshEpisodes bool   `json:"refreshepisodes"`
	Title           string `json:"title"`
}

// GetMoviesRequest 获取电影请求参数
type GetMoviesRequest struct {
	Filter     *Filter  `json:"filter"`
	Limit      *Limits  `json:"limits"`
	Properties []string `json:"properties"`
}

// GetMoviesResponse 获取电影返回参数
type GetMoviesResponse struct {
	Limits LimitsResult   `json:"limits"`
	Movies []MovieDetails `json:"movies"`
}

// GetTVShowsRequest 获取电视剧请求参数
type GetTVShowsRequest struct {
	Filter     *Filter  `json:"filter"`
	Limit      *Limits  `json:"limits"`
	Properties []string `json:"properties"`
}

// GetTVShowsResponse 获取电视剧返回参数
type GetTVShowsResponse struct {
	Limits  LimitsResult     `json:"limits"`
	TvShows []*TvShowDetails `json:"tvshows"`
}

// ScanRequest 扫描视频媒体库请求参数
type ScanRequest struct {
	Directory   string `json:"directory"`
	ShowDialogs bool   `json:"showdialogs"`
}

func NewVideoLibrary() *VideoLibrary {
	vlOnce.Do(func() {
		vl = &VideoLibrary{
			scanLimiter:   NewLimiter(300),
			refreshMovie:  NewLimiter(300),
			refreshTVShow: NewLimiter(300),
		}
	})
	return vl
}

// Scans the video sources for new library items
// TODO 异步
func (vl *VideoLibrary) Scan(req *ScanRequest) bool {
	if !vl.scanLimiter.take() {
		return false
	}

	if req == nil {
		req = &ScanRequest{Directory: "", ShowDialogs: false}
	}
	_, err := request(&JsonRpcRequest{
		Method: "VideoLibrary.Scan",
		Params: req,
	})
	return err == nil
}

// GetMovies Retrieve all movies
func (vl *VideoLibrary) GetMovies(req *GetMoviesRequest) *GetMoviesResponse {
	body, err := request(&JsonRpcRequest{Method: "VideoLibrary.GetMovies", Params: req})
	if len(body) == 0 {
		return nil
	}

	resp := &JsonRpcResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		panic(err)
	}

	if resp != nil && resp.Result != nil {
		jsonBytes, _ := json.Marshal(resp.Result)

		moviesResp := &GetMoviesResponse{}
		_ = json.Unmarshal(jsonBytes, moviesResp)

		return moviesResp
	}

	return nil
}

// RefreshMovie Refresh the given movie in the library
func (vl *VideoLibrary) RefreshMovie(req *RefreshMovieRequest) bool {
	if !vl.refreshMovie.take() {
		return false
	}

	_, err := request(&JsonRpcRequest{
		Method: "VideoLibrary.RefreshMovie",
		Params: req,
	})
	return err == nil
}

// GetTVShows Retrieve all tv shows
func (vl *VideoLibrary) GetTVShows(req *GetTVShowsRequest) *GetTVShowsResponse {
	body, err := request(&JsonRpcRequest{Method: "VideoLibrary.GetTVShows", Params: req})
	if len(body) == 0 {
		return nil
	}

	resp := &JsonRpcResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		panic(err)
	}

	if resp != nil && resp.Result != nil {
		jsonBytes, _ := json.Marshal(resp.Result)

		moviesResp := &GetTVShowsResponse{}
		_ = json.Unmarshal(jsonBytes, moviesResp)

		return moviesResp
	}

	return nil
}

// RefreshTVShow Refresh the given tv show in the library
func (vl *VideoLibrary) RefreshTVShow(req *RefreshTVShowRequest) bool {
	if !vl.refreshTVShow.take() {
		return false
	}

	_, err := request(&JsonRpcRequest{
		Method: "VideoLibrary.RefreshTVShow",
		Params: req,
	})
	return err == nil
}
