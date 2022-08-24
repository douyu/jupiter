// Copyright 2022 Douyu
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

package sentinel

import (
	"context"
	"time"

	"github.com/alibaba/sentinel-golang/ext/datasource"
	"github.com/alibaba/sentinel-golang/util"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/pkg/errors"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type etcdv3DataSource struct {
	datasource.Base
	propertyKey         string
	lastUpdatedRevision int64
	client              *clientv3.Client
	// cancel is the func, call cancel will stop watching on the propertyKey
	cancel context.CancelFunc
	// closed indicate whether continuing to watch on the propertyKey
	closed util.AtomicBool
}

func newDataSource(client *clientv3.Client, key string, handlers ...datasource.PropertyHandler) (*etcdv3DataSource, error) {
	if client == nil {
		return nil, errors.New("The etcdv3 client is nil.")
	}
	ds := &etcdv3DataSource{
		client:      client,
		propertyKey: key,
	}
	for _, h := range handlers {
		ds.AddPropertyHandler(h)
	}
	return ds, nil
}

func (s *etcdv3DataSource) Initialize() error {
	err := s.doReadAndUpdate()
	if err != nil {
		xlog.Jupiter().Warn("Fail to update data for key when execute Etcdv3DataSource.Initialize()",
			zap.String("propertyKey", s.propertyKey), zap.Error(err))
	}
	go util.RunWithRecover(s.watch)
	return nil
}

func (s *etcdv3DataSource) ReadSource() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := s.client.Get(ctx, s.propertyKey)
	if err != nil {
		return nil, errors.Errorf("Fail to get value for property key[%s]", s.propertyKey)
	}
	if resp.Count == 0 {
		return nil, errors.Errorf("The key[%s] is not existed in etcd server.", s.propertyKey)
	}
	s.lastUpdatedRevision = resp.Header.GetRevision()
	xlog.Jupiter().Info("[Etcdv3] Get the newest data for key", zap.String("propertyKey", s.propertyKey),
		zap.Int64("revision", resp.Header.GetRevision()), zap.ByteString("value", resp.Kvs[0].Value))
	return resp.Kvs[0].Value, nil
}

func (s *etcdv3DataSource) doReadAndUpdate() error {

	src, err := s.ReadSource()
	if err != nil {
		xlog.Jupiter().Warn("etcd ReadSource failed", xlog.FieldErr(err))
		return err
	}
	xlog.Jupiter().Debug("etcd ReadSource", xlog.FieldEvent("doReadAndUpdate"), xlog.FieldAddrAny(string(src)))

	return s.Handle(src)
}

func (s *etcdv3DataSource) processWatchResponse(resp *clientv3.WatchResponse) {
	if resp == nil {
		xlog.Jupiter().Debug("etcd ReadSource", xlog.FieldEvent("processWatchResponse read nil"))
		return
	}
	if resp.CompactRevision > s.lastUpdatedRevision {
		s.lastUpdatedRevision = resp.CompactRevision
	}
	if resp.Header.GetRevision() > s.lastUpdatedRevision {
		s.lastUpdatedRevision = resp.Header.GetRevision()
	}

	if err := resp.Err(); err != nil {
		xlog.Jupiter().Error("etcd ReadSource", xlog.FieldEvent("Watch on etcd endpoints err"), xlog.FieldErr(err))
		return
	}

	for _, ev := range resp.Events {
		if ev.Type == mvccpb.PUT {
			err := s.doReadAndUpdate()
			if err != nil {
				xlog.Jupiter().Error("etcd ReadSource", xlog.FieldEvent("Fail to execute doReadAndUpdate for PUT"), xlog.FieldErr(err))
			}
		}
		if ev.Type == mvccpb.DELETE {
			xlog.Jupiter().Debug("etcd ReadSource", xlog.FieldEvent("processWatchResponse delete"))

			updateErr := s.Handle(nil)
			if updateErr != nil {
				xlog.Jupiter().Error("etcd ReadSource", xlog.FieldEvent("Fail to execute doReadAndUpdate for DELETE"), xlog.FieldErr(updateErr))
			}
		}
	}
}

func (s *etcdv3DataSource) watch() {
	// Add watch for propertyKey from lastUpdatedRevision updated after Initializing
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	rch := s.client.Watch(ctx, s.propertyKey, clientv3.WithCreatedNotify(), clientv3.WithRev(s.lastUpdatedRevision))

	for {
		for resp := range rch {
			s.processWatchResponse(&resp)
		}
		// Stop watching if datasource had been closed.
		if s.closed.Get() {
			xlog.Jupiter().Info("etcd ReadSource", xlog.FieldEvent("watch closed detected"))
			return
		}
		time.Sleep(time.Duration(1) * time.Second)
		ctx, cancel = context.WithCancel(context.Background())
		s.cancel = cancel
		if s.lastUpdatedRevision > 0 {
			rch = s.client.Watch(ctx, s.propertyKey, clientv3.WithRev(s.lastUpdatedRevision+1))
		} else {
			rch = s.client.Watch(ctx, s.propertyKey)
		}
	}
}

func (s *etcdv3DataSource) Close() error {
	// stop to watch property key.
	s.closed.Set(true)
	s.cancel()

	return nil
}
