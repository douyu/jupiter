// Copyright 2020 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etcdv3

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/douyu/jupiter/pkg/ecode"

	"github.com/coreos/etcd/clientv3"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/xlog"
)

type etcdv3Registry struct {
	client *etcdv3.Client
	lease  clientv3.LeaseID
	kvs    sync.Map
	*Config
}

func newETCDRegistry(config *Config) *etcdv3Registry {
	if config.logger == nil {
		config.logger = xlog.DefaultLogger
	}
	config.logger = config.logger.With(xlog.FieldMod(ecode.ModRegistryETCD), xlog.FieldAddrAny(config.Config.Endpoints))
	res := &etcdv3Registry{
		client: config.Config.Build(),
		Config: config,
		kvs:    sync.Map{},
	}
	return res
}

// RegisterService ...
func (e *etcdv3Registry) RegisterService(ctx context.Context, info *server.ServiceInfo) error {
	opOptions := make([]clientv3.OpOption, 0)
	if e.lease != 0 {
		opOptions = append(opOptions, clientv3.WithLease(e.lease), clientv3.WithSerializable())
	}

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.ReadTimeout)
		defer cancel()
	}

	key := fmt.Sprintf("/%s/%s/providers/%s://%s:%d", e.Prefix, info.Name, info.Scheme, info.IP, info.Port)
	val, err := json.Marshal(info)

	if err != nil {
		return err
	}

	_, err = e.client.Put(ctx, key, string(val), opOptions...)
	if err != nil {
		e.logger.Error("register service", xlog.FieldErrKind(ecode.ErrKindRegisterErr), xlog.FieldErr(err), xlog.FieldKeyAny(key), xlog.FieldValueAny(info))
		return err
	}
	// xdebug.PrintKVWithPrefix("registry", "register key", key)
	e.logger.Info("register service", xlog.FieldKeyAny(key), xlog.FieldValueAny(info))
	e.kvs.Store(key, val)
	return err
}

// DeregisterService ...
func (e *etcdv3Registry) DeregisterService(ctx context.Context, info *server.ServiceInfo) error {
	key := fmt.Sprintf("/%s/%s/providers/%s://%s:%d", e.Prefix, info.Name, info.Scheme, info.IP, info.Port)
	return e.deregister(ctx, key)
}

func (e *etcdv3Registry) deregister(ctx context.Context, key string) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, e.ReadTimeout)
		defer cancel()
	}
	_, err := e.client.Delete(ctx, key)
	if err == nil {
		e.kvs.Delete(key)
	}
	return err
}

// Close ...
func (e *etcdv3Registry) Close() error {
	var wg sync.WaitGroup
	e.kvs.Range(func(k, v interface{}) bool {
		wg.Add(1)
		go func(k interface{}) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err := e.deregister(ctx, k.(string))
			if err != nil {
				e.logger.Error("deregister service", xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err), xlog.FieldErr(err), xlog.FieldKeyAny(k), xlog.FieldValueAny(v))
			} else {
				e.logger.Info("deregister service", xlog.FieldKeyAny(k), xlog.FieldValueAny(v))
			}
			cancel()
		}(k)
		return true
	})
	wg.Wait()

	if e.lease > 0 {
		// revoke 有一些延迟，考虑直接删除
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		_, err := e.client.Revoke(ctx, e.lease)
		cancel()
		return err
	}
	return nil
}
