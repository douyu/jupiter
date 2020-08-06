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

package signals

import (
	"os"
	"syscall"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func kill(sig os.Signal) {
	pro, _ := os.FindProcess(os.Getpid())
	pro.Signal(sig)
}
func TestShutdownSIGQUIT(t *testing.T) {
	quit := make(chan struct{})
	Convey("test shutdown signal by SIGQUIT", t, func(c C) {
		fn := func(grace bool) {
			c.So(grace, ShouldEqual, false)
			close(quit)
		}
		Shutdown(fn)
		kill(syscall.SIGQUIT)
		<-quit
	})
}

// func TestShutdownSIGINT(t *testing.T) {
// 	quit := make(chan struct{})
// 	Convey("test shutdown signal by SIGINT", t, func(c C) {
// 		fn := func(grace bool) {
// 			c.So(grace, ShouldEqual, true)
// 			close(quit)
// 		}
// 		Shutdown(fn)
// 		kill(syscall.SIGINT)
// 		<-quit
// 	})
// }
