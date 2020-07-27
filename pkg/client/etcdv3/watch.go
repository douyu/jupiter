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
	"sync"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/xlog"
)

// Watch A watch only tells the latest revision
type Watch struct {
	revision  int64
	cancel    context.CancelFunc
	eventChan chan *clientv3.Event
	lock      *sync.RWMutex
	logger    *xlog.Logger

	incipientKVs []*mvccpb.KeyValue
}

// C ...
func (w *Watch) C() chan *clientv3.Event {
	return w.eventChan
}

// IncipientKeyValues incipient key and values
func (w *Watch) IncipientKeyValues() []*mvccpb.KeyValue {
	return w.incipientKVs
}

// NewWatch ...
func (client *Client) WatchPrefix(ctx context.Context, prefix string) (*Watch, error) {
	resp, err := client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	var w = &Watch{
		revision:     resp.Header.Revision,
		eventChan:    make(chan *clientv3.Event, 100),
		incipientKVs: resp.Kvs,
	}

	xgo.Go(func() {
		ctx, cancel := context.WithCancel(context.Background())
		w.cancel = cancel
		rch := client.Client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify(), clientv3.WithRev(w.revision))
		for {
			for n := range rch {
				if n.CompactRevision > w.revision {
					w.revision = n.CompactRevision
				}
				if n.Header.GetRevision() > w.revision {
					w.revision = n.Header.GetRevision()
				}
				if err := n.Err(); err != nil {
					xlog.Error(ecode.MsgWatchRequestErr, xlog.FieldErrKind(ecode.ErrKindRegisterErr), xlog.FieldErr(err), xlog.FieldAddr(prefix))
					continue
				}
				for _, ev := range n.Events {
					select {
					case w.eventChan <- ev:
					default:
						xlog.Error("watch etcd with prefix", xlog.Any("err", "block event chan, drop event message"))
					}
				}
			}
			ctx, cancel := context.WithCancel(context.Background())
			w.cancel = cancel
			if w.revision > 0 {
				rch = client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify(), clientv3.WithRev(w.revision))
			} else {
				rch = client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify())
			}
		}
	})

	return w, nil
}

// Close close watch
func (w *Watch) Close() error {
	if w.cancel != nil {
		w.cancel()
	}
	return nil
}
