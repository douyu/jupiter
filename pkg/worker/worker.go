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

package worker

import (
	"fmt"

	"github.com/douyu/jupiter/pkg/component"
)

// Worker could scheduled by jupiter or customized scheduler
type Worker interface {
	Run() error
	Stop() error
}

var _ component.Component = &WorkerComponent{}

type WorkerComponent struct {
	Worker
}

func (c WorkerComponent) Start(stopCh <-chan struct{}) error {
	var errCh = make(chan error)
	go func() {
		fmt.Println("before work...")
		errCh <- c.Run()
		fmt.Println("after work...")
	}()
	go func() {
		select {
		case <-stopCh:
			fmt.Println("stop...")
			c.Stop()
		case <-errCh:
			fmt.Println("err occur...")
		}
		fmt.Println("close err ch")
		close(errCh)
	}()
	return nil
}

func (c WorkerComponent) ShouldBeLeader() bool {
	return false
}
