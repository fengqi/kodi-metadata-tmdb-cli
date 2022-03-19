package tmdb

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"os"
)

type TvEpisodeRequest struct {
	ApiKey           string `json:"api_key"`
	Language         string `json:"language"`
	AppendToResponse string `json:"append_to_response"`
}

type TvEpisodeDetail struct {
	AirDate        string       `json:"air_date"`
	Crew           []Crew       `json:"crew"`
	GuestStars     []GuestStars `json:"guest_stars"`
	Name           string       `json:"name"`
	Overview       string       `json:"overview"`
	Id             int          `json:"id"`
	ProductionCode string       `json:"production_code"`
	SeasonNumber   int          `json:"season_number"`
	EpisodeNumber  int          `json:"episode_number"`
	StillPath      string       `json:"still_path"`
	VoteAverage    float32      `json:"vote_average"`
	VoteCount      int          `json:"vote_count"`
	FromCache      bool         `json:"from_cache"`
}

type Crew struct {
	Id          int    `json:"id"`
	CreditId    string `json:"credit_id"`
	Name        string `json:"name"`
	Department  string `json:"department"`
	Job         string `json:"job"`
	ProfilePath string `json:"profile_path"`
}

type GuestStars struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	CreditId    string `json:"credit_id"`
	Character   string `json:"character"`
	Order       int    `json:"order"`
	ProfilePath string `json:"profile_path"`
}

func (t *tmdb) GetTvEpisodeDetail(tvId, season, episode int) (*TvEpisodeDetail, error) {
	utils.Logger.DebugF("get tv episode detail from tmdb: %d %d-%d", tvId, season, episode)

	if tvId <= 0 || season <= 0 || episode <= 0 {
		return nil, nil
	}

	api := fmt.Sprintf(ApiTvEpisode, tvId, season, episode)
	req := map[string]string{
		"append_to_response": "",
	}

	body, err := t.request(api, req)
	if err != nil {
		utils.Logger.ErrorF("read tmdb response: %s err: %v", api, err)
		return nil, err
	}

	tvResp := &TvEpisodeDetail{}
	err = json.Unmarshal(body, tvResp)
	if err != nil {
		utils.Logger.ErrorF("parse tmdb response: %s err: %v", api, err)
		return nil, err
	}

	return tvResp, err
}

func (r *TvEpisodeRequest) ToQuery() string {
	return fmt.Sprintf(
		"api_key=%s&language=%s&append_to_response=%s",
		r.ApiKey,
		r.Language,
		r.AppendToResponse,
	)
}

// SaveToCache 保存单集详情到文件
func (d *TvEpisodeDetail) SaveToCache(file string) {
	if d.Id == 0 || d.Name == "" {
		return
	}

	utils.Logger.InfoF("save episode detail to: %s", file)

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		utils.Logger.ErrorF("save to episode file: %s err: %v", file, err)
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
		utils.Logger.ErrorF("save to episode, marshal err: %v", err)
		return
	}

	_, err = f.Write(bytes)
}
