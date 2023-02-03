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

package xecho

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/douyu/jupiter/pkg/util/xerror"
	helloworldv1 "github.com/douyu/jupiter/proto/helloworld/v1"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func hander(ctx context.Context, req *helloworldv1.SayHelloRequest) (*helloworldv1.SayHelloResponse, error) {
	if req.Name != "bob" {
		return &helloworldv1.SayHelloResponse{
			Error: uint32(xerror.InvalidArgument.GetEcode()),
			Msg:   "invalid name",
			Data:  &helloworldv1.SayHelloResponse_Data{},
		}, nil
	}

	return &helloworldv1.SayHelloResponse{
		Msg: "",
		Data: &helloworldv1.SayHelloResponse_Data{
			Name: "hello bob",
		},
	}, nil
}

func TestGRPCProxyWrapper(t *testing.T) {
	type args struct {
		req    *http.Request
		header map[string]string
	}
	tests := []struct {
		name       string
		args       args
		wantErr    error
		wantBody   string
		wantHeader http.Header
	}{
		{
			name: "case 1: post with json",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{\"name\":\"bob\"}")),
				header: map[string]string{
					"Content-Type": "application/json",
				},
			},
			wantErr: nil,
			wantHeader: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
			wantBody: `{"error":0,"msg":"","data":{"name":"hello bob","ageNumber":"0","sex":0,"metadata":{}}}`,
		},
		{
			name: "case 2: get with query",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/?name=bob", nil),
			},
			wantErr:  nil,
			wantBody: `{"error":0,"msg":"","data":{"name":"hello bob","ageNumber":"0","sex":0,"metadata":{}}}`,
		},
		{
			name: "case 3: post with form",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("name=bob")),
				header: map[string]string{
					"Content-Type": "application/x-www-form-urlencoded",
				},
			},
			wantErr: nil,
			wantHeader: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
			wantBody: `{"error":0,"msg":"","data":{"name":"hello bob","ageNumber":"0","sex":0,"metadata":{}}}`,
		},
		{
			name: "case 4: post with query",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/?name=bob", nil),
				header: map[string]string{
					"Content-Type": "application/x-www-form-urlencoded",
				},
			},
			wantErr: nil,
			wantHeader: http.Header{
				"Content-Type": []string{"application/json; charset=utf-8"},
			},
			wantBody: `{"error":3,"msg":"invalid name","data":{"name":"","ageNumber":"0","sex":0,"metadata":{}}}`,
		},
		{
			name: "case 5: json without content-type",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/?name=query", bytes.NewBufferString("{\"name\":\"json\"}")),
				header: map[string]string{
					"Content-Type": "application/x-www-form-urlencoded",
				},
			},
			wantErr:  nil,
			wantBody: `{"error":3,"msg":"invalid name","data":{"name":"","ageNumber":"0","sex":0,"metadata":{}}}`,
		},
		{
			name: "case 6: form without content-type",
			args: args{
				req: httptest.NewRequest(http.MethodPost, "/?name=query", bytes.NewBufferString("name=form")),
				header: map[string]string{
					"Content-Type": "application/json",
				},
			},
			wantErr:  nil,
			wantBody: "{\"error\":3,\"msg\":\"bad request\",\"data\":{}}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			for k, v := range tt.args.header {
				tt.args.req.Header.Add(k, v)
			}

			rec := httptest.NewRecorder()
			c := e.NewContext(tt.args.req, rec)

			assert.Equal(t, tt.wantErr, GRPCProxyWrapper(hander)(c))
			if tt.wantHeader != nil {
				assert.Equal(t, tt.wantHeader, rec.HeaderMap)
			}

			// protojson does not generate frozen json, so
			var rm json.RawMessage = rec.Body.Bytes()
			data2, err := json.Marshal(rm)

			assert.Nil(t, err)
			assert.Equal(t, tt.wantBody, string(data2))
		})
	}
}
