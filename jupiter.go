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
	"net/http"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/conf"
	file_datasource "github.com/douyu/jupiter/pkg/datasource/file"
	http_datasource "github.com/douyu/jupiter/pkg/datasource/http"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/flag"
	"github.com/douyu/jupiter/pkg/govern"
	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/sentinel"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/trace"
	"github.com/douyu/jupiter/pkg/trace/jaeger"
	"github.com/douyu/jupiter/pkg/util/xcolor"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/util/xstring"
	"github.com/douyu/jupiter/pkg/worker"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/automaxprocs/maxprocs"
	"golang.org/x/sync/errgroup"
)

// Application is the framework's instance, it contains the servers, workers, client and configuration settings.
// Create an instance of Application, by using &Application{}
type Application struct {
	servers []server.Server
	workers []worker.Worker
	logger  *xlog.Logger

	stopOnce    sync.Once
	initOnce    sync.Once
	startupOnce sync.Once

	registerer registry.Registry

	signalHooker func(*Application)
	defers       []func() error

	governor *http.Server
}

// initialize application
func (app *Application) initialize() {
	app.initOnce.Do(func() {
		app.servers = make([]server.Server, 0)
		app.workers = make([]worker.Worker, 0)
		app.logger = xlog.JupiterLogger
		app.signalHooker = hookSignals
		app.defers = []func() error{}
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
		)()
	})
	return
}

func (app *Application) Startup(fns ...func() error) error {
	app.initialize()
	if err := app.startup(); err != nil {
		return err
	}
	return xgo.SerialUntilError(fns...)()
}

func (app *Application) Defer(fns ...func() error) {
	app.initialize()
	if app.defers == nil {
		app.defers = make([]func() error, 0)
	}
	app.defers = append(app.defers, fns...)
}

func (app *Application) Serve(s server.Server) error {
	app.servers = append(app.servers, s)
	return nil
}

func (app *Application) Schedule(w worker.Worker) error {
	app.workers = append(app.workers, w)
	return nil
}

// SetRegistry set customize registry
func (app *Application) SetRegistry(reg registry.Registry) {
	app.registerer = reg
}

// SetSignalHooker set signal hooker
func (app *Application) SetSignalHooker(hook func(*Application)) {
	app.signalHooker = hook
}

// SetGovernor governor
// 127.0.0.1:9990 as default governor addr
func (app *Application) SetGovernor(addr string) {
	app.governor = &http.Server{
		Handler: govern.DefaultServeMux,
		Addr:    addr,
	}
}

// Run run application
func (app *Application) Run() error {
	defer app.clean()
	if app.signalHooker == nil {
		app.signalHooker = hookSignals
	}
	if app.governor == nil {
		app.governor = &http.Server{
			Handler: govern.DefaultServeMux,
			Addr:    "127.0.0.1:9990", // 默认治理端口
		}
	}

	if app.registerer == nil {
		app.registerer = registry.Nop{}
	}

	app.signalHooker(app)

	// start govern
	var eg errgroup.Group
	eg.Go(app.startGovernor)
	eg.Go(app.startServers)
	eg.Go(app.startWorkers)
	return eg.Wait()
}

// Stop application immediately after necessary cleanup
func (app *Application) Stop() (err error) {
	app.beforeStop()
	app.stopOnce.Do(func() {
		err = app.registerer.Close()
		if err != nil {
			app.logger.Error("stop register close err", xlog.FieldMod(ecode.ModApp), xlog.FieldErr(err))
		}
		err = app.governor.Close()
		if err != nil {
			app.logger.Error("stop governor close err", xlog.FieldMod(ecode.ModApp), xlog.FieldErr(err))
		}
		var eg errgroup.Group
		for _, s := range app.servers {
			s := s
			eg.Go(s.Stop)
		}
		for _, w := range app.workers {
			w := w
			eg.Go(w.Stop)
		}
		err = eg.Wait()
	})
	return
}

