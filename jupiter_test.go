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
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/douyu/jupiter/pkg/server/xgrpc"
	"google.golang.org/grpc"

	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/util/xcycle"
	"github.com/douyu/jupiter/pkg/util/xdefer"
	"github.com/douyu/jupiter/pkg/worker"
	"github.com/douyu/jupiter/pkg/xlog"
	. "github.com/smartystreets/goconvey/convey"
)

type testServer struct {
	ServeBlockTime time.Duration
	ServeErr       error

	StopBlockTime time.Duration
	StopErr       error

	GstopBlockTime time.Duration
	GstopErr       error
}

func (s *testServer) Serve() error {
	time.Sleep(s.ServeBlockTime)
	return s.ServeErr
}
func (s *testServer) Stop() error {
	time.Sleep(s.StopBlockTime)
	return s.StopErr
}
func (s *testServer) GracefulStop(ctx context.Context) error {
	time.Sleep(s.GstopBlockTime)
	return s.GstopErr
}
func (s *testServer) Info() *server.ServiceInfo {
	return &server.ServiceInfo{}
}
func TestApplication_Run_1(t *testing.T) {
	Convey("test application run serve", t, func(c C) {
		srv := &testServer{
			ServeErr: errors.New("when server call serve error"),
		}
		app := &Application{}
		app.initialize()
		err := app.Serve(srv)
		So(err, ShouldBeNil)
		go func() {
			// make sure Serve() is called
			time.Sleep(time.Millisecond * 1500)
			err = app.Stop()
			c.So(err, ShouldBeNil)
		}()
		err = app.Run()
		So(err, ShouldEqual, srv.ServeErr)
	})
	Convey("test application run serve block", t, func(c C) {
		srv := &testServer{
			ServeBlockTime: time.Second,
			ServeErr:       errors.New("when server call serve error"),
		}
		app := &Application{}
		app.initialize()
		err := app.Serve(srv)
		So(err, ShouldBeNil)
		go func() {
			// make sure Serve() is called
			time.Sleep(time.Millisecond * 1500)
			err = app.Stop()
			c.So(err, ShouldBeNil)
		}()
		err = app.Run()
		So(err, ShouldEqual, srv.ServeErr)
	})
	Convey("test application run stop", t, func(c C) {
		srv := &testServer{
			ServeBlockTime: time.Second * 2,
			StopBlockTime:  time.Second,
			StopErr:        errors.New("when server call stop error"),
		}
		app := &Application{}
		app.initialize()
		err := app.Serve(srv)
		So(err, ShouldBeNil)
		go func() {
			// make sure Serve() is called
			time.Sleep(time.Millisecond * 1500)
			err = app.Stop()
			c.So(err, ShouldBeNil)
		}()
		err = app.Run()
		So(err, ShouldEqual, srv.StopErr)
	})
}

func TestApplication_initialize(t *testing.T) {
	Convey("test application initialize", t, func() {
		app := &Application{}
		app.initialize()
		So(app.servers, ShouldNotBeNil)
		So(app.workers, ShouldNotBeNil)
		So(app.logger, ShouldNotBeNil)
		So(app.afterStop, ShouldNotBeNil)
		So(app.beforeStop, ShouldNotBeNil)
		So(app.cycle, ShouldNotBeNil)
	})
}

func TestApplication_Startup(t *testing.T) {
	//Convey("test application startup error", t, func() {
	//	app := &Application{}
	//	startUpErr := errors.New("throw startup error")
	//	err := app.Startup(func() error {
	//		return startUpErr
	//	})
	//	So(err,ShouldEqual,startUpErr)
	//})
	//
	//Convey("test application startup nil", t, func() {
	//	app := &Application{}
	//	err := app.Startup(func() error {
	//		return nil
	//	})
	//	So(err,ShouldBeNil)
	//})
}

func TestApplication_Defer(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	type args struct {
		fns []func() error
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			app.Defer(tt.args.fns...)
		})
	}
}

type stopInfo struct {
	state bool
}

func (info *stopInfo) Stop() error {
	info.state = true
	return nil
}

func TestApplication_BeforeStop(t *testing.T) {
	Convey("test application before stop", t, func() {
		si := &stopInfo{}
		app := &Application{}
		app.BeforeStop(si.Stop)
		err := app.Run()
		So(err, ShouldBeNil)
		So(si.state, ShouldEqual, false)

		err = app.Stop()
		So(err, ShouldBeNil)
		So(si.state, ShouldEqual, true)
	})
}

