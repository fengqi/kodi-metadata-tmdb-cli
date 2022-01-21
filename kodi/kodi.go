package kodi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"io"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	jsonRpc  = ""
	timeout  = 1
	username = "kodi"
	password = ""
)

type JsonRpcRequest struct {
	Id      string      `json:"id"`
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
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

func InitKodi(c config.KodiConfig) {
	jsonRpc = c.JsonRpc
	timeout = c.Timeout
	username = c.Username
	password = c.Password
}

func Ping() bool {
	_, err := request(&JsonRpcRequest{})
	return err == nil
}

func request(rpcReq *JsonRpcRequest) ([]byte, error) {
	if rpcReq.JsonRpc == "" {
		rpcReq.JsonRpc = "2.0"
	}

	if rpcReq.Id == "" {
		rpcReq.Id = time.Now().String()
	}

	jsonBytes, err := json.Marshal(rpcReq)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, jsonRpc, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(username, password)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout:   time.Duration(timeout) * time.Second,
		Transport: http.DefaultTransport,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	return ioutil.ReadAll(resp.Body)
}
