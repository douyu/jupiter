# jupiter接入xxl-job分布式调度中心模块
example: /example/xxl-job
用法如下：
1. 增加配置:
以下是配置说明
```toml
[xxl]
  [xxl.job]
    [xxl.job.executor]
      address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名
      access_token = "jupiter-token"    # 注册xxl-job执行器需要的token信息
      appname = "jupiter-xxl-job-demo"  # 启动执行器的名称
      port = "59000"                    # 启动执行器的服务端口
      log_dir = "./"                    # 执行器产生的日志文件目录
	  # 以下的配置建议使用默认
	  #host = ""                        # 启动执行器的主机。默认通过ip.get注册。确保xxl-job能调度该地址
	  debug = true						# 是否开启debug模式
	  switch = true 					# 执行器的总开关
	  registry_group = "EXECUTOR"		# 执行器组
	  #timeout = 2000					# 接口超时时间
```

2. 示例用法:
```go
func main() {
	eng := NewEngine()
	eng.RegisterHooks(hooks.Stage_AfterStop, func() {
		fmt.Println("exit jupiter app ...")
	})
	if err := eng.Run(); err != nil {
		log.Fatal(err)
	}
}

type Engine struct {
	jupiter.Application
}

func NewEngine() *Engine {
	eng := &Engine{}
	if err := eng.Startup(
		xgo.ParallelWithError(
			eng.startXxlJob,
		),
	); err != nil {
		xlog.Panic("startup engine", xlog.Any("err", err))
	}
	return eng
}

func (eng *Engine) startXxlJob() error {
	executor := xxl.StdNewExecutor()
	// 注册定时任务
	executor.RegXJob(
		NewTest(),
		NewTest2(),
	)
	eng.Executor(executor)
	return nil
}


// =======以下为示例任务test.go=========
type Test struct{}

func NewTest() *Test {
	return &Test{}
}

// 任务名称
func (t *Test) GetJobName() string {
	return "test"
}

// xxl-job 分布式调度任务执行函数
func (t *Test) Run(ctx context.Context, param *executor.RunReq) (msg string, err error) {
	//使用xxl-logger日志即可在xxl-job平台上看到日志
	logger.Info(param.LogID, "start run...")
	logger.Info(param.LogID, fmt.Sprintf("job param is: %s", param.ExecutorParams))
	fmt.Println("test job has been executed")
	return "success", nil
}
```

```
注：调度器中的停止某项任务是停止调度接下来的任务。如果用中途中止某个正在进行中的定时任务，请使用ctx.Done()
```