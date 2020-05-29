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

// ServiceInfo ...
type ServiceInfo struct {
	Name      string
	Scheme    string
	IP        string
	Port      int
	Weight    float64
	Enable    bool
	Healthy   bool
	Metadata  map[string]string
	Region    string
	Zone      string
	GroupName string
}

// Label ...
func (si ServiceInfo) Label() string {
	return fmt.Sprintf("%s://%s:%d", si.Scheme, si.IP, si.Port)
}

// Server ...
type Server interface {
	Serve() error
	Stop() error
	GracefulStop(ctx context.Context) error
	Info() *ServiceInfo
}
