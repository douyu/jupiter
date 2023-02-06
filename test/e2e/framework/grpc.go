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
	"context"
	"reflect"

	"github.com/douyu/jupiter/pkg/client/grpc"
	"github.com/imdario/mergo"
	"github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
	ggrpc "google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type GRPCTestCase struct {
	Conf   *grpc.Config
	Method string
	Args   proto.Message

	ExpectError    error
	ExpectMetadata metadata.MD
	ExpectReply    proto.Message
}

// RunGRPCTestCase runs a test case against the given handler.
func RunGRPCTestCase(gtc GRPCTestCase) {
	ginkgoT := ginkgo.GinkgoT()

	err := mergo.Merge(gtc.Conf, grpc.DefaultConfig())
	assert.Nil(ginkgoT, err)

	clientConn, err := gtc.Conf.Build()
	assert.Nil(ginkgoT, err)

	reply := reflect.New(reflect.TypeOf(gtc.ExpectReply).Elem())
	metadata := metadata.New(nil)

	err = clientConn.Invoke(context.Background(), gtc.Method, gtc.Args, reply.Interface(),
		ggrpc.Header(&metadata))
	assert.Equal(ginkgoT, gtc.ExpectError, err,
		"expected: %s\nactually: %s", gtc.ExpectError, err)

	assert.True(ginkgoT, proto.Equal(gtc.ExpectReply, reply.Interface().(proto.Message)),
		"expected: %s\nactually: %s", gtc.ExpectReply, reply.Interface().(proto.Message))

	if gtc.ExpectMetadata != nil {
		assert.Equal(ginkgoT, gtc.ExpectMetadata, metadata,
			"expected: %s\nactually: %s", gtc.ExpectMetadata, metadata)
	}

	assert.Nil(ginkgoT, clientConn.Close())
}
