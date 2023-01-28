package kodi

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
)

type Files struct{}

// GetSources 获取kodi添加的媒体源
func (f *Files) GetSources(media string) *GetSourcesResponse {
	body, _ := Rpc.request(&JsonRpcRequest{
		Method: "Files.GetSources",
		Params: GetSourcesRequest{
			Method: media,
		},
	})

	resp := &JsonRpcResponse{}
	err := json.Unmarshal(body, resp)
	if err != nil {
		utils.Logger.WarningF("parse Files.GetSources response err: %v", err)
		return nil
	}

	if resp != nil && resp.Result != nil {
		jsonBytes, _ := json.Marshal(resp.Result)
		sourcesResp := &GetSourcesResponse{}
		_ = json.Unmarshal(jsonBytes, sourcesResp)

		return sourcesResp
	}

	return nil
}
