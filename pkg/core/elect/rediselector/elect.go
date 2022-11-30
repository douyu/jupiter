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

package rediselector

import "github.com/douyu/jupiter/pkg/core/elect"

var _ elect.LeaderElector = &redisLeaderElector{}

type redisLeaderElector struct {
	callbacks []elect.LeaderElectCallback
}

func New() *redisLeaderElector {
	return &redisLeaderElector{
		callbacks: make([]elect.LeaderElectCallback, 0),
	}
}

func (r *redisLeaderElector) IsLeader() bool {
	return false
}

func (r *redisLeaderElector) AddCallbacks(callbacks ...elect.LeaderElectCallback) {

}

func (r *redisLeaderElector) Start(stop <-chan struct{}) {
}
