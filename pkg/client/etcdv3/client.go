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
	"crypto/tls"
	"crypto/x509"
	"github.com/coreos/etcd/clientv3/concurrency"
	"io/ioutil"
	"strings"
	"time"

	"github.com/douyu/jupiter/pkg/ecode"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/douyu/jupiter/pkg/xlog"
	grpcprom "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

// Client ...
type Client struct {
	*clientv3.Client
	config *Config
}

// New ...
func newClient(config *Config) *Client {
	conf := clientv3.Config{
		Endpoints:            config.Endpoints,
		DialTimeout:          config.ConnectTimeout,
		DialKeepAliveTime:    10 * time.Second,
		DialKeepAliveTimeout: 3 * time.Second,
		DialOptions: []grpc.DialOption{
			grpc.WithBlock(),
			grpc.WithUnaryInterceptor(grpcprom.UnaryClientInterceptor),
			grpc.WithStreamInterceptor(grpcprom.StreamClientInterceptor),
		},
		AutoSyncInterval: config.AutoSyncInterval,
	}

	config.logger = config.logger.With(xlog.FieldAddrAny(config.Endpoints))

	if config.Endpoints == nil {
		config.logger.Panic("client etcd endpoints empty", xlog.FieldMod(ecode.ModClientETCD), xlog.FieldValueAny(config))
	}

	if !config.Secure {
		conf.DialOptions = append(conf.DialOptions, grpc.WithInsecure())
	}

	if config.BasicAuth {
		conf.Username = config.UserName
		conf.Password = config.Password
	}

	tlsEnabled := false
	tlsConfig := &tls.Config{
		InsecureSkipVerify: false,
	}

	if config.CaCert != "" {
		certBytes, err := ioutil.ReadFile(config.CaCert)
		if err != nil {
			config.logger.Panic("parse CaCert failed", xlog.Any("err", err))
		}

		caCertPool := x509.NewCertPool()
		ok := caCertPool.AppendCertsFromPEM(certBytes)

		if ok {
			tlsConfig.RootCAs = caCertPool
		}
		tlsEnabled = true
	}

	if config.CertFile != "" && config.KeyFile != "" {
		tlsCert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			config.logger.Panic("load CertFile or KeyFile failed", xlog.Any("config", config), xlog.Any("err", err))
		}
		tlsConfig.Certificates = []tls.Certificate{tlsCert}
		tlsEnabled = true
	}

	if tlsEnabled {
		conf.TLS = tlsConfig
	}

	client, err := clientv3.New(conf)

	if err != nil {
		config.logger.Panic("client etcd start panic", xlog.FieldMod(ecode.ModClientETCD), xlog.FieldErrKind(ecode.ErrKindAny), xlog.FieldErr(err), xlog.FieldValueAny(config))
	}

	cc := &Client{
		Client: client,
		config: config,
	}

	config.logger.Info("dial etcd server")
	return cc
}

// GetKeyValue queries etcd key, returns mvccpb.KeyValue
func (client *Client) GetKeyValue(ctx context.Context, key string) (kv *mvccpb.KeyValue, err error) {
	rp, err := client.Client.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	if len(rp.Kvs) > 0 {
		return rp.Kvs[0], nil
	}

	return
}

// GetPrefix get prefix
func (client *Client) GetPrefix(ctx context.Context, prefix string) (map[string]string, error) {
	var (
		vars = make(map[string]string)
	)

	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return vars, err
	}

	for _, kv := range resp.Kvs {
		vars[string(kv.Key)] = string(kv.Value)
	}

	return vars, nil
}

// DelPrefix 按前缀删除
func (client *Client) DelPrefix(ctx context.Context, prefix string) (deleted int64, err error) {
	resp, err := client.Delete(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return 0, err
	}
	return resp.Deleted, err
}

// GetValues queries etcd for keys prefixed by prefix.
func (client *Client) GetValues(ctx context.Context, keys ...string) (map[string]string, error) {
	var (
		firstRevision = int64(0)
		vars          = make(map[string]string)
		maxTxnOps     = 128
		getOps        = make([]string, 0, maxTxnOps)
	)

	doTxn := func(ops []string) error {
		txnOps := make([]clientv3.Op, 0, maxTxnOps)

		for _, k := range ops {
			txnOps = append(txnOps, clientv3.OpGet(k,
				clientv3.WithPrefix(),
				clientv3.WithSort(clientv3.SortByKey, clientv3.SortDescend),
				clientv3.WithRev(firstRevision)))
		}

		result, err := client.Txn(ctx).Then(txnOps...).Commit()
		if err != nil {
			return err
		}
		for i, r := range result.Responses {
			originKey := ops[i]
			originKeyFixed := originKey
			if !strings.HasSuffix(originKeyFixed, "/") {
				originKeyFixed = originKey + "/"
			}
			for _, ev := range r.GetResponseRange().Kvs {
				k := string(ev.Key)
				if k == originKey || strings.HasPrefix(k, originKeyFixed) {
					vars[string(ev.Key)] = string(ev.Value)
				}
			}
		}
		if firstRevision == 0 {
			firstRevision = result.Header.GetRevision()
		}
		return nil
	}
	for _, key := range keys {
		getOps = append(getOps, key)
		if len(getOps) >= maxTxnOps {
			if err := doTxn(getOps); err != nil {
				return vars, err
			}
			getOps = getOps[:0]
		}
	}
	if len(getOps) > 0 {
		if err := doTxn(getOps); err != nil {
			return vars, err
		}
	}
	return vars, nil
}

//GetLeaseSession 创建租约会话
func (client *Client) GetLeaseSession(ctx context.Context, opts ...concurrency.SessionOption) (leaseSession *concurrency.Session, err error) {
	return concurrency.NewSession(client.Client, opts...)
}
