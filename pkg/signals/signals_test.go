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
	Convey("test shutdown signal is SIGQUIT", t, func(c C) {
		fn := func(grace bool) {
			c.So(grace, ShouldEqual, false)
			if grace != false {
				t.Fatal("SIGQUIT should be not grace")
			}
		}
		Shutdown(fn)
		kill(syscall.SIGQUIT)
	})
}
func TestShutdownSIGINT(t *testing.T) {
	Convey("test shutdown signal is SIGINT", t, func(c C) {
		fn := func(grace bool) {
			if grace != true {
				t.Fatal("SIGINT should be grace")
			}
		}
		Shutdown(fn)
		kill(syscall.SIGQUIT)
	})
}
