package redigo

import (
	"git.dz11.com/vega/invokerx/logger"
	"github.com/go-redis/redis/v8"
)

// StdNewStub read standard config and return a stub client
func StdNewStub(name string) *redis.Client {
	return NewStub(StdStubConfig(name))
}

func NewStub(name string, config StubOptions) *redis.Client {
	if _, ok := stubInstances.Load(name); ok {
		logger.Panicd("duplicated stub", ilog.FieldName(name))
	}
}
