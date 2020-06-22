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
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/douyu/jupiter/pkg/ecode"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/xlog"
	jsoniter "github.com/json-iterator/go"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/naming"
	"google.golang.org/grpc/resolver"
)

// etcdResolver implement grpc resolve.Builder
type etcdResolver struct {
	client *etcdv3.Client
	*Config
}

func newETCDResolver(config *Config) *etcdResolver {
	if config.logger == nil {
		config.logger = xlog.DefaultLogger
	}
	config.logger = config.logger.With(xlog.FieldMod("resolver.etcd"), xlog.FieldAddrAny(config.Config.Endpoints))
	res := &etcdResolver{
		client: config.Config.Build(),
		Config: config,
	}
	return res
}

// Build ...
func (r *etcdResolver) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOption) (resolver.Resolver, error) {
	r.logger.Info("build etcd resolver", xlog.String("watchFullKey", "/"+r.Prefix+"/"+target.Endpoint))
	go r.watch(cc, r.Prefix, target.Endpoint)
	return r, nil
}

// Scheme ...
func (r etcdResolver) Scheme() string {
	return "etcd"
}

// ResolveNow ...
func (r etcdResolver) ResolveNow(rn resolver.ResolveNowOption) {
	r.logger.Info("resolve now")
}

// close closes the resolver.
func (r etcdResolver) Close() {
	r.client.Close()
	r.logger.Info("close")
}

func (r *etcdResolver) deleteAddrList(al *AddressList, prefix string, kvs ...*mvccpb.KeyValue) {
	al.mtx.Lock()
	defer al.mtx.Unlock()
	for _, kv := range kvs {
		var addr = strings.TrimPrefix(string(kv.Key), prefix)
		if strings.HasPrefix(addr, "providers/grpc://") {
			// 解析服务注册键
			addr = strings.TrimPrefix(addr, "providers/")
			if addr == "" {
				continue
			}
			host, _, err := parseURI(addr)
			if err != nil {
				r.logger.Error("parse uri", xlog.FieldErrKind(ecode.ErrKindUriErr), xlog.FieldErr(err), xlog.FieldKey(string(kv.Key)))
			}
			delete(al.regInfos, host)
		}

		if strings.HasPrefix(addr, "configurators/grpc://") {
			// 解析服务配置键
			addr = strings.TrimPrefix(addr, "configurators/")
			if addr == "" {
				continue
			}
			host, _, err := parseURI(addr)
			if err != nil {
				r.logger.Error("parse uri", xlog.FieldErrKind(ecode.ErrKindUriErr), xlog.FieldErr(err), xlog.FieldKey(string(kv.Key)))
			}
			delete(al.cfgInfos, host)
		}

		if isIPPort(addr) {
			// 直接删除addr 因为Delete操作的value值为空
			delete(al.cfgInfos, addr)
			delete(al.regInfos, addr)
		}
	}
}

func isIPPort(addr string) bool {
	_, _, err := net.SplitHostPort(addr)
	return err == nil
}

func (r *etcdResolver) updateAddrList(al *AddressList, prefix string, kvs ...*mvccpb.KeyValue) {
	al.mtx.Lock()
	defer al.mtx.Unlock()
	for _, kv := range kvs {
		var addr = strings.TrimPrefix(string(kv.Key), prefix)
		switch {
		// 解析服务注册键
		case strings.HasPrefix(addr, "providers/grpc://"):
			addr = strings.TrimPrefix(addr, "providers/")
			host, meta, err := parseURI(addr)
			if err != nil {
				r.logger.Error("parse uri", xlog.FieldErrKind(ecode.ErrKindUriErr), xlog.FieldErr(err), xlog.FieldKey(string(kv.Key)))
				continue
			}
			if sm, err := parseValue(kv.Value); err == nil {
				for key, val := range sm.Metadata {
					meta.Set(key, val)
				}
			}
			al.regInfos[host] = &meta
		case strings.HasPrefix(addr, "configurators/grpc://"):
			addr = strings.TrimPrefix(addr, "configurators/")
			host, meta, err := parseURI(addr)
			if err != nil {
				r.logger.Error("parse uri", xlog.FieldErrKind(ecode.ErrKindUriErr), xlog.FieldErr(err), xlog.FieldKey(string(kv.Key)))
				continue
			}
			if sm, err := parseValue(kv.Value); err == nil {
				for key, val := range sm.Metadata {
					meta.Set(key, val)
				}
			}
			al.cfgInfos[host] = &meta
		case isIPPort(addr): // v1 协议
			var meta naming.Update
			if err := jsoniter.Unmarshal(kv.Value, &meta); err != nil {
				r.logger.Error("unmarshal metadata", xlog.FieldErrKind(ecode.ErrKindUnmarshalConfigErr), xlog.FieldErr(err), xlog.FieldKey(string(kv.Key)), xlog.FieldValue(string(kv.Value)))
				continue
			}

			if _, ok := al.cfgInfos[addr]; !ok {
				al.cfgInfos[addr] = &url.Values{}
			}

			// 解析value
			switch meta.Op {
			case naming.Add:
				al.regInfos[addr] = &url.Values{}
				al.cfgInfos[addr].Set("enable", "true")
				al.cfgInfos[addr].Set("weight", "100")
			case naming.Delete:
				al.regInfos[addr] = &url.Values{}
				al.cfgInfos[addr].Set("enable", "false")
				al.cfgInfos[addr].Set("weight", "0")
			}

		}
	}
	r.logger.Info("update addr list", xlog.Any("reg_info", al.regInfos), xlog.Any("cfg_info", al.cfgInfos))
}

