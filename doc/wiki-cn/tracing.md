# 链路追踪

得益于统一的拦截器架构，jupiter为主要模块内置了链路追踪。本文通过一个试验，介绍jupiter的链路追踪.


### 启动jaeger集群

docker环境下，可以很容器启动一个jaeger集群:

```
docker run -d --name jaeger   -e COLLECTOR_ZIPKIN_HTTP_PORT=9411   -p 5775:5775/udp   -p 6831:6831/udp   -p 6832:6832/udp   -p 5778:5778   -p 16686:16686   -p 14268:14268   -p 9411:9411   jaegertracing/all-in-one:1.6
```

### 启动服务

jupiter内置了一个简单的demo，可以直接启动:

```
cd github.com/douyu/jupiter
make demo
```

### 启动客户端

```
cd github.com/douyu/jupiter/example/grpc/direct/direct-client
go run main.go --config=config.toml
```

### 查看链路

打开浏览器，输入:
```
http://127.0.0.1:16686
```
即可看到链路信息


## 状态

目前，支持http-server/grpc-server/grpc-client/gorm四种链路。未来，还将释放http-client/redigo/bigcache等链路采集。

> gorm链路，请参考: github.com/douyu/jupiter/example/client/gorm
