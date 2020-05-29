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

package compound

import (
	"context"

	registry2 "github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/server"
	"golang.org/x/sync/errgroup"
)

type compoundRegistry struct {
	registries []registry2.Registry
}

// RegisterService ...
func (c compoundRegistry) RegisterService(ctx context.Context, bean *server.ServiceInfo) error {
	var eg errgroup.Group
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			return registry.RegisterService(ctx, bean)
		})
	}
	return eg.Wait()
}

// DeregisterService ...
func (c compoundRegistry) DeregisterService(ctx context.Context, bean *server.ServiceInfo) error {
	var eg errgroup.Group
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			return registry.DeregisterService(ctx, bean)
		})
	}
	return eg.Wait()
}

// Close ...
func (c compoundRegistry) Close() error {
	var eg errgroup.Group
	for _, registry := range c.registries {
		registry := registry
		eg.Go(func() error {
			return registry.Close()
		})
	}
	return eg.Wait()
}

// New ...
func New(registries ...registry2.Registry) registry2.Registry {
	return compoundRegistry{
		registries: registries,
	}
}