// AddressList ...
type AddressList struct {
	serverName string

	// TODO 2019/9/10 gorexlv: need lock
	regInfos map[string]*url.Values // 注册信息
	cfgInfos map[string]*url.Values // 配置信息

	mtx sync.RWMutex
}

// List ...
func (al *AddressList) List() []resolver.Address {
	// TODO 2019/9/10 gorexlv:
	addrs := make([]resolver.Address, 0)
	al.mtx.RLock()
	defer al.mtx.RUnlock()
	for addr, values := range al.regInfos {
		metadata := *values
		address := resolver.Address{
			Addr:       addr,
			ServerName: al.serverName,
			Attributes: attributes.New(),
		}
		if infos, ok := al.cfgInfos[addr]; ok {
			for cfg := range *infos {
				metadata.Set(cfg, infos.Get(cfg))
				// address.Attributes.WithValues(cfg, infos.Get(cfg))
			}
		}
		if enable := metadata.Get("enable"); enable != "" && enable != "true" {
			continue
		}
		if metadata.Get("weight") == "" {
			metadata.Set("weight", "100")
		}
		// group: 客户端配置的分组，默认为default
		// metadata.Get("group"): 服务端配置的分组
		// if metadata.Get("group") != "" && metadata.Get("group") != group {
		// 	continue
		// }
		address.Metadata = &metadata
		addrs = append(addrs, address)
	}
	return addrs
}

func escape(raw string) string {
	return strings.Replace(raw, "$", "%24", -1)
}

func (r *etcdResolver) watch(cc resolver.ClientConn, prefix string, serviceName string) {
	cli := r.client
	target := fmt.Sprintf("/%s/%s/", prefix, serviceName)
	for {
		var al = &AddressList{
			serverName: serviceName,
			regInfos:   make(map[string]*url.Values),
			cfgInfos:   make(map[string]*url.Values),
			mtx:        sync.RWMutex{},
		}

		getResp, err := cli.Get(context.Background(), target, clientv3.WithPrefix())
		if err != nil {
			r.logger.Error(ecode.MsgWatchRequestErr, xlog.FieldErrKind(ecode.ErrKindRequestErr), xlog.FieldErr(err), xlog.FieldAddr(target))
			time.Sleep(time.Second * 5)
			continue
		}

		r.updateAddrList(al, target, getResp.Kvs...)

		cc.UpdateState(resolver.State{
			Addresses: al.List(),
		})

		// 处理配置键变更, reg_key
		// 处理注册键变更, cfg_key
		ctx, cancel := context.WithCancel(context.Background())
		rch := cli.Watch(ctx, target, clientv3.WithPrefix(), clientv3.WithCreatedNotify())
		for n := range rch {
			for _, ev := range n.Events {
				switch ev.Type {
				// 添加或者更新
				case mvccpb.PUT:
					r.updateAddrList(al, target, ev.Kv)
				// 硬删除
				case mvccpb.DELETE:
					r.deleteAddrList(al, target, ev.Kv)
				}
			}

			cc.UpdateState(resolver.State{
				Addresses: al.List(),
			})
		}

		cancel()
	}
}

func parseURI(uri string) (host string, meta url.Values, err error) {
	uri, err = url.PathUnescape(uri)
	if err != nil {
		return
	}
	if strings.Index(uri, "://") > 0 {
		u, e := url.Parse(uri)
		if e != nil || u == nil {
			return "", nil, e
		}
		host = u.Host
		meta = u.Query()
		meta.Set("scheme", u.Scheme)
		return
	}
	return uri, meta, nil
}

func parseValue(raw []byte) (*server.ServiceInfo, error) {
	var meta server.ServiceInfo
	if err := jsoniter.Unmarshal(raw, &meta); err != nil {
		return nil, err
	}

	return &meta, nil
}
