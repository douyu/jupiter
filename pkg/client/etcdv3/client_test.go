package etcdv3

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.etcd.io/etcd/pkg/mock/mockserver"
)

func startMockServer() {
	ms, err := mockserver.StartMockServers(1)
	if err != nil {
		log.Fatal(err)
	}

	if err := ms.StartAt(0); err != nil {
		log.Fatal(err)
	}
}

func TestMain(m *testing.M) {
	go startMockServer()
}

func Test_GetKeyValue(t *testing.T) {
	config := DefaultConfig()
	config.Endpoints = []string{"localhost:0"}
	config.TTL = 5
	etcdCli, err := newClient(config)
	assert.Nil(t, err)

	ctx := context.TODO()

	leaseSession, err := etcdCli.GetLeaseSession(ctx, concurrency.WithTTL(int(config.TTL)))
	assert.Nil(t, err)
	defer leaseSession.Close()
	fmt.Printf("111=%+v\n", 111)

	_, err = etcdCli.Client.KV.Put(ctx, "/test/key", "{...}", clientv3.WithLease(leaseSession.Lease()))
	assert.Nil(t, err)

	keyValue, err := etcdCli.GetKeyValue(ctx, "/test/key")
	assert.Nil(t, err)

	assert.Equal(t, string(keyValue.Value), "{...}")
}

func Test_MutexLock(t *testing.T) {
	config := DefaultConfig()
	config.Endpoints = []string{"localhost:0"}
	config.TTL = 10
	etcdCli, err := newClient(config)
	assert.Nil(t, err)

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
