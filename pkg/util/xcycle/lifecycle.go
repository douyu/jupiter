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
)

//Cycle ..
type Cycle struct {
	mu uint32
	sync.Once
	wg   *sync.WaitGroup
	quit chan error
}

//NewCycle new a cycle life
func NewCycle() *Cycle {
	return &Cycle{
		mu:   0,
		wg:   &sync.WaitGroup{},
		quit: make(chan error),
	}
}

//Run a new goroutine
func (c *Cycle) Run(fn func() error) {
	c.wg.Add(1)
	go func(c *Cycle) {
		defer c.wg.Done()
		if err := fn(); err != nil {
			c.quit <- err
		}
	}(c)

}

//Done block and return a chan error
func (c *Cycle) Done() <-chan error {
	go func() {
		c.wg.Wait()
		c.Close()
	}()
	return c.quit
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
func (c *Cycle) Wait() <-chan error {
	return c.quit
}
