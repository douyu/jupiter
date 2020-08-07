package etcdv3

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_GetKeyValue(t *testing.T) {
	config := DefaultConfig()
	config.Endpoints = []string{"127.0.0.1:2379"}
	config.TTL = 5
	etcdCli := newClient(config)

	ctx := context.TODO()

	leaseSession, err := etcdCli.GetLeaseSession(ctx, concurrency.WithTTL(int(config.TTL)))
	assert.Nil(t, err)
	defer leaseSession.Close()

	_, err = etcdCli.Client.KV.Put(ctx, "/test/key", "{...}", clientv3.WithLease(leaseSession.Lease()))
	assert.Nil(t, err)

	keyValue, err := etcdCli.GetKeyValue(ctx, "/test/key")
	assert.Nil(t, err)

	assert.Equal(t, string(keyValue.Value), "{...}")
}

func Test_MutexLock(t *testing.T) {
	config := DefaultConfig()
	config.Endpoints = []string{"127.0.0.1:2379"}
	config.TTL = 10
	etcdCli := newClient(config)

	etcdMutex1, err := etcdCli.NewMutex("/test/lock",
		concurrency.WithTTL(int(config.TTL)))
	assert.Nil(t, err)

	err = etcdMutex1.Lock(time.Second * 1)
	assert.Nil(t, err)
	defer etcdMutex1.Unlock()

	// Grab the lock
	etcdMutex, err := etcdCli.NewMutex("/test/lock",
		concurrency.WithTTL(int(config.TTL)))
	assert.Nil(t, err)
	defer etcdMutex.Unlock()

	err = etcdMutex.Lock(time.Second * 1)
	assert.NotNil(t, err)
}
