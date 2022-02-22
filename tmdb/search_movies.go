package tmdb

import (
	"encoding/json"
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io/ioutil"
	"net/http"
	"strconv"
)

type SearchMoviesResponse struct {
	Page         int                    `json:"page"`
	TotalResults int                    `json:"total_results"`
	TotalPages   int                    `json:"total_pages"`
	Results      []*SearchMoviesResults `json:"results"`
}

type SearchMoviesResults struct {
	PosterPath       string  `json:"poster_path"`
	Adult            bool    `json:"adult"`
	Overview         string  `json:"overview"`
	ReleaseDate      string  `json:"release_date"`
	GenreIds         []int   `json:"genre_ids"`
	Id               int     `json:"id"`
	OriginalTitle    string  `json:"original_title"`
	OriginalLanguage string  `json:"original_language"`
	Title            string  `json:"title"`
	BackdropPath     string  `json:"backdrop_path"`
	Popularity       float32 `json:"popularity"`
	VoteCount        int     `json:"vote_count"`
	Video            bool    `json:"video"`
	VoteAverage      float32 `json:"vote_average"`
}

func SearchMovie(chsTitle, engTitle string, year int) (*SearchMoviesResults, error) {
	utils.Logger.InfoF("search: %s or %s %d from tmdb", chsTitle, engTitle, year)

	strYear := strconv.Itoa(year)
	searchComb := make([]map[string]string, 0)

	if chsTitle != "" {
		// chs + year
		if year > 0 {
			searchComb = append(searchComb, map[string]string{
				"api_key":       getApiKey(),
				"language":      getLanguage(),
				"query":         chsTitle,
				"page":          "1",
				"include_adult": "true",
				//"region": "US",
				"year":                 strYear,
				"primary_release_year": strYear,
			})
		}
		// chs
		searchComb = append(searchComb, map[string]string{
			"api_key":       getApiKey(),
			"language":      getLanguage(),
			"query":         chsTitle,
			"page":          "1",
			"include_adult": "true",
			//"region": "US",
		})
	}

	if engTitle != "" {
		// eng + year
		if year > 0 {
			searchComb = append(searchComb, map[string]string{
				"api_key":       getApiKey(),
				"language":      getLanguage(),
				"query":         engTitle,
				"page":          "1",
				"include_adult": "true",
				//"region": "US",
				"year":                 strYear,
				"primary_release_year": strYear,
			})
		}
		// eng
		searchComb = append(searchComb, map[string]string{
			"api_key":       getApiKey(),
			"language":      getLanguage(),
			"query":         engTitle,
			"page":          "1",
			"include_adult": "true",
			//"region": "US",
		})
	}

	if len(searchComb) == 0 {
		return nil, errors.New("title empty")
	}

	moviesResp := &SearchMoviesResponse{}
	for _, req := range searchComb {
		api := host + apiSearchMovie + "?" + utils.StringMapToQuery(req)
		utils.Logger.DebugF("request tmdb: %s", api)

		resp, err := http.Get(api)
		if err != nil {
			utils.Logger.WarningF("search shows err: %v", err)
			continue
		}

		body, err := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			utils.Logger.ErrorF("read tmdb response err: %v", err)
			continue
		}

		err = json.Unmarshal(body, moviesResp)
		if err != nil {
			utils.Logger.ErrorF("parse tmdb response err: %v", err)
			continue
		}

		if len(moviesResp.Results) > 0 {
			utils.Logger.InfoF("search movies: %s %d result count: %d, use: %v", chsTitle, year, len(moviesResp.Results), moviesResp.Results[0])
			return moviesResp.Results[0], nil
		}
	}

	return nil, errors.New("search not found")
}
