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
	"testing"

	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

func Test_ServiceInfo(t *testing.T) {

	info1 := ServiceInfo{
		Name:    "main",
		Address: "127.0.0.1:1234",
		Weight:  100,
		Kind:    constant.ServiceProvider,
		Group:   "g1",
	}

	info2 := ServiceInfo{
		Name:    "main",
		Address: "127.0.0.1:1234",
		Weight:  100,
		Kind:    constant.ServiceProvider,
		Group:   "g1",
	}

	info3 := ServiceInfo{
		Name:    "main",
		Address: "127.0.0.1:1235",
		Weight:  100,
		Kind:    constant.ServiceProvider,
		Group:   "g1",
	}

	var (
		address1, address2, address3 resolver.Address
	)

	address1.Addr = info1.Address
	address1.Attributes = attributes.New(constant.KeyServiceInfo, info1)

	address2.Addr = info2.Address
	address2.Attributes = attributes.New(constant.KeyServiceInfo, info2)

	address3.Addr = info3.Address
	address3.Attributes = attributes.New(constant.KeyServiceInfo, info3)

	// the Equal method will check the info which added to attributes,
	// two attributes with the same content are Equal.
	if !address1.Equal(address2) {
		t.Fatalf("%+v.Equals(%+v) = false; want true", address1, address2)
	}
	if !address2.Equal(address1) {
		t.Fatalf("%+v.Equals(%+v) = false; want true", address2, address1)
	}

	if address1.Equal(address3) {
		t.Fatalf("%+v.Equals(%+v) = true; want false", address1, address3)
	}

	if address3.Equal(address1) {
		t.Fatalf("%+v.Equals(%+v) = true; want false", address3, address1)
	}

}

// Reproduce the pannic problem of  issue  #293
// The structure #ServiceInfo# inside the test case
// does not implement the equal method, so the
// comparison at runtime results in panic
func TestNotImplementEqual(t *testing.T) {

	// The previous structure:  Equal method is not implemented
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
		Group string `json:"group"`
	}

	info1 := ServiceInfo{
		Name:    "main",
		Address: "127.0.0.1:1234",
		Weight:  100,
		Kind:    constant.ServiceProvider,
		Group:   "g1",
	}

	info2 := ServiceInfo{
		Name:    "main",
		Address: "127.0.0.1:1234",
		Weight:  100,
		Kind:    constant.ServiceProvider,
		Group:   "g1",
	}

	var (
		address1, address2 resolver.Address
	)

	// Attributes as above
	address1.Addr = info1.Address
	address1.Attributes = attributes.New(constant.KeyServiceInfo, info1)

	address2.Addr = info2.Address
	address2.Attributes = attributes.New(constant.KeyServiceInfo, info2)

	assert.Panics(t, func() {
		// This will cause panic
		address1.Equal(address2)
	})
}

func TestServer(t *testing.T) {
	server := ApplyOptions(
		WithScheme("grpc"),
		WithAddress("127.0.0.1"),
		WithKind(constant.ServiceGovernor),
		WithMetaData("zone", "wh"),
	)

	assert.Equal(t, "grpc:server.test:v1:unknown-mode/127.0.0.1", server.RegistryName())
	assert.Equal(t, "grpc:server.test:v1:unknown-mode/", server.ServicePrefix())
	assert.Equal(t, "grpc://127.0.0.1", server.Label())
}
