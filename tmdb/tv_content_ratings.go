package tmdb

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"fmt"
)

// TV 内容分级

type TvContentRatings struct {
	Id      int                      `json:"id"`
	Results []TvContentRatingsResult `json:"results"`
}

type TvContentRatingsResult struct {
	ISO31661 string `json:"iso_3166_1"`
	Rating   string `json:"rating"`
}

func (t *tmdb) GetTvContentRatings(tvId int) (*TvContentRatings, error) {
	utils.Logger.DebugF("get tv content ratings from tmdb: %d", tvId)

	api := fmt.Sprintf(ApiTvContentRatings, tvId)
	req := map[string]string{}

	body, err := t.request(api, req)
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
