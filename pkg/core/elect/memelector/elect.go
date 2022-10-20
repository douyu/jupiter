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

package memelector

import "github.com/douyu/jupiter/pkg/core/elect"

type noopLeaderElector struct {
	alwaysLeader bool
	callbacks    []elect.LeaderElectCallback
}

var _ elect.LeaderElector = &noopLeaderElector{}

func NewAlwaysLeaderElector() elect.LeaderElector {
	return &noopLeaderElector{
		alwaysLeader: true,
	}
}

func NewNeverLeaderElector() elect.LeaderElector {
	return &noopLeaderElector{
		alwaysLeader: false,
	}
}

func (n *noopLeaderElector) AddCallbacks(callbacks ...elect.LeaderElectCallback) {
	n.callbacks = append(n.callbacks, callbacks...)
}

func (n *noopLeaderElector) IsLeader() bool {
	return n.alwaysLeader
}

func (n *noopLeaderElector) Start(stop <-chan struct{}) {
	if n.alwaysLeader {
		for _, callback := range n.callbacks {
			callback(elect.CallbackPhasePostStarted)
		}
	}
}
