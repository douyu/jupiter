package xhttp

import (
	"net/http"
	"strings"

	"github.com/douyu/jupiter/pkg/util/xerror"
	"github.com/labstack/echo/v4"
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
	acceptEncoding := c.Request().Header.Get(HeaderAcceptEncoding)

	var msg proto.Message
	switch obj := i.(type) {
	case proto.Message:
		msg = obj
	case error:
		c.Response().Header().Set(HeaderHRPCErr, "true")
		c.Response().Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)

		return c.JSON(http.StatusOK, xerror.Convert(obj))
	}

	// protobuf output
	if strings.Contains(acceptEncoding, MIMEApplicationProtobuf) {
		c.Response().Header().Set(HeaderContentType, MIMEApplicationProtobuf)
		c.Response().WriteHeader(code)
		bs, _ := proto.Marshal(msg)
		_, err := c.Response().Write(bs)
		return err
	}
	// json output
	c.Response().Header().Set(HeaderContentType, MIMEApplicationJSONCharsetUTF8)
	c.Response().WriteHeader(code)

	body, err := pbjson.Marshal(msg)
	if err != nil {
		return err
	}

	_, err = c.Response().Write(body)
	return err
}
