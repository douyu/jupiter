package testproto

import (
	"net/http"

	"github.com/douyu/jupiter/pkg/util/xhttp"
	"github.com/labstack/echo/v4"
)

type HTTPConverter struct {
	impl   GreeterServiceServer
	binder *xhttp.ProtoBinder
}

func NewHTTPConverter(srv GreeterServiceServer) *HTTPConverter {
	return &HTTPConverter{
		impl:   srv,
		binder: new(xhttp.ProtoBinder),
	}
}

func (ins *HTTPConverter) SayHelloPath() string {
	return "/sayhello"
}

func (ins *HTTPConverter) SayHello() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(SayHelloRequest)
		err := ins.binder.Bind(req, c)
		if err != nil {
			return xhttp.ProtoError(c, http.StatusBadRequest, xhttp.ErrBadRequest)
		}

		res, err := ins.impl.SayHello(c.Request().Context(), req)
		if err != nil {
			return xhttp.ProtoError(c, http.StatusBadRequest, err)
		}

		return xhttp.ProtoJSON(c, http.StatusOK, res)
	}
}
