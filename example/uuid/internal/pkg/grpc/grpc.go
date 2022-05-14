package grpc

import (
	"context"

	grpcclient "github.com/douyu/jupiter/pkg/client/grpc"
	"github.com/google/wire"
	"google.golang.org/grpc"
)

// 本文件负责grpc下面各种ProviderSet的注册
var (
	ProviderSet = wire.NewSet(
		NewUuid,
	)
)

type Uuid struct {
	cc grpc.ClientConnInterface
}

func NewUuid() UuidInterface {
	return &Uuid{
		cc: grpcclient.StdConfig("uuid").Build(),
	}
}

func (s *Uuid) GetUuidBySnowflake(ctx context.Context) error {
	return nil
}

func (s *Uuid) GetUuidByGoogleUUIDV4(ctx context.Context) error {
	return nil
}
