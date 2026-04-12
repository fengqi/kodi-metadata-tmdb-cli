package tmdb

import (
	"fengqi/kodi-metadata-tmdb-cli/common/httpx"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

var Api *Tmdb
var HttpClient *http.Client

const (
	ApiSearchTv           = "/3/search/tv"
	ApiSearchMovie        = "/3/search/movie"
	ApiTvDetail           = "/3/tv/%d"
	ApiTvEpisode          = "/3/tv/%d/season/%d/episode/%d"
	ApiTvAggregateCredits = "/3/tv/%d/aggregate_credits"
	ApiTvContentRatings   = "/3/tv/%d/content_ratings"
	ApiTvEpisodeGroup     = "/3/tv/episode_group/%s"
	ApiMovieDetail        = "/3/movie/%d"
)

func InitTmdb() {
	HttpClient = httpx.NewClient(config.Tmdb.Proxy, 0)
	Api = &Tmdb{
		apiHost:   config.Tmdb.ApiHost,
		apiKey:    config.Tmdb.ApiKey,
		imageHost: config.Tmdb.ImageHost,
		language:  config.Tmdb.Language,
		rating:    config.Tmdb.Rating,
	}
}

// GetImageW500 压缩后的图片
func (t *Tmdb) GetImageW500(path string) string {
	if path == "" {
		return ""
	}
	return Api.imageHost + "/t/p/w500" + path
}

// GetImageOriginal 原始图片
func (t *Tmdb) GetImageOriginal(path string) string {
	if path == "" {
		return ""
	}
	return Api.imageHost + "/t/p/original" + path
}

func (t *Tmdb) request(api string, args map[string]string) ([]byte, error) {
	if args == nil {
		args = make(map[string]string, 0)
	}

	args["api_key"] = t.apiKey
	args["language"] = t.language

	api = t.apiHost + api + "?" + utils.StringMapToQuery(args)
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

// DownloadFile 下载文件, 提供网址和目的地
func DownloadFile(url string, filename string) error {
	if info, err := os.Stat(filename); err == nil && info.Size() > 0 {
		return nil
	}

	utils.Logger.InfoF("download %s to %s", url, filename)

	resp, err := HttpClient.Get(url)
	if err != nil {
		utils.Logger.ErrorF("download: %s err: %v", url, err)
		return err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.Logger.WarningF("download file, close body err: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		utils.Logger.ErrorF("download: %s status code failed: %d", url, resp.StatusCode)
		return nil
	}

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		utils.Logger.ErrorF("download: %s open_file %s err: %v", url, filename, err)
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			utils.Logger.WarningF("download file, close file %s err: %v", filename, err)
		}
	}(f)

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		utils.Logger.ErrorF("save content to image: %s err: %v", filename, err)
		return err
	}

	return nil
}
