package tmdb

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// GetMovieDetail 获取电影详情
func GetMovieDetail(id int) (*MovieDetail, error) {
	utils.Logger.DebugF("get movie detail from tmdb: %d", id)

	req := map[string]string{
		"api_key":            getApiKey(),
		"language":           getLanguage(),
		"append_to_response": "credits",
	}

	api := host + fmt.Sprintf(apiMovieDetail, id) + "?" + utils.StringMapToQuery(req)
	utils.Logger.DebugF("request tmdb: %s", api)

	resp, err := http.Get(api)
	if err != nil {
		utils.Logger.ErrorF("request tmdb: %s err: %v", api, err)
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		utils.Logger.ErrorF("read tmdb response: %s err: %v", api, err)
		return nil, err
	}

	detail := &MovieDetail{}
	err = json.Unmarshal(body, detail)
	if err != nil {
		utils.Logger.ErrorF("parse tmdb response: %s err: %v", api, err)
		return nil, err
	}

	return detail, err
}

// SaveToCache 保存剧集详情到文件
func (d *MovieDetail) SaveToCache(file string) {
	if d.Id == 0 || d.Title == "" {
		return
	}

	utils.Logger.InfoF("save movie detail to: %s", file)

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		utils.Logger.ErrorF("save movie to cache, open_file err: %v", err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	bytes, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		utils.Logger.ErrorF("save movie to cache, marshal struct err: %v", err)
		return
	}

	_, err = f.Write(bytes)
}
