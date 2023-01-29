package kodi

import (
	"encoding/json"
)

type XBMC struct{}

// GetInfoBooleans Retrieve info booleans about Kodi and the system
// https://kodi.wiki/view/List_of_boolean_conditions
func (x *XBMC) GetInfoBooleans(booleans []string) map[string]bool {
	bytes, err := Rpc.request(&JsonRpcRequest{
		Method: "XBMC.GetInfoBooleans",
		Params: &GetInfoBooleansRequest{
			Booleans: booleans,
		},
	})

	if err != nil {
		return nil
	}

	resp := &JsonRpcResponse{}
	err = json.Unmarshal(bytes, resp)
	if err != nil {
		return nil
	}

	if resp != nil && resp.Result != nil {
		jsonBytes, _ := json.Marshal(resp.Result)
		info := make(map[string]bool, 0)
		_ = json.Unmarshal(jsonBytes, &info)
		return info
	}
	return nil
}

// GetInfoLabels 获取Kodi和系统相关信息
// https://kodi.wiki/view/JSON-RPC_API/v13#XBMC.GetInfoLabels
// https://kodi.wiki/view/InfoLabels
func (x *XBMC) GetInfoLabels(labels []string) map[string]interface{} {
	// TODO
	return nil
}
