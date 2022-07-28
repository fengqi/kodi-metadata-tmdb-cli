package tmdb

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"os"
)

// TvDetailsRequest
// Append To Response: https://developers.themoviedb.org/3/getting-started/append-to-response
type TvDetailsRequest struct {
	ApiKey           string `json:"api_key"`            // api_key, required
	Language         string `json:"language"`           // ISO 639-1, optional, default en-US
	TvId             int    `json:"tv_id"`              // tv id, required
	AppendToResponse string `json:"append_to_response"` // optional
}

type TvDetail struct {
	Id                   int                   `json:"id"`
	Name                 string                `json:"name"`
	BackdropPath         string                `json:"backdrop_path"`
	CreatedBy            []CreatedBy           `json:"created_by"`
	EpisodeRunTime       []int                 `json:"episode_run_time"`
	FirstAirDate         string                `json:"first_air_date"`
	LastAirDate          string                `json:"last_air_date"`
	Genres               []Genre               `json:"genres"`
	Homepage             string                `json:"homepage"`
	InProduction         bool                  `json:"in_production"`
	Languages            []string              `json:"languages"`
	LastEpisodeToAir     LastEpisodeToAir      `json:"last_episode_to_air"`
	NextEpisodeToAir     NextEpisodeToAir      `json:"next_episode_to_air"`
	Networks             []Network             `json:"networks"`
	NumberOfEpisodes     int                   `json:"number_of_episodes"`
	NumberOfSeasons      int                   `json:"number_of_seasons"`
	OriginCountry        []string              `json:"origin_country"`
	OriginalLanguage     string                `json:"original_language"`
	OriginalName         string                `json:"original_name"`
	Overview             string                `json:"overview"`
	Popularity           float32               `json:"popularity"`
	PosterPath           string                `json:"poster_path"`
	ProductionCompanies  []ProductionCompany   `json:"production_companies"`
	ProductionCountries  []ProductionCountry   `json:"production_countries"`
	Seasons              []Season              `json:"seasons"`
	SpokenLanguages      []SpokenLanguage      `json:"spoken_languages"`
	Status               string                `json:"status"`
	Tagline              string                `json:"tagline"`
	Type                 string                `json:"type"`
	VoteAverage          float32               `json:"vote_average"`
	VoteCount            int                   `json:"vote_count"`
	AggregateCredits     *TvAggregateCredits   `json:"aggregate_credits"`
	ContentRatings       *TvContentRatings     `json:"content_ratings"`
	TvEpisodeGroupDetail *TvEpisodeGroupDetail `json:"tv_episode_group_detail"`
	FromCache            bool                  `json:"from_cache"`
}

type Genre struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type Network struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	LogoPath      string `json:"logo_path"`
	OriginCountry string `json:"origin_country"`
}

type CreatedBy struct {
	Id          int    `json:"id"`
	CreditId    string `json:"credit_id"`
	Name        string `json:"name"`
	Gender      int    `json:"gender"`
	ProfilePath string `json:"profile_path"`
}

type LastEpisodeToAir struct {
	Id             int     `json:"id"`
	AirDate        string  `json:"air_date"`
	EpisodeNumber  int     `json:"episode_number"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	ProductionCode string  `json:"production_code"`
	SeasonNumber   int     `json:"season_number"`
	StillPath      string  `json:"still_path"`
	VoteAverage    float32 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
}

type NextEpisodeToAir struct {
	Id             int     `json:"id"`
	AirDate        string  `json:"air_date"`
	EpisodeNumber  int     `json:"episode_number"`
	Name           string  `json:"name"`
	Overview       string  `json:"overview"`
	ProductionCode string  `json:"production_code"`
	SeasonNumber   int     `json:"season_number"`
	StillPath      string  `json:"still_path"`
	VoteAverage    float32 `json:"vote_average"`
	VoteCount      int     `json:"vote_count"`
}

type ProductionCompany struct {
	Id            int    `json:"id"`
	LogoPath      string `json:"logo_path"`
	Name          string `json:"name"`
	OriginCountry string `json:"origin_country"`
}

type ProductionCountry struct {
	Iso31661 string `json:"iso_3166_1"`
	Name     string `json:"name"`
}

type Season struct {
	Id           int    `json:"id"`
	AirDate      string `json:"air_date"`
	EpisodeCount int    `json:"episode_count"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	PosterPath   string `json:"poster_path"`
	SeasonNumber int    `json:"season_number"`
}

type SpokenLanguage struct {
	EnglishName string `json:"english_name"`
	Iso6391     string `json:"iso_639_1"`
	Name        string `json:"name"`
}

func (t *tmdb) GetTvDetail(id int) (*TvDetail, error) {
	utils.Logger.DebugF("get tv detail from tmdb: %d", id)

	api := fmt.Sprintf(ApiTvDetail, id)
	req := map[string]string{
		"append_to_response": "aggregate_credits,content_ratings",
	}

	body, err := t.request(api, req)
	if err != nil {
		utils.Logger.ErrorF("read tmdb response err: %v", err)
		return nil, err
	}

	tvResp := &TvDetail{}
	err = json.Unmarshal(body, tvResp)
	if err != nil {
		utils.Logger.ErrorF("parse tmdb response err: %v", err)
		return nil, err
	}

	return tvResp, err
}

func (r *TvDetailsRequest) ToQuery() string {
	return fmt.Sprintf(
		"api_key=%s&language=%s&append_to_response=%s",
		r.ApiKey,
		r.Language,
		r.AppendToResponse,
	)
}

// SaveToCache 保存剧集详情到文件
func (d *TvDetail) SaveToCache(file string) {
	if d.Id == 0 || d.Name == "" {
		return
	}

	utils.Logger.InfoF("save tv detail to: %s", file)

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		utils.Logger.ErrorF("save tv to cache, open_file err: %v", err)
		return
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			utils.Logger.WarningF("save tv to cache, close file err: %v", err)
		}
	}(f)

	bytes, err := json.MarshalIndent(d, "", "    ")
	if err != nil {
		utils.Logger.ErrorF("save tv to cache, marshal struct errr: %v", err)
		return
	}

	_, err = f.Write(bytes)
}
