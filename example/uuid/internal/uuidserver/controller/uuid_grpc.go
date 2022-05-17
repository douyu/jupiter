package controller

import (
	"context"

	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
	uuidv1 "uuid/gen/api/go/uuid/v1"
	"uuid/internal/uuidserver/service"
)

type UuidGrpc struct {
	uuid *service.Uuid
}

func NewUUuidGrpcController(uuid *service.Uuid) *UuidGrpc {
	return &UuidGrpc{
		uuid: uuid,
	}
}

func (u *UuidGrpc) GetUuidBySnowflake(ctx context.Context, req *uuidv1.GetUuidBySnowflakeRequest) (*uuidv1.GetUuidBySnowflakeRequestResponse, error) {
	res, err := u.uuid.GetUuidBySnowflake(ctx, req)
	if err != nil {
		xlog.Error("getUuidBySnowflake failed", zap.Error(err), zap.Any("res", res), zap.Any("req", req))
		return nil, err
	}

	return res, nil
}

func (u *UuidGrpc) GetUuidByGoogleUUIDV4(ctx context.Context, req *uuidv1.GetUuidByGoogleUUIDV4Request) (*uuidv1.GetUuidByGoogleUUIDV4Response, error) {
	res, err := u.uuid.GetUuidByGoogleUUIDV4(ctx, req)
	if err != nil {
		xlog.Error("getUuidByGoogleUUIDV4 failed", zap.Error(err), zap.Any("res", res), zap.Any("req", req))
		return nil, err
	}

	return res, nil
}
