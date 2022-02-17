// Copyright 2021 rex lv
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

package pgelector

import (
	"context"
	"sync/atomic"
	"time"

	"cirello.io/pglock"
	"github.com/douyu/jupiter/pkg/elect"
	"github.com/douyu/jupiter/pkg/xlog"
)

var _logger = xlog.DefaultLogger.With(xlog.FieldMod("pgelector"))

// postgresLeaderElector implements leader election using PostgreSQL DB.
// pglock does not rely on timestamps, which eliminates the problem of clock skews, but the cost is that first leader election can happen only after lease duration
// pglock does optimistic locking under the hood, the alternative would be to use pg_advisory_lock
type postgresLeaderElector struct {
	leader      int32
	lockClient  *pglock.Client
	callbacks   []elect.LeaderElectCallback
	lockName    string
	backoffTime time.Duration
}

var _ elect.LeaderElector = &postgresLeaderElector{}

func New(lockClient *pglock.Client, lockName string, backoffTime time.Duration) elect.LeaderElector {
	return &postgresLeaderElector{
		lockClient:  lockClient,
		lockName:    lockName,
		backoffTime: backoffTime,
	}
}

func (p *postgresLeaderElector) Start(stop <-chan struct{}) {
	_logger.Info("starting Leader Elector")
	ctx, cancelFn := context.WithCancel(context.Background())
	go func() {
		<-stop
		_logger.Info("stopping Leader Elector")
		cancelFn()
	}()

	for {
		_logger.Info("waiting for lock")
		if err := p.lockClient.Do(ctx, p.lockName, func(ctx context.Context, lock *pglock.Lock) error {
			p.leaderAcquired()
			<-ctx.Done()
			p.leaderLost()
			return nil
		}); err != nil {
			_logger.Errorw(err.Error(), "error waiting for lock")
		}

		select {
		case <-stop:
			break
		default:
		}

		time.Sleep(p.backoffTime)
	}
	// _logger.Info("Leader Elector stopped")
}

func (p *postgresLeaderElector) leaderAcquired() {
	p.setLeader(true)
	for _, callback := range p.callbacks {
		callback(elect.CallbackPhasePostStarted)
	}
}

func (p *postgresLeaderElector) leaderLost() {
	p.setLeader(false)
	for _, callback := range p.callbacks {
		callback(elect.CallbackPhasePostStopped)
	}
}

func (p *postgresLeaderElector) AddCallbacks(callbacks ...elect.LeaderElectCallback) {
	p.callbacks = append(p.callbacks, callbacks...)
}

func (p *postgresLeaderElector) setLeader(leader bool) {
	var value int32 = 0
	if leader {
		value = 1
	}
	atomic.StoreInt32(&p.leader, value)
}

func (p *postgresLeaderElector) IsLeader() bool {
	return atomic.LoadInt32(&(p.leader)) == 1
}
