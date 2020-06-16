# jupiter

## 脚手架介绍
### jupiter脚手架工具集
1. 快速生成模板项目
2. 基于proto文件生成pb.go
3. 基于proto文件生成服务端实现

# go version
 GO >= 1.13

## 脚手架工具获取
windows :  
```shell script
set GO111MODULE=on
set GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```
linux :  
```shell script
export GO111MODULE=on
export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
```
`go get -u -v github.com/douyu/jupiter/tools/jupiter`
* windows 用户:  
  会在${GOPATH}/src/bin 目录下生成jupiter.exe 文件,若想方便的在任何地方使用jupiter命令,请将该
  命令配置在系统的环境变量中
* Linux 用户:  
* Mac 用户 :


## 如何使用
* jupiter -h  
```shell script
E:\go\goworkspace\src\jupiter-demo\cmd>jupiter -h
NAME:
   jupiter - jupiter tools

USAGE:
   jupiter [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
   new, n     Create Jupiter template project
   protoc, p  jupiter protoc tools
   help, h    Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```
* jupiter new -h

```shell script
E:\go\goworkspace\src\jupiter-demo\cmd>jupiter new -h
NAME:
   jupiter new - Create Jupiter template project

USAGE:

jupiter [commands|flags]

The commands & flags are:
  new     Create Jupiter template project

  -d      Build the specified directory for the template project

Examples:
  # Build the specified directory for the template project
  jupiter new (your project name) -d (project dir)
```
* jupiter protoc -h  
```shell script
E:\go\goworkspace\src\jupiter-demo\cmd>jupiter protoc -h
NAME:
   jupiter protoc - jupiter protoc tools

USAGE:

jupiter [commands|flags]

The commands & flags are:
  protoc        jupiter protoc tools
  -g,--grpc     whether to generate GRPC code
  -s,--server   whether to generate grpc server code
  -f,--file     path of proto file
  -o,--out      path of code generation
  -p,--prefix   prefix(current project name)
Examples:
   # Generate GRPC code from the Proto file
   # -f: Proto file address -o: Code generation path -g: Whether to generate GRPC code
   jupiter protoc -f ./pb/hello/hello.proto -o ./pb/hello -g
   # According to the proto file, generate the server implementation
   # -f: Proto file address -o: Code generation path -p:prefix(Current project name) -g: Whether to generate Server code
   jupiter protoc -f ./pb/hello/hello.proto -o ./internal/app/grpc -p jupiter-demo -s
```
## 开始实战 
 接下来我们会一步一步的带着大家从无到有开发jupiter应用!(gopher Let's go)
### 快速创建jupiter模板项目
```shell script
cd ${GOPATH}/src
jupiter new jupiter-demo -d ./
```
命令解释: 
* new :创建jupiter模板项目
* jupiter-demo: 项目名称
* -d: 生成项目所在路径

然后就会在${GOPATH}/src 下生成jupiter-demo 项目  
项目目录结构:
```go
build                           编译目录
cmd                             应用启动目录
config                          应用配置目录
internal
├─app                           应用目录
│  ├─engine                     
│  │  ├─engine.go               核心编排引擎(启动HTTP,GRPC,JOB等服务)
│  ├─grpc                       grpc服务实现目录
│  ├─handler                    控制器目录（接收用户请求）              
│  │  ├─user.go                 控制器文件
│  ├─model                      model目录（定义持久层结构体）
│  │  ├─db
│  │  │  ├─user.go
│  │  ├─init.go                 初始化全局数据库句柄
│  ├─service                    service层
│  │  ├─user                    模块
│  │  │  ├─impl  
│  │  │  │  ├─mysqlImpl.go      实现
│  │  │  ├─repository.go        service 接口
│  │  ├─init.go
pb                              proto文件
sql                             sql脚本
.gitignore
go.mod
Makefile
```
### 运行jupiter应用
* 数据库环境准备
1. 请先在mysql中创建test库，然后将sql中的user.sql在该库中执行
2. 修改mysql配置: 在config/config.toml中 找到 `[jupiter.mysql.test]`
   将`dsn` 改成你所在的环境
* 下载go.mod依赖
*  编译运行  
   在项目跟目录下执行以下命令
```shell script
cd cmd
go run main.go --config=../config/config.toml
```
   接下来你将会看到以下日志  
```shell script
E:\go\goworkspace\src\jupiter-demo\cmd>go run main.go --config=../config/config.toml

   (_)_   _ _ __ (_) |_ ___ _ __
   | | | | | '_ \| | __/ _ \ '__|
   | | |_| | |_) | | ||  __/ |
  _/ |\__,_| .__/|_|\__\___|_|
 |__/      |_|

 Welcome to jupiter, starting application ...

1592274902      INFO    load local file                         {"mod": "config", "addr": "../config/config.toml"}
1592274902      INFO    auto max procs                          {"mod": "proc", "procs": 4}
1592274902      INFO    set global tracer                       {"mod": "trace"}
1592274902      INFO    add job                                 {"mod": "worker.cron", "name": "jupiter-demo/internal/app/engine.(*Engine).execJob-fm"}

?[33m[2020-06-16 10:35:02]?[0m ?[35m[info] replacing callback `gorm:delete` from E:/go/goworkspace/src/jupiter-demo/vendor/github.com/douyu/jupiter/pkg/store/gorm/orm.go:118?[0m ?[31;1m
?[0m

?[33m[2020-06-16 10:35:02]?[0m ?[35m[info] replacing callback `gorm:update` from E:/go/goworkspace/src/jupiter-demo/vendor/github.com/douyu/jupiter/pkg/store/gorm/orm.go:118?[0m ?[31;1m
?[0m

?[33m[2020-06-16 10:35:02]?[0m ?[35m[info] replacing callback `gorm:create` from E:/go/goworkspace/src/jupiter-demo/vendor/github.com/douyu/jupiter/pkg/store/gorm/orm.go:118?[0m ?[31;1m
?[0m

?[33m[2020-06-16 10:35:02]?[0m ?[35m[info] replacing callback `gorm:query` from E:/go/goworkspace/src/jupiter-demo/vendor/github.com/douyu/jupiter/pkg/store/gorm/orm.go:118?[0m ?[31;1m ?
[0m

?[33m[2020-06-16 10:35:02]?[0m ?[35m[info] replacing callback `gorm:row_query` from E:/go/goworkspace/src/jupiter-demo/vendor/github.com/douyu/jupiter/pkg/store/gorm/orm.go:118?[0m ?[31;
1m ?[0m
1592274902      INFO    client mysql start                      {"mod": "gorm", "addr": "127.0.0.1:3306", "name": "test"}
1592274902      INFO    run worker                              {"mod": "worker.cron", "number of scheduled jobs": 1}
1592274902      INFO    cron start                              {"mod": "worker.cron"}
1592274902      INFO    cron schedule                           {"mod": "worker.cron", "now": 1592274902, "entry": 1, "next": 1592274902}
1592274902      INFO    cron wake                               {"mod": "worker.cron", "now": 1592274902}
1592274902      INFO    cron run                                {"mod": "worker.cron", "now": 1592274902, "entry": 1, "next": 1592274912}
1592274902      INFO    start servers                           {"mod": "app", "addr": "http://127.0.0.1:20105"}
1592274902      INFO    start servers                           {"mod": "app", "addr": "grpc://127.0.0.1:20102"}
1592274902      INFO    start governor                          {"mod": "app", "addr": "http://127.0.0.1:9990"}
1592274902      INFO    exec job                                {"info": "print info"}
1592274902      WARN    exec job                                {"warn": "print warning"}
1592274902      INFO    echo add route                          {"mod": "server.echo", "method": "GET", "path": "/jupiter"}
1592274902      INFO    echo add route                          {"mod": "server.echo", "method": "GET", "path": "/api/user/:id"}
1592274902      INFO    echo add route                          {"mod": "server.echo", "method": "GET", "path": "/grpc/get"}
1592274902      INFO    echo add route                          {"mod": "server.echo", "method": "POST", "path": "/grpc/post"}
1592274902      INFO    run job                                 {"mod": "worker.cron", "name": "jupiter-demo/internal/app/engine.(*Engine).execJob-fm"}
⇨ http server started on 127.0.0.1:20105
1592274912      INFO    cron wake                               {"mod": "worker.cron", "now": 1592274912}
1592274912      INFO    cron run                                {"mod": "worker.cron", "now": 1592274912, "entry": 1, "next": 1592274922}
1592274912      INFO    exec job                                {"info": "print info"}
1592274912      WARN    exec job                                {"warn": "print warning"}
1592274912      INFO    run job                                 {"mod": "worker.cron", "name": "jupiter-demo/internal/app/engine.(*Engine).execJob-fm"}
1592274922      INFO    cron wake                               {"mod": "worker.cron", "now": 1592274922}
1592274922      INFO    cron run                                {"mod": "worker.cron", "now": 1592274922, "entry": 1, "next": 1592274932}
1592274922      INFO    exec job                                {"info": "print info"}
1592274922      WARN    exec job                                {"warn": "print warning"}
1592274922      INFO    run job                                 {"mod": "worker.cron", "name": "jupiter-demo/internal/app/engine.(*Engine).execJob-fm"}
1592274932      INFO    cron wake                               {"mod": "worker.cron", "now": 1592274932}
1592274932      INFO    cron run                                {"mod": "worker.cron", "now": 1592274932, "entry": 1, "next": 1592274942}
1592274932      INFO    exec job                                {"info": "print info"}
1592274932      WARN    exec job                                {"warn": "print warning"}
1592274932      INFO    run job                                 {"mod": "worker.cron", "name": "jupiter-demo/internal/app/engine.(*Engine).execJob-fm"}
```
模板项目中 我们默认帮你开启了HTTP,GRPC,JOB等服务。当然你也可以根据自己的需求进行取舍  
接下来打开浏览器输入:  
http://127.0.0.1:20105/jupiter  
你将会看到  `welcome to jupiter`

在开始运行项目前我们不是配置过数据库环境么? 接下来我们可以验证查询是否正常:  
打开浏览器输入: http://127.0.0.1:20105/api/user/1  
你将会看到  `{"id":1,"username":"admin","password":"123456","nickname":"rose","address":"WUHAN"}`

如果以上操作你都正常完成，恭喜你，你已经完成了基础操作，接下来我们将介绍脚手架工具的 protoc模块

### protoc 工具集介绍

jupiter protoc工具集目前可用于根据pb文件中的proto文件，生成pb.go，以及生成服务端实现，
你只需要去完善对应服务端实现的业务逻辑即可。
* 环境准备
  请确保当前系统安装protoc 编译工具  
  go版本的 Protobuf 编译器插件,工具内部采用的是性能更好的 protoc-gen-gofast ,若检测到当前系统安装会自动进行安装
* 实操演示  
脚手架中默认提供的 `pb/hello/hello.proto` 文件,以下操作都是基于该文件  
进入到项目根目录下，执行以下命令   
1. 生成pb.go 文件
```shell script
 jupiter protoc -f ./pb/hello/hello.proto -o ./pb/hello -g
```
命令解释:  
-f : proto文件路径  
-o :生成文件目录(建议生成目录与proto文件同级，方便管理)  
-g :是否生成pb.go 文件   

正常情况下会在 `./pb/hello/`中 生成 `hello.pb.go` 文件

2. 生成grpc服务端实现
```
jupiter protoc -f ./pb/hello/hello.proto -o ./internal/app/grpc -p jupiter-demo -s
```
命令解释:    
-f : proto文件路径  
-o : grpc服务端代码生成文件目录(建议生成至./internal/app/grpc 中)  
-p : grpc服务端代码实现依赖pb.go的前缀(请使用当前项目名)  
-s : 是否生成grpc服务端实现代码  

正常情况下会在 `./internal/app/grpc/hello` 目录中 生成 `helloServiceServer.go` 文件  
请注意: 我们在grpc包中会默认以proto的package的包名来管理 grpc服务端实现代码,从而让代码结构更加清晰  
接下来 你只需要去该文件中完成对应的业务逻辑实现即可  

快速测试grpc服务端代码:  
框架内部我们支持 代理http到grpc控制器,并且在engine.go中我们也提供了相关的实例
```go
	//support proxy for http to grpc controller
	g := greeter.Greeter{}
	group2 := server.Group("/grpc")
	group2.GET("/get", xecho.GRPCProxyWrapper(g.SayHello))
	group2.POST("/post", xecho.GRPCProxyWrapper(g.SayHello))
```
你可以参照该示例,将自己的grpc服务端实现注册在路由中，然后通过HTTP请求即可访问 


