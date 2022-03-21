package xxl

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/douyu/jupiter/pkg/executor"
	"github.com/douyu/jupiter/pkg/executor/xxl/logger"
	"github.com/douyu/jupiter/pkg/util/xdebug"
)

var MaxQueueSize = 1

type JobExecutor struct {
	opts    Options
	address string
	regList *taskList //注册任务列表
	runList *taskList //正在执行任务列表
	quit    chan os.Signal
}

// 创建执行器
func StdNewExecutor(opts ...Option) executor.Executor {
	return newJobExecutor(opts...)
}

// 创建Job执行器
func newJobExecutor(opts ...Option) *JobExecutor {
	// 加载默认配置
	options := DefaultConfig()
	// 加载参数式配置
	for _, o := range opts {
		o(options)
	}
	executor := options.Build()
	// 校验配置项，如果校验失败将会panic
	return executor
}

// 执行器启动
func (e *JobExecutor) Run() (err error) {
	// 总开关
	if !e.opts.Switch {
		return
	}
	// 初始化日志路径
	_ = logger.InitLogPath(e.opts.LogDir)
	// 创建路由器
	mux := http.NewServeMux()
	// 设置路由规则
	mux.HandleFunc("/run", e.handlerWithAuth(e.runTaskHandler))
	mux.HandleFunc("/kill", e.handlerWithAuth(e.killTaskHandler))
	mux.HandleFunc("/log", e.handlerWithAuth(e.taskLogHandler))
	mux.HandleFunc("/heartbeat", e.handlerWithAuth(e.heartBeatHandler))
	mux.HandleFunc("/idle", e.handlerWithAuth(e.idleHandler))
	// 创建服务器
	server := &http.Server{
		Addr:         e.address,
		WriteTimeout: time.Second * 3,
		Handler:      mux,
	}
	xdebug.PrettyKVWithPrefix("[Executor]", "start xxl-job golang executor server at", e.address)
	// 1.初始化系统消息hook，用于反注册
	e.quit = make(chan os.Signal)
	e.handlePending()
	// 2.启动服务器
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	// 3.向xxl-job任务中心注册任务
	go e.registry()
	return nil
}

// 注册执行器任务
func (e *JobExecutor) RegXJob(jobs ...executor.XJob) {
	for _, j := range jobs {
		e.regTask(j.GetJobName(), j.Run)
	}
}

// 注册任务
func (e *JobExecutor) regTask(pattern string, task TaskFunc) {
	var t = &Task{}
	t.fn = task
	e.regList.Set(pattern, t)
}

// 运行一个任务
func (e *JobExecutor) runTaskHandler(writer http.ResponseWriter, request *http.Request) {
	req, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	// 1.解析参数
	param := &executor.RunReq{}
	err := json.Unmarshal(req, &param)
	if err != nil {
		_, _ = writer.Write(returnCall(param, http.StatusInternalServerError, "params err"))
		log.Println("参数解析错误:" + string(req))
		return
	}
	logger.Info(param.LogID, "任务["+Int64ToStr(param.JobID)+":"+param.ExecutorHandler+"]接收到调度请求")
	if !e.regList.Exists(param.ExecutorHandler) {
		_, _ = writer.Write(returnCall(param, http.StatusInternalServerError, "Task not registered"))
		logger.Info(param.LogID, "任务["+Int64ToStr(param.JobID)+":"+param.ExecutorHandler+"]没有注册")
		return
	}
	//TODO:考虑使用sync.Pool建立一个context pool进行优化
	ctx := context.Background()
	ctx = context.WithValue(ctx, logger.DefaultLogIDKey, param.LogID)
	// 2.从注册任务列表重，获取任务基础信息并封装参数，context等信息
	task := e.regList.Get(param.ExecutorHandler)
	task.lock.Lock()
	timeout := time.Duration(param.ExecutorTimeout) * time.Second
	task.Id = param.JobID
	task.Name = param.ExecutorHandler
	task.Param = param
	task.timeout = timeout
	task.ctx = ctx
	task.lock.Unlock()
	// 如果已经有任务在运行，根据阻塞策略做不同处理
	if e.runList.Exists(Int64ToStr(task.GetId())) {
		switch param.ExecutorBlockStrategy {
		case CoverEarly:
			oldTask := e.runList.Get(Int64ToStr(task.GetId()))
			if oldTask != nil {
				// 移除进行中的任务
				oldTask.Cancel()
				e.runList.Del(Int64ToStr(task.GetId()))
			}
			e.runList.Set(Int64ToStr(task.GetId()), task)
			go task.Run(ctx, func(ctx context.Context, status int, msg string) error {
				e.callback(task, status, msg, false)
				return nil
			})
			_, _ = writer.Write(returnGeneral())
			return
		case SerialExecution:
			// 串行执行，将任务放到pending list中
			id := task.GetId()
			err := e.runList.EnqueuePending(Int64ToStr(id))
			if err != nil && err.Error() == OverLimit {
				_, _ = writer.Write(returnCall(param, http.StatusInternalServerError,
					fmt.Sprintf("There are %v tasks running", MaxQueueSize)))
				return
			}
			_, _ = writer.Write(returnGeneral())
			return
		case DiscardLaterNoAlarm:
			// 丢弃本次调度，正常返回
			_, _ = writer.Write(returnCall(param, http.StatusOK, "There are tasks running"))
			logger.Info(param.LogID, "任务["+Int64ToStr(param.JobID)+":"+param.ExecutorHandler+"]已经在运行了:"+param.ExecutorHandler)
			return
		case DiscardLater:
			// 丢弃本次调度，返回500，admin报警
			_, _ = writer.Write(returnCall(param, http.StatusInternalServerError, "There are tasks running"))
			logger.Info(param.LogID, "任务["+Int64ToStr(param.JobID)+":"+param.ExecutorHandler+"]已经在运行了:"+param.ExecutorHandler)
			return
		}
	}
	// 任务未开始，将任务加入调度列表
	e.runList.Set(Int64ToStr(task.GetId()), task)
	switch param.ExecutorBlockStrategy {
	case CoverEarly:
		// 立即执行
		go task.Run(ctx, func(ctx context.Context, status int, msg string) error {
			e.callback(task, status, msg, false)
			return nil
		})
	case DiscardLater, DiscardLaterNoAlarm:
		// 立即执行
		go task.Run(ctx, func(ctx context.Context, status int, msg string) error {
			// 执行完成删除 runlist 中对应的task
			e.callback(task, status, msg, true)
			return nil
		})
	case SerialExecution:
		// 放入队列，延迟执行
		task.lock.Lock()
		task.lock.Unlock()
		id := task.GetId()
		_ = e.runList.EnqueuePending(Int64ToStr(id))
	}
	_, _ = writer.Write(returnGeneral())
}

