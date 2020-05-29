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

package registry

import (
	"context"
	"io"

	"github.com/douyu/jupiter/pkg/server"
)

// ServerInstance ...
type ServerInstance struct {
	Scheme string
	IP     string
	Port   int
	Labels map[string]string
}

// Registry register/deregister service
// registry impl should control rpc timeout
type Registry interface {
	RegisterService(context.Context, *server.ServiceInfo) error
	DeregisterService(context.Context, *server.ServiceInfo) error
	io.Closer
}

// Nop registry, used for local development/debugging
type Nop struct{}

// RegisterService ...
func (n Nop) RegisterService(context.Context, *server.ServiceInfo) error { return nil }

// DeregisterService ...
func (n Nop) DeregisterService(context.Context, *server.ServiceInfo) error { return nil }

// Close ...
func (n Nop) Close() error { return nil }

// Configuration ...
type Configuration struct {
}

// Rule ...
type Rule struct {
	Target  string
	Pattern string
}
