package tmdb

import (
	"encoding/json"
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"strconv"
)

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
func (t *tmdb) SearchShows(chsTitle, engTitle string, year int) (*SearchResults, error) {
	utils.Logger.InfoF("search: %s or %s %d from tmdb", chsTitle, engTitle, year)

	strYear := strconv.Itoa(year)
	searchComb := make([]map[string]string, 0)

	if chsTitle != "" {
		if year > 0 {
			searchComb = append(searchComb, map[string]string{
				"query":               chsTitle,
				"page":                "1",
				"include_adult":       "true",
				"first_air_date_year": strYear,
			})
		}
		searchComb = append(searchComb, map[string]string{
			"query":         chsTitle,
			"page":          "1",
			"include_adult": "true",
		})
	}

	if engTitle != "" {
		if year > 0 {
			searchComb = append(searchComb, map[string]string{
				"query":               engTitle,
				"page":                "1",
				"include_adult":       "true",
				"first_air_date_year": strYear,
			})
		}
		searchComb = append(searchComb, map[string]string{
			"query":         engTitle,
			"page":          "1",
			"include_adult": "true",
		})
	}

	if len(searchComb) == 0 {
		return nil, errors.New("title empty")
	}

	tvResp := &SearchTvResponse{}
	for _, req := range searchComb {
		body, err := t.request(ApiSearchTv, req)
		if err != nil {
			utils.Logger.ErrorF("read tmdb response err: %v", err)
			continue
		}

		err = json.Unmarshal(body, tvResp)
		if err != nil {
			utils.Logger.ErrorF("parse tmdb response err: %v", err)
			continue
		}

		if len(tvResp.Results) > 0 {
			utils.Logger.InfoF("search tv: %s %d result count: %d, use: %v", chsTitle, year, len(tvResp.Results), tvResp.Results[0])
			return tvResp.Results[0], nil
		}
	}

	return nil, errors.New("search tv not found")
}
