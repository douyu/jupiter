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
	_ "github.com/douyu/jupiter/pkg/datasource/apollo"
	"github.com/douyu/jupiter/pkg/xlog"
)

// apollo: http://106.54.227.205/config.html#/appid=jupiter&env=DEV&cluster=default
// account/password: apollo/admin

//  go run main.go --config="apollo://106.54.227.205:8080?appId=jupiter&cluster=default&namespaceName=application&key=jupiter-test"
func main() {
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
	xlog.Info("people info", xlog.String("name", peopleName), xlog.String("type", "onelineByApollo"))
	return nil
}
