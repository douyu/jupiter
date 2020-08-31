package mongo

import (
	"time"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/globalsign/mgo"
)

// StdNew ...
func StdNew(name string, opts ...interface{}) *mgo.Session {
	return New(name, StdConfig(name))
}

/*
DB: 返回name定义的mysql DB handler
name: 唯一名称
opts: Open Option, 用于覆盖配置文件中定义的配置
example: DB := DB("StdConfig", orm.RawConfig("jupiter.mongodb.StdConfig"))
*/
func New(name string, config Config) *mgo.Session {
	if _, ok := _instances.Load(name); ok {
		_logger.Panic("duplicated new", xlog.FieldName(name), xlog.FieldExtMessage(config))
	}

	session, err := mgo.Dial(config.DSN)
	if err != nil {
		_logger.Panic("dial mongo", xlog.FieldName(name), xlog.FieldAddr(config.DSN), xlog.Any("error", err))
	}

	if config.SocketTimeout == time.Duration(0) {
		_logger.Panic("invalid config", xlog.FieldName(name), xlog.FieldExtMessage("socketTimeout"))
	}

	if config.PoolLimit == 0 {
		_logger.Panic("invalid config", xlog.FieldName(name), xlog.FieldExtMessage("poolLimit"))
	}

	session.SetSocketTimeout(config.SocketTimeout)
	session.SetPoolLimit(config.PoolLimit)

	_instances.Store(name, session)
	return session
}
