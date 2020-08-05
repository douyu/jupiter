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
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/server/xecho"
	"github.com/douyu/jupiter/pkg/xlog"
	"time"
)

//  go run main.go --config=config.toml --watch=true
func main() {
	app := NewEngine()
	if err := app.Run(); err != nil {
		panic(err)
	}
}

type Engine struct {
	jupiter.Application
}

type People struct {
	Name string
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.fileWatch,
		eng.serveHTTP,
	); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}

	return eng
}

func (eng *Engine) serveHTTP() error {
	server := xecho.StdConfig("http").Build()
	return eng.Serve(server)
}

func (s *Engine) fileWatch() error {
	xlog.DefaultLogger = xlog.StdConfig("default").Build()
	p := People{}
	conf.OnChange(func(config *conf.Configuration) {
		err := config.UnmarshalKey("people", &p)
		if err != nil {
			panic(err.Error())
		}
	})

	go func() {
		// 循环打印配置
		for {
			time.Sleep(1 * time.Second)
			xlog.Info("people info", xlog.String("name", p.Name), xlog.String("type", "structByFileWatch"))
		}
	}()
	return nil
}
