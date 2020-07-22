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

type jobLock struct {
	client     *etcdv3.Client
	kv         clientv3.KV
	lease      clientv3.Lease
	cancelFunc context.CancelFunc // 用于终止自动续租
	isLocked   bool               // 是否上锁成功
	leaseId    clientv3.LeaseID
}

func newETCDXcron(config *Config) {
	if config.logger == nil {
		config.logger = xlog.DefaultLogger
	}
	config.logger = config.logger.With(xlog.FieldMod(ecode.ModXcronETCD), xlog.FieldAddrAny(config.Config.Endpoints))
	config.jLock = &jobLock{
		client: config.Config.Build(),
	}

	config.jLock.kv = clientv3.NewKV(config.jLock.client.Client)
	config.jLock.lease = clientv3.NewLease(config.jLock.client.Client)

	return
}

// 尝试上锁
func (jLock *jobLock) TryLock(lockKey string) (err error) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		leaseId        clientv3.LeaseID
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		txnResp        *clientv3.TxnResponse
	)

	if leaseGrantResp, err = jLock.lease.Grant(context.TODO(), 5); err != nil {
		return
	}

	cancelCtx, cancelFunc = context.WithCancel(context.TODO())

	leaseId = leaseGrantResp.ID

	if keepRespChan, err = jLock.lease.KeepAlive(cancelCtx, leaseId); err != nil {
		goto FAIL
	}

	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		for {
			select {
			case keepResp = <-keepRespChan: // 自动续租的应答
				if keepResp == nil {
					goto END
				}
			}
		}
	END:
	}()

	txn = jLock.kv.Txn(context.TODO())

	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))

	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	if !txnResp.Succeeded { // 锁被占用
		err = errors.New("lock is occupied")
		goto FAIL
	}

	jLock.leaseId = leaseId
	jLock.cancelFunc = cancelFunc
	jLock.isLocked = true
	return

FAIL:
	cancelFunc()
	jLock.lease.Revoke(context.TODO(), leaseId)
	return
}

func (jLock *jobLock) Unlock() {
	if jLock.isLocked {
		jLock.cancelFunc()
		jLock.lease.Revoke(context.TODO(), jLock.leaseId)
	}
}
