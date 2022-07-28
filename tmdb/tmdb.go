package tmdb

import (
	"context"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"golang.org/x/net/proxy"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

var Api *tmdb
var HttpClient *http.Client

const (
	ApiSearchTv           = "/search/tv"
	ApiSearchMovie        = "/search/movie"
	ApiTvDetail           = "/tv/%d"
	ApiTvEpisode          = "/tv/%d/season/%d/episode/%d"
	ApiTvAggregateCredits = "/tv/%d/aggregate_credits"
	ApiTvContentRatings   = "/tv/%d/content_ratings"
	ApiTvEpisodeGroup     = "/tv/episode_group/%s"
	ApiMovieDetail        = "/movie/%d"
	ImageW500             = "https://image.tmdb.org/t/p/w500"     // 压缩后的
	ImageOriginal         = "https://image.tmdb.org/t/p/original" // 原始文件
)

func InitTmdb(config *config.TmdbConfig) {
	HttpClient = getHttpClient(config.Proxy)
	Api = &tmdb{
		host:     config.Host,
		key:      config.ApiKey,
		language: config.Language,
		rating:   config.Rating,
	}
}

func (t *tmdb) request(api string, args map[string]string) ([]byte, error) {
	if args == nil {
		args = make(map[string]string, 0)
	}

	args["api_key"] = t.key
	args["language"] = t.language

	api = t.host + api + "?" + utils.StringMapToQuery(args)
	resp, err := HttpClient.Get(api)
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

// 支持 http 和 socks5 代理
func getHttpClient(proxyConnect string) *http.Client {
	proxyUrl, err := url.Parse(proxyConnect)
	if err != nil || proxyConnect == "" {
		return http.DefaultClient
	}

	if proxyUrl.Scheme == "http" || proxyUrl.Scheme == "https" {
		_ = os.Setenv("HTTP_PROXY", proxyConnect)
		_ = os.Setenv("HTTPS_PROXY", proxyConnect)

		return http.DefaultClient
	}

	if proxyUrl.Scheme == "socks5" || proxyUrl.Scheme == "socks5h" {
		dialContext := func(ctx context.Context, network, addr string) (net.Conn, error) {
			dialer := &net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}

			proxyDialer, err := proxy.FromURL(proxyUrl, dialer)
			if err != nil {
				utils.Logger.WarningF("tmdb new proxy dialer err: %v\n", err)
				return dialer.Dial(network, addr)
			}

			return proxyDialer.Dial(network, addr)
		}

		transport := http.DefaultTransport.(*http.Transport)
		transport.DialContext = dialContext
		return &http.Client{
			Transport: transport,
		}
	}

	return http.DefaultClient
}
