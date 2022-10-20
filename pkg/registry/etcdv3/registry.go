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
	"net"
	"strings"
	"sync"
	"time"

	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/registry"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
)

type etcdv3Registry struct {
	client *etcdv3.Client
	kvs    sync.Map
	*Config
	cancel   context.CancelFunc
	rmu      *sync.RWMutex
	sessions map[string]*concurrency.Session
}

const (
	servicePrefix = "%s:%s:%s:%s/"
	// schema:appname:version:mode/host:port
	registerService = "%s:%s:%s:%s/%s"
)

func newETCDRegistry(config *Config) (*etcdv3Registry, error) {
	if config.logger == nil {
		config.logger = xlog.Jupiter()
	}
	config.logger = config.logger.With(xlog.FieldMod(ecode.ModRegistryETCD), xlog.FieldAddrAny(config.Config.Endpoints))
	etcdv3Client, err := config.Config.Build()
	if err != nil {
		config.logger.Error("create etcdv3 client", xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err))
		return nil, err
	}
	reg := &etcdv3Registry{
		client:   etcdv3Client,
		Config:   config,
		kvs:      sync.Map{},
		rmu:      &sync.RWMutex{},
		sessions: make(map[string]*concurrency.Session),
	}
	return reg, nil
}

func (reg *etcdv3Registry) Kind() string { return "etcdv3" }

// RegisterService register service to registry
func (reg *etcdv3Registry) RegisterService(ctx context.Context, info *server.ServiceInfo) error {
	err := reg.registerBiz(ctx, info)
	if err != nil {
		return err
	}
	return reg.registerMetric(ctx, info)
}

// UnregisterService unregister service from registry
func (reg *etcdv3Registry) UnregisterService(ctx context.Context, info *server.ServiceInfo) error {
	return reg.unregister(ctx, reg.registerKey(info))
}

// ListServices list service registered in registry with name `name`
func (reg *etcdv3Registry) ListServices(ctx context.Context, name string, scheme string) (services []*server.ServiceInfo, err error) {
	target := fmt.Sprintf(servicePrefix, scheme, name, "v1", conf.GetString("app.mode"))
	getResp, getErr := reg.client.Get(ctx, target, clientv3.WithPrefix())
	if getErr != nil {
		reg.logger.Error(ecode.MsgWatchRequestErr, xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(getErr), xlog.FieldAddr(target))
		return nil, getErr
	}

	for _, kv := range getResp.Kvs {
		var service server.ServiceInfo
		if err := json.Unmarshal(kv.Value, &service); err != nil {
			reg.logger.Warn("invalid service", xlog.FieldErr(err))
			continue
		}
		services = append(services, &service)
	}

	return
}

// WatchServices watch service change event, then return address list
func (reg *etcdv3Registry) WatchServices(ctx context.Context, name string, scheme string) (chan registry.Endpoints, error) {
	prefix := fmt.Sprintf(servicePrefix, scheme, name, "v1", conf.GetString("jupiter.mode"))
	watch, err := reg.client.WatchPrefix(context.Background(), prefix)
	if err != nil {
		return nil, err
	}

	var addresses = make(chan registry.Endpoints, 10)
	var al = &registry.Endpoints{
		Nodes:           make(map[string]server.ServiceInfo),
		RouteConfigs:    make(map[string]registry.RouteConfig),
		ConsumerConfigs: make(map[string]registry.ConsumerConfig),
		ProviderConfigs: make(map[string]registry.ProviderConfig),
	}

	for _, kv := range watch.IncipientKeyValues() {
		updateAddrList(al, prefix, scheme, kv)
	}

	// var snapshot registry.Endpoints
	// xstruct.CopyStruct(al, &snapshot)
	addresses <- *al.DeepCopy()

	xgo.Go(func() {
		for event := range watch.C() {
			switch event.Type {
			case mvccpb.PUT:
				updateAddrList(al, prefix, scheme, event.Kv)
			case mvccpb.DELETE:
				deleteAddrList(al, prefix, scheme, event.Kv)
			}

			// var snapshot registry.Endpoints
			// xstruct.CopyStruct(al, &snapshot)
			out := al.DeepCopy()
			select {
			// case addresses <- snapshot:
			case addresses <- *out:
			default:
				xlog.Jupiter().Warn("invalid")
			}
		}
	})

	return addresses, nil
}

func (reg *etcdv3Registry) unregister(ctx context.Context, key string) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, reg.ReadTimeout)
		defer cancel()
	}

	if err := reg.delSession(key); err != nil {
		return err
	}

	_, err := reg.client.Delete(ctx, key)
	if err == nil {
		reg.kvs.Delete(key)
	}
	return err
}

// Close ...
func (reg *etcdv3Registry) Close() error {
	if reg.cancel != nil {
		reg.cancel()
	}
	var wg sync.WaitGroup
	reg.kvs.Range(func(k, v interface{}) bool {
		wg.Add(1)
		go func(k interface{}) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err := reg.unregister(ctx, k.(string))
			if err != nil {
				reg.logger.Error("unregister service", xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err), xlog.FieldErr(err), xlog.FieldKeyAny(k), xlog.FieldValueAny(v))
			} else {
				reg.logger.Info("unregister service", xlog.FieldKeyAny(k), xlog.FieldValueAny(v))
			}
			cancel()
		}(k)
		return true
	})
	wg.Wait()
	return nil
}

