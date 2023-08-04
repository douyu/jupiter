package redis

import (
	"context"
	"errors"
	"math/rand"

	"github.com/go-redis/redis/v8"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/core/singleton"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/xlog"
)

type Client struct {
	master *redis.Client
	slave  []*redis.Client
	config *Config
}

func (ins *Client) CmdOnMaster() *redis.Client {
	if ins.master == nil {
		ins.config.logger.Panic("redis:no master for "+ins.config.name, xlog.FieldExtMessage(ins.config))
	}
	return ins.master
}
func (ins *Client) CmdOnSlave() *redis.Client {
	if len(ins.slave) == 0 {
		ins.config.logger.Panic("redis:no slave for "+ins.config.name, xlog.FieldExtMessage(ins.config))
	}
	return ins.slave[rand.Intn(len(ins.slave))]
}

// Singleton returns a singleton client conn.
func (config *Config) Singleton() (*Client, error) {
	if val, ok := singleton.Load(constant.ModuleClientRedis, config.name); ok && val != nil {
		return val.(*Client), nil
	}

	cc, err := config.Build()
	if err != nil {
		xlog.Jupiter().Error("build redis client failed", zap.Error(err))
		return nil, err
	}
	singleton.Store(constant.ModuleClientRedis, config.name, cc)
	return cc, nil
}

// MustSingleton panics when error found.
func (config *Config) MustSingleton() *Client {
	return lo.Must(config.Singleton())
}

// MustBuild panics when error found.
func (config *Config) MustBuild() *Client {
	return lo.Must(config.Build())
}

// Build ..
func (config *Config) Build() (*Client, error) {
	ins := new(Client)
	var err error
	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint("redis's config: "+config.name, config)
	}
	if config.Master.Addr != "" {
		addr, user, pass := getUsernameAndPassword(config.Master.Addr)
		ins.master, err = config.build(addr, user, pass)
		if err != nil {
			return ins, err
		}
	}
	if len(config.Slaves.Addr) > 0 {
		ins.slave = []*redis.Client{}
		for _, slave := range config.Slaves.Addr {
			addr, user, pass := getUsernameAndPassword(slave)
			cli, err := config.build(addr, user, pass)
			if err != nil {
				return ins, err
			}
			ins.slave = append(ins.slave, cli)
		}
	}

	if ins.master == nil && len(ins.slave) == 0 {
		return ins, errors.New("no master or slaves for " + config.name)
	}
	return ins, nil
}

func (config *Config) build(addr, user, pass string) (*redis.Client, error) {

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

	if config.EnableSentinel {
		stubClient.AddHook(sentinelInterceptor(config.name, addr, config, config.logger))
	}

	if err := stubClient.Ping(context.Background()).Err(); err != nil {
		if config.OnDialError == "panic" {
			config.logger.Panic("redis stub client start err: " + err.Error())
		}
		config.logger.Error("redis stub client start err", xlog.FieldErr(err))
		return nil, err
	}

	instances.Store(config.name, &storeRedis{
		ClientStub: stubClient,
	})
	return stubClient, nil
}
