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

type CallbackPhase int

const (
	CallbackPhasePostStarted CallbackPhase = 1
	CallbackPhasePostStopped CallbackPhase = 2
)

type LeaderElectCallback func(CallbackPhase)

type LeaderElector interface {
	Start(stop <-chan struct{})
	IsLeader() bool
	AddCallbacks(...LeaderElectCallback)
}
