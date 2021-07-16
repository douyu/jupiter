// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/go-resty/resty/v2"
)

/*
基于http的配置轮询的配置获取
*/
type yaseeDataSource struct {
	lastRevision int64
	enableWatch  bool
	client       *resty.Client
	addr         string
	changed      chan struct{}
	data         string
}

// default client resp struct
type yaseeRes struct {
	Code int        `json:"code"`
	Msg  string     `json:"msg"`
	Data ConfigData `json:"data"`
}

// ConfigData ...
type ConfigData struct {
	Content      string `json:"content"`
	LastRevision int64  `json:"last_revision"`
}

// NewDataSource ...
func NewDataSource(addr string, enableWatch bool) *yaseeDataSource {
	yasee := &yaseeDataSource{
		client:      resty.New(),
		addr:        addr,
		changed:     make(chan struct{}),
		enableWatch: enableWatch,
	}
	if enableWatch {
		xgo.Go(yasee.watch)
	}
	return yasee
}

// ReadConfig ...
func (y *yaseeDataSource) ReadConfig() ([]byte, error) {
	// 检查watch 如果watch为真，走长轮询逻辑
	switch y.enableWatch {
	case true:
		if y.data == "" {
			content, err := y.getConfigInner(y.addr, y.enableWatch)
			return []byte(content), err
		}
		return []byte(y.data), nil
	default:
		content, err := y.getConfigInner(y.addr, y.enableWatch)
		return []byte(content), err
	}
}

// IsConfigChanged ...
func (y *yaseeDataSource) IsConfigChanged() <-chan struct{} {
	return y.changed
}

// Close ...
func (y *yaseeDataSource) Close() error {
	close(y.changed)
	return nil
}

func (y *yaseeDataSource) watch() {
	for {
		resp, err := y.client.R().SetQueryParam("watch", strconv.FormatBool(y.enableWatch)).Get(y.addr)
		// client get err
		if err != nil {
			xlog.Error("yaseeDataSource", xlog.String("listenConfig curl err", err.Error()))
			time.Sleep(time.Second * 1)
			continue
		}
		if resp.StatusCode() != 200 {
			xlog.Error("yaseeDataSource", xlog.String("listenConfig status err", resp.Status()))
			time.Sleep(time.Second * 1)
			continue
		}
		var yaseeRes yaseeRes
		if err := json.Unmarshal(resp.Body(), &yaseeRes); err != nil {
			y.updateConfig(string(resp.Body()), 0)
			time.Sleep(time.Second * 1)
			continue
		}
		// default code != 200 means not change
		if yaseeRes.Code != 200 {
			xlog.Info("yaseeDataSource", xlog.Int64("code", int64(yaseeRes.Code)))
			time.Sleep(time.Second * 1)
			continue
		}
		y.updateConfig(yaseeRes.Data.Content, yaseeRes.Data.LastRevision)
	}
}

func (y *yaseeDataSource) updateConfig(data string, version int64) {
	select {
	case y.changed <- struct{}{}:
		// record the config change data
		y.data = data
		y.lastRevision = version
		xlog.Info("yaseeDataSource", xlog.String("change", data))
	default:
	}
}

func (y *yaseeDataSource) getConfigInner(addr string, enableWatch bool) (string, error) {
	urlParse, err := url.Parse(addr)
	if err != nil {
		return "", fmt.Errorf("config addr is wrong, err:%v", err.Error())
	}
	appName := urlParse.Query().Get("name")
	appEnv := urlParse.Query().Get("env")
	target := urlParse.Query().Get("target")
	port := urlParse.Query().Get("port")
	commonKey := fmt.Sprintf("%s-%s-%s-%s", appName, appEnv, target, port)
	if commonKey == "" {
		return "", fmt.Errorf("config check key is null")
	}

	content, err := y.getConfig(addr, enableWatch)
	if err != nil {
		content, err = ReadConfigFromFile(commonKey, "configCacheDir")
		if err != nil {
			return "", errors.New("read config from both server and cache fail")
		}
		return "", err
	}
	WriteConfigToFile(commonKey, "configCacheDir", content)
	return content, nil
}

func (y *yaseeDataSource) getConfig(addr string, enableWatch bool) (string, error) {
	resp, err := y.client.SetDebug(false).R().SetQueryParam("watch", strconv.FormatBool(enableWatch)).Get(addr)
	if err != nil {
		return "", errors.New("get config err")
	}
	if resp.StatusCode() != 200 {
		return "", fmt.Errorf("get config reply err code:%v", resp.Status())
	}
	configRes := yaseeRes{}
	if err := json.Unmarshal(resp.Body(), &configRes); err != nil {
		return string(resp.Body()), nil
	}
	if configRes.Code != 200 {
		return "", fmt.Errorf("get config reply err code:%v", resp.Status())
	}
	return configRes.Data.Content, nil
}
