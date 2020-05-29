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

package xecho

import (
	"context"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/examples/helloworld/helloworld"
)

type GreeterMock struct{}

var (
	requestJSON = `
	{"name":"test"}
`
)

func (mock GreeterMock) SayHello(context context.Context, request *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{
		Message: "hello",
	}, nil
}
func TestGRPCProxyWrapper(t *testing.T) {
	e := echo.New()
	mock := GreeterMock{}
	req := httptest.NewRequest(echo.POST, "/", strings.NewReader(requestJSON))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	wrapper := GRPCProxyWrapper(mock.SayHello)
	err := wrapper(c)
	assert.Nil(t, err)
	assert.Equal(t, rec.Result().StatusCode, 200)
}
