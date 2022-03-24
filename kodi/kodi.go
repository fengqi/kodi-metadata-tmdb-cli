package kodi

import (
	"bytes"
	"encoding/json"
	"errors"
	"fengqi/kodi-metadata-tmdb-cli/config"
	"fengqi/kodi-metadata-tmdb-cli/utils"
	"io"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

var Rpc *JsonRpc

func InitKodi(c config.KodiConfig) {
	Rpc = &JsonRpc{
		config: c,
		queue:  make(map[string]*JsonRpcRequest, 0),
		lock:   &sync.RWMutex{},
	}
}

func (r *JsonRpc) AddTask(name string, req *JsonRpcRequest) bool {
	if !r.config.Enable {
		return false
	}

	r.lock.Lock()
	defer r.lock.Unlock()

	if _, ok := r.queue[name]; !ok {
		r.queue[name] = req
	}

	return true
}

func (r *JsonRpc) RunNotify() {
	if !r.config.Enable {
		return
	}

	task := func() {
		if len(r.queue) == 0 {
			return
		}

		if !r.Ping() {
			return
		}

		r.lock.RLock()
		defer r.lock.RUnlock()

		utils.Logger.DebugF("kodi request queue size: %d", len(r.queue))
		for k, req := range r.queue {
			_, err := r.request(req)
			if err != nil {
				utils.Logger.ErrorF("request kodi %s err: %v", req.Method, err)
				continue
			}

			delete(r.queue, k)

			time.Sleep(time.Second * 30)
		}
	}

	ticker := time.NewTicker(time.Second * 60)
	for {
		select {
		case <-ticker.C:
			task()
		}
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
	utils.Logger.InfoF("request kodi: %s", rpcReq.Method)

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

	req, err := http.NewRequest(http.MethodPost, r.config.JsonRpc, bytes.NewReader(jsonBytes))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(r.config.Username, r.config.Password)
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout:   time.Duration(r.config.Timeout) * time.Second,
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
