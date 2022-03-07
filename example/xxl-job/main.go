// Copyright 2022 Douyu
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
	"fmt"
	"log"

	"github.com/douyu/jupiter"
	"github.com/douyu/jupiter/pkg/executor"
	"github.com/douyu/jupiter/pkg/executor/xxl"
	"github.com/douyu/jupiter/pkg/executor/xxl/logger"
	"github.com/douyu/jupiter/pkg/hooks"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/xlog"
)

func main() {
	eng := NewEngine()
	eng.RegisterHooks(hooks.Stage_AfterStop, func() {
		fmt.Println("exit jupiter app ...")
	})
	if err := eng.Run(); err != nil {
		log.Fatal(err)
	}
}

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		xgo.ParallelWithError(
			eng.startXxlJob,
		),
	); err != nil {
		xlog.Panic("startup engine", xlog.Any("err", err))
	}
	return eng
}

func (eng *Engine) startXxlJob() error {
	executor := xxl.StdNewExecutor(xxl.ExecutorHost("127.0.0.1"))
	executor.RegXJob(
		NewTest(),
		NewTest2(),
	)
	eng.Executor(executor)
	return nil
}

// =======以下为示例任务test.go=========
type Test struct{}

func NewTest() *Test {
	return &Test{}
}

// 任务名称
func (t *Test) GetJobName() string {
	return "test"
}

// xxl-job 分布式调度任务执行函数
func (t *Test) Run(ctx context.Context, param *executor.RunReq) (msg string, err error) {
	//使用xxl-logger日志即可在xxl-job平台上看到日志
	logger.Info(param.LogID, "start run...")
	logger.Info(param.LogID, fmt.Sprintf("job param is: %s", param.ExecutorParams))
	fmt.Println("test job has been executed")
	return "success", nil
}

// =======以下为第二个示例任务test2.go=========
type Test2 struct{}

func NewTest2() *Test2 {
	return &Test2{}
}

// 任务名称
func (t *Test2) GetJobName() string {
	return "test2"
}

// xxl-job 分布式调度任务执行函数
func (t *Test2) Run(ctx context.Context, param *executor.RunReq) (msg string, err error) {
	return "success", nil
}
