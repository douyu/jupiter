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

package xcycle

import (
	"sync"
	"sync/atomic"

	"golang.org/x/sync/errgroup"
)

//Cycle ..
type Cycle struct {
	mu uint32
	sync.Once
	eg   *errgroup.Group
	quit chan struct{}
}

//NewCycle new a cycle life
func NewCycle() *Cycle {
	return &Cycle{
		mu:   0,
		eg:   &errgroup.Group{},
		quit: make(chan struct{}),
	}
}

//Run a new goroutine
func (c *Cycle) Run(fn func() error) {
	c.eg.Go(fn)
}

//Done block and return a chan error
func (c *Cycle) Done() <-chan error {
	errCh := make(chan error)
	go func() {
		if err := c.eg.Wait(); err != nil {
			errCh <- err
		}
		close(errCh)
	}()
	return errCh
}

//DoneAndClose ..
func (c *Cycle) DoneAndClose() {
	<-c.Done()
	c.Close()
}

//Close ..
func (c *Cycle) Close() {
	if c.mu == 0 {
		if atomic.CompareAndSwapUint32(&c.mu, 0, 1) {
			close(c.quit)
		}
	}
}

// Wait blocked for a life cycle
func (c *Cycle) Wait() <-chan struct{} {
	return c.quit
}
