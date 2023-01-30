package kodi

import (
	"encoding/json"
	"fengqi/kodi-metadata-tmdb-cli/utils"
)

type Files struct{}

// GetSources 获取kodi添加的媒体源
func (f *Files) GetSources(media string) []*FileSource {
	body, _ := Rpc.request(&JsonRpcRequest{
		Method: "Files.GetSources",
		Params: GetSourcesRequest{
			Media: media,
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
		err = json.Unmarshal(jsonBytes, sourcesResp)
		if err == nil {
			return sourcesResp.Sources
		}
	}

	return nil
}
