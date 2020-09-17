package mongox

import (
	"context"
	"time"

	"github.com/douyu/jupiter/pkg/xlog"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

/*
DB: 返回name定义的mysql DB handler
name: 唯一名称
opts: Open Option, 用于覆盖配置文件中定义的配置
example: DB := DB("StdConfig", orm.RawConfig("jupiter.mongodb.StdConfig"))
*/

func newSession(config Config) *mongo.Client {

	// check config param
	isConfigErr(config)

	mps := uint64(config.PoolLimit)

	clientOpts := options.Client()
	clientOpts.MaxPoolSize = &mps
	clientOpts.SocketTimeout = &config.SocketTimeout
	client, err := mongo.Connect(context.Background(), clientOpts.ApplyURI(config.DSN))
	if err != nil {
		_logger.Panic("dial mongo", xlog.FieldAddr(config.DSN), xlog.Any("error", err))
	}

	_instances.Store(config.Name, client)
	return client
}

func isConfigErr(config Config) {
	if config.SocketTimeout == time.Duration(0) {
		_logger.Panic("invalid config", xlog.FieldExtMessage("socketTimeout"))
	}

	if config.PoolLimit == 0 {
		_logger.Panic("invalid config", xlog.FieldExtMessage("poolLimit"))
	}
}
