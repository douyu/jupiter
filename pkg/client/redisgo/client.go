package redisgo

import (
	"context"
	"math/rand"

	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/singleton"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/go-redis/redis/v8"
)

type instance struct {
	master *redis.Client
	slave  []*redis.Client
	config *Config
}

func (ins *instance) CmdOnMaster() *redis.Client {
	if ins.master == nil {
		ins.config.logger.Panic("redisgo:no master for "+ins.config.name, xlog.FieldExtMessage(ins.config))
	}
	return ins.master
}
func (ins *instance) CmdOnSlave() *redis.Client {
	if len(ins.slave) == 0 {
		ins.config.logger.Panic("redisgo:no slave for "+ins.config.name, xlog.FieldExtMessage(ins.config))
	}
	return ins.slave[rand.Intn(len(ins.slave))]
}

// Singleton 单例模式
func (config *Config) Singleton() *instance {
	if val, ok := singleton.Load(constant.ModuleClientRedis, config.name); ok && val != nil {
		return val.(*instance)
	}

	cc := config.Build()
	singleton.Store(constant.ModuleClientRedis, config.name, cc)
	return cc
}

// Build ..
func (config *Config) Build() *instance {
	ins := new(instance)
	if config.Master.Addr != "" {
		addr, user, pass := getUsernameAndPassword(config.Master.Addr)
		ins.master = config.build(addr, user, pass)
	}
	if len(config.Slaves.Addr) > 0 {
		ins.slave = []*redis.Client{}
		for _, slave := range config.Slaves.Addr {
			addr, user, pass := getUsernameAndPassword(slave)
			ins.slave = append(ins.slave, config.build(addr, user, pass))
		}
	}

	if ins.master == nil && len(ins.slave) == 0 {
		config.logger.Panic("redisgo:no master or slaves for "+config.name, xlog.FieldExtMessage(config))
	}
	return ins
}

func (config *Config) build(addr, user, pass string) *redis.Client {

	stubClient := redis.NewClient(&redis.Options{
		Addr:         addr,
		Username:     user,
		Password:     pass,
		DB:           config.DB,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  config.IdleTimeout,
	})
	stubClient.AddHook(fixedInterceptor(config.name, addr, config, config.logger))
	if config.EnableMetricInterceptor {
		stubClient.AddHook(metricInterceptor(config.name, addr, config, config.logger))
	}
	if config.Debug {
		stubClient.AddHook(debugInterceptor(config.name, addr, config, config.logger))
	}
	if config.EnableTraceInterceptor {
		stubClient.AddHook(traceInterceptor(config.name, addr, config, config.logger))
	}
	if config.EnableAccessLogInterceptor {
		stubClient.AddHook(accessInterceptor(config.name, addr, config, config.logger))
	}

	if err := stubClient.Ping(context.Background()).Err(); err != nil {
		if config.OnDialError == "panic" {
			config.logger.Panic("redisgo stub client start err: " + err.Error())
		}
		config.logger.Error("redisgo stub client start err", xlog.FieldErr(err))
	}

	instances.Store(config.name, &storeRedis{
		ClientStub: stubClient,
	})
	return stubClient
}
