package kodi

import (
	"encoding/json"
	"sync"
)

var (
	vlOnce sync.Once
	vl     *VideoLibrary
)

type VideoLibrary struct{}

type RefreshMovieRequest struct {
	MovieId   int    `json:"movieid"`
	IgnoreNfo bool   `json:"ignorenfo"`
	Title     string `json:"title"`
}

type GetMoviesRequest struct {
	Filter     *Filter  `json:"filter"`
	Limit      *Limits  `json:"limits"`
	Properties []string `json:"properties"`
}

type GetMoviesResponse struct {
	Id      string                  `json:"id"`
	JsonRpc string                  `json:"jsonrpc"`
	Result  GetMoviesResponseResult `json:"result"`
}

type GetMoviesResponseResult struct {
	Limits LimitsResult    `json:"limits"`
	Movies []*MovieDetails `json:"movies"`
}

type ScanRequest struct {
	Directory   string `json:"directory"`
	ShowDialogs bool   `json:"showdialogs"`
}

func NewVideoLibrary() *VideoLibrary {
	vlOnce.Do(func() {
		vl = &VideoLibrary{}
	})
	return vl
}

func (vl *VideoLibrary) Scan(req *ScanRequest) bool {
	if req == nil {
		req = &ScanRequest{Directory: "", ShowDialogs: false}
	}
	_, err := request(&JsonRpcRequest{
		Method: "VideoLibrary.Scan",
		Params: req,
	})
	return err == nil
}

func (vl *VideoLibrary) GetMovies(req *GetMoviesRequest) *GetMoviesResponse {
	body, err := request(&JsonRpcRequest{Method: "VideoLibrary.GetMovies", Params: req})
	if len(body) == 0 {
		return nil
	}

	resp := &GetMoviesResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		panic(err)
	}

	return resp
}

// RefreshMovie Refresh the given movie in the library
func (vl *VideoLibrary) RefreshMovie(req *RefreshMovieRequest) bool {
	_, err := request(&JsonRpcRequest{
		Method: "VideoLibrary.RefreshMovie",
		Params: req,
	})
	return err == nil
}
