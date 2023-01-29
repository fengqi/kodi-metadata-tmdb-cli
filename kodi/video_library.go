package kodi

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
)

// Scans the video sources for new library items
func (vl *VideoLibrary) Scan(req *ScanRequest) bool {
	if !vl.scanLimiter.take() {
		return false
	}

	if req == nil {
		req = &ScanRequest{Directory: "", ShowDialogs: false}
	}

	return Rpc.AddTask("scan video library", &JsonRpcRequest{
		Method: "VideoLibrary.Scan",
		Params: req,
	})
}

// GetMovies Retrieve all movies
func (vl *VideoLibrary) GetMovies(req *GetMoviesRequest) *GetMoviesResponse {
	body, err := Rpc.request(&JsonRpcRequest{Method: "VideoLibrary.GetMovies", Params: req})
	if len(body) == 0 {
		return nil
	}

	resp := &JsonRpcResponse{}
	err = json.Unmarshal(body, resp)
	if err != nil {
		return nil
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

	return Rpc.AddTask(fmt.Sprintf("refresh movie %d", req.MovieId), &JsonRpcRequest{
		Method: "VideoLibrary.RefreshMovie",
		Params: req,
	})
}

// GetTVShowsByField 自定义根据字段搜索，例如：GetTVShowsByField("title", "is", "去有风的地方")
func (vl *VideoLibrary) GetTVShowsByField(field, operator, value string) *GetTVShowsResponse {
	req := &GetTVShowsRequest{
		Filter: &Filter{
			Field:    field,
			Operator: operator,
			Value:    value,
		},
		Limit: &Limits{
			Start: 0,
			End:   5,
		},
		Properties: []string{"title", "originaltitle", "year", "file"},
	}

	body, err := Rpc.request(&JsonRpcRequest{Method: "VideoLibrary.GetTVShows", Params: req})
	if err != nil {
		utils.Logger.WarningF("GetTVShowsByField(%s, %s, %s) err: %v", field, operator, value, err)
		return nil
	}

	resp := &JsonRpcResponse{}
	_ = json.Unmarshal(body, resp)
	if resp != nil && resp.Result != nil {
		jsonBytes, _ := json.Marshal(resp.Result)

		moviesResp := &GetTVShowsResponse{}
		_ = json.Unmarshal(jsonBytes, moviesResp)

		return moviesResp
	}

	return nil
}

// GetTVShows Retrieve all tv shows
func (vl *VideoLibrary) GetTVShows(req *GetTVShowsRequest) *GetTVShowsResponse {
	body, err := Rpc.request(&JsonRpcRequest{Method: "VideoLibrary.GetTVShows", Params: req})
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

	return Rpc.AddTask(fmt.Sprintf("refresh tvshow %d", req.TvShowId), &JsonRpcRequest{
		Method: "VideoLibrary.RefreshTVShow",
		Params: req,
	})
}

// RefreshEpisode 刷新剧集信息
// https://kodi.wiki/view/JSON-RPC_API/v13#VideoLibrary.RefreshEpisode
func (vl *VideoLibrary) RefreshEpisode(episodeId int) bool {
	_, err := Rpc.request(&JsonRpcRequest{
		Method: "VideoLibrary.RefreshEpisode",
		Params: &RefreshEpisodeRequest{
			EpisodeId: episodeId,
		},
	})

	return err == nil
}

// Clean 清理资料库
func (vl *VideoLibrary) Clean(req *CleanRequest) bool {
	if req == nil {
		req = &CleanRequest{Directory: "", ShowDialogs: false, Content: "video"}
	}

	return Rpc.AddTask("clean video library", &JsonRpcRequest{
		Method: "VideoLibrary.Clean",
		Params: req,
	})
}

// GetEpisodes 获取电视剧剧集列表
func (vl *VideoLibrary) GetEpisodes(tvShowId, season int, filter *Filter) ([]*Episode, error) {
	bytes, err := Rpc.request(&JsonRpcRequest{
		Method: "VideoLibrary.GetEpisodes",
		Params: &GetEpisodesRequest{
			TvShowId: tvShowId,
			Season:   season,
			Filter:   filter,
		},
	})

	resp := &JsonRpcResponse{}
	err = json.Unmarshal(bytes, resp)
	if err != nil || resp.Result == nil {
		return nil, err
	}

	jsonBytes, _ := json.Marshal(resp.Result)
	episodesResp := &GetEpisodesResponse{}
	err = json.Unmarshal(jsonBytes, episodesResp)
	if err != nil {
		return nil, err
	}

	return episodesResp.Episodes, nil
}
