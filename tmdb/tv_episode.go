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

func GetTvEpisodeDetail(tvId, season, episode int) (*TvEpisodeDetail, error) {
	utils.Logger.DebugF("get tv episode detail from tmdb: %d %d-%d", tvId, season, episode)

	if tvId <= 0 || season <= 0 || episode <= 0 {
		return nil, nil
	}

	req := &TvEpisodeRequest{
		ApiKey:           getApiKey(),
		Language:         getLanguage(),
		AppendToResponse: "",
	}

	api := fmt.Sprintf(host+apiTvEpisode, tvId, season, episode)
	api = api + "?" + req.ToQuery()
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

	if resp.StatusCode != 200 {
		utils.Logger.ErrorF("request tmdb status failed: %d err: %v", resp.StatusCode, err)
		return nil, nil
	}

	body, err := ioutil.ReadAll(resp.Body)
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
