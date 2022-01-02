package tmdb

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io"
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

func SearchMovie(title string, year int) (*SearchMoviesResults, error) {
	utils.Logger.InfoF("search: %s %d from tmdb", title, year)

	req := map[string]string{
		"api_key":       getApiKey(),
		"language":      getLanguage(),
		"query":         title,
		"page":          "1",
		"include_adult": "true",
		//"region": "US",
	}

	if year > 0 {
		req["year"] = strconv.Itoa(year)
		req["primary_release_year"] = strconv.Itoa(year)
	}

	api := host + apiSearchMovie + "?" + utils.StringMapToQuery(req)
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

	moviesResp := &SearchMoviesResponse{}
	err = json.Unmarshal(body, moviesResp)
	if err != nil {
		utils.Logger.ErrorF("parse tmdb response err: %v", err)
		return nil, err
	}

	if len(moviesResp.Results) == 0 {
		return nil, nil
	}

	if len(moviesResp.Results) > 0 {
		utils.Logger.InfoF("search movies: %s %d result count: %d, use: %v", title, year, len(moviesResp.Results), moviesResp.Results[0])
	}

	return moviesResp.Results[0], nil
}
