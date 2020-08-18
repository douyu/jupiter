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
	"os/signal"
	"syscall"
)

//Shutdown suport twice signal must exit
func Shutdown(stop func(grace bool)) {
	sig := make(chan os.Signal, 2)
	signal.Notify(
		sig,
		shutdownSignals...,
	)
	go func() {
		s := <-sig
		go stop(s != syscall.SIGQUIT)
		<-sig
		os.Exit(128 + int(s.(syscall.Signal))) // second signal. Exit directly.
	}()
}
