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

var (
	jsonRpc  = ""
	timeout  = 1
	username = "kodi"
	password = ""
	RpcQueue *RequestQueue
)

func InitKodi(c config.KodiConfig) {
	jsonRpc = c.JsonRpc
	timeout = c.Timeout
	username = c.Username
	password = c.Password
	RpcQueue = &RequestQueue{
		queue: make(map[string]*JsonRpcRequest, 0),
		lock:  &sync.RWMutex{},
	}
	go RpcQueue.notify()
}

func (q *RequestQueue) addTask(name string, req *JsonRpcRequest) bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	if _, ok := q.queue[name]; !ok {
		q.queue[name] = req
	}

	return true
}

func (q *RequestQueue) notify() {
	task := func() {
		if len(q.queue) == 0 {
			return
		}

		if !Ping() {
			return
		}

		q.lock.RLock()
		defer q.lock.RUnlock()

		utils.Logger.DebugF("kodi request queue size: %d", len(q.queue))
		for k, req := range q.queue {
			resp, err := request(req)
			if err != nil {
				panic(err)
			}

			delete(q.queue, k)
			utils.Logger.DebugF("req kodi: %s", resp)
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

func Ping() bool {
	_, err := request(&JsonRpcRequest{Method: "JSONRPC.Ping"})
	if err != nil {
		utils.Logger.WarningF("ping kodi err: %v", err)
	}
	return err == nil
}

// 发送json rpc请求
func request(rpcReq *JsonRpcRequest) ([]byte, error) {
	utils.Logger.DebugF("request kodi: %s", rpcReq.Method)

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
