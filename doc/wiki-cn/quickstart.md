# 快速开始

jupiter 提供了脚手架以便快速的创建应用:

1. 安装脚手架工具
```
go get -u github.com/douyu/jupiter/tools/jupiter
```
2. 创建一个项目
``` bash
jupiter new demo
```


## 创建Application

1. 创建一个Application

```golang
type MyApplication struct {
    jupiter.Application
}

func main() {
    app := &MyApplication{}
    app.Startup()
    app.Run()
}
```

2. 为Application添加一个http server

```golang
func(app *MyApplication) startHTTPServer() error {
    server := xecho.DefaultConfig().Build()
    server.GET("ping",func(ctx echo.Context) error {
        return ctx.JSON(200, "pong")
    })
    return app.Serve(server)
}

func main() {
    app := &MyApplication{}
    app.Startup(
        app.startHTTPServer,
    )
    app.Run()
}
```

2. 为Application添加一个grpc server


```golang
func(app *MyApplication) startGRPCServer() error {
    server := xgrpc.DefaultConfig().Build()
    helloworld.RegisterGreeterServer(server.Server, new(greeter.Greeter))
    returnapp.Serve(server)
}

func main() {
    app := &MyApplication{}
    app.Startup(
        app.startGRPCServer,
    )
    app.Run()
}
```

3. 为Application添加一个cron job

```golang
func(app *MyApplication) startCronJobs() error {
    config := xcron.DefaultConfig()
    // 如果要支持秒调度(默认为分钟调度)
    config.WithSeconds = true
    // 如果并发执行了，后一个任务延迟(-1：则跳过后一个任务)
    config.ConcurrentDelay = time.Second * 3
    // 立即执行
	config.ImmediatelyRun = true

    cron := config.Build()
	cron.AddFunc("*/1 * * * *", func() {
        log.Println("hello")
    })
	return app.Schedule(cron)
}

func main() {
    app := &MyApplication{}
    app.Startup(
        app.startCronJobs,
    )
    app.Run()
}
```

4. 为Application配置一个注册中心

```golang
func main() {
    app := &MyApplication{}
    app.Startup()
    app.SetRegistry(
        etcdv3.DefaultConfig().Build(),
    )
    app.Run()
}
```

如果要多注册:

``` golang
import (
    "github.com/douyu/jupiter/pkg/registry/etcdv3"
)
func main() {
    app := &MyApplication{}
    app.Startup()
    app.SetRegistry( // 多注册中心
		compound.New(
			etcdv3.StdConfig("bj01").Build(),  // 读取配置文件中 jupiter.etcdv3.bj01的配置，自动初始化一个etcd registry
			etcdv3.StdConfig("bj02").Build(),  // 读取配置文件中 jupiter.etcdv3.bj02的配置，自动初始化一个etcd registry
		),
	)
    app.Run()
}
```

5. 为Application配置一个治理入口
```golang
func main() {
    app := &MyApplication{}
    app.Startup()
    app.SetGovernor("127.0.0.1:9990") // 默认为127.0.0.1:9990
    app.Run()
}

// curl http://127.0.0.1:9091/all 获取所有治理路径
```

6. 创建一个grpc客户端

直连服务器:
```golang
func dialServer() {
    config := DefaultConfig()
    config = config.WithDialOption()
    config.Address = "127.0.0.1:9990"
    client := helloworld.NewGreeterClient(config.Build())
    rep, err := client.SayHello(context.Background(), &helloworld.HelloRequest{
        Name: "hi",
    })
}
```
通过etcd的注册服务器
```golang
// 1. 注册一个resolver
func init() {
	resolver.Register("etcd", etcdv3.StdConfig("wh").Build())
}
// 2. 通过resolver获取服务器地址
func dialServer() {
    config := grpc.DefaultConfig()
    config = config.WithDialOption()
    config.Address = "etcd:///server-name" // 注意必须填写etcd:///
    client := helloworld.NewGreeterClient(config.Build())
    rep, err := client.SayHello(context.Background(), &helloworld.HelloRequest{
        Name: "hi",
    })
}
```