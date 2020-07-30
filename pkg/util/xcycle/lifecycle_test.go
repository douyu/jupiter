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
	"testing"
	"time"
)

//TestCycleDone
func TestCycleDone(t *testing.T) {
	state := "init"
	c := NewCycle()
	c.Run(func() error {
		time.Sleep(time.Microsecond * 100)
		return nil
	})
	go func() {
		state = "done"
		<-c.Done()
	}()
	go func() {
		time.Sleep(time.Microsecond * 200)
		state = "close"
		c.Close()
	}()
	<-c.Wait()

	want := "done"
	if state != want {
		t.Errorf("TestCycleDone error want: %v, ret: %v\r\n", want, state)
	}
}

//TestCycleClose
func TestCycleClose(t *testing.T) {
	state := "init"
	c := NewCycle()
	c.Run(func() error {
		time.Sleep(time.Microsecond * 100)
		return nil
	})
	go func() {
		time.Sleep(time.Microsecond * 200)
		<-c.Done()
		state = "done"
	}()
	go func() {
		c.Close()
		state = "close"
	}()
	<-c.Wait()
	want := "close"
	if state != want {
		t.Errorf("TestCycleClose error want: %v, ret: %v\r\n", want, state)
	}
}
