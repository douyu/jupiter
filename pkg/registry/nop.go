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

package registry

import (
	"context"

	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/xlog"
)

// Nop registry, used for local development/debugging
type Local struct{}

// ListServices ...
func (n Local) ListServices(ctx context.Context, s string, s2 string) ([]*server.ServiceInfo, error) {
	panic("implement me")
}

// WatchServices ...
func (n Local) WatchServices(ctx context.Context, s string, s2 string) (chan Endpoints, error) {
	panic("implement me")
}

// RegisterService ...
func (n Local) RegisterService(ctx context.Context, si *server.ServiceInfo) error {
	xlog.Jupiter().Info("register service locally", xlog.FieldMod("registry"), xlog.FieldName(si.Name), xlog.FieldAddr(si.Label()))
	return nil
}

// UnregisterService ...
func (n Local) UnregisterService(ctx context.Context, si *server.ServiceInfo) error {
	xlog.Jupiter().Info("unregister service locally", xlog.FieldMod("registry"), xlog.FieldName(si.Name), xlog.FieldAddr(si.Label()))
	return nil
}

// Close ...
func (n Local) Close() error { return nil }

// Close ...
func (n Local) Kind() string { return "local" }
