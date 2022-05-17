package controller

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/douyu/jupiter/pkg/util/xerror"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	uuidv1 "uuid/gen/api/go/uuid/v1"
	"uuid/internal/uuidserver/service"
)

type UuidHTTP struct {
	uuid *service.Uuid
}

func NewUuidHTTPController(uuid *service.Uuid) *UuidHTTP {
	return &UuidHTTP{
		uuid: uuid,
	}
}

func (s *UuidHTTP) GetUuidBySnowflake(c echo.Context) error {
	nodeId, err := strconv.Atoi(c.QueryParam("nodeId"))
	if err != nil {
		xlog.Error("getUuidBySnowflake failed", zap.Error(err))
		return c.JSON(http.StatusOK, fmt.Errorf("nodeId invalid err:%v", err))
	}

	req := &uuidv1.GetUuidBySnowflakeRequest{
		NodeId: int32(nodeId),
	}

	res, err := s.uuid.GetUuidBySnowflake(c.Request().Context(), req)
	if err != nil {
		xlog.Error("getUuidBySnowflake failed", zap.Error(err), zap.Any("res", res), zap.Any("req", req))
		return c.JSON(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, xerror.OK.WithData(res))
}

func (s *UuidHTTP) GetUuidByGoogleUUIDV4(c echo.Context) error {
	req := &uuidv1.GetUuidByGoogleUUIDV4Request{}

	res, err := s.uuid.GetUuidByGoogleUUIDV4(c.Request().Context(), req)
	if err != nil {
		xlog.Error("getUuidByGoogleUUIDV4 failed", zap.Error(err), zap.Any("res", res), zap.Any("req", req))
		return c.JSON(http.StatusOK, err)
	}

	return c.JSON(http.StatusOK, xerror.OK.WithData(res))
}
