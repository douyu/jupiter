package xxl

import (
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"sync"
	"sync/atomic"
	"time"

	"github.com/douyu/jupiter/pkg/executor"
	"github.com/douyu/jupiter/pkg/executor/xxl/logger"
	"github.com/douyu/jupiter/pkg/metric"
)

// 任务执行函数
type TaskFunc func(ctx context.Context, param *executor.RunReq) (msg string, err error)

// 任务执行完后通知回调函数
type CallbackFunc func(ctx context.Context, status int, msg string) error

type TaskResult struct {
	taskResultType int
	msg            string
	err            error
}

const (
	TaskResultTypeDone = iota
	TaskResultTypeFailed
	TaskResultTypeCancel
	TaskResultTypeTimeout
	TaskResultTypePanic
)

// 任务定义
type Task struct {
	Id        int64            // 任务id
	Name      string           // 任务名
	Param     *executor.RunReq // 参数
	StartTime int64            // 开始时间
	EndTime   int64            // 结束时间

	running int32 // 是否在运行状态
	timeout time.Duration
	lock    sync.RWMutex       // 任务锁
	cancel  context.CancelFunc // 任务结束句柄
	fn      TaskFunc           // 任务执行函数

	ctx context.Context // 异步任务缓存调用时的context
}

// 任务执行序
func (t *Task) Run(ctx context.Context, cb CallbackFunc) {
	t.Trace("开始任务")
	msg := "执行成功"
	var beg = time.Now()
	if t.fn != nil {
		// 创建任务取消句柄
		fmt.Println("设置running前")
		atomic.StoreInt32(&t.running, 1)
		t.lock.Lock()
		t.StartTime = time.Now().Unix()
		if t.timeout > 0 {
			ctx, t.cancel = context.WithTimeout(ctx, t.timeout)
		} else {
			ctx, t.cancel = context.WithCancel(ctx)
		}
		t.lock.Unlock()

		done := make(chan TaskResult)
		go func() {
			defer func() {
				// recover panic
				if err := recover(); err != nil {
					log.Println(t.Info()+" panic: ", err)
					debug.PrintStack() // 堆栈跟踪
					done <- TaskResult{TaskResultTypePanic, "panic", fmt.Errorf("%v", err)}
					close(done)
				}
			}()
			info, err := t.fn(ctx, t.GetParam())
			if err == nil {
				done <- TaskResult{TaskResultTypeDone, info, nil}
			} else {
				done <- TaskResult{TaskResultTypeFailed, info, err}
			}
			close(done)
		}()

		select {
		case r := <-done:
			atomic.StoreInt32(&t.running, 0)
			if r.err != nil {
				if r.taskResultType == TaskResultTypePanic {
					t.Trace("执行异常退出")
				} else {
					t.Trace("任务失败")
				}
				_ = cb(ctx, r.taskResultType, r.err.Error())
				metric.JobHandleCounter.Inc("xxl-job", t.Name, "over fail")
				return
			}
			if r.msg != "" {
				msg = r.msg
			}
			t.Trace("执行完成")
			_ = cb(ctx, TaskResultTypeDone, msg)
			metric.JobHandleCounter.Inc("xxl-job", t.Name, "over suc")
		case <-ctx.Done():
			atomic.StoreInt32(&t.running, 0)
			if d, ok := ctx.Deadline(); ok && time.Now().After(d) {
				t.Trace("任务超时")
				_ = cb(ctx, TaskResultTypeTimeout, "任务超时")
			} else {
				t.Trace("任务取消")
				_ = cb(ctx, TaskResultTypeCancel, "任务取消")
			}
			metric.JobHandleCounter.Inc("xxl-job", t.Name, "over cancel")
		}
		metric.ServerHandleHistogram.Observe(time.Since(beg).Seconds(), "xxl-job", t.Name, "client")
	}
}

// 取消任务
func (t *Task) Cancel() {
	t.lock.RLock()
	defer t.lock.RUnlock()
	if t.cancel != nil {
		t.cancel()
	}
}

// 更新任务timeout
func (t *Task) SetTimeout(timeout time.Duration) {
	t.lock.Lock()
	defer t.lock.Unlock()
	t.timeout = timeout
}

// 判断任务是否在运行中
func (t *Task) IsRunning() bool {
	running := atomic.LoadInt32(&t.running)
	return running == 1
}

// 获取参数
func (t *Task) GetParam() *executor.RunReq {
	t.lock.RLocker().Lock()
	defer t.lock.RLocker().Unlock()
	r := *t.Param
	return &r
}

// 获取ID
func (t *Task) GetId() int64 {
	t.lock.RLocker().Lock()
	defer t.lock.RLocker().Unlock()
	return t.Id
}

// 获取Name
func (t *Task) GetName() string {
	t.lock.RLocker().Lock()
	defer t.lock.RLocker().Unlock()
	return t.Name
}

// 获取context
func (t *Task) GetContext() context.Context {
	t.lock.RLocker().Lock()
	defer t.lock.RLocker().Unlock()
	return t.ctx
}

// 任务信息
func (t *Task) Info() string {
	if t == nil {
		log.Fatal("非法任务")
	}

	param := "nil"
	if t.Param != nil {
		param = t.GetParam().ExecutorParams
	}
	return fmt.Sprintf("任务ID[%v]任务名称[%v]参数：%v", t.GetId(), t.GetName(), param)
}

// 任务跟踪
func (t *Task) Trace(step string) {
	if t == nil {
		log.Println("非法任务")
	} else {
		logId := int64(0)
		p := t.GetParam()
		if p != nil {
			logId = p.LogID
		}
		logger.Info(logId, t.Info()+"; "+step)
	}
}
