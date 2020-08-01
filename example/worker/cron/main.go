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
	"fmt"
	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/worker/xcron"
	"github.com/douyu/jupiter/pkg/xlog"
	"time"
)

func main() {
	eng := NewEngine()
	if err := eng.Run(); err != nil {
		xlog.Error(err.Error())
	}
}

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}

	if err := eng.Startup(
		eng.startJobs,
	); err != nil {
		xlog.Panic("startup engine", xlog.Any("err", err))
	}

	return eng
}

func (eng *Engine) startJobs() error {
	cron := xcron.StdConfig("test").Build()
	cron.Schedule(xcron.Every(time.Second*10), xcron.FuncJob(eng.execJob))
	return eng.Schedule(cron)
}

func (eng *Engine) execJob() error {
	xlog.Info("info job")
	xlog.Warn("warn job")
	fmt.Println("run job")
	return nil
}
