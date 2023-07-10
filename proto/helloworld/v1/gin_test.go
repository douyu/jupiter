package helloworldv1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGinGreeterService_SayHello_0(t *testing.T) {
	type fields struct {
		server       GreeterServiceGinServer
		createRouter func() *gin.Engine
	}
	type args struct {
		createReq func() *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantErr    bool
		wantRes    string
		wantHeader http.Header
	}{
		{
			name: "case 1: post with form",
			fields: fields{
				server: new(FooServer),
				createRouter: func() *gin.Engine {
					server := gin.Default()
					return server
				},
			},
			args: args{
				createReq: func() *http.Request {
					req := httptest.NewRequest(
						"POST", "http://localhost/v1/helloworld.Greeter/SayHello",
						bytes.NewBufferString("name=bob"))
					req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
					return req
				},
			},
			wantErr:    false,
			wantRes:    `{"error":0,"msg":"","data":{"name":"bob","ageNumber":0,"sex":0,"metadata":null}}`,
			wantHeader: http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
		},
		{
			name: "case 2: post with json",
			fields: fields{
				server: new(FooServer),
				createRouter: func() *gin.Engine {
					server := gin.Default()
					server.Use(errorHandler)
					return server
				},
			},
			args: args{
				createReq: func() *http.Request {
					req := httptest.NewRequest(
						"POST", "http://localhost/v1/helloworld.Greeter/SayHello",
						bytes.NewBufferString(`{"name":"needErr"}`))
					req.Header.Add("Content-Type", "application/json")
					return req
				},
			},
			wantErr:    false,
			wantRes:    "{\"error\":500,\"msg\":\"rpc error: code = DataLoss desc = error foo\",\"data\":{}}",
			wantHeader: http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
		},
		{
			name: "case 3: invalid content type",
			fields: fields{
				server: new(FooServer),
				createRouter: func() *gin.Engine {
					server := gin.Default()
					server.Use(errorHandler)
					return server
				},
			},
			args: args{
				createReq: func() *http.Request {
					req := httptest.NewRequest(
						"POST", "http://localhost/v1/helloworld.Greeter/SayHello",
						bytes.NewBufferString("name=bob"))
					req.Header.Add("Content-Type", "application/json")
					return req
				},
			},
			wantErr:    false,
			wantRes:    "{\"error\":500,\"msg\":\"invalid character 'a' in literal null (expecting 'u')\",\"data\":{}}",
			wantHeader: http.Header{"Content-Type": []string{"application/json; charset=utf-8"}},
		},
		{
			name: "case 4: not found",
			fields: fields{
				server: new(FooServer),
				createRouter: func() *gin.Engine {
					server := gin.Default()
					server.Use(errorHandler)
					return server
				},
			},
			args: args{
				createReq: func() *http.Request {
					req := httptest.NewRequest(
						"POST", "http://localhost/v1/helloworld.Greeter/NotFound",
						bytes.NewBufferString("name=bob"))
					req.Header.Add("Content-Type", "application/json")
					return req
				},
			},
			wantErr:    false,
			wantRes:    "404 page not found",
			wantHeader: http.Header{"Content-Type": []string{"text/plain"}},
		},
		{
			name: "case 5: bind param for GET",
			fields: fields{
				server: new(FooServer),
				createRouter: func() *gin.Engine {
					server := gin.Default()
					return server
				},
			},
			args: args{
				createReq: func() *http.Request {
					req := httptest.NewRequest(
						"GET", "http://localhost/v1/helloworld.Greeter/SayHello/bob",
						nil)
					return req
				},
			},
			wantErr: false,
			wantRes: `{"error":0,"msg":"","data":{"name":"bob","ageNumber":0,"sex":0,"metadata":null}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := tt.fields.createRouter()
			RegisterGreeterServiceGinServer(router, tt.fields.server)

			res := httptest.NewRecorder()
			router.ServeHTTP(res, tt.args.createReq())

			var body bytes.Buffer
			err := json.Compact(&body, res.Body.Bytes())
			// 如果Compact失败，则说明不是json格式
			if err != nil {
				body.WriteString(res.Body.String())
			}
			assert.Equal(t, tt.wantRes, body.String())

			if len(tt.wantHeader) > 0 {
				assert.Equal(t, tt.wantHeader, res.Header())
			}
		})
	}
}

func errorHandler(c *gin.Context) {
	c.Next()

	for _, err := range c.Errors {
		switch err.Err {
		default:
			c.JSON(-1, struct {
				Error int      `json:"error"`
				Msg   string   `json:"msg"`
				Data  struct{} `json:"data"`
			}{
				Error: http.StatusInternalServerError,
				Msg:   err.Error(),
				Data:  struct{}{},
			})
		}
		return
	}
}
