package redisgo

import (
	"context"

	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/singleton"

	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/go-redis/redis/v8"
)

func (config *Config) BuildStub() *redis.Client {
	stubClient := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Username:     config.Username,
		Password:     config.Password,
		DB:           config.DB,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  config.IdleTimeout,
	})
	stubClient.AddHook(fixedInterceptor(config.name, config, config.logger))
	if config.EnableMetricInterceptor {
		stubClient.AddHook(metricInterceptor(config.name, config, config.logger))
	}
	if config.Debug {
		stubClient.AddHook(debugInterceptor(config.name, config, config.logger))
	}
	if config.EnableTraceInterceptor {
		stubClient.AddHook(traceInterceptor(config.name, config, config.logger))
	}
	if config.EnableAccessLogInterceptor {
		stubClient.AddHook(accessInterceptor(config.name, config, config.logger))
	}

	if err := stubClient.Ping(context.Background()).Err(); err != nil {
		if config.OnDialError == "panic" {
			config.logger.Panic("redis stub start err", xlog.FieldErr(err))
		}
		config.logger.Error("redis stub start err", xlog.FieldErr(err))
	}

	instances.Store(config.name, &storeRedis{
		ClientStub: stubClient,
	})
	return stubClient
}

func (config *Config) BuildCluster() *redis.ClusterClient {
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        config.Addrs,
		Password:     config.Password,
		MaxRetries:   config.MaxRetries,
		ReadOnly:     config.ReadOnly,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  config.IdleTimeout,
	})
	clusterClient.AddHook(fixedInterceptor(config.name, config, config.logger))
	if config.EnableMetricInterceptor {
		clusterClient.AddHook(metricInterceptor(config.name, config, config.logger))
	}
	if config.Debug {
		clusterClient.AddHook(debugInterceptor(config.name, config, config.logger))
	}
	if config.EnableTraceInterceptor {
		clusterClient.AddHook(traceInterceptor(config.name, config, config.logger))
	}
	if config.EnableAccessLogInterceptor {
		clusterClient.AddHook(accessInterceptor(config.name, config, config.logger))
	}

	if err := clusterClient.Ping(context.Background()).Err(); err != nil {
		if config.OnDialError == "panic" {
			config.logger.Panic("redis cluster client start err", xlog.FieldErr(err))
		}
		config.logger.Error("redis cluster client start err", xlog.FieldErr(err))
	}
	instances.Store(config.name, &storeRedis{
		ClientCluster: clusterClient,
	})
	return clusterClient
}

func (config *Config) StubSingleton() *redis.Client {
	if val, ok := singleton.Load(constant.ModuleClientRedisStub, config.name); ok && val != nil {
		return val.(*redis.Client)
	}

	cc := config.BuildStub()
	singleton.Store(constant.ModuleClientGrpc, config.name, cc)
	return cc
}

func (config *Config) ClusterSingleton() *redis.ClusterClient {
	if val, ok := singleton.Load(constant.ModuleClientRedisCluster, config.name); ok && val != nil {
		return val.(*redis.ClusterClient)
	}

	cc := config.BuildCluster()
	singleton.Store(constant.ModuleClientGrpc, config.name, cc)
	return cc
}
