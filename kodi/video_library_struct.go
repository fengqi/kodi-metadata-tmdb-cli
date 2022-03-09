package kodi

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

type CleanRequest struct {
	ShowDialogs bool   `json:"showdialogs"`
	Content     string `json:"content"`
	Directory   string `json:"directory"`
}
