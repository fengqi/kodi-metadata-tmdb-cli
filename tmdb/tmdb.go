package tmdb

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io"
	"io/ioutil"
	"net/http"
)

var Api *tmdb

const (
	ApiSearchTv           = "/search/tv"
	ApiSearchMovie        = "/search/movie"
	ApiTvDetail           = "/tv/%d"
	ApiTvEpisode          = "/tv/%d/season/%d/episode/%d"
	ApiTvAggregateCredits = "/tv/%d/aggregate_credits"
	ApiTvContentRatings   = "/tv/%d/content_ratings"
	ApiMovieDetail        = "/movie/%d"
	ImageW500             = "https://image.tmdb.org/t/p/w500"     // 压缩后的
	ImageOriginal         = "https://image.tmdb.org/t/p/original" // 原始文件
)

func InitTmdb(config *config.Config) {
	Api = &tmdb{
		host:     "https://api.themoviedb.org/3",
		key:      config.ApiKey,
		language: config.Language,
		rating:   "US",
	}
}

func (t *tmdb) request(api string, args map[string]string) ([]byte, error) {
	if args == nil {
		args = make(map[string]string, 0)
	}

	args["api_key"] = t.key
	args["language"] = t.language

	api = t.host + api + "?" + utils.StringMapToQuery(args)
	resp, err := http.Get(api)
	if err != nil {
		utils.Logger.ErrorF("request tmdb: %s err: %v", api, err)
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.Logger.WarningF("request tmdb close body err: %v", err)
		}
	}(resp.Body)

	return ioutil.ReadAll(resp.Body)
}