// 删除一个任务
func (e *JobExecutor) killTaskHandler(writer http.ResponseWriter, request *http.Request) {
	req, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	param := &killReq{}
	_ = json.Unmarshal(req, &param)
	if !e.runList.Exists(Int64ToStr(param.JobID)) {
		_, _ = writer.Write(returnKill(param, http.StatusInternalServerError))
		log.Println("任务[" + Int64ToStr(param.JobID) + "]没有运行")
		return
	}
	task := e.runList.Get(Int64ToStr(param.JobID))
	task.Cancel()
	e.runList.Del(Int64ToStr(param.JobID))
	_, _ = writer.Write(returnGeneral())
}

//任务日志
func (e *JobExecutor) taskLogHandler(writer http.ResponseWriter, request *http.Request) {
	data, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	req := &logReq{}
	_ = json.Unmarshal(data, &req)
	fmt.Println("查看日志请求：" + string(data))
	_, _ = writer.Write(returnLog(req, http.StatusOK))
}

func (e *JobExecutor) idleHandler(writer http.ResponseWriter, request *http.Request) {
	data, _ := ioutil.ReadAll(request.Body)
	defer request.Body.Close()
	req := &idleReq{}
	_ = json.Unmarshal(data, &req)

	// 执行队列为空，返回空闲
	t := e.runList.Get(Int64ToStr(req.JobID))
	if t == nil {
		_, _ = writer.Write(returnIdle(e.address, true))
		return
	}
	// 执行队列为空，返回空闲
	_, q := e.runList.GetPending(Int64ToStr(req.JobID))
	if len(q) == 0 && !t.IsRunning() {
		_, _ = writer.Write(returnIdle(e.address, true))
		return
	}
	_, _ = writer.Write(returnIdle(e.address, false))
}

// 心跳信息
func (e *JobExecutor) heartBeatHandler(writer http.ResponseWriter, request *http.Request) {
	_, _ = ioutil.ReadAll(request.Body)
	defer request.Body.Close()

	_, _ = writer.Write(returnHeatBeat(e.address))
}

//注册执行器到调度中心
func (e *JobExecutor) registry() {
	t := time.NewTimer(time.Second * 0) //初始立即执行
	defer t.Stop()
	req := &Registry{
		RegistryGroup: e.opts.RegistryGroup, // EXECUTOR  CONTAINER
		RegistryKey:   e.opts.RegistryKey,
		RegistryValue: "http://" + e.address,
	}
	param, err := json.Marshal(req)
	if err != nil {
		log.Fatal("执行器注册信息解析失败:" + err.Error())
	}
	for {
		<-t.C
		t.Reset(time.Second * time.Duration(20)) //20秒心跳防止过期
		func() {
			result, err := e.post("/api/registry", string(param))
			if err != nil {
				log.Println("执行器注册失败:" + err.Error())
				return
			}
			if result == nil {
				log.Println("执行器注册失败")
				return
			}
			defer result.Body.Close()
			body, err := ioutil.ReadAll(result.Body)
			if err != nil {
				log.Println("执行器注册失败:" + err.Error())
				return
			}
			res := &res{}
			_ = json.Unmarshal(body, &res)
			if result.StatusCode != http.StatusOK || res.Msg != nil {
				log.Println("执行器注册失败:" + string(body))
				return
			}
		}()
	}
}

