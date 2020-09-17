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

package jupiter

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"sync"

	"github.com/douyu/jupiter/pkg/server/governor"
	job "github.com/douyu/jupiter/pkg/worker/xjob"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/conf"

	//go-lint
	_ "github.com/douyu/jupiter/pkg/datasource/file"
	_ "github.com/douyu/jupiter/pkg/datasource/http"
	"github.com/douyu/jupiter/pkg/datasource/manager"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/sentinel"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/signals"
	"github.com/douyu/jupiter/pkg/trace"
	"github.com/douyu/jupiter/pkg/trace/jaeger"
	"github.com/douyu/jupiter/pkg/util/xcolor"
	"github.com/douyu/jupiter/pkg/util/xcycle"
	"github.com/douyu/jupiter/pkg/util/xdefer"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/worker"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/automaxprocs/maxprocs"
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
	servers      []server.Server
	workers      []worker.Worker
	jobs         map[string]job.Runner
	logger       *xlog.Logger
	registerer   registry.Registry
	hooks        map[uint32]*xdefer.DeferStack
	configParser conf.Unmarshaller
	disableMap   map[Disable]bool
}

//New new a Application
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
		app.servers = make([]server.Server, 0)
		app.workers = make([]worker.Worker, 0)
		app.jobs = make(map[string]job.Runner)
		app.logger = xlog.JupiterLogger
		app.configParser = toml.Unmarshal
		app.disableMap = make(map[Disable]bool)
		//private method
		app.initHooks(StageBeforeStop, StageAfterStop)
		//public method
		app.SetRegistry(registry.Nop{}) //default nop without registry
	})
}

// start up application
// By default the startup composition is:
// - parse config, watch, version flags
// - load config
// - init default biz logger, jupiter frame logger
// - init procs
func (app *Application) startup() (err error) {
	app.startupOnce.Do(func() {
		err = xgo.SerialUntilError(
			app.parseFlags,
			app.printBanner,
			app.loadConfig,
			app.initLogger,
			app.initMaxProcs,
			app.initTracer,
			app.initSentinel,
			app.initGovernor,
		)()
	})
	return
}

//Startup ..
func (app *Application) Startup(fns ...func() error) error {
	app.initialize()
	if err := app.startup(); err != nil {
		return err
	}
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

// SetRegistry set customize registry
func (app *Application) SetRegistry(reg registry.Registry) {
	app.registerer = reg
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

	app.waitSignals() //start signal listen task in goroutine
	defer app.clean()

	// todo jobs not graceful
	app.startJobs()

	// start servers and govern server
	app.cycle.Run(app.startServers)
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
		app.runHooks(StageBeforeStop)

		if app.registerer != nil {
			err = app.registerer.Close()
			if err != nil {
				app.logger.Error("stop register close err", xlog.FieldMod(ecode.ModApp), xlog.FieldErr(err))
			}
		}
		//stop servers
		app.smu.RLock()
		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(s.Stop)
			}(s)
		}
		app.smu.RUnlock()

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
		app.runHooks(StageBeforeStop)

		if app.registerer != nil {
			err = app.registerer.Close()
			if err != nil {
				app.logger.Error("stop register close err", xlog.FieldMod(ecode.ModApp), xlog.FieldErr(err))
			}
		}
		//stop servers
		app.smu.RLock()
		for _, s := range app.servers {
			func(s server.Server) {
				app.cycle.Run(func() error {
					return s.GracefulStop(ctx)
				})
			}(s)
		}
		app.smu.RUnlock()

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

func (app *Application) initGovernor() error {
	if app.isDisable(DisableDefaultGovernor) {
		app.logger.Info("defualt governor disable", xlog.FieldMod(ecode.ModApp))
		return nil
	}

	config := governor.StdConfig("governor")
	if !config.Enable {
		return nil
	}
	return app.Serve(config.Build())
}

