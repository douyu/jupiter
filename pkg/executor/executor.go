package executor

import (
	"context"
	"net/http"
	"sync"

	"github.com/douyu/jupiter/pkg/util/xdebug"
	"github.com/douyu/jupiter/pkg/util/xgo"
	"github.com/douyu/jupiter/pkg/xlog"
)

var _instances = sync.Map{}

type HttpHandler func(http.ResponseWriter, *http.Request)

type Executor interface {
	RegXJob(jobs ...XJob) // 注册带参数的XJob
	Run() error           // 执行器启动
	GetAddress() string   // 获取执行器服务地址信息
	GracefulStop()        // 执行器优雅退出，向调度中心取消注册
	Stop()                // 执行器退出，向调度中心取消注册
}

// XJob 定时任务接口
type XJob interface {
	Run(ctx context.Context, param *RunReq) (msg string, err error)
	GetJobName() string
}

// 触发任务请求参数
type RunReq struct {
	JobID                 int64  `json:"jobId"`                 // 任务ID
	ExecutorHandler       string `json:"executorHandler"`       // 任务标识
	ExecutorParams        string `json:"executorParams"`        // 任务参数
	ExecutorBlockStrategy string `json:"executorBlockStrategy"` // 任务阻塞策略
	ExecutorTimeout       int64  `json:"executorTimeout"`       // 任务超时时间，单位秒，大于零时生效
	LogID                 int64  `json:"logId"`                 // 本次调度日志ID
	LogDateTime           int64  `json:"logDateTime"`           // 本次调度日志时间
	GlueType              string `json:"glueType"`              // 任务模式，可选值参考 com.xxl.job.core.glue.GlueTypeEnum
	GlueSource            string `json:"glueSource"`            // GLUE脚本代码
	GlueUpdateTime        int64  `json:"glueUpdatetime"`        // GLUE脚本更新时间，用于判定脚本是否变更以及是否需要刷新
	BroadcastIndex        int64  `json:"broadcastIndex"`        // 分片参数：当前分片
	BroadcastTotal        int64  `json:"broadcastTotal"`        // 分片参数：总分片
}

// 注册执行器到jupiter
func Register(address string, e Executor) {
	if _, ok := _instances.Load(address); ok {
		panic("duplicate executor address " + address)
	}
	xdebug.PrettyKVWithPrefix("[jupiter]", "add executor for", address)
	_instances.Store(address, e)
}

// Run ...
func Run() error {
	var executors = make([]func(), 0)
	_instances.Range(func(key, val interface{}) bool {
		address := key.(string)
		if executor, ok := val.(Executor); ok {
			executors = append(executors, func() {
				xlog.Info("xxl-job executor run begin for ", xlog.FieldName(address))
				defer xlog.Info("xxl-job executor run exit for ", xlog.FieldName(address))
				_ = executor.Run()
			})
		}
		return true
	})

	xgo.Parallel(executors...)()
	return nil
}

// Stop ...
func Stop() error {
	var executors = make([]func(), 0)
	_instances.Range(func(key, val interface{}) bool {
		address := key.(string)
		if stopper, ok := val.(interface{ Stop() }); ok {
			executors = append(executors, func() {
				xlog.Info("xxl-job executor stop for ", xlog.FieldName(address))
				defer xlog.Info("xxl-job executor exit for ", xlog.FieldName(address))
				stopper.Stop()
			})
		}
		return true
	})

	xgo.Parallel(executors...)()
	return nil
}

// GracefulStop ...
func GracefulStop() error {
	var executors = make([]func(), 0)
	_instances.Range(func(key, val interface{}) bool {
		address := key.(string)
		if stopper, ok := val.(interface{ GracefulStop() }); ok {
			executors = append(executors, func() {
				xlog.Info("xxl-job executor stop begin for ", xlog.FieldName(address))
				defer xlog.Info("xxl-job executor stop exit for", xlog.FieldName(address))
				stopper.GracefulStop()
			})
		}
		return true
	})

	xgo.Parallel(executors...)()
	return nil
}
