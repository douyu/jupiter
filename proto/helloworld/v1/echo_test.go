package helloworldv1

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/douyu/jupiter/pkg/core/encoding"
	v4 "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestGreeterService_SayHello_0(t *testing.T) {
	type fields struct {
		server       GreeterServiceEchoServer
		createRouter func() *v4.Echo
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
				createRouter: func() *v4.Echo {
					echo := v4.New()
					echo.HTTPErrorHandler = defaultErrorHandler
					echo.JSONSerializer = new(encoding.ProtoJsonSerializer)
					return echo
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
			wantRes:    `{"error":0,"msg":"","data":{"name":"bob","ageNumber":"0","sex":0,"metadata":{}}}`,
			wantHeader: http.Header{"Content-Type": []string{"application/json; charset=UTF-8"}},
		},
		{
			name: "case 2: post with json",
			fields: fields{
				server: new(FooServer),
				createRouter: func() *v4.Echo {
					echo := v4.New()
					echo.HTTPErrorHandler = defaultErrorHandler
					return echo
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
			wantHeader: http.Header{"Content-Type": []string{"application/json; charset=UTF-8"}},
		},
		{
			name: "case 3: invalid content type",
			fields: fields{
				server: new(FooServer),
				createRouter: func() *v4.Echo {
					echo := v4.New()
					echo.HTTPErrorHandler = defaultErrorHandler
					return echo
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
			wantRes:    "{\"error\":500,\"msg\":\"code=400, message=Syntax error: offset=2, error=invalid character 'a' in literal null (expecting 'u'), internal=invalid character 'a' in literal null (expecting 'u')\",\"data\":{}}",
			wantHeader: http.Header{"Content-Type": []string{"application/json; charset=UTF-8"}},
		},
		{
			name: "case 4: not found",
			fields: fields{
				server: new(FooServer),
				createRouter: func() *v4.Echo {
					echo := v4.New()
					echo.HTTPErrorHandler = defaultErrorHandler
					return echo
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
			wantRes:    "{\"error\":500,\"msg\":\"code=404, message=Not Found\",\"data\":{}}",
			wantHeader: http.Header{"Content-Type": []string{"application/json; charset=UTF-8"}},
		},
		{
			name: "case 5: bind param for GET",
			fields: fields{
				server: new(FooServer),
				createRouter: func() *v4.Echo {
					echo := v4.New()
					echo.HTTPErrorHandler = defaultErrorHandler
					return echo
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
			RegisterGreeterServiceEchoServer(router, tt.fields.server)

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

var defaultErrorHandler = func(err error, c v4.Context) {
	c.JSON(http.StatusOK, struct {
		Error int      `json:"error"`
		Msg   string   `json:"msg"`
		Data  struct{} `json:"data"`
	}{
		Error: http.StatusInternalServerError,
		Msg:   err.Error(),
		Data:  struct{}{},
	})
}
