package redis

import (
	"strings"

	xredis "github.com/douyu/jupiter/pkg/client/redis"
	"github.com/google/wire"
)

var (
	ProviderSet = wire.NewSet(
		NewRedis,
	)

	redisNodeIdKey = "jupiter.uuid.node"
)

type Redis struct {
	*xredis.Redis
}

func NewRedis() *Redis {
	return &Redis{
		xredis.StdRedisConfig("uuid").Build(),
	}
}

// todo 采用 redis 的原子操作，这里目前是为了先实现功能
func (r *Redis) GetNodeId() (int64, error) {
	nodeID, err := r.Redis.Client.Get(redisNodeIdKey).Int64()
	if err != nil && !strings.Contains(err.Error(), "redis: nil") {
		return 0, err
	}

	nodeID++

	r.Redis.Client.Set(redisNodeIdKey, nodeID, 0)

	return nodeID, nil
}
