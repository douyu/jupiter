package etcdv3

import (
	"context"
	"testing"
	"time"

	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/stretchr/testify/assert"
)

func Test_etcdv3Registry(t *testing.T) {
	etcdConfig := etcdv3.DefaultConfig()
	etcdConfig.Endpoints = []string{"127.0.0.1:2379"}
	registry := newETCDRegistry(&Config{
		Config:      etcdConfig,
		ReadTimeout: time.Second * 10,
		Prefix:      "jupiter",
		logger: xlog.DefaultLogger,
	})

	assert.Nil(t, registry.RegisterService(context.Background(), &server.ServiceInfo{
		Name:       "service_1",
		Scheme:     "grpc",
		Address:         "10.10.10.1:9091",
		Weight:     40,
		Enable:     true,
		Healthy:    true,
		Metadata: map[string]string{},
		Region:     "default",
		Zone:       "default",
		Deployment: "default",
	}))

	services, err := registry.ListServices(context.Background(), "service_1", "grpc")
	assert.Nil(t, err)
	t.Logf("services: %+v\n", services[0])
	assert.Equal(t, 1, len(services))
	assert.Equal(t, "10.10.10.1:9091", services[0].Address)

	go func() {
		si := &server.ServiceInfo{
			Name:       "service_1",
			Scheme:     "grpc",
			Weight:     40,
			Address:         "10.10.10.1:9092",
			Enable:     true,
			Healthy:    true,
			Metadata: map[string]string{},
			Region:     "default",
			Zone:       "default",
			Deployment: "default",
		}
		time.Sleep(time.Second)
		assert.Nil(t, registry.RegisterService(context.Background(), si))
		assert.Nil(t, registry.UnregisterService(context.Background(), si))
	}()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		services, eventChan, err := registry.WatchServices(ctx, "service_1", "grpc")
		assert.Nil(t, err)
		assert.Equal(t, 1, len(services))
		for msg := range eventChan {
			t.Logf("watch service: %+v\n", msg)
			assert.Equal(t, "10.10.10.2:9092", msg)
		}
	}()

	time.Sleep(time.Second * 3)
	cancel()
	_ = registry.Close()
	time.Sleep(time.Second * 1)
}