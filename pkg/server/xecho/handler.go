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
	"bytes"
	"net/http"
	"reflect"
	"strings"

	"github.com/codegangsta/inject"
	"github.com/douyu/jupiter/pkg/util/xerror"
	"github.com/labstack/echo/v4"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// pbjson is a protojson.MarshalOptions with some default options.
var pbjson = protojson.MarshalOptions{
	Multiline:       false,
	UseEnumNumbers:  true,
	EmitUnpopulated: true,
}

// ProtoError ...
func ProtoError(c echo.Context, code int, e error) error {
	return ProtoJSON(c, code, xerror.Convert(e))
}

// ProtoJSON sends a Protobuf JSON response with status code and data.
func ProtoJSON(c echo.Context, code int, i interface{}) error {
	var acceptEncoding = c.Request().Header.Get(HeaderAcceptEncoding)

	var m proto.Message
	switch msg := i.(type) {
	case proto.Message:
		m = msg
	case error:
		c.Response().Header().Set(HeaderHRPCErr, "true")
		c.Response().Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)

		return c.JSON(http.StatusOK, xerror.Convert(i.(error)))
	}

	// protobuf output
	if strings.Contains(acceptEncoding, MIMEApplicationProtobuf) {
		c.Response().Header().Set(HeaderContentType, MIMEApplicationProtobuf)
		c.Response().WriteHeader(code)
		bs, _ := proto.Marshal(m)
		_, err := c.Response().Write(bs)
		return err
	}
	// json output
	c.Response().Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(code)

	body, err := pbjson.Marshal(m)
	if err != nil {
		return err
	}

	_, err = c.Response().Write(body)
	return err
}

var grpcBinder = &ProtoBinder{}

// GRPCProxyWrapper ...
func GRPCProxyWrapper(h interface{}) echo.HandlerFunc {
	t := reflect.TypeOf(h)
	if t.Kind() != reflect.Func {
		panic("reflect error: handler must be func")
	}

	return func(c echo.Context) error {

		var req = reflect.New(t.In(1).Elem()).Interface()
		if err := grpcBinder.Bind(req, c); err != nil {
			return ProtoError(c, http.StatusBadRequest, errBadRequest)
		}

		var md = metadata.MD{}
		for k, vs := range c.Request().Header {
			for _, v := range vs {
				bs := bytes.TrimFunc([]byte(v), func(r rune) bool {
					return r == '\n' || r == '\r' || r == '\000'
				})
				md.Append(k, string(bs))
			}
		}

		ctx := metadata.NewIncomingContext(c.Request().Context(), md)
		var inj = inject.New()
		inj.Map(ctx)
		inj.Map(req)
		vs, err := inj.Invoke(h)
		if err != nil {
			return ProtoError(c, http.StatusInternalServerError, errMicroInvoke)
		}
		if len(vs) != 2 {
			return ProtoError(c, http.StatusInternalServerError, errMicroInvokeLen)
		}
		repV, errV := vs[0], vs[1]
		if !errV.IsNil() || repV.IsNil() {
			if e, ok := errV.Interface().(error); ok {
				return ProtoError(c, http.StatusOK, e)
			}
			return ProtoError(c, http.StatusInternalServerError, errMicroInvokeInvalid)
		}
		if !repV.IsValid() {
			return ProtoError(c, http.StatusInternalServerError, errMicroResInvalid)
		}
		return ProtoJSON(c, http.StatusOK, repV.Interface())
	}
}
