package etcdv3

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

func Test_GetKeyValue(t *testing.T) {
	config := DefaultConfig()
	config.Endpoints = []string{"localhost:2379"}
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
	config.Endpoints = []string{"localhost:2379"}
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

func TestClient_GetPrefix(t *testing.T) {
	type fields struct {
		Client *clientv3.Client
		config *Config
	}
	type args struct {
		ctx    context.Context
		prefix string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "test getprefix",
			fields: fields{
				Client: StdConfig("default").MustSingleton().Client,
				config: StdConfig("default"),
			},
			args: args{
				ctx:    context.Background(),
				prefix: "/test/getprefix",
			},
			want: map[string]string{
				"/test/getprefix/1": "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			client := &Client{
				Client: tt.fields.Client,
				config: tt.fields.config,
			}

			_, err := client.Put(tt.args.ctx, "/test/getprefix/1", "1")
			assert.Nil(t, err)

			got, err := client.GetPrefix(tt.args.ctx, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetPrefix() = %v, want %v", got, tt.want)
			}

		})
	}
}

func TestClient_DelPrefix(t *testing.T) {
	type fields struct {
		Client *clientv3.Client
		config *Config
	}
	type args struct {
		ctx    context.Context
		prefix string
	}
	tests := []struct {
		name        string
		fields      fields
		args        args
		wantDeleted int64
		wantErr     bool
	}{
		{
			name: "test operation",
			fields: fields{
				Client: StdConfig("default").MustSingleton().Client,
				config: StdConfig("default"),
			},
			args: args{
				ctx:    context.Background(),
				prefix: "/test/delprefix",
			},
			wantDeleted: 0,
			wantErr:     false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				Client: tt.fields.Client,
				config: tt.fields.config,
			}
			gotDeleted, err := client.DelPrefix(tt.args.ctx, tt.args.prefix)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.DelPrefix() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotDeleted != tt.wantDeleted {
				t.Errorf("Client.DelPrefix() = %v, want %v", gotDeleted, tt.wantDeleted)
			}
		})
	}
}

func TestClient_GetValues(t *testing.T) {
	type fields struct {
		Client *clientv3.Client
		config *Config
	}
	type args struct {
		ctx  context.Context
		keys []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]string
		wantErr bool
	}{
		{
			name: "test operation",
			fields: fields{
				Client: StdConfig("default").MustSingleton().Client,
				config: StdConfig("default"),
			},
			args: args{
				ctx:  context.Background(),
				keys: []string{"/test/getvalues/1"},
			},
			want: map[string]string{
				"/test/getvalues/1": "1",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := &Client{
				Client: tt.fields.Client,
				config: tt.fields.config,
			}

			_, err := client.Put(tt.args.ctx, "/test/getvalues/1", "1")
			assert.Nil(t, err)

			got, err := client.GetValues(tt.args.ctx, tt.args.keys...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Client.GetValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Client.GetValues() = %v, want %v", got, tt.want)
			}
		})
	}
}
