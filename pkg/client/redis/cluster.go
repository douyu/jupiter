package redis

import (
	"context"
	"errors"
	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/core/singleton"
	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/go-redis/redis/v8"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

type ClusterClient struct {
	cluster *redis.ClusterClient
	config  *Config
}

func (ins *ClusterClient) CmdOnCluster() *redis.ClusterClient {
	if ins.cluster == nil {
		ins.config.logger.Panic("redis:no cluster for "+ins.config.name, xlog.FieldExtMessage(ins.config))
	}
	return ins.cluster
}

// ClusterSingleton returns a singleton client conn.
func (config *Config) ClusterSingleton() (*ClusterClient, error) {
	if val, ok := singleton.Load(constant.ModuleClusterRedis, config.name); ok && val != nil {
		return val.(*ClusterClient), nil
	}

	cc, err := config.BuildCluster()
	if err != nil {
		xlog.Jupiter().Error("build redis cluster client failed", zap.Error(err))
		return nil, err
	}
	singleton.Store(constant.ModuleClusterRedis, config.name, cc)
	return cc, nil
}

// MustClusterSingleton panics when error found.
func (config *Config) MustClusterSingleton() *ClusterClient {
	return lo.Must(config.ClusterSingleton())
}

// MustClusterBuild panics when error found.
func (config *Config) MustClusterBuild() *ClusterClient {
	return lo.Must(config.BuildCluster())
}

// BuildCluster ..
func (config *Config) BuildCluster() (*ClusterClient, error) {
	ins := new(ClusterClient)
	var err error
	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint("redis's config: "+config.name, config)
	}

	if len(config.Addr) > 0 {
		ins.cluster, err = config.buildCluster()
		if err != nil {
			return ins, err
		}
		return ins, nil
	}

	if ins.cluster == nil {
		return ins, errors.New("no cluster for " + config.name)
	}
	return ins, nil
}

func (config *Config) buildCluster() (*redis.ClusterClient, error) {
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        config.Addr,
		Username:     config.Username,
		Password:     config.Password,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		IdleTimeout:  config.IdleTimeout,
	})

	for _, addr := range config.Addr {
		clusterClient.AddHook(fixedInterceptor(config.name, addr, config, config.logger))
		if config.EnableMetricInterceptor {
			clusterClient.AddHook(metricInterceptor(config.name, addr, config, config.logger))
		}
		if config.Debug {
			clusterClient.AddHook(debugInterceptor(config.name, addr, config, config.logger))
		}
		if config.EnableTraceInterceptor {
			clusterClient.AddHook(traceInterceptor(config.name, addr, config, config.logger))
		}
		if config.EnableAccessLogInterceptor {
			clusterClient.AddHook(accessInterceptor(config.name, addr, config, config.logger))
		}

		if config.EnableSentinel {
			clusterClient.AddHook(sentinelInterceptor(config.name, addr, config, config.logger))
		}
	}

	clusterClient.Ping(context.Background())
	if err := clusterClient.Ping(context.Background()).Err(); err != nil {
		if config.OnDialError == "panic" {
			config.logger.Panic("redis cluster client start err: " + err.Error())
		}
		config.logger.Error("redis cluster client start err", xlog.FieldErr(err))
		return nil, err
	}

	instances.Store(config.name, &storeRedis{
		ClientCluster: clusterClient,
	})
	return clusterClient, nil
}
