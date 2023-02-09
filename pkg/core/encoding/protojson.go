package encoding

import (
	"encoding/json"
	"fmt"
	"net/http"

	v4 "github.com/labstack/echo/v4"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

type ProtoJsonSerializer struct {
	v4.JSONSerializer
}

func (s *ProtoJsonSerializer) Serialize(c v4.Context, i interface{}, indent string) error {
	var (
		err  error
		data []byte
	)
	switch obj := i.(type) {
	case proto.Message:
		data, err = protojson.MarshalOptions{
			EmitUnpopulated: true,
			UseEnumNumbers:  true,
		}.Marshal(obj)
	default:
		data, err = json.Marshal(obj)
	}

	if err != nil {
		return err
	}

	return c.JSONBlob(http.StatusOK, data)
}

func (s *ProtoJsonSerializer) Deserialize(c v4.Context, i interface{}) error {
	err := json.NewDecoder(c.Request().Body).Decode(i)
	if ute, ok := err.(*json.UnmarshalTypeError); ok {
		return v4.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Unmarshal type error: expected=%v, got=%v, field=%v, offset=%v", ute.Type, ute.Value, ute.Field, ute.Offset)).SetInternal(err)
	} else if se, ok := err.(*json.SyntaxError); ok {
		return v4.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("Syntax error: offset=%v, error=%v", se.Offset, se.Error())).SetInternal(err)
	}
	return err
}
