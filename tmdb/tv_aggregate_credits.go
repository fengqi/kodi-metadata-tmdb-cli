package tmdb

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type TvAggregateCredits struct {
	Id   string   `json:"id"`
	Cast []TvCast `json:"cast"`
	Crew []TvCrew `json:"crew"`
}

type TvCast struct {
	Adult              bool    `json:"adult"`
	Gender             int     `json:"gender"`
	Id                 int     `json:"id"`
	KnownForDepartment string  `json:"known_for_department"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name"`
	Popularity         float32 `json:"popularity"`
	ProfilePath        string  `json:"profile_path"`
	Roles              []Role  `json:"roles"`
	TotalEpisodeCount  int     `json:"total_episode_count"`
	Order              int     `json:"order"`
}

type Role struct {
	CreditId     string `json:"credit_id"`
	Character    string `json:"character"`
	EpisodeCount int    `json:"episode_count"`
}

type Job struct {
	CreditId     string `json:"credit_id"`
	Job          string `json:"job"`
	EpisodeCount int    `json:"episode_count"`
}

type TvCrew struct {
	Adult              bool    `json:"adult"`
	Gender             int     `json:"gender"`
	Id                 int     `json:"id"`
	KnownForDepartment string  `json:"known_for_department"`
	Name               string  `json:"name"`
	OriginalName       string  `json:"original_name"`
	Popularity         float32 `json:"popularity"`
	ProfilePath        string  `json:"profile_path"`
	Jobs               []Job   `json:"jobs"`
	Department         string  `json:"department"`
	TotalEpisodeCount  int     `json:"total_episode_count"`
}

func GetTvAggregateCredits(tvId int) (*TvAggregateCredits, error) {
	utils.Logger.DebugF("get tv aggregate credits from tmdb: %d", tvId)

	m := map[string]string{
		"api_key":  getApiKey(),
		"language": getLanguage(),
	}

	api := host + fmt.Sprintf(apiTvAggregateCredits, tvId) + "?" + utils.StringMapToQuery(m)
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

	credits := &TvAggregateCredits{}
	err = json.Unmarshal(body, credits)
	if err != nil {
		utils.Logger.ErrorF("parse tmdb response: %s err: %v", api, err)
		return nil, err
	}

	return credits, err
}
