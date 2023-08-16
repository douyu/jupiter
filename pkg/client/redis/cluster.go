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
	redis.ClusterClient
}

// ClusterSingleton returns a singleton client conn.
func (config *ClusterOptions) ClusterSingleton() (*ClusterClient, error) {
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
func (config *ClusterOptions) MustClusterSingleton() *ClusterClient {
	return lo.Must(config.ClusterSingleton())
}

// MustClusterBuild panics when error found.
func (config *ClusterOptions) MustClusterBuild() *ClusterClient {
	return lo.Must(config.BuildCluster())
}

// BuildCluster ..
func (config *ClusterOptions) BuildCluster() (*ClusterClient, error) {
	if xdebug.IsDevelopmentMode() {
		xdebug.PrettyJsonPrint("redis's config: "+config.name, config)
	}

	if len(config.Addr) <= 0 {
		return nil, errors.New("cluster redis addr is empty")
	}

	ins, err := config.buildCluster()
	if err != nil {
		return nil, err
	}
	return ins, nil
}

func (config *ClusterOptions) buildCluster() (*ClusterClient, error) {
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

	cfg := &config.Config
	for _, addr := range config.Addr {
		clusterClient.AddHook(fixedInterceptor(config.name, addr, cfg, config.logger))
		if config.EnableMetricInterceptor {
			clusterClient.AddHook(metricInterceptor(config.name, addr, cfg, config.logger))
		}
		if config.Debug {
			clusterClient.AddHook(debugInterceptor(config.name, addr, cfg, config.logger))
		}
		if config.EnableTraceInterceptor {
			clusterClient.AddHook(traceInterceptor(config.name, addr, cfg, config.logger))
		}
		if config.EnableAccessLogInterceptor {
			clusterClient.AddHook(accessInterceptor(config.name, addr, cfg, config.logger))
		}

		if config.EnableSentinel {
			clusterClient.AddHook(sentinelInterceptor(config.name, addr, cfg, config.logger))
		}
	}

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
	return &ClusterClient{ClusterClient: *clusterClient}, nil
}