//执行器注册摘除
func (e *JobExecutor) registryRemove() {
	t := time.NewTimer(time.Second * 0) //初始立即执行
	defer t.Stop()
	req := &Registry{
		RegistryGroup: "EXECUTOR",
		RegistryKey:   e.opts.RegistryKey,
		RegistryValue: "http" + e.address,
	}
	param, err := json.Marshal(req)
	if err != nil {
		log.Println("执行器摘除失败:" + err.Error())
		return
	}
	res, err := e.post("/api/registryRemove", string(param))
	if err != nil {
		log.Println("执行器摘除失败:" + err.Error())
		return
	}
	if res == nil {
		log.Println("执行器摘除失败: 无效返回")
		return
	}
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	xdebug.PrettyKVWithPrefix("[Executor]", "executor stop success: post /api/registryRemove", string(body))
}

// 执行器优雅地退出
func (e *JobExecutor) GracefulStop() {
	e.registryRemove()
}

// 执行器退出
func (e *JobExecutor) Stop() {
	e.registryRemove()
}

//回调任务列表
func (e *JobExecutor) callback(task *Task, code int, msg string, delete bool) {
	if delete {
		e.runList.Del(Int64ToStr(task.Id))
	}
	status := http.StatusOK
	if code != TaskResultTypeDone {
		status = http.StatusInternalServerError
	}
	resp, err := e.post("/api/callback", string(returnCall(task.GetParam(), int64(status), msg)))
	if err != nil {
		fmt.Println(err)
		return
	}
	if resp == nil {
		fmt.Println("任务回调失败")
		return
	}
	defer resp.Body.Close()
	fmt.Println("任务回调成功")
}

// 回调请求
func (e *JobExecutor) post(action, body string) (resp *http.Response, err error) {
	url := e.opts.ServerAddr + action

	request, err := http.NewRequest("POST", url, strings.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json;charset=UTF-8")
	request.Header.Set("XXL-JOB-ACCESS-TOKEN", e.opts.AccessToken)
	client := http.Client{
		Timeout: e.opts.Timeout,
	}
	return client.Do(request)
}

// 处理串行消息
func (e *JobExecutor) handlePending() {
	go func() {
		// 任务主循环
		// 这里使用time.NewTicker而非time.Tick. https://stackoverflow.com/questions/38856959/go-time-tick-vs-time-newticker
		ticker := time.NewTicker(10 * time.Millisecond)
		for {
			select {
			case <-ticker.C:
				// 10ms 主动检查一次队列兜底
				e.runList.KillIdleTasks()
				hasMore := e.runPendingTasks()
				if !hasMore {
					return
				}
				runtime.Gosched()
			case <-e.quit:
				return
			}
		}
	}()
}

// 执行所有的pending task
// TODO：每次检查所有任务队列，存在惊群效应，可以优化
func (e *JobExecutor) runPendingTasks() bool {
	// 获取所有pending的任务的队列
	tasks := e.runList.GetAllPending()
	for _, task := range tasks {
		// 已经有任务在运行, 则跳过
		if task.IsRunning() {
			continue
		}
		// non-blocking polling
		select {
		case t := <-task.pending:
			if t != nil {
				ctx := t.GetContext()
				go task.Run(ctx, func(ctx context.Context, status int, msg string) error {
					e.callback(t, status, msg, false)
					return nil
				})
			}
		case <-e.quit:
			return false
		default:
			continue
		}
	}
	return true
}

// 中间件
type HttpHandler func(http.ResponseWriter, *http.Request)

// 鉴权中间件
func (e *JobExecutor) handlerWithAuth(next HttpHandler) HttpHandler {
	return func(writer http.ResponseWriter, request *http.Request) {
		if !e.verifyAuth(request) {
			_, _ = writer.Write(returnAuthError())
			return
		}
		next(writer, request)
	}
}

// 鉴权函数
func (e *JobExecutor) verifyAuth(request *http.Request) bool {
	if request == nil || request.Header == nil {
		return false
	}
	token := request.Header.Get("XXL-JOB-ACCESS-TOKEN")
	return token == e.opts.AccessToken
}

// GetAddress
func (e *JobExecutor) GetAddress() string {
	return e.address
}
