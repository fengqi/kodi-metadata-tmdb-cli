package kodi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var Rpc *JsonRpc
var httpClient *http.Client

func InitKodi(config *config.KodiConfig) {
	Rpc = &JsonRpc{
		config:       config,
		refreshQueue: make(map[string]struct{}, 0),
		scanQueue:    make(map[string]struct{}, 0),
		refreshLock:  &sync.RWMutex{},
		scanLock:     &sync.RWMutex{},
		VideoLibrary: &VideoLibrary{
			scanLimiter:   NewLimiter(300),
			refreshMovie:  NewLimiter(60),
			refreshTVShow: NewLimiter(60),
		},
		Files: &Files{},
		XBMC:  &XBMC{},
	}

	go Rpc.ConsumerRefreshTask()
	go Rpc.ConsumerScanTask()

	httpClient = &http.Client{
		Timeout:   time.Duration(config.Timeout) * time.Second,
		Transport: &http.Transport{},
	}
}

func (r *JsonRpc) Ping() bool {
	_, err := r.request(&JsonRpcRequest{Method: "JSONRPC.Ping"})
	if err != nil {
		utils.Logger.WarningF("ping kodi err: %v", err)
	}
	return err == nil
}

// 发送json rpc请求
func (r *JsonRpc) request(rpcReq *JsonRpcRequest) ([]byte, error) {
	if rpcReq.JsonRpc == "" {
		rpcReq.JsonRpc = "2.0"
	}

	if rpcReq.Id == "" {
		rpcReq.Id = strconv.FormatInt(time.Now().UnixNano(), 10)
	}

	jsonBytes, err := json.Marshal(rpcReq)
	utils.Logger.DebugF("request kodi: %s", jsonBytes)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, r.config.JsonRpc, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(r.config.Username, r.config.Password)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			utils.Logger.WarningF("request kodi closeBody err: %v", err)
		}
	}(resp.Body)

	if resp.StatusCode != 200 {
		return nil, errors.New(resp.Status)
	}

	return io.ReadAll(resp.Body)
}
