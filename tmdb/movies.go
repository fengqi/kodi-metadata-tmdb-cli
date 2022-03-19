package tmdb

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"os"
)

// GetMovieDetail 获取电影详情
func (t *tmdb) GetMovieDetail(id int) (*MovieDetail, error) {
	utils.Logger.DebugF("get movie detail from tmdb: %d", id)

	api := fmt.Sprintf(ApiMovieDetail, id)
	req := map[string]string{
		"append_to_response": "credits",
	}

	body, err := t.request(api, req)
	if err != nil {
		utils.Logger.ErrorF("get movie detail err: %d %v", id, err)
		return nil, err
	}

	detail := &MovieDetail{}
	err = json.Unmarshal(body, detail)
	if err != nil {
		utils.Logger.ErrorF("parse movie detail err: %d %v", id, err)
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
