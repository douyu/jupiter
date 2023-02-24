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
// See the License for the specific language governing permissions and
// limitations under the License.

package application

import (
	"context"
	"fmt"
	"sync"
	"time"

	//go-lint
	_ "github.com/douyu/jupiter/pkg/conf/datasource/etcdv3"
	_ "github.com/douyu/jupiter/pkg/conf/datasource/file"
	_ "github.com/douyu/jupiter/pkg/conf/datasource/http"
	_ "github.com/douyu/jupiter/pkg/core/autoproc"
	_ "github.com/douyu/jupiter/pkg/core/rocketmq"
	_ "github.com/douyu/jupiter/pkg/core/xgrpclog"
	_ "github.com/douyu/jupiter/pkg/registry/etcdv3"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/core/component"
	"github.com/douyu/jupiter/pkg/core/ecode"
	"github.com/douyu/jupiter/pkg/core/hooks"
	"github.com/douyu/jupiter/pkg/core/signals"
	"github.com/douyu/jupiter/pkg/executor"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/util/xcycle"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/worker"
	job "github.com/douyu/jupiter/pkg/worker/xjob"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/fatih/color"
	"golang.org/x/sync/errgroup"
)

// Application is the framework's instance, it contains the servers, workers, client and configuration settings.
// Create an instance of Application, by using &Application{}
type Application struct {
	cycle    *xcycle.Cycle
	smu      *sync.RWMutex
	initOnce sync.Once
	// startupOnce  sync.Once
	stopOnce sync.Once
	servers  []server.Server
	workers  []worker.Worker
	jobs     map[string]job.Runner
	logger   *xlog.Logger
	// hooks        map[uint32]*xdefer.DeferStack
	configParser conf.Unmarshaller
	disableMap   map[Disable]bool
	HideBanner   bool
	stopped      chan struct{}
	components   []component.Component
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

// run hooks
func (app *Application) runHooks(stage hooks.Stage) {
	hooks.Do(stage)
}

// RegisterHooks register a stage Hook
func (app *Application) RegisterHooks(stage hooks.Stage, fns ...func()) {
	hooks.Register(stage, fns...)
}

// initialize application
func (app *Application) initialize() {
	app.initOnce.Do(func() {
		//assign
		app.cycle = xcycle.NewCycle()
		app.smu = &sync.RWMutex{}
		app.servers = make([]server.Server, 0)
		app.workers = make([]worker.Worker, 0)
		app.jobs = make(map[string]job.Runner)
		app.logger = xlog.Jupiter()
		app.configParser = toml.Unmarshal
		app.disableMap = make(map[Disable]bool)
		app.stopped = make(chan struct{})
		app.components = make([]component.Component, 0)
		//private method

		_ = app.parseFlags()
		_ = app.printBanner()
		// app.initLogger()
	})
}

// func (app *Application) initLogger() {
// 	xgrpclog.SetLogger(xlog.Jupiter())
// 	rocketmq.SetLogger(xlog.Jupiter())
// }

// // start up application
// // By default the startup composition is:
// // - parse config, watch, version flags
// // - load config
// // - init default biz logger, jupiter frame logger
// // - init procs
// func (app *Application) startup() (err error) {
// 	app.startupOnce.Do(func() {
// 		err = xgo.SerialUntilError(
// 			app.parseFlags,
// 			// app.printBanner,
// 			// app.loadConfig,
// 			// app.initLogger,
// 			// app.initMaxProcs,
// 			// app.initTracer,
// 			// app.initSentinel,
// 			// app.initGovernor,
// 		)()
// 	})
// 	return
// }

// Startup ..
func (app *Application) Startup(fns ...func() error) error {
	app.initialize()
	// if err := app.startup(); err != nil {
	// 	return err
	// }
	return xgo.SerialUntilError(fns...)()
}

// Defer ..
// Deprecated: use AfterStop instead
// func (app *Application) Defer(fns ...func() error) {
// 	app.AfterStop(fns...)
// }

// BeforeStop hook
// Deprecated: use RegisterHooks instead
// func (app *Application) BeforeStop(fns ...func() error) {
// 	app.RegisterHooks(StageBeforeStop, fns...)
// }

// AfterStop hook
// Deprecated: use RegisterHooks instead
// func (app *Application) AfterStop(fns ...func() error) {
// 	app.RegisterHooks(StageAfterStop, fns...)
// }

// Serve start server
func (app *Application) Serve(s ...server.Server) error {
	app.smu.Lock()
	defer app.smu.Unlock()
	app.servers = append(app.servers, s...)
	return nil
}

// Schedule ..
func (app *Application) Schedule(w worker.Worker) error {
	app.workers = append(app.workers, w)
	return nil
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

// Executor ...
func (app *Application) Executor(e executor.Executor) {
	executor.Register(e.GetAddress(), e)
}

// SetRegistry set customize registry
func (app *Application) SetRegistry(reg registry.Registry) {
	registry.DefaultRegisterer = reg
}

// SetGovernor set governor addr (default 127.0.0.1:0)
// Deprecated
//func (app *Application) SetGovernor(addr string) {
//	app.governorAddr = addr
//}

// Run run application
func (app *Application) Run(servers ...server.Server) error {
	app.smu.Lock()
	app.servers = append(app.servers, servers...)
	app.smu.Unlock()

	hooks.Do(hooks.Stage_BeforeRun)

	app.waitSignals() //start signal listen task in goroutine
	defer app.clean()

	// todo jobs not graceful
	_ = app.startJobs()

	// start servers and govern server
	app.cycle.Run(app.startServers)
	// start workers
	app.cycle.Run(app.startWorkers)
	// start executors
	app.cycle.Run(app.startExecutors)
	//blocking and wait quit
	if err := <-app.cycle.Wait(); err != nil {
		app.logger.Error("jupiter shutdown with error", xlog.FieldMod(ecode.ModApp), xlog.FieldErr(err))
		return err
	}
	app.logger.Info("shutdown jupiter, bye!", xlog.FieldMod(ecode.ModApp))
	return nil
}

// clean after app quit
func (app *Application) clean() {
	_ = xlog.Default().Sync()
	_ = xlog.Jupiter().Sync()
}

// Stop application immediately after necessary cleanup
func (app *Application) Stop() (err error) {
	app.stopOnce.Do(func() {
		app.stopped <- struct{}{}
		app.runHooks(hooks.Stage_BeforeStop)
		//stop servers
		for _, s := range app.servers {
			func(s server.Server) {
				app.smu.RLock()
				// unregister before stop
				e := registry.DefaultRegisterer.UnregisterService(context.Background(), s.Info())
				if e != nil {
					app.logger.Error("exit server", xlog.FieldMod(ecode.ModApp), xlog.FieldEvent("stop"), xlog.FieldName(s.Info().Name), xlog.FieldAddr(s.Info().Label()), xlog.FieldErr(err))
				}
				app.logger.Info("exit server", xlog.FieldMod(ecode.ModApp), xlog.FieldEvent("stop"), xlog.FieldName(s.Info().Name), xlog.FieldAddr(s.Info().Label()))

				app.cycle.Run(s.Stop)
				app.smu.RUnlock()
			}(s)
		}
		//stop workers
		for _, w := range app.workers {
			func(w worker.Worker) {
				app.cycle.Run(w.Stop)
			}(w)
		}
		app.cycle.Run(executor.Stop)

		<-app.cycle.Done()
		// run hook
		app.runHooks(hooks.Stage_AfterStop)
		app.cycle.Close()
	})
	return
}

// GracefulStop application after necessary cleanup
func (app *Application) GracefulStop(ctx context.Context) (err error) {
	app.stopOnce.Do(func() {
		app.stopped <- struct{}{}
		app.runHooks(hooks.Stage_BeforeStop)
		//stop servers
		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(func() error {
					app.smu.RLock()
					defer app.smu.RUnlock()
					// unregister before graceful stop
					e := registry.DefaultRegisterer.UnregisterService(ctx, s.Info())
					if e != nil {
						app.logger.Error("exit server", xlog.FieldMod(ecode.ModApp), xlog.FieldEvent("graceful stop"), xlog.FieldName(s.Info().Name), xlog.FieldAddr(s.Info().Label()), xlog.FieldErr(err))
					}
					app.logger.Info("exit server", xlog.FieldMod(ecode.ModApp), xlog.FieldEvent("graceful stop"), xlog.FieldName(s.Info().Name), xlog.FieldAddr(s.Info().Label()))

					return s.GracefulStop(ctx)
				})
			}(s)
		}
		//stop workers
		for _, w := range app.workers {
			func(w worker.Worker) {
				app.cycle.Run(w.Stop)
			}(w)
		}
		// stop executor
		app.cycle.Run(executor.GracefulStop)
		<-app.cycle.Done()
		// run hooks
		app.runHooks(hooks.Stage_AfterStop)
		app.cycle.Close()
	})
	return err
}

