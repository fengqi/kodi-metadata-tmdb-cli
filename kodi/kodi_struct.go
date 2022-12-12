package kodi

import (
	"fengqi/kodi-metadata-tmdb-cli/config"
	"sync"
)

type JsonRpc struct {
	config       *config.KodiConfig
	queue        map[string]*JsonRpcRequest
	lock         *sync.RWMutex
	VideoLibrary *VideoLibrary
}

// JsonRpcRequest JsonRpc 请求参数
type JsonRpcRequest struct {
	Id      string      `json:"id"`
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

// JsonRpcResponse JsonRpc 返回参数
type JsonRpcResponse struct {
	Id      string                 `json:"id"`
	JsonRpc string                 `json:"jsonrpc"`
	Result  map[string]interface{} `json:"result"`
}

type Limits struct {
	Start int `json:"start"`
	End   int `json:"end"`
}

type LimitsResult struct {
	Start int `json:"start"`
	End   int `json:"end"`
	Total int `json:"total"`
}

type Sort struct {
	Order         string `json:"order"`
	Method        string `json:"method"`
	IgnoreArticle bool   `json:"ignorearticle"`
}

type Filter struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}