// GracefulStop application after necessary cleanup
func (app *Application) GracefulStop(ctx context.Context) (err error) {
	app.beforeStop()
	app.stopOnce.Do(func() {
		err = app.registerer.Close()
		if err != nil {
			app.logger.Error("graceful stop register close err", xlog.FieldMod(ecode.ModApp), xlog.FieldErr(err))
		}
		err = app.governor.Close()
		if err != nil {
			app.logger.Error("graceful stop governor close err", xlog.FieldMod(ecode.ModApp), xlog.FieldErr(err))
		}
		if err != nil {

		}
		var eg errgroup.Group
		for _, s := range app.servers {
			s := s
			eg.Go(func() error {
				return s.GracefulStop(ctx)
			})
		}
		err = eg.Wait()
	})
	return err
}

func (app *Application) beforeStop() {
	// todo(gorexlv): before stop hooks
	app.logger.Info("leaving jupiter, bye....", xlog.FieldMod(ecode.ModApp))
}

func (app *Application) startGovernor() (err error) {
	app.logger.Info("start governor", xlog.FieldMod(ecode.ModApp), xlog.FieldAddr("http://"+app.governor.Addr))
	defer func() {
		if err != nil {
			app.logger.Panic("stop governor", xlog.FieldMod(ecode.ModApp), xlog.FieldErr(err), xlog.FieldAddr("http://"+app.governor.Addr))
		}
	}()
	err = app.governor.ListenAndServe()
	if err == http.ErrServerClosed {
		app.logger.Info("stop governor", xlog.FieldMod(ecode.ModApp), xlog.FieldAddr("http://"+app.governor.Addr))
		return nil
	}

	return err
}

func (app *Application) startServers() error {
	var eg errgroup.Group
	xgo.ParallelWithErrorChan()
	// start multi servers
	for _, s := range app.servers {
		s := s
		eg.Go(func() (err error) {
			_ = app.registerer.RegisterService(context.TODO(), s.Info())
			defer app.registerer.DeregisterService(context.TODO(), s.Info())
			app.logger.Info("start servers", xlog.FieldMod(ecode.ModApp), xlog.FieldAddr(s.Info().Label()), xlog.Any("scheme", s.Info().Scheme))
			defer app.logger.Info("exit server", xlog.FieldMod(ecode.ModApp), xlog.FieldErr(err), xlog.FieldAddr(s.Info().Label()))
			return s.Serve()
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

func (app *Application) parseFlags() error {
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
	return flag.Parse()
}

func (app *Application) clean() {
	for i := len(app.defers) - 1; i >= 0; i-- {
		fn := app.defers[i]
		if err := fn(); err != nil {
			xlog.Error("clean.defer", xlog.String("func", xstring.FunctionName(fn)))
		}
	}
	_ = xlog.DefaultLogger.Flush()
	_ = xlog.JupiterLogger.Flush()
}

func (app *Application) loadConfig() error {
	var (
		watchConfig = flag.Bool("watch")
		configAddr  = flag.String("config")
	)

	if configAddr == "" {
		app.logger.Warn("no config ...")
		return nil
	}
	switch {
	case strings.HasPrefix(configAddr, "http://"),
		strings.HasPrefix(configAddr, "https://"):
		provider := http_datasource.NewDataSource(configAddr, watchConfig)
		if err := conf.LoadFromDataSource(provider, toml.Unmarshal); err != nil {
			app.logger.Panic("load remote config", xlog.FieldMod(ecode.ModConfig), xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err))
		}
		app.logger.Info("load remote config", xlog.FieldMod(ecode.ModConfig), xlog.FieldAddr(configAddr))
	default:
		provider := file_datasource.NewDataSource(configAddr, watchConfig)

		if err := conf.LoadFromDataSource(provider, toml.Unmarshal); err != nil {
			app.logger.Panic("load local file", xlog.FieldMod(ecode.ModConfig), xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err))
		}
		app.logger.Info("load local file", xlog.FieldMod(ecode.ModConfig), xlog.FieldAddr(configAddr))
	}
	return nil
}

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

func (app *Application) initTracer() error {
	// init tracing component jaeger
	if conf.Get("jupiter.trace.jaeger") != nil {
		var config = jaeger.RawConfig("jupiter.trace.jaeger")
		trace.SetGlobalTracer(config.Build())
	}
	return nil
}

func (app *Application) initSentinel() error {
	// init reliability component sentinel
	if conf.Get("jupiter.reliability.sentinel") != nil {
		app.logger.Info("init reliability component sentinel")
		return sentinel.RawConfig("jupiter.reliability.sentinel").
			InitSentinelCoreComponent()
	}
	return nil
}

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
