package tmdb

import "fengqi/kodi-metadata-tmdb-cli/config"

var (
	apiKey   = ""
	language = "zh-CN"
	host     = "https://api.themoviedb.org/3"

	apiSearchTv           = "/search/tv"
	apiSearchMovie        = "/search/movie"
	apiTvDetail           = "/tv" // todo 改成下面那种的占位符
	apiTvEpisode          = "/tv/%d/season/%d/episode/%d"
	apiTvAggregateCredits = "/tv/%d/aggregate_credits"
	apiTvContentRatings   = "/tv/%d/content_ratings"
	apiMovieDetail        = "/movie/%d"

	ImageW500     = "https://image.tmdb.org/t/p/w500"     // 压缩后的
	ImageOriginal = "https://image.tmdb.org/t/p/original" // 原始文件
)

func InitTmdb(config *config.Config) {
	apiKey = config.ApiKey
	language = config.Language
}

func getApiKey() string {
	return apiKey
}

func getLanguage() string {
	return language
}
