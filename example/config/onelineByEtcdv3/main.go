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
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/conf"
	_ "github.com/douyu/jupiter/pkg/datasource/etcdv3"
	"github.com/douyu/jupiter/pkg/xlog"
)

//  go run main.go --config=etcdv3://10.0.101.68:2379?key=test

var configText = `
[people]
    name = "jupiter"
[jupiter.logger.default]
    debug = true
    enableConsole = true
[jupiter.server.governor]
    enable = false
    host = "0.0.0.0"
    port = 9246
`

func initTestData() {
	etcdCfg := clientv3.Config{
		Endpoints: []string{"127.0.0.1:2379"},
	}
	cli, _ := clientv3.New(etcdCfg)
	cli.Put(context.Background(), "test", configText)
}

func main() {
	initTestData()
	app := NewEngine()
	if err := app.Run(); err != nil {
		panic(err)
	}
}

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		eng.printConfig,
	); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
	return eng
}

func (s *Engine) printConfig() error {
	xlog.DefaultLogger = xlog.StdConfig("default").Build()
	peopleName := conf.GetString("people.name")
	xlog.Info("people info", xlog.String("name", peopleName), xlog.String("type", "onelineByEtcdv3"))
	return nil
}
