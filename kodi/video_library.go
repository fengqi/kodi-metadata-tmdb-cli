package kodi

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
)

// Scan 扫描媒体库
func (vl *VideoLibrary) Scan(directory string, showDialogs bool) bool {
	_, err := Rpc.request(&JsonRpcRequest{
		Method: "VideoLibrary.Scan",
		Params: &ScanRequest{
			Directory:   directory,
			ShowDialogs: showDialogs,
		},
	})

	return err == nil
}

// IsScanning 是否正在扫描
func (vl *VideoLibrary) IsScanning() bool {
	info := Rpc.XBMC.GetInfoBooleans([]string{"Library.IsScanningVideo"})
	return info != nil && info["Library.IsScanningVideo"]
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
func (vl *VideoLibrary) RefreshMovie(movieId int) bool {

	_, err := Rpc.request(&JsonRpcRequest{
		Method: "VideoLibrary.RefreshMovie",
		Params: &RefreshMovieRequest{
			MovieId:   movieId,
			IgnoreNfo: false,
		},
	})

	return err == nil
}

// GetTVShowsByField 自定义根据字段搜索，例如：GetTVShowsByField("title", "is", "去有风的地方")
// TODO 根据名字获取可能有重复，可等待后续使用uniqueId搜索：https://github.com/xbmc/xbmc/pull/22498
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
func (vl *VideoLibrary) RefreshTVShow(tvShowId int) bool {
	_, err := Rpc.request(&JsonRpcRequest{
		Method: "VideoLibrary.RefreshTVShow",
		Params: &RefreshTVShowRequest{
			TvShowId:        tvShowId,
			IgnoreNfo:       false,
			RefreshEpisodes: false,
		},
	})

	return err == nil
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
func (vl *VideoLibrary) Clean(directory string, showDialogs bool) bool {
	_, err := Rpc.request(&JsonRpcRequest{
		Method: "VideoLibrary.Clean",
		Params: &CleanRequest{
			Directory:   directory,
			ShowDialogs: showDialogs,
			Content:     "video",
		},
	})

	return err == nil
}

// GetEpisodes 获取电视剧剧集列表
func (vl *VideoLibrary) GetEpisodes(tvShowId, season int, filter *Filter) ([]*Episode, error) {
	bytes, err := Rpc.request(&JsonRpcRequest{
		Method: "VideoLibrary.GetEpisodes",
		Params: &GetEpisodesRequest{
			TvShowId:   tvShowId,
			Season:     season,
			Filter:     filter,
			Properties: EpisodeFields,
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

// SetEpisodeDetails 更新剧集详情
// https://kodi.wiki/view/JSON-RPC_API/v12#VideoLibrary.SetEpisodeDetails
func (vl *VideoLibrary) SetEpisodeDetails(episodeId int, params map[string]interface{}) bool {
	params["episodeid"] = episodeId

	bytes, err := Rpc.request(&JsonRpcRequest{
		Method: "VideoLibrary.SetEpisodeDetails",
		Params: params,
	})

	resp := &JsonRpcResponse{}
	_ = json.Unmarshal(bytes, resp)

	return err == nil
}

// SetMovieDetails 更新电影信息
// https://kodi.wiki/view/JSON-RPC_API/v12#VideoLibrary.SetMovieDetails
func (vl *VideoLibrary) SetMovieDetails(movieId int, params map[string]interface{}) bool {
	params["movieid"] = movieId

	bytes, err := Rpc.request(&JsonRpcRequest{
		Method: "VideoLibrary.SetMovieDetails",
		Params: params,
	})

	resp := &JsonRpcResponse{}
	_ = json.Unmarshal(bytes, resp)

	return err == nil
}
