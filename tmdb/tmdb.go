package tmdb

import (
	"fengqi/kodi-metadata-tmdb-cli/common/httpx"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"io"

	"net/http"
	"os"
	"strconv"
	"time"
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
	HttpClient = httpx.NewClient(config.Tmdb.Proxy, config.Tmdb.TimeoutSeconds)
	Api = &Tmdb{
		apiHost:    config.Tmdb.ApiHost,
		apiKey:     config.Tmdb.ApiKey,
		imageHost:  config.Tmdb.ImageHost,
		language:   config.Tmdb.Language,
		rating:     config.Tmdb.Rating,
		retryCount: config.Tmdb.RetryCount,
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

// statusError 非成功状态码的错误，携带状态码和 429 时的 Retry-After
type statusError struct {
	code       int
	retryAfter time.Duration
}

func (e *statusError) Error() string {
	return fmt.Sprintf("request tmdb status code: %d", e.code)
}

// request 请求 TMDB 接口
// 网络错误、429、5xx 在 retry_count 范围内重试，指数退避（1s/2s/4s...封顶32s）
// 429 优先等待 Retry-After（封顶10s），避免单消费者协程被长时间挂起；其余 4xx 不重试直接失败
func (t *Tmdb) request(api string, args map[string]string) ([]byte, error) {
	if args == nil {
		args = make(map[string]string, 0)
	}

	args["api_key"] = t.apiKey
	args["language"] = t.language

	api = t.apiHost + api + "?" + utils.StringMapToQuery(args)

	for attempt := 0; ; attempt++ {
		body, err := t.requestOnce(api)
		if err == nil {
			return body, nil
		}

		// 可重试判定：网络错误、429、5xx
		statusErr, ok := err.(*statusError)
		retryable := !ok || statusErr.code == http.StatusTooManyRequests || statusErr.code >= 500
		if attempt >= t.retryCount || !retryable {
			return nil, err
		}

		wait := time.Second << uint(min(attempt, 5))
		if ok && statusErr.code == http.StatusTooManyRequests && statusErr.retryAfter > 0 {
			wait = min(statusErr.retryAfter, 10*time.Second)
		}

		utils.Logger.WarningF("request tmdb err: %v, retry %d/%d after %s", err, attempt+1, t.retryCount, wait)
		time.Sleep(wait)
	}
}

// requestOnce 发起一次请求，非 200 状态码返回 statusError
func (t *Tmdb) requestOnce(api string) ([]byte, error) {
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

	if resp.StatusCode != http.StatusOK {
		statusErr := &statusError{code: resp.StatusCode}
		if retryAfter := resp.Header.Get("Retry-After"); retryAfter != "" {
			if sec, err := strconv.Atoi(retryAfter); err == nil && sec > 0 {
				statusErr.retryAfter = time.Duration(sec) * time.Second
			}
		}
		return nil, statusErr
	}

	return io.ReadAll(resp.Body)
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
