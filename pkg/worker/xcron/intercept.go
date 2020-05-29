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

package xcron

import (
	"sync"
	"time"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/robfig/cron/v3"
)

// delayIfStillRunning serializes jobs, delaying subsequent runs until the
// previous one is complete. Jobs running after a delay of more than a minute
// have the delay logged at Info.
func delayIfStillRunning(logger *xlog.Logger) JobWrapper {
	return func(j Job) Job {
		var mu sync.Mutex
		return cron.FuncJob(func() {
			start := time.Now()
			mu.Lock()
			defer mu.Unlock()
			if dur := time.Since(start); dur > time.Minute {
				logger.Info("cron delay", xlog.String("duration", dur.String()))
			}
			j.Run()
		})
	}
}

// skipIfStillRunning skips an invocation of the Job if a previous invocation is
// still running. It logs skips to the given logger at Info level.
func skipIfStillRunning(logger *xlog.Logger) JobWrapper {
	var ch = make(chan struct{}, 1)
	ch <- struct{}{}
	return func(j Job) Job {
		return cron.FuncJob(func() {
			select {
			case v := <-ch:
				j.Run()
				ch <- v
			default:
				logger.Info("cron skip")
			}
		})
	}
}
