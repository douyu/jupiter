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

package server

import (
	"context"
	"fmt"
)

// ServiceConfigurator represents service configurator
type ServiceConfigurator struct {
	Routes []Route
}

// ServiceInfo represents service info
type ServiceInfo struct {
	Name       string
	Scheme     string
	IP         string
	Port       int
	Address    string
	Weight     float64
	Enable     bool
	Healthy    bool
	Metadata   map[string]string
	Region     string
	Zone       string
	Deployment string
	Services   map[string]*Service `json:"services" toml:"services"`
}

type Service struct {
	Namespace string            `json:"namespace" toml:"namespace"`
	Name      string            `json:"name" toml:"name"`
	Labels    map[string]string `json:"labels" toml:"labels"`
	Methods   []string          `json:"methods" toml:"methods"`
}

// Label ...
func (si ServiceInfo) Label() string {
	return fmt.Sprintf("%s://%s", si.Scheme, si.Address)
}

// Server ...
type Server interface {
	Serve() error
	Stop() error
	GracefulStop(ctx context.Context) error
	Info() *ServiceInfo
}

// Route ...
type Route struct {
	// 权重组，按照
	WeightGroups []WeightGroup
	// 方法名
	Method string
}

// WeightGroup ...
type WeightGroup struct {
	Group  string
	Weight int
}
