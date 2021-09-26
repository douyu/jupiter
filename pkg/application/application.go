// Copyright 2021 rex lv
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
// See the License for the specific language governing permissions and // limitations under the License.

package application

import (
	"context"
	"fmt"
	"sync"

	"go.uber.org/multierr"

	_ "github.com/douyu/jupiter/internal/autoproc"
	_ "github.com/douyu/jupiter/internal/banner"

	"github.com/douyu/jupiter/pkg/component"
	"github.com/douyu/jupiter/pkg/container"
	"github.com/douyu/jupiter/pkg/governor"
	job "github.com/douyu/jupiter/pkg/worker/xjob"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"

	//go-lint
	_ "github.com/douyu/jupiter/pkg/conf/datasource/file"
	_ "github.com/douyu/jupiter/pkg/conf/datasource/http"
	_ "github.com/douyu/jupiter/pkg/registry/etcdv3"

	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/signals"
	"github.com/douyu/jupiter/pkg/util/xcycle"
	"github.com/douyu/jupiter/pkg/util/xdefer"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/worker"
	"github.com/douyu/jupiter/pkg/xlog"
	"golang.org/x/sync/errgroup"
)

const (
	//StageAfterStop after app stop
	StageAfterStop uint32 = iota + 1
	//StageBeforeStop before app stop
	StageBeforeStop
)

// Application is the framework's instance, it contains the servers, workers, client and configuration settings.
// Create an instance of Application, by using &Application{}
type Application struct {
	cycle        *xcycle.Cycle
	smu          *sync.RWMutex
	initOnce     sync.Once
	startupOnce  sync.Once
	stopOnce     sync.Once
	workers      []worker.Worker
	jobs         map[string]job.Runner
	logger       *xlog.Logger
	hooks        map[uint32]*xdefer.DeferStack
	configParser conf.Unmarshaller
	disableMap   map[Disable]bool
	HideBanner   bool
	stopped      chan struct{}

	// generic components
	componentHeap *container.PriorityQueue
}

// New create a new Application instance
func New(fns ...func() error) (*Application, error) {
	app := &Application{}
	if err := app.Startup(fns...); err != nil {
		return nil, err
	}
	return app, nil
}

func DefaultApp() *Application {
	app := &Application{}
	app.initialize()
	return app
}

//init hooks
func (app *Application) initHooks(hookKeys ...uint32) {
	app.hooks = make(map[uint32]*xdefer.DeferStack, len(hookKeys))
	for _, k := range hookKeys {
		app.hooks[k] = xdefer.NewStack()
	}
}

//run hooks
func (app *Application) runHooks(k uint32) {
	hooks, ok := app.hooks[k]
	if ok {
		hooks.Clean()
	}
}

//RegisterHooks register a stage Hook
func (app *Application) RegisterHooks(k uint32, fns ...func() error) error {
	hooks, ok := app.hooks[k]
	if ok {
		hooks.Push(fns...)
		return nil
	}
	return fmt.Errorf("hook stage not found")
}

// initialize application
func (app *Application) initialize() {
	app.initOnce.Do(func() {
		//assign
		app.cycle = xcycle.NewCycle()
		app.smu = &sync.RWMutex{}
		app.workers = make([]worker.Worker, 0)
		app.jobs = make(map[string]job.Runner)
		app.logger = xlog.JupiterLogger
		app.configParser = toml.Unmarshal
		app.disableMap = make(map[Disable]bool)
		app.stopped = make(chan struct{})
		//app.components = make([]component.Component, 0)
		app.componentHeap = container.NewPriorityQueue()
		//private method
		app.initHooks(StageBeforeStop, StageAfterStop)

		app.componentHeap.Push(governor.StdConfig("governor").MustBuild(), 1)
		flag.Parse()
	})
}

//Startup ..
func (app *Application) Startup(fns ...func() error) error {
	app.initialize()
	//app.initGovernor()
	// if err := app.startup(); err != nil {
	// 	return err
	// }
	return xgo.SerialUntilError(fns...)()
}

// Serve start server
func (app *Application) Serve(servers ...server.Server) error {
	for _, svr := range servers {
		app.componentHeap.Push(registry.ServerComponent{Server: svr}, 10)
	}
	return nil
}

// Schedule ..
func (app *Application) Schedule(w worker.Worker) error {
	return app.AddComponent(15, worker.WorkerComponent{Worker: w})
}

func (app *Application) AddComponent(priority int, comps ...component.Component) error {
	var errs error
	for _, comp := range comps {
		errs = multierr.Append(errs, app.componentHeap.Push(comp, priority))
	}
	return errs
}

