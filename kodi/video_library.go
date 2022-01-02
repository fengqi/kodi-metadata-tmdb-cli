package kodi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type JsonRpcRequest struct {
	Id      string      `json:"id"`
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type VideoLibraryScan struct {
	Description string      `json:"description"`
	Params      []Parameter `json:"params"`
	Permission  string      `json:"permission"`
	Returns     string      `json:"returns"`
	Type        string      `json:"type"`
}

type Parameter struct {
	Type    string `json:"type"`
	Name    string `json:"name"`
	Default string `json:"default"`
}

func SendRequest(v interface{}) {
	req := &JsonRpcRequest{
		JsonRpc: "2.0",
		Id:      "1",
		Method:  "",
	}

	jsonBytes, err := json.Marshal(req)
	if err != nil {
		//
	}

	resp, err := http.Post("http://192.168.50.142/jsonrpc", "", bytes.NewReader(jsonBytes))
	if err != nil {
		//
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			//
		}
	}(resp.Body)

	//
}
