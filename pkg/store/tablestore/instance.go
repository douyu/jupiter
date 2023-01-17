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

package tstore

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/aliyun/aliyun-tablestore-go-sdk/tablestore"
	prome "github.com/douyu/jupiter/pkg/core/metric"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/util/xstring"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/fatih/color"
)

func newTs(config *Config) *tablestore.TableStoreClient {
	tsConfig := &tablestore.TableStoreConfig{
		RetryTimes:   config.RetryTimes,
		MaxRetryTime: config.MaxRetryTime,
		HTTPTimeout: tablestore.HTTPTimeout{
			ConnectionTimeout: config.ConnectionTimeout,
			RequestTimeout:    config.RequestTimeout,
		},
		MaxIdleConnections: config.MaxIdleConnections,
		Transport: &TsRoundTripper{
			name:   config.Name,
			config: *config,
			Transport: http.Transport{
				MaxIdleConnsPerHost: config.MaxIdleConnections,
				DialContext: (&net.Dialer{
					Timeout: config.ConnectionTimeout,
				}).DialContext,
			},
		},
	}
	return tablestore.NewClientWithConfig(config.EndPoint, config.Instance, config.AccessKeyId, config.AccessKeySecret, config.SecurityToken, tsConfig)
}

type TsRoundTripper struct {
	http.Transport
	name   string
	config Config
}

func (h *TsRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	var (
		addr    = req.URL.String()
		hostURL = strings.TrimLeft(req.URL.Path, "/")
		method  = req.Method
		beg     = time.Now()
	)

	resp, err := h.Transport.RoundTrip(req)
	var cost = time.Since(beg)

	// 指标采集
	if err != nil {
		prome.ClientHandleCounter.WithLabelValues(prome.TypeTableStore, h.name, method, hostURL, "error").Inc()
	} else {
		prome.ClientHandleCounter.WithLabelValues(prome.TypeTableStore, h.name, method, hostURL, resp.Status).Inc()
	}
	prome.ClientHandleHistogram.WithLabelValues(prome.TypeTableStore, h.name, method, hostURL).Observe(cost.Seconds())

	statusCode := -1
	if resp != nil {
		statusCode = resp.StatusCode
	}

	if xdebug.IsDevelopmentMode() {
		prefix := fmt.Sprintf("[%s:%s]", h.name, addr)
		fmt.Printf("%-50s => %s\n", color.GreenString(prefix), color.GreenString("Send: "+method+" | "+xstring.PrettyJson(req.Body)))
	}

	// 访问日志
	if err != nil {
		xlog.Jupiter().Error("access",
			xlog.FieldErr(err),
			xlog.FieldMethod(method),
			xlog.FieldAddr(addr),
			xlog.FieldCode(int32(statusCode)),
			xlog.FieldCost(cost),
		)
	} else {
		if h.config.EnableAccessLog {
			xlog.Jupiter().Info("access",
				xlog.FieldMethod(method),
				xlog.FieldAddr(addr),
				xlog.FieldCost(cost),
				xlog.FieldCode(int32(statusCode)),
			)
		}
	}

	if h.config.SlowThreshold > time.Duration(0) {
		// 慢日志
		if cost > h.config.SlowThreshold {
			xlog.Jupiter().Error("slow",
				xlog.FieldErr(errSlowCommand),
				xlog.FieldMethod(method),
				xlog.FieldCost(cost),
				xlog.FieldAddr(addr),
				xlog.FieldCode(int32(statusCode)),
			)
		}
	}
	return resp, err
}
