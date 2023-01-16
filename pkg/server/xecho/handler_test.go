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
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/douyu/jupiter/proto/testproto/v1"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

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
			wantBody: "{\"error\":0,\"msg\":\"\",\"data\":{\"name\":\"hello bob\",\"ageNumber\":\"0\"}}",
		},
		{
			name: "case 2: get with query",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/?name=bob", nil),
			},
			wantErr:  nil,
			wantBody: "{\"error\":0,\"msg\":\"\",\"data\":{\"name\":\"hello bob\",\"ageNumber\":\"0\"}}",
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
			wantBody: "{\"error\":0,\"msg\":\"\",\"data\":{\"name\":\"hello bob\",\"ageNumber\":\"0\"}}",
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
			wantBody: "{\"error\":3,\"msg\":\"invalid name\",\"data\":{\"name\":\"\",\"ageNumber\":\"0\"}}",
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
			wantBody: "{\"error\":3,\"msg\":\"invalid name\",\"data\":{\"name\":\"\",\"ageNumber\":\"0\"}}",
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

			impl := new(impl)
			assert.Equal(t, tt.wantErr, GRPCProxyWrapper(impl.SayHello)(c))
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

func TestHTTPConverter(t *testing.T) {
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
			wantBody: "{\"error\":0,\"msg\":\"\",\"data\":{\"name\":\"hello bob\",\"ageNumber\":\"0\"}}",
		},
		{
			name: "case 2: get with query",
			args: args{
				req: httptest.NewRequest(http.MethodGet, "/?name=bob", nil),
			},
			wantErr:  nil,
			wantBody: "{\"error\":0,\"msg\":\"\",\"data\":{\"name\":\"hello bob\",\"ageNumber\":\"0\"}}",
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
			wantBody: "{\"error\":0,\"msg\":\"\",\"data\":{\"name\":\"hello bob\",\"ageNumber\":\"0\"}}",
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
			wantBody: "{\"error\":3,\"msg\":\"invalid name\",\"data\":{\"name\":\"\",\"ageNumber\":\"0\"}}",
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
			wantBody: "{\"error\":3,\"msg\":\"invalid name\",\"data\":{\"name\":\"\",\"ageNumber\":\"0\"}}",
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

			impl := new(impl)
			httpConverter := testproto.NewHTTPConverter(impl)
			assert.Equal(t, tt.wantErr, httpConverter.SayHello()(c))
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

func BenchmarkHTTP(b *testing.B) {

	b.Run("HTTP with reflect", func(b *testing.B) {
		server := echo.New()
		impl := new(impl)

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{\"name\":\"bob\"}"))
			req.Header.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			c := server.NewContext(req, rec)

			GRPCProxyWrapper(impl.SayHello)(c)
		}
	})

	b.Run("HTTP without reflect", func(b *testing.B) {
		server := echo.New()
		httpConvert := testproto.NewHTTPConverter(new(impl))
		sayHello := httpConvert.SayHello()

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{\"name\":\"bob\"}"))
			req.Header.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			c := server.NewContext(req, rec)

			sayHello(c)
		}
	})

	b.Run("HTTP with protojson", func(b *testing.B) {
		server := echo.New()

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{\"name\":\"bob\"}"))
			req.Header.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			c := server.NewContext(req, rec)

			echoHandler(c)
		}
	})

	b.Run("HTTP with json", func(b *testing.B) {
		server := echo.New()

		for i := 0; i < b.N; i++ {
			req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBufferString("{\"name\":\"bob\"}"))
			req.Header.Add("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			c := server.NewContext(req, rec)

			echoJsonHandler(c)
		}
	})
}