func (app *Application) startServers() error {
	var eg errgroup.Group
	// start multi servers
	for _, s := range app.servers {
		s := s
		eg.Go(func() (err error) {
			_ = app.registerer.RegisterService(context.TODO(), s.Info())
			defer app.registerer.UnregisterService(context.TODO(), s.Info())
			app.logger.Info("start server", xlog.FieldMod(ecode.ModApp), xlog.FieldEvent("init"), xlog.FieldName(s.Info().Name), xlog.FieldAddr(s.Info().Label()), xlog.Any("scheme", s.Info().Scheme))
			defer app.logger.Info("exit server", xlog.FieldMod(ecode.ModApp), xlog.FieldEvent("exit"), xlog.FieldName(s.Info().Name), xlog.FieldErr(err), xlog.FieldAddr(s.Info().Label()))
			err = s.Serve()
			return
		})
	}
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

//parseFlags init
func (app *Application) parseFlags() error {
	if app.isDisable(DisableParserFlag) {
		app.logger.Info("parseFlags disable", xlog.FieldMod(ecode.ModApp))
		return nil
	}
	flag.Register(&flag.StringFlag{
		Name:    "config",
		Usage:   "--config",
		EnvVar:  "JUPITER_CONFIG",
		Default: "",
		Action:  func(name string, fs *flag.FlagSet) {},
	})

	flag.Register(&flag.BoolFlag{
		Name:    "watch",
		Usage:   "--watch, watch config change event",
		Default: false,
		EnvVar:  "JUPITER_CONFIG_WATCH",
	})

	flag.Register(&flag.BoolFlag{
		Name:    "version",
		Usage:   "--version, print version",
		Default: false,
		Action: func(string, *flag.FlagSet) {
			pkg.PrintVersion()
			os.Exit(0)
		},
	})

	flag.Register(&flag.StringFlag{
		Name:    "host",
		Usage:   "--host, print host",
		Default: "127.0.0.1",
		Action:  func(string, *flag.FlagSet) {},
	})
	return flag.Parse()
}

//loadConfig init
func (app *Application) loadConfig() error {
	if app.isDisable(DisableLoadConfig) {
		app.logger.Info("load config disable", xlog.FieldMod(ecode.ModConfig))
		return nil
	}

	var configAddr = flag.String("config")
	provider, err := manager.NewDataSource(configAddr)
	if err != manager.ErrConfigAddr {
		if err != nil {
			app.logger.Panic("data source: provider error", xlog.FieldMod(ecode.ModConfig), xlog.FieldErr(err))
		}

		if err := conf.LoadFromDataSource(provider, app.configParser); err != nil {
			app.logger.Panic("data source: load config", xlog.FieldMod(ecode.ModConfig), xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err))
		}
	} else {
		app.logger.Info("no config... ", xlog.FieldMod(ecode.ModConfig))
	}
	return nil
}

//initLogger init
func (app *Application) initLogger() error {
	if conf.Get("jupiter.logger.default") != nil {
		xlog.DefaultLogger = xlog.RawConfig("jupiter.logger.default").Build()
	}
	xlog.DefaultLogger.AutoLevel("jupiter.logger.default")

	if conf.Get("jupiter.logger.jupiter") != nil {
		xlog.JupiterLogger = xlog.RawConfig("jupiter.logger.jupiter").Build()
	}
	xlog.JupiterLogger.AutoLevel("jupiter.logger.jupiter")

	return nil
}

//initTracer init
func (app *Application) initTracer() error {
	// init tracing component jaeger
	if conf.Get("jupiter.trace.jaeger") != nil {
		var config = jaeger.RawConfig("jupiter.trace.jaeger")
		trace.SetGlobalTracer(config.Build())
	}
	return nil
}

//initSentinel init
func (app *Application) initSentinel() error {
	// init reliability component sentinel
	if conf.Get("jupiter.reliability.sentinel") != nil {
		app.logger.Info("init sentinel")
		return sentinel.RawConfig("jupiter.reliability.sentinel").Build()
	}
	return nil
}

//initMaxProcs init
func (app *Application) initMaxProcs() error {
	if maxProcs := conf.GetInt("maxProc"); maxProcs != 0 {
		runtime.GOMAXPROCS(maxProcs)
	} else {
		if _, err := maxprocs.Set(); err != nil {
			app.logger.Panic("auto max procs", xlog.FieldMod(ecode.ModProc), xlog.FieldErrKind(ecode.ErrKindAny), xlog.FieldErr(err))
		}
	}
	app.logger.Info("auto max procs", xlog.FieldMod(ecode.ModProc), xlog.Int64("procs", int64(runtime.GOMAXPROCS(-1))))
	return nil
}

func (app *Application) isDisable(d Disable) bool {
	b, ok := app.disableMap[d]
	if !ok {
		return false
	}
	return b
}

//printBanner init
func (app *Application) printBanner() error {
	const banner = `
   (_)_   _ _ __ (_) |_ ___ _ __
   | | | | | '_ \| | __/ _ \ '__|
   | | |_| | |_) | | ||  __/ |
  _/ |\__,_| .__/|_|\__\___|_|
 |__/      |_|
 
 Welcome to jupiter, starting application ...
`
	fmt.Println(xcolor.Green(banner))
	return nil
}
