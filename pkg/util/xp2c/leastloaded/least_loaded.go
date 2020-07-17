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

package leastloaded

import (
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/douyu/jupiter/pkg/util/xp2c"
	"google.golang.org/grpc/balancer"
)

type leastLoadedNode struct {
	item     interface{}
	inflight int64
}

type leastLoaded struct {
	items []*leastLoadedNode
	mu    sync.Mutex
	rand  *rand.Rand
}

func New() xp2c.P2c {
	return &leastLoaded{
		items: make([]*leastLoadedNode, 0),
		rand:  rand.New(rand.NewSource(time.Now().Unix())),
	}
}

func (p *leastLoaded) Add(item interface{}) {
	p.items = append(p.items, &leastLoadedNode{item: item})
}

func (p *leastLoaded) Next() (interface{}, func(balancer.DoneInfo)) {
	var sc, backsc *leastLoadedNode

	switch len(p.items) {
	case 0:
		return nil, func(balancer.DoneInfo) {}
	case 1:
		sc = p.items[0]
	default:
		// rand needs lock
		p.mu.Lock()
		a := p.rand.Intn(len(p.items))
		b := p.rand.Intn(len(p.items) - 1)
		p.mu.Unlock()

		if b >= a {
			b = b + 1
		}
		sc, backsc = p.items[a], p.items[b]

		// choose the least loaded item based on inflight
		if sc.inflight > backsc.inflight {
			sc, backsc = backsc, sc
		}
	}

	atomic.AddInt64(&sc.inflight, 1)

	return sc.item, func(balancer.DoneInfo) {
		atomic.AddInt64(&sc.inflight, -1)
	}
}