// waitSignals wait signal
func (app *Application) waitSignals() {
	app.logger.Info("init listen signal", xlog.FieldMod(ecode.ModApp), xlog.FieldEvent("init"))
	signals.Shutdown(func(grace bool) { //when get shutdown signal
		if grace {
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_ = app.GracefulStop(ctx)
		} else {
			_ = app.Stop()
		}
	})
}

func (app *Application) startServers() error {
	var eg errgroup.Group
	var ctx, cancel = context.WithTimeout(context.Background(), 3*time.Second)
	go func() {
		<-app.stopped
		cancel()
	}()
	// start multi servers
	app.smu.Lock()
	for _, s := range app.servers {
		s := s
		eg.Go(func() (err error) {
			time.AfterFunc(time.Second, func() {
				_ = registry.DefaultRegisterer.RegisterService(ctx, s.Info())
				app.logger.Info("start server", xlog.FieldMod(ecode.ModApp), xlog.FieldEvent("init"), xlog.FieldName(s.Info().Name), xlog.FieldAddr(s.Info().Label()), xlog.Any("scheme", s.Info().Scheme))
			})
			err = s.Serve()
			return
		})
	}

	app.smu.Unlock()
	return eg.Wait()
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

// start executor
func (app *Application) startExecutors() error {
	return executor.Run()
}

// parseFlags init
func (app *Application) parseFlags() error {
	if app.isDisable(DisableParserFlag) {
		app.logger.Info("parseFlags disable", xlog.FieldMod(ecode.ModApp))
		return nil
	}

	return flag.Parse()
}

func (app *Application) isDisable(d Disable) bool {
	b, ok := app.disableMap[d]
	if !ok {
		return false
	}
	return b
}

// printBanner init
func (app *Application) printBanner() error {
	if app.HideBanner {
		return nil
	}

	if xdebug.IsTestingMode() {
		return nil
	}

	const banner = `
   (_)_   _ _ __ (_) |_ ___ _ __
   | | | | | '_ \| | __/ _ \ '__|
   | | |_| | |_) | | ||  __/ |
  _/ |\__,_| .__/|_|\__\___|_|
 |__/      |_|

 Welcome to jupiter, starting application ...
`
	fmt.Println(color.GreenString(banner))
	return nil
}
