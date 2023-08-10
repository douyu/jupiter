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
				assert.Equal(t, "no cluster for default", r.(string))
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
		config.Addr = []string{"127.0.0.1:7000", "127.0.0.1:7001", "127.0.0.1:7002"}
		client, err := config.BuildCluster()
		assert.Nil(t, err)
		assert.NotNil(t, client)
		err = client.Ping(context.Background()).Err()
		if err != nil {
			t.Errorf("Test_Cluster ping err %v", err)
		}
	})

}
