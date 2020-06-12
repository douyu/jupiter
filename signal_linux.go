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
	"os"
	"os/signal"
	"syscall"
	"time"
)

func hookSignals(app *Application) {
	sigChan := make(chan os.Signal)
	signal.Notify(
		sigChan,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGSTOP,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
		syscall.SIGKILL,
	)

	go func() {
		var sig os.Signal
		for {
			sig = <-sigChan
			switch sig {
			case syscall.SIGQUIT, syscall.SIGSTOP, syscall.SIGUSR1:
				_ = app.Stop() // graceful stop
			case syscall.SIGHUP:
				_ = app.GracefulStop(context.TODO())
			case syscall.SIGINT, syscall.SIGKILL, syscall.SIGUSR2, syscall.SIGTERM:
				_ = app.Stop() // terminate now
			}
			time.Sleep(time.Second * 3)
		}
	}()
}
