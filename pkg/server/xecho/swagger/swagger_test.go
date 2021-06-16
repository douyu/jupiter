package swagger

import (
	"testing"

	"github.com/labstack/echo/v4"
)

func TestRoute(t *testing.T) {
	r := &echo.Route{}
	Route(r).
		Doc("document").
		PathParam("name", "user's name").
		

}