func TestApplication_AfterStop(t *testing.T) {
	Convey("test application after stop", t, func() {
		si := &stopInfo{}
		app := &Application{}
		app.AfterStop(si.Stop)
		err := app.Run()
		So(err, ShouldBeNil)
		So(si.state, ShouldEqual, true)
	})
}

func TestApplication_Serve(t *testing.T) {
	Convey("test application serve throw wrong ip", t, func(c C) {
		app := &Application{}
		grpcConfig := xgrpc.DefaultConfig()
		grpcConfig.Port = 0
		app.initialize()
		err := app.Serve(grpcConfig.Build())
		So(err, ShouldBeNil)
		go func() {
			// make sure Serve() is called
			time.Sleep(time.Millisecond * 1500)
			err = app.Stop()
			c.So(err, ShouldBeNil)
		}()
		err = app.Run()
		So(err, ShouldEqual, grpc.ErrServerStopped)

	})
}

func TestApplication_Schedule(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	type args struct {
		w worker.Worker
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			if err := app.Schedule(tt.args.w); (err != nil) != tt.wantErr {
				t.Errorf("Application.Schedule() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_SetRegistry(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	type args struct {
		reg registry.Registry
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			app.SetRegistry(tt.args.reg)
		})
	}
}

func TestApplication_Run(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			if err := app.Run(); (err != nil) != tt.wantErr {
				t.Errorf("Application.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_Stop(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			if err := app.Stop(); (err != nil) != tt.wantErr {
				t.Errorf("Application.Stop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_GracefulStop(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			if err := app.GracefulStop(tt.args.ctx); (err != nil) != tt.wantErr {
				t.Errorf("Application.GracefulStop() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_waitSignals(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			app.waitSignals()
		})
	}
}

func TestApplication_initGovernor(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			if err := app.initGovernor(); (err != nil) != tt.wantErr {
				t.Errorf("Application.initGovernor() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_startServers(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			if err := app.startServers(); (err != nil) != tt.wantErr {
				t.Errorf("Application.startServers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_startWorkers(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			if err := app.startWorkers(); (err != nil) != tt.wantErr {
				t.Errorf("Application.startWorkers() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_parseFlags(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			if err := app.parseFlags(); (err != nil) != tt.wantErr {
				t.Errorf("Application.parseFlags() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_clean(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	tests := []struct {
		name   string
		fields fields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			app.clean()
		})
	}
}

func TestApplication_loadConfig(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			if err := app.loadConfig(); (err != nil) != tt.wantErr {
				t.Errorf("Application.loadConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_initLogger(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			if err := app.initLogger(); (err != nil) != tt.wantErr {
				t.Errorf("Application.initLogger() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_initTracer(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			if err := app.initTracer(); (err != nil) != tt.wantErr {
				t.Errorf("Application.initTracer() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_initSentinel(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			if err := app.initSentinel(); (err != nil) != tt.wantErr {
				t.Errorf("Application.initSentinel() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_initMaxProcs(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			if err := app.initMaxProcs(); (err != nil) != tt.wantErr {
				t.Errorf("Application.initMaxProcs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestApplication_printBanner(t *testing.T) {
	type fields struct {
		cycle       *xcycle.Cycle
		stopOnce    sync.Once
		initOnce    sync.Once
		startupOnce sync.Once
		afterStop   *xdefer.DeferStack
		beforeStop  *xdefer.DeferStack
		defers      []func() error
		servers     []server.Server
		workers     []worker.Worker
		logger      *xlog.Logger
		registerer  registry.Registry
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{
				cycle:       tt.fields.cycle,
				stopOnce:    tt.fields.stopOnce,
				initOnce:    tt.fields.initOnce,
				startupOnce: tt.fields.startupOnce,
				afterStop:   tt.fields.afterStop,
				beforeStop:  tt.fields.beforeStop,
				defers:      tt.fields.defers,
				servers:     tt.fields.servers,
				workers:     tt.fields.workers,
				logger:      tt.fields.logger,
				registerer:  tt.fields.registerer,
			}
			if err := app.printBanner(); (err != nil) != tt.wantErr {
				t.Errorf("Application.printBanner() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
