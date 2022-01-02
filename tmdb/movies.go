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

type MovieDetail struct {
	Adult               bool                `json:"adult"`
	BackdropPath        string              `json:"backdrop_path"`
	BelongsToCollection BelongsToCollection `json:"belongs_to_collection"`
	Budget              int                 `json:"budget"`
	Genres              []Genre             `json:"genres"`
	Homepage            string              `json:"homepage"`
	Id                  int                 `json:"id"`
	ImdbId              string              `json:"imdb_id"`
	OriginalLanguage    string              `json:"original_language"`
	OriginalTitle       string              `json:"original_title"`
	Overview            string              `json:"overview"`
	Popularity          float32             `json:"popularity"`
	PosterPath          string              `json:"poster_path"`
	ProductionCompanies []ProductionCompany `json:"production_companies"`
	ProductionCountries []ProductionCountry `json:"production_countries"`
	ReleaseDate         string              `json:"release_date"`
	Revenue             int                 `json:"revenue"`
	Runtime             int                 `json:"runtime"`
	SpokenLanguages     []SpokenLanguage    `json:"spoken_languages"`
	Status              string              `json:"status"`
	Tagline             string              `json:"tagline"`
	Title               string              `json:"title"`
	Video               bool                `json:"video"`
	VoteAverage         float32             `json:"vote_average"`
	VoteCount           int                 `json:"vote_count"`
	Credits             *Credit             `json:"credits"`
}

type BelongsToCollection struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	PosterPath   string `json:"poster_path"`
	BackdropPath string `json:"backdrop_path"`
}

// GetMovieDetail 获取电影详情
func GetMovieDetail(id int) (*MovieDetail, error) {
	utils.Logger.DebugF("get movie detail from tmdb: %d", id)

	req := map[string]string{
		"api_key":            getApiKey(),
		"language":           getLanguage(),
		"append_to_response": "credits",
	}

	api := host + fmt.Sprintf(apiMovieDetail, id) + "?" + utils.StringMapToQuery(req)
	utils.Logger.DebugF("request tmdb: %s", api)

	resp, err := http.Get(api)
	if err != nil {
		utils.Logger.ErrorF("request tmdb: %d err: %v", api, err)
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

	detail := &MovieDetail{}
	err = json.Unmarshal(body, detail)
	if err != nil {
		utils.Logger.ErrorF("parse tmdb response: %s err: %v", api, err)
		return nil, err
	}

	return detail, err
}

// SaveToCache 保存剧集详情到文件
func (d *MovieDetail) SaveToCache(file string) {
	if d.Id == 0 || d.Title == "" {
		return
	}

	utils.Logger.InfoF("save movie detail to: %s", file)

	f, err := os.OpenFile(file, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		utils.Logger.ErrorF("save movie to cache, open_file err: %v", err)
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
		utils.Logger.ErrorF("save movie to cache, marshal struct err: %v", err)
		return
	}

	_, err = f.Write(bytes)
}
