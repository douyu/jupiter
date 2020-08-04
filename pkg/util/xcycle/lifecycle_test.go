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
	"fmt"
	"testing"
	"time"
)

//TestCycleDone
func TestCycleDone(t *testing.T) {
	state := "init"
	c := NewCycle()
	c.Run(func() error {
		time.Sleep(time.Microsecond * 5)
		return nil
	})
	go func() {
		select {
		case <-c.Done():
			state = "done"
		case <-time.After(time.Millisecond):
			state = "close"
		}
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
		time.Sleep(time.Millisecond * 100)
		return nil
	})
	go func() {
		select {
		case <-c.Done():
			state = "done"
		case <-time.After(time.Millisecond):
			state = "close"
		}
		c.Close()
	}()
	<-c.Wait()
	want := "close"
	if state != want {
		t.Errorf("TestCycleClose error want: %v, ret: %v\r\n", want, state)
	}
}

func TestCycleDoneAndClose(t *testing.T) {
	ch := make(chan string, 2)
	state := "init"
	c := NewCycle()
	c.Run(func() error {
		time.Sleep(time.Microsecond * 100)
		return nil
	})
	go func() {
		c.DoneAndClose()
		ch <- "close"
	}()
	<-c.Wait()
	want := "close"
	state = <-ch
	if state != want {
		t.Errorf("TestCycleClose error want: %v, ret: %v\r\n", want, state)
	}
}
func TestCycleWithError(t *testing.T) {
	c := NewCycle()
	c.Run(func() error {
		return fmt.Errorf("run error")
	})
	err := <-c.Wait()
	want := fmt.Errorf("run error")
	if err.Error() != want.Error() {
		t.Errorf("TestCycleClose error want: %v, ret: %v\r\n", want.Error(), err.Error())
	}
}
