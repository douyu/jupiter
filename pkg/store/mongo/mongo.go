package mongo

import (
	"time"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/globalsign/mgo"
)

/*
DB: 返回name定义的mysql DB handler
name: 唯一名称
opts: Open Option, 用于覆盖配置文件中定义的配置
example: DB := DB("StdConfig", orm.RawConfig("jupiter.mongodb.StdConfig"))
*/
func newSession(config Config) *mgo.Session {
	session, err := mgo.Dial(config.DSN)
	if err != nil {
		_logger.Panic("dial mongo", xlog.FieldAddr(config.DSN), xlog.Any("error", err))
	}

	if config.SocketTimeout == time.Duration(0) {
		_logger.Panic("invalid config", xlog.FieldExtMessage("socketTimeout"))
	}

	if config.PoolLimit == 0 {
		_logger.Panic("invalid config", xlog.FieldExtMessage("poolLimit"))
	}

	session.SetSocketTimeout(config.SocketTimeout)
	session.SetPoolLimit(config.PoolLimit)

	return session
}
