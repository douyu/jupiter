package redis

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_Cluster(t *testing.T) {
	config := DefaultClusterOption()
	t.Run("should panic when addr nil", func(t *testing.T) {
		var client *ClusterClient
		defer func() {
			if r := recover(); r != nil {
				assert.Equal(t, "cluster redis addr is empty", r.(string))
				assert.Nil(t, client)
			}
		}()
		client = config.MustClusterSingleton()
		assert.NotNil(t, client)
	})
	t.Run("should not panic when dial err", func(t *testing.T) {
		config.Addr = []string{"1.1.1.1"}
		config.OnDialError = "error"
		client, err := config.BuildCluster()
		assert.NotNil(t, err)
		assert.Nil(t, client)

	})
	t.Run("normal start", func(t *testing.T) {
		config.Addr = []string{"localhost:7111"}
		config.name = "test"
		client, err := config.BuildCluster()
		if err != nil {
			t.Fatalf("Failed to build cluster: %v", err)
		}
		if client == nil {
			t.Fatal("Client is nil")
		}
		err = client.Ping(context.Background()).Err()
		assert.Nil(t, err)
	})

}
