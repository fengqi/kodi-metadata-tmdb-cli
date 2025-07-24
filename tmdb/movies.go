package tmdb

import (
	"encoding/json"
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"os"
)

// GetMovieDetail 获取电影详情
func (t *Tmdb) GetMovieDetail(id int) (*MovieDetail, error) {
	utils.Logger.DebugF("get movie detail from tmdb: %d", id)

	api := fmt.Sprintf(ApiMovieDetail, id)
	req := map[string]string{
		"append_to_response":     "credits,releases,images",
		"include_image_language": "zh,en,null",
	}

	body, err := t.request(api, req)
	if err != nil {
		return nil, err
	}

	detail := &MovieDetail{}
	err = json.Unmarshal(body, detail)
	if err != nil {
		return nil, err
	}

	return detail, err
}

// SaveToCache 保存剧集详情到文件
func (d *MovieDetail) SaveToCache(file string) error {
	if d.Id == 0 || d.Title == "" {
		return errors.New("id or title empty")
	}

	utils.Logger.InfoF("save movie detail to: %s", file)

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			utils.Logger.WarningF("save movie to cache, close file err: %v", err)
		}
	}(f)

	bytes, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		return err
	}

	_, err = f.Write(bytes)
	return err
}
