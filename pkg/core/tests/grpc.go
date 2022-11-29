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

package tests

import (
	"context"
	"reflect"
	"time"

	"github.com/onsi/ginkgo/v2"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/proto"
)

type GRPCTestCase struct {
	Addr    string
	Timeout time.Duration
	Method  string
	Args    proto.Message

	ExpectError    error
	ExpectMetadata metadata.MD
	ExpectReply    proto.Message
}

// RunGRPCTestCase runs a test case against the given handler.
func RunGRPCTestCase(gtc GRPCTestCase) {
	ginkgoT := ginkgo.GinkgoT()

	if gtc.Timeout == 0 {
		gtc.Timeout = time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), gtc.Timeout)
	defer cancel()

	clientConn, err := grpc.DialContext(ctx, gtc.Addr,
		grpc.WithBlock(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.Nil(ginkgoT, err)

	reply := reflect.New(reflect.TypeOf(gtc.ExpectReply).Elem())
	metadata := metadata.New(nil)

	err = clientConn.Invoke(ctx, gtc.Method, gtc.Args, reply.Interface(),
		grpc.Header(&metadata))
	assert.Equal(ginkgoT, gtc.ExpectError, err)

	assert.True(ginkgoT, proto.Equal(gtc.ExpectReply, reply.Interface().(proto.Message)),
		"expected: %s\nactually: %s", gtc.ExpectReply, reply.Interface().(proto.Message))

	if gtc.ExpectMetadata != nil {
		assert.Equal(ginkgoT, gtc.ExpectMetadata, metadata,
			"expected: %s\nactually: %s", gtc.ExpectMetadata, metadata)
	}

	assert.Nil(ginkgoT, clientConn.Close())
}
