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
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/douyu/jupiter/pkg/xlog"
)

// Watch A watch only tells the latest revision
type Watch struct {
	revision  int64
	client    *Client
	cancel    context.CancelFunc
	eventChan chan *clientv3.Event
	lock      *sync.RWMutex
	logger    *xlog.Logger
}

// C ...
func (w *Watch) C() chan *clientv3.Event {
	return w.eventChan
}

func (w *Watch) update(resp *clientv3.WatchResponse) {
	if resp.CompactRevision > w.revision {
		w.revision = resp.CompactRevision
	} else if resp.Header.GetRevision() > w.revision {
		w.revision = resp.Header.GetRevision()
	}

	if err := resp.Err(); err != nil {
		w.logger.Error("handle watch update", xlog.Any("err", err))
		return
	}

	for _, event := range resp.Events {
		select {
		case w.eventChan <- event:
		default:
			w.logger.Warn("handle watch block", xlog.Int64("revision", w.revision), xlog.Any("kv", event.Kv))
		}
	}
}

// NewWatch ...
func (client *Client) NewWatch(prefix string) (*Watch, error) {
	var (
		ctx, cancel = context.WithCancel(context.Background())
		watcher     = &Watch{
			client:    client,
			revision:  0,
			cancel:    cancel,
			eventChan: make(chan *clientv3.Event, 100),
			lock:      &sync.RWMutex{},
			logger:    client.config.logger,
		}
	)

	go func() {
		rch := client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify())
		for {
			for resp := range rch {
				watcher.update(&resp)
			}

			time.Sleep(time.Duration(1) * time.Second)
			if watcher.revision > 0 {
				rch = client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify())
			} else {
				rch = client.Watch(ctx, prefix, clientv3.WithPrefix(), clientv3.WithCreatedNotify())
			}
		}
	}()

	return watcher, nil
}
