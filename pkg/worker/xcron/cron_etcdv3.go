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

package xcron

import (
	"context"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/douyu/jupiter/pkg/client/etcdv3"
	"github.com/douyu/jupiter/pkg/ecode"
	"github.com/douyu/jupiter/pkg/xlog"
)

const (
	// 任务锁
	WorkerLockDir = "/xcron/lock/"
)

type jobLock struct {
	client     *etcdv3.Client
	cancelFunc context.CancelFunc
	isLocked   bool // 是否上锁成功
}

func newETCDXcron(config *Config) {
	if config.logger == nil {
		config.logger = xlog.DefaultLogger
	}
	config.logger = config.logger.With(xlog.FieldMod(ecode.ModXcronETCD), xlog.FieldAddrAny(config.Config.Endpoints))
	config.jLock = &jobLock{
		client: config.Config.Build(),
	}

	return
}

// 尝试上锁
func (jLock *jobLock) TryLock(ctx context.Context, lockKey string) (err error) {
	var (
		txn     clientv3.Txn
		txnResp *clientv3.TxnResponse
	)

	_, jLock.cancelFunc = context.WithCancel(ctx)

	leaseID, err := jLock.client.GetLease(ctx)
	if err != nil {
		goto FAIL
	}

	txn = jLock.client.KV.Txn(ctx)

	// add lock prefix
	lockKey = WorkerLockDir + lockKey

	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseID))).
		Else(clientv3.OpGet(lockKey))

	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	if !txnResp.Succeeded { // 锁被占用
		err = errors.New("lock is occupied")
		goto FAIL
	}

	jLock.isLocked = true
	return

FAIL:
	jLock.cancelFunc()
	return
}

func (jLock *jobLock) Unlock() {
	if jLock.isLocked {
		jLock.cancelFunc()
	}
}