func (reg *etcdv3Registry) registerMetric(ctx context.Context, info *server.ServiceInfo) error {
	if info.Kind != constant.ServiceGovernor {
		return nil
	}

	metric := "/prometheus/job/%s/%s/%s"

	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, reg.ReadTimeout)
		defer cancel()
	}

	val := info.Address
	key := fmt.Sprintf(metric, info.Name, pkg.HostName(), val)

	opOptions := make([]clientv3.OpOption, 0)
	// opOptions = append(opOptions, clientv3.WithSerializable())
	if ttl := reg.Config.ServiceTTL.Seconds(); ttl > 0 {
		//todo ctx without timeout for same as service life?
		sess, err := reg.getSession(key, concurrency.WithTTL(int(ttl)))
		if err != nil {
			return err
		}
		opOptions = append(opOptions, clientv3.WithLease(sess.Lease()))
	}
	_, err := reg.client.Put(ctx, key, val, opOptions...)
	if err != nil {
		reg.logger.Error("register service", xlog.FieldErrKind(ecode.ErrKindRegisterErr), xlog.FieldErr(err), xlog.FieldKeyAny(key), xlog.FieldValueAny(info))
		return err
	}

	reg.logger.Info("register service", xlog.FieldKeyAny(key), xlog.FieldValueAny(val))
	reg.kvs.Store(key, val)
	return nil

}
func (reg *etcdv3Registry) registerBiz(ctx context.Context, info *server.ServiceInfo) error {
	if _, ok := ctx.Deadline(); !ok {
		var readCancel context.CancelFunc
		ctx, readCancel = context.WithTimeout(ctx, reg.ReadTimeout)
		defer readCancel()
	}

	key := reg.registerKey(info)
	val := reg.registerValue(info)

	opOptions := make([]clientv3.OpOption, 0)
	// opOptions = append(opOptions, clientv3.WithSerializable())
	if ttl := reg.Config.ServiceTTL.Seconds(); ttl > 0 {
		//todo ctx without timeout for same as service life?
		sess, err := reg.getSession(key, concurrency.WithTTL(int(ttl)))
		if err != nil {
			return err
		}
		opOptions = append(opOptions, clientv3.WithLease(sess.Lease()))
	}
	_, err := reg.client.Put(ctx, key, val, opOptions...)
	if err != nil {
		reg.logger.Error("register service", xlog.FieldErrKind(ecode.ErrKindRegisterErr), xlog.FieldErr(err), xlog.FieldKeyAny(key), xlog.FieldValueAny(info))
		return err
	}
	reg.logger.Info("register service", xlog.FieldKeyAny(key), xlog.FieldValueAny(val))
	reg.kvs.Store(key, val)
	return nil
}

func (reg *etcdv3Registry) getSession(k string, opts ...concurrency.SessionOption) (*concurrency.Session, error) {
	reg.rmu.RLock()
	sess, ok := reg.sessions[k]
	reg.rmu.RUnlock()
	if ok {
		return sess, nil
	}
	sess, err := concurrency.NewSession(reg.client.Client, opts...)
	if err != nil {
		return sess, err
	}
	reg.rmu.Lock()
	reg.sessions[k] = sess
	reg.rmu.Unlock()
	return sess, nil
}

func (reg *etcdv3Registry) delSession(k string) error {
	if ttl := reg.Config.ServiceTTL.Seconds(); ttl > 0 {
		reg.rmu.RLock()
		sess, ok := reg.sessions[k]
		reg.rmu.RUnlock()
		if ok {
			reg.rmu.Lock()
			delete(reg.sessions, k)
			reg.rmu.Unlock()
			if err := sess.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

func (reg *etcdv3Registry) registerKey(info *server.ServiceInfo) string {
	return getServiceKey(info)
}

func (reg *etcdv3Registry) registerValue(info *server.ServiceInfo) string {
	update := Update{
		Op:   Add,
		Addr: info.Address,
	}

	val, _ := json.Marshal(update)

	return string(val)
}

func deleteAddrList(al *registry.Endpoints, prefix, scheme string, kvs ...*mvccpb.KeyValue) {
	for _, kv := range kvs {
		var addr = strings.TrimPrefix(string(kv.Key), prefix)

		if isIPPort(addr) {
			// 直接删除addr 因为Delete操作的value值为空
			delete(al.Nodes, addr)
			delete(al.RouteConfigs, addr)
		}
	}
}

func updateAddrList(al *registry.Endpoints, prefix, scheme string, kvs ...*mvccpb.KeyValue) {
	for _, kv := range kvs {
		var addr = strings.TrimPrefix(string(kv.Key), prefix)
		if isIPPort(addr) {
			var meta Update
			if err := json.Unmarshal(kv.Value, &meta); err != nil {
				xlog.Jupiter().Error("unmarshal meta", xlog.FieldErr(err),
					xlog.FieldExtMessage("value", string(kv.Value), "key", string(kv.Key)))
				continue
			}

			switch meta.Op {
			case Add:
				al.Nodes[addr] = server.ServiceInfo{
					Address: addr,
				}
			case Delete:
				delete(al.Nodes, addr)
			}
		}
	}
}

// getServiceKey ..
func getServiceKey(s *server.ServiceInfo) string {
	return fmt.Sprintf(registerService, s.Scheme, s.Name, "v1", conf.GetString("jupiter.mode"), s.Address)
}

func isIPPort(addr string) bool {
	_, _, err := net.SplitHostPort(addr)
	return err == nil
}
