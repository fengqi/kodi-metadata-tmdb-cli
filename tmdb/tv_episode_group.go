package tmdb

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"os"
	"sort"
)

type TvEpisodeGroupDetail struct {
	Id           string           `json:"id"`
	Name         string           `json:"name"`
	Type         int              `json:"type"`
	Network      Network          `json:"network"`
	GroupCount   int              `json:"group_count"`
	EpisodeCount int              `json:"episode_count"`
	Description  string           `json:"description"`
	Groups       []TvEpisodeGroup `json:"groups"`
	FromCache    bool             `json:"from_cache"`
}

type TvEpisodeGroup struct {
	Id       string                  `json:"id"`
	Name     string                  `json:"name"`
	Order    int                     `json:"order"`
	Episodes []TvEpisodeGroupEpisode `json:"episodes"`
	Locked   bool                    `json:"locked"`
}

type TvEpisodeGroupEpisode struct {
	AirDate        string  `json:"air_date"`
	EpisodeNumber  int     `json:"episode_number"`
	Id             int     `json:"id"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	ProductionCode string  `json:"production_code"`
	SeasonNumber   int     `json:"season_number"`
	ShowId         int     `json:"show_id"`
	StillPath      string  `json:"still_path"`
	VoteAverage    float32 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
	Order          int     `json:"order"`
}

type TvEpisodeGroupEpisodeWrapper struct {
	episodes []TvEpisodeGroupEpisode
	by       func(l, r *TvEpisodeGroupEpisode) bool
}

func (ew TvEpisodeGroupEpisodeWrapper) Len() int {
	return len(ew.episodes)
}
func (ew TvEpisodeGroupEpisodeWrapper) Swap(i, j int) {
	ew.episodes[i], ew.episodes[j] = ew.episodes[j], ew.episodes[i]
}
func (ew TvEpisodeGroupEpisodeWrapper) Less(i, j int) bool {
	return ew.by(&ew.episodes[i], &ew.episodes[j])
}

func (d TvEpisodeGroup) SortEpisode() {
	sort.Sort(TvEpisodeGroupEpisodeWrapper{d.Episodes, func(l, r *TvEpisodeGroupEpisode) bool {
		return l.Order < r.Order
	}})
}

func (t *tmdb) GetTvEpisodeGroupDetail(groupId string) (*TvEpisodeGroupDetail, error) {
	utils.Logger.DebugF("get tv episode group detail from tmdb: %s", groupId)

	if groupId == "" {
		return nil, nil
	}

	api := fmt.Sprintf(ApiTvEpisodeGroup, groupId)
	req := map[string]string{
		//"append_to_response": "",
	}

	body, err := t.request(api, req)
	if err != nil {
		utils.Logger.ErrorF("read tmdb response: %s err: %v", api, err)
		return nil, err
	}

	tvResp := &TvEpisodeGroupDetail{}
	err = json.Unmarshal(body, tvResp)
	if err != nil {
		utils.Logger.ErrorF("parse tmdb response: %s err: %v", api, err)
		return nil, err
	}

	return tvResp, err
}

func (d TvEpisodeGroupDetail) SaveToCache(file string) {
	if d.Id == "" || d.Name == "" {
		return
	}

	utils.Logger.InfoF("save tv episode group detail to: %s", file)

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		utils.Logger.ErrorF("save tv episode group detail to cache, open_file err: %v", err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			utils.Logger.WarningF("save tv episode group detail to cache, close file err: %v", err)
		}
	}(f)

	bytes, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		utils.Logger.ErrorF("save  tv episode group detail to cache, marshal struct err: %v", err)
		return
	}

	_, err = f.Write(bytes)
}
