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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/xlog"
)

type etcdv3Registry struct {
	client *etcdv3.Client
	lease  clientv3.LeaseID
	kvs    sync.Map
	*Config
	cancel context.CancelFunc
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

// RegisterService register service to registry
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

	key := e.providerKey(info)
	val := e.providerValue(info)

	_, err := e.client.Put(ctx, key, string(val), opOptions...)
	if err != nil {
		e.logger.Error("register service", xlog.FieldErrKind(ecode.ErrKindRegisterErr), xlog.FieldErr(err), xlog.FieldKeyAny(key), xlog.FieldValueAny(info))
		return err
	}

	e.logger.Info("register service", xlog.FieldKeyAny(key), xlog.FieldValueAny(info))
	e.kvs.Store(key, val)
	return err
}

// UnregisterService unregister service from registry
func (e *etcdv3Registry) UnregisterService(ctx context.Context, info *server.ServiceInfo) error {
	return e.unregister(ctx, e.providerKey(info))
}

// ListServices list service registered in registry with name `name`
func (e *etcdv3Registry) ListServices(ctx context.Context, name string, scheme string) (services []*server.ServiceInfo, err error) {
	target := fmt.Sprintf("/%s/%s/providers/%s://", e.Prefix, name, scheme)
	getResp, getErr := e.client.Get(ctx, target, clientv3.WithPrefix())
	if getErr != nil {
		e.logger.Error(ecode.MsgWatchRequestErr, xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(getErr), xlog.FieldAddr(target))
		return nil, getErr
	}

	for _, kv := range getResp.Kvs {
		var service server.ServiceInfo
		if err := json.Unmarshal(kv.Value, &service); err != nil {
			e.logger.Warnf("invalid service", xlog.FieldErr(err))
			continue
		}
		services = append(services, &service)
	}

	return
}

func (e *etcdv3Registry) filterKey(name string, scheme string, key []byte) bool {
	// standard key format: /{prefix}/{app_name}/providers/{scheme}://{ip}:{port}
	// or: /{prefix}/{app_name}/configurators/{scheme}:///
	// or: /{prefix}/{app_name}/consumers/{scheme}://
	target := fmt.Sprintf("/%s/%s/", e.Prefix, name)
	return bytes.HasPrefix(key, []byte(target+"/providers/"+scheme+"://")) ||
		bytes.HasPrefix(key, []byte(target+"/configurators/"+scheme+"://")) ||
		bytes.HasPrefix(key, []byte(target+"/consumers/"+scheme+"://"))
}

// WatchServices list services then watch service change event
func (e *etcdv3Registry) WatchServices(ctx context.Context, name string, scheme string) (services []*server.ServiceInfo, messages chan *registry.EventMessage, err error) {
	target := fmt.Sprintf("/%s/%s/", e.Prefix, name)
	messages = make(chan *registry.EventMessage, 8)
	getResp, getErr := e.client.Get(ctx, target, clientv3.WithPrefix())
	if getErr != nil {
		e.logger.Error(ecode.MsgWatchRequestErr, xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(getErr), xlog.FieldAddr(target))
		return nil, nil, getErr
	}

	lastRevision := getResp.Header.Revision

	for _, kv := range getResp.Kvs {
		if !e.filterKey(name, scheme, kv.Key) {
			continue
		}
		var service server.ServiceInfo
		if err := json.Unmarshal(kv.Value, &service); err != nil {
			e.logger.Warnf("invalid service", xlog.FieldErr(err))
			continue
		}
		services = append(services, &service)
	}

	xgo.Go(func() {
		ctx, cancel := context.WithCancel(context.Background())
		e.cancel = cancel

		rch := e.client.Watch(ctx, target, clientv3.WithPrefix(), clientv3.WithCreatedNotify(), clientv3.WithRev(lastRevision))
		for {
			for n := range rch {
				if n.CompactRevision > lastRevision {
					lastRevision = n.CompactRevision
				}
				if n.Header.GetRevision() > lastRevision {
					lastRevision = n.Header.GetRevision()
				}
				if err := n.Err(); err != nil {
					e.logger.Error(ecode.MsgWatchRequestErr, xlog.FieldErrKind(ecode.ErrKindRegisterErr), xlog.FieldErr(err), xlog.FieldAddr(target))
					continue
				}
				for _, ev := range n.Events {
					if !e.filterKey(name, scheme, ev.Kv.Key) {
						continue
					}
					msg, err := extractEventMessage(target, ev)
					if err != nil {
						e.logger.Error(ecode.MsgWatchRequestErr, xlog.FieldErrKind(ecode.ErrKindRegisterErr), xlog.FieldErr(err), xlog.FieldAddr(target))
						continue
					}
					messages <- msg
				}
			}
			ctx, cancel := context.WithCancel(context.Background())
			e.cancel = cancel
			if lastRevision > 0 {
				rch = e.client.Watch(ctx, target, clientv3.WithPrefix(), clientv3.WithCreatedNotify(), clientv3.WithRev(lastRevision))
			} else {
				rch = e.client.Watch(ctx, target, clientv3.WithPrefix(), clientv3.WithCreatedNotify())
			}
		}
	})
	return services, messages, nil
}

func extractEventMessage(target string, event *clientv3.Event) (*registry.EventMessage, error) {
	var em = &registry.EventMessage{
		Event: registry.EventUnknown,
		Kind:  registry.KindUnknown,
	}
	registryKey, err := ToRegistryKey(string(event.Kv.Key))
	if err != nil {
		return nil, err
	}

	em.Name = registryKey.AppName
	em.Address = registryKey.Host
	em.Scheme = registryKey.Scheme
	em.Kind = registryKey.Kind

	switch event.Type {
	case mvccpb.PUT:
		em.Event = registry.EventUpdate
		if em.Kind == registry.KindProvider {
			var configuration registry.Configuration
			if err := json.Unmarshal(event.Kv.Value, &configuration); err != nil { return nil, err }
			em.Message = configuration
		}
		if em.Kind == registry.KindConfigurator {
			var serviceInfo server.ServiceInfo
			if err := json.Unmarshal(event.Kv.Value, &serviceInfo); err != nil {
				return nil, err
			}
			em.Message = serviceInfo
		}
	case mvccpb.DELETE:
		em.Event = registry.EventDelete
	}

	return em, nil
}

func (e *etcdv3Registry) unregister(ctx context.Context, key string) error {
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
	if e.cancel != nil {
		e.cancel()
	}
	var wg sync.WaitGroup
	e.kvs.Range(func(k, v interface{}) bool {
		wg.Add(1)
		go func(k interface{}) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err := e.unregister(ctx, k.(string))
			if err != nil {
				e.logger.Error("unregister service", xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err), xlog.FieldErr(err), xlog.FieldKeyAny(k), xlog.FieldValueAny(v))
			} else {
				e.logger.Info("unregister service", xlog.FieldKeyAny(k), xlog.FieldValueAny(v))
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

func (e *etcdv3Registry) providerKey(info *server.ServiceInfo) string {
	return fmt.Sprintf("/%s/%s/providers/%s://%s", e.Prefix, info.Name, info.Scheme, info.Address)
}

func (e *etcdv3Registry) providerValue(info *server.ServiceInfo) string {
	val, _ := json.Marshal(info)
	return string(val)
}
