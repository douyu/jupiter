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

package main

import (
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/server/xgoframe"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

func main() {
	eng := NewEngine()
	if err := eng.Run(); err != nil {
		xlog.Panic(err.Error())
	}
}

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.serveHTTP,
	); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
	return eng
}

// HTTP地址
func (eng *Engine) serveHTTP() error {
	server := xgoframe.StdConfig("http").Build()
	server.BindHandler("/hello", func(r *ghttp.Request) {
		_ = r.Response.WriteJson("Hello GoFrame")
	})
	server.BindHandler("/", func(r *ghttp.Request) {
		_ = r.Response.WriteJson(g.Map{
			"id":   1,
			"name": "hello,jupiter",
		})
	})
	server.BindHandler("/panic", func(r *ghttp.Request) {
		panic("it is a test for panic")
	})
	return eng.Serve(server)
}
