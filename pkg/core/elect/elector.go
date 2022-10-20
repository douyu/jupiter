// Copyright 2021 rex lv
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

package elect

import (
	"sync"

	"github.com/douyu/jupiter/pkg/core/component"
)

var _ component.Manager = &electorComponent{}

type electorComponent struct {
	components    []component.Component
	leaderElector LeaderElector
}

func NewComponent() *electorComponent {
	return &electorComponent{
		components: make([]component.Component, 0),
	}
}

func (e *electorComponent) Start(stop <-chan struct{}) error {
	errCh := make(chan error)
	e.startNonLeaderComponents(stop, errCh)
	e.startLeaderComponents(stop, errCh)

	select {
	case <-stop:
		return nil
	case err := <-errCh:
		return err
	}
}

func (e *electorComponent) AddComponent(components ...component.Component) error {
	e.components = append(e.components, components...)
	return nil
}

func (e *electorComponent) ShouldBeLeader() bool {
	return false
}

func (e *electorComponent) startNonLeaderComponents(stop <-chan struct{}, errCh chan error) {
	for _, item := range e.components {
		if !item.ShouldBeLeader() {
			go func(c component.Component) {
				if err := c.Start(stop); err != nil {
					errCh <- err
				}
			}(item)
		}
	}
}

func (e *electorComponent) startLeaderComponents(stop <-chan struct{}, errCh chan error) {
	var mutex sync.Mutex
	stopCh := make(chan struct{})
	closeCh := func() {
		mutex.Lock()
		defer mutex.Unlock()
		select {
		case <-stopCh:
		default:
			close(stopCh)
		}
	}

	e.leaderElector.AddCallbacks(
		func(phase CallbackPhase) {
			if phase != CallbackPhasePostStarted {
				return
			}
			mutex.Lock()
			defer mutex.Unlock()
			stopCh = make(chan struct{})
			for _, item := range e.components {
				if item.ShouldBeLeader() {
					go func(c component.Component) {
						if err := c.Start(stopCh); err != nil {
							errCh <- err
						}
					}(item)
				}
			}
		},
		func(phase CallbackPhase) {
			if phase != CallbackPhasePostStopped {
				return
			}
			closeCh()
		},
	)

	go e.leaderElector.Start(stop)
	go func() {
		<-stop
		closeCh()
	}()
}
