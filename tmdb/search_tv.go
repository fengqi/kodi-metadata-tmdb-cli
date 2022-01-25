package tmdb

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

type SearchTvRequest struct {
	ApiKey           string `json:"api_key"`             // api_key, required
	Language         string `json:"language"`            // ISO 639-1, optional, default en-US
	Page             int    `json:"page"`                // page, 1-1000, optional, default 1
	Query            string `json:"query"`               // query text, required, URI encoded
	IncludeAdult     bool   `json:"include_adult"`       // include adult, optional
	FirstAirDateYear int    `json:"first_air_date_year"` // year, optional
}

type SearchTvResponse struct {
	Page         int              `json:"page"`
	TotalResults int              `json:"total_results"`
	TotalPages   int              `json:"total_pages"`
	Results      []*SearchResults `json:"results"`
}

type SearchResults struct {
	Id               int      `json:"id"`
	PosterPath       string   `json:"poster_path"`
	Popularity       float32  `json:"popularity"`
	BackdropPath     string   `json:"backdrop_path"`
	VoteAverage      float32  `json:"vote_average"`
	Overview         string   `json:"overview"`
	FirstAirDate     string   `json:"first_air_date"`
	OriginCountry    []string `json:"origin_country"`
	GenreIds         []int    `json:"genre_ids"`
	OriginalLanguage string   `json:"original_language"`
	VoteCount        int      `json:"vote_count"`
	Name             string   `json:"name"`
	OriginalName     string   `json:"original_name"`
}

type Response struct {
	Success       bool   `json:"success"`
	StatusCode    int    `json:"status_code"`
	StatusMessage string `json:"status_message"`
}

// SearchShows 搜索tmdb
func SearchShows(title string, year int) (*SearchResults, error) {
	utils.Logger.InfoF("search: %s %d from tmdb", title, year)

	req := &SearchTvRequest{
		ApiKey:           getApiKey(),
		Query:            title,
		FirstAirDateYear: year,
		Page:             1,
		IncludeAdult:     true,
		Language:         getLanguage(),
	}

	api := host + apiSearchTv + "?" + req.ToQuery()
	utils.Logger.DebugF("request tmdb: %s", api)

	resp, err := http.Get(api)
	if err != nil {
		utils.Logger.WarningF("search shows err: %v", err)
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
		utils.Logger.ErrorF("read tmdb response err: %v", err)
		return nil, err
	}

	tvResp := &SearchTvResponse{}
	err = json.Unmarshal(body, tvResp)
	if err != nil {
		utils.Logger.ErrorF("parse tmdb response err: %v", err)
		return nil, err
	}

	if len(tvResp.Results) == 0 {
		return nil, nil
	}

	if len(tvResp.Results) > 0 {
		utils.Logger.InfoF("search tv: %s %d result count: %d, use: %v", title, year, len(tvResp.Results), tvResp.Results[0])
	}

	return tvResp.Results[0], nil
}

func (r *SearchTvRequest) ToQuery() string {
	return fmt.Sprintf(
		"api_key=%s&language=%s&page=%d&include_adult=%t&query=%s&first_air_date_year=%d",
		r.ApiKey,
		r.Language,
		r.Page,
		r.IncludeAdult,
		url.QueryEscape(r.Query),
		r.FirstAirDateYear,
	)
}
