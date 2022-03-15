package kodi

import (
	"encoding/json"
	"fmt"
	"sync"
)

var (
	vlOnce sync.Once
	vl     *VideoLibrary
)

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
func (vl *VideoLibrary) Scan(req *ScanRequest) bool {
	if !vl.scanLimiter.take() {
		return false
	}

	if req == nil {
		req = &ScanRequest{Directory: "", ShowDialogs: false}
	}

	return RpcQueue.addTask("scan video library", &JsonRpcRequest{
		Method: "VideoLibrary.Scan",
		Params: req,
	})
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

	return RpcQueue.addTask(fmt.Sprintf("refresh movie %d", req.MovieId), &JsonRpcRequest{
		Method: "VideoLibrary.RefreshMovie",
		Params: req,
	})
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

	return RpcQueue.addTask(fmt.Sprintf("refresh tvshow %d", req.TvShowId), &JsonRpcRequest{
		Method: "VideoLibrary.RefreshTVShow",
		Params: req,
	})
}

// Clean 清理资料库
func (vl *VideoLibrary) Clean(req *CleanRequest) bool {
	if req == nil {
		req = &CleanRequest{Directory: "", ShowDialogs: false, Content: "video"}
	}

	return RpcQueue.addTask("clean video library", &JsonRpcRequest{
		Method: "VideoLibrary.Clean",
		Params: req,
	})
}
