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

package component

import "github.com/douyu/jupiter/pkg/metric"

type Component interface {
	// Start blocks until the channel is closed or an error occurs.
	// The component will stop running when the channel is closed.
	Start(<-chan struct{}) error

	ShouldBeLeader() bool
}

var _ Component = ComponentFunc(nil)

type ComponentFunc func(<-chan struct{}) error

func (f ComponentFunc) Start(stop <-chan struct{}) error {
	return f(stop)
}

func (f ComponentFunc) ShouldBeLeader() bool {
	return false
}

// Component manager, aggregate multiple components to one
type Manager interface {
	Component
	AddComponent(...Component) error
}

// Component builder, build component with injecting govern plugin
type Builder interface {
	WithComponentManager(Manager) Builder
	WithMetrics(metric.Metrics) Builder
	Build() Component
}
