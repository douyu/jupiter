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
	"time"

	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/xlog"
)

//  go run main.go --config=config.toml
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
		eng.printLogger,
	); err != nil {
		xlog.Panic("startup", xlog.Any("err", err))
	}
	return eng
}

func (s *Engine) printLogger() error {
	xlog.DefaultLogger = xlog.StdConfig("default").Build()
	go func() {
		for {
			xlog.Info("logger info", xlog.String("gopher", "jupiter1"), xlog.String("type", "file"))
			xlog.Debug("logger debug", xlog.String("gopher", "jupiter2"), xlog.String("type", "file"))
			time.Sleep(1 * time.Second)
		}
	}()
	return nil
}
