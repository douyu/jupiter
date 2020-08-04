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

	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/constant"
)

type Option func(c *ServiceInfo)

// ServiceConfigurator represents service configurator
type ConfigInfo struct {
	Routes []Route
}

// ServiceInfo represents service info
type ServiceInfo struct {
	Name     string               `json:"name"`
	AppID    string               `json:"appId"`
	Scheme   string               `json:"scheme"`
	Address  string               `json:"address"`
	Weight   float64              `json:"weight"`
	Enable   bool                 `json:"enable"`
	Healthy  bool                 `json:"healthy"`
	Metadata map[string]string    `json:"metadata"`
	Region   string               `json:"region"`
	Zone     string               `json:"zone"`
	Kind     constant.ServiceKind `json:"kind"`
	// Deployment 部署组: 不同组的流量隔离
	// 比如某些服务给内部调用和第三方调用，可以配置不同的deployment,进行流量隔离
	Deployment string `json:"deployment"`
	// Group 流量组: 流量在Group之间进行负载均衡
	Group    string              `json:"group"`
	Services map[string]*Service `json:"services" toml:"services"`
}

// Service ...
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

func ApplyOptions(options ...Option) ServiceInfo {
	info := defaultServiceInfo()
	for _, option := range options {
		option(&info)
	}
	return info
}

func WithMetaData(key, value string) Option {
	return func(c *ServiceInfo) {
		c.Metadata[key] = value
	}
}

func WithScheme(scheme string) Option {
	return func(c *ServiceInfo) {
		c.Scheme = scheme
	}
}

func WithAddress(address string) Option {
	return func(c *ServiceInfo) {
		c.Address = address
	}
}

func WithKind(kind constant.ServiceKind) Option {
	return func(c *ServiceInfo) {
		c.Kind = kind
	}
}

func defaultServiceInfo() ServiceInfo {
	si := ServiceInfo{
		Name:       pkg.Name(),
		AppID:      pkg.AppID(),
		Weight:     100,
		Enable:     true,
		Healthy:    true,
		Metadata:   make(map[string]string),
		Region:     pkg.AppRegion(),
		Zone:       pkg.AppZone(),
		Kind:       0,
		Deployment: "",
		Group:      "",
	}
	si.Metadata["appMode"] = pkg.AppMode()
	si.Metadata["appHost"] = pkg.AppHost()
	si.Metadata["startTime"] = pkg.StartTime()
	si.Metadata["buildTime"] = pkg.BuildTime()
	si.Metadata["appVersion"] = pkg.AppVersion()
	si.Metadata["jupiterVersion"] = pkg.JupiterVersion()
	return si
}
