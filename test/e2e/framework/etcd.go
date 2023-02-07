// Copyright 2022 Douyu
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

package framework

import (
	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/registry/etcdv3"
	"github.com/imdario/mergo"
	"github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
)

type ETCDTestCase struct {
	Conf *etcdv3.Config
	DoFn func(registry.Registry) (interface{}, error)

	ExpectError error
	ExpectReply interface{}
}

// RunETCDTestCase runs a test case against the given handler.
func RunETCDTestCase(tc ETCDTestCase) {
	ginkgoT := ginkgo.GinkgoT()

	err := mergo.Merge(tc.Conf, etcdv3.DefaultConfig())
	assert.Nil(ginkgoT, err)

	clientConn := tc.Conf.MustBuild()

	reply, err := tc.DoFn(clientConn)
	assert.Equal(ginkgoT, tc.ExpectError, err,
		"expected: %s\nactually: %s", tc.ExpectError, err)

	assert.Equal(ginkgoT, tc.ExpectReply, reply,
		"expected: %s\nactually: %s", tc.ExpectReply, reply)

	assert.Nil(ginkgoT, clientConn.Close())
}
