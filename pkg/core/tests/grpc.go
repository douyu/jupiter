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
	"google.golang.org/protobuf/proto"
)

type GRPCTestCase struct {
	Addr    string
	Timeout time.Duration
	Method  string
	Args    interface{}

	ExpectError error
	ExpectReply interface{}
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
	err = clientConn.Invoke(ctx, gtc.Method, gtc.Args, reply.Interface())
	assert.Equal(ginkgoT, gtc.ExpectError, err)

	replyData, err := proto.Marshal(reply.Interface().(proto.Message))
	assert.Nil(ginkgoT, err)

	expectReplyData, err := proto.Marshal(gtc.ExpectReply.(proto.Message))
	assert.Nil(ginkgoT, err)

	assert.Equal(ginkgoT, string(expectReplyData), string(replyData))

	assert.Nil(ginkgoT, clientConn.Close())
}
