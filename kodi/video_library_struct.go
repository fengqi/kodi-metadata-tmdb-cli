package kodi

var EpisodeFields = []string{
	"title",
	"plot",
	"votes",
	"rating",
	"writer",
	"firstaired",
	"playcount",
	"runtime",
	"director",
	"productioncode",
	"season",
	"episode",
	"originaltitle",
	"showtitle",
	"cast",
	"streamdetails",
	"lastplayed",
	"fanart",
	"thumbnail",
	"file",
	"resume",
	"tvshowid",
	"dateadded",
	"uniqueid",
	"art",
	"specialsortseason",
	"specialsortepisode",
	"userrating",
	"seasonid",
	"ratings",
}

var MovieFields = []string{
	"title",
	"genre",
	"year",
	"rating",
	"director",
	"trailer",
	"tagline",
	"plot",
	"plotoutline",
	"originaltitle",
	"lastplayed",
	"playcount",
	"writer",
	"studio",
	"mpaa",
	"cast",
	"country",
	"imdbnumber",
	"runtime",
	"set",
	"showlink",
	"streamdetails",
	"top250",
	"votes",
	"fanart",
	"thumbnail",
	"file",
	"sorttitle",
	"resume",
	"setid",
	"dateadded",
	"tag",
	"art",
	"userrating",
	"ratings",
	"premiered",
	"uniqueid",
}

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

// CleanRequest 清理资料库
type CleanRequest struct {
	ShowDialogs bool   `json:"showdialogs"`
	Content     string `json:"content"`
	Directory   string `json:"directory"`
}

type TvShowDetails struct {
	TvShowId      int    `json:"tvshowid"`
	Title         string `json:"title"`
	OriginalTitle string `json:"originaltitle"`
	File          string `json:"file"`
}

type MovieDetails struct {
	MovieId       int      `json:"movieid"`
	Title         string   `json:"title"`
	OriginalTitle string   `json:"originaltitle"`
	Label         string   `json:"label"`
	Year          int      `json:"year"`
	PlayCount     int      `json:"playcount"`
	LastPlayed    string   `json:"lastplayed"`
	UniqueId      UniqueId `json:"uniqueid"`
}

type UniqueId struct {
	Imdb string `json:"imdb"`
	Tvdb string `json:"tvdb"`
	Tmdb string `json:"tmdb"`
}

type GetEpisodesRequest struct {
	TvShowId   int      `json:"tvshowid"`
	Season     int      `json:"season"`
	Filter     *Filter  `json:"filter"`
	Properties []string `json:"properties"`
}

type GetEpisodesResponse struct {
	Limits   LimitsResult `json:"limits"`
	Episodes []*Episode   `json:"episodes"`
}

type Episode struct {
	EpisodeId  int    `json:"episodeid"`
	TvShowId   int    `json:"tvshowid"`
	Label      string `json:"label"`
	LastPlayed string `json:"lastplayed"`
	PlayCount  int    `json:"playcount"`
	Episode    int    `json:"episode"`
}

type RefreshEpisodeRequest struct {
	EpisodeId int `json:"episodeid"`
}
