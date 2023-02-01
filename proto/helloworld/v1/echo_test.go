package helloworldv1

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	v4 "github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
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
					echo.JSONSerializer = new(protoJsonSerializer)
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
			err := json.Compact(&body, []byte(res.Body.String()))
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

type protoJsonSerializer struct {
	v4.JSONSerializer
}

func (s *protoJsonSerializer) Serialize(c v4.Context, i interface{}, indent string) error {
	var (
		err  error
		data []byte
	)
	switch i.(type) {
	case proto.Message:
		data, err = protojson.MarshalOptions{
			EmitUnpopulated: true,
			UseEnumNumbers:  true,
		}.Marshal(i.(proto.Message))
	default:
		data, err = json.Marshal(i)
	}

	if err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, data)
}

func (s *protoJsonSerializer) Deserialize(c v4.Context, i interface{}) error {
	err := json.NewDecoder(c.Request().Body).Decode(i)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return v4.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
	} else if se, ok := err.(*json.SyntaxError); ok {
		return v4.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
	}
	return err
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
