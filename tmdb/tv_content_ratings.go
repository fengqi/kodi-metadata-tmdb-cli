package tmdb

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// TV 内容分级

type TvContentRatings struct {
	Id      int                      `json:"id"`
	Results []TvContentRatingsResult `json:"results"`
}

type TvContentRatingsResult struct {
	Iso31661 string `json:"iso_3166_1"`
	Rating   string `json:"rating"`
}

func GetTvContentRatings(tvId int) (*TvContentRatings, error) {
	utils.Logger.DebugF("get tv content ratings from tmdb: %d", tvId)

	m := map[string]string{
		"api_key":  getApiKey(),
		"language": getLanguage(),
	}

	api := host + fmt.Sprintf(apiTvContentRatings, tvId) + "?" + utils.StringMapToQuery(m)
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

	ratings := &TvContentRatings{}
	err = json.Unmarshal(body, ratings)
	if err != nil {
		utils.Logger.ErrorF("parse tmdb response: %s err: %v", api, err)
		return nil, err
	}

	return ratings, err
}