// Job ..
func (app *Application) Job(runner job.Runner) error {
	namedJob, ok := runner.(interface{ GetJobName() string })
	// job runner must implement GetJobName
	if !ok {
		return nil
	}
	jobName := namedJob.GetJobName()
	if flag.Bool("disable-job") {
		app.logger.Info("jupiter disable job", xlog.FieldName(jobName))
		return nil
	}

	// start job by name
	jobFlag := flag.String("job")
	if jobFlag == "" {
		app.logger.Error("jupiter jobs flag name empty", xlog.FieldName(jobName))
		return nil
	}

	if jobName != jobFlag {
		app.logger.Info("jupiter disable jobs", xlog.FieldName(jobName))
		return nil
	}
	app.logger.Info("jupiter register job", xlog.FieldName(jobName))
	app.jobs[jobName] = runner
	return nil
}

// SetRegistry set customize registry
// Deprecated, please use registry.DefaultRegisterer instead.
func (app *Application) SetReistry(reg registry.Registry) {
	registry.DefaultRegisterer = reg
}

// SetGovernor set governor addr (default 127.0.0.1:0)
// Deprecated
//func (app *Application) SetGovernor(addr string) {
//	app.governorAddr = addr
//}

func (app *Application) Start(stop <-chan struct{}) error {
	return nil
}

// Run run application
func (app *Application) Run() error {
	app.waitSignals() //start signal listen task in goroutine
	defer app.clean()

	// todo jobs not graceful
	app.startJobs()

	// app.components = append(app.components, governor.StdConfig("governor").MustBuild())

	for {
		if app.componentHeap.Len() > 0 {
			var comp, err = app.componentHeap.Pop()
			if err != nil {
				panic(err)
			}
			c := comp.(component.Component)
			xlog.Infof("start component %s", c.Name())
			if err = c.Start(app.stopped); err != nil {
				xlog.Errorf("start component %s failed %+v", c.Name(), err)
			} else {
				xlog.Infof("start component %s succeeded", c.Name())
			}
		}
	}

	// for _, component := range app.components {
	//component.Start(app.stopped)
	//}

	// start servers and govern server
	// app.cycle.Run(app.startServers)
	// start workers
	app.cycle.Run(app.startWorkers)

	//blocking and wait quit
	if err := <-app.cycle.Wait(); err != nil {
		app.logger.Error("jupiter shutdown with error", xlog.FieldMod(ecode.ModApp), xlog.FieldErr(err))
		return err
	}
	app.logger.Info("shutdown jupiter, bye!", xlog.FieldMod(ecode.ModApp))
	return nil
}

//clean after app quit
func (app *Application) clean() {
	_ = xlog.DefaultLogger.Flush()
	_ = xlog.JupiterLogger.Flush()
}

// Stop application immediately after necessary cleanup
func (app *Application) Stop() (err error) {
	app.stopOnce.Do(func() {
		close(app.stopped)
		app.runHooks(StageBeforeStop)

		//stop workers
		for _, w := range app.workers {
			func(w worker.Worker) {
				app.cycle.Run(w.Stop)
			}(w)
		}
		<-app.cycle.Done()
		app.runHooks(StageAfterStop)
		app.cycle.Close()
	})
	return
}

// GracefulStop application after necessary cleanup
func (app *Application) GracefulStop(ctx context.Context) (err error) {
	app.stopOnce.Do(func() {
		close(app.stopped)
		app.runHooks(StageBeforeStop)

		//stop workers
		for _, w := range app.workers {
			func(w worker.Worker) {
				app.cycle.Run(w.Stop)
			}(w)
		}
		<-app.cycle.Done()
		app.runHooks(StageAfterStop)
		app.cycle.Close()
	})
	return err
}

// waitSignals wait signal
func (app *Application) waitSignals() {
	app.logger.Info("init listen signal", xlog.FieldMod(ecode.ModApp), xlog.FieldEvent("init"))
	signals.Shutdown(func(grace bool) { //when get shutdown signal
		//todo: support timeout
		if grace {
			app.GracefulStop(context.TODO())
		} else {
			app.Stop()
		}
	})
}

func (app *Application) startWorkers() error {
	var eg errgroup.Group
	// start multi workers
	for _, w := range app.workers {
		w := w
		eg.Go(func() error {
			return w.Run()
		})
	}
	return eg.Wait()
}

// todo handle error
func (app *Application) startJobs() error {
	if len(app.jobs) == 0 {
		return nil
	}
	var jobs = make([]func(), 0)
	//warp jobs
	for name, runner := range app.jobs {
		jobs = append(jobs, func() {
			app.logger.Info("job run begin", xlog.FieldName(name))
			defer app.logger.Info("job run end", xlog.FieldName(name))
			// runner.Run panic 错误在更上层抛出
			runner.Run()
		})
	}
	xgo.Parallel(jobs...)()
	return nil
}
