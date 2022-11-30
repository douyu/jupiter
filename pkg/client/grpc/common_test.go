package grpc

import (
	"bytes"
	"context"
	"net"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/registry/etcdv3"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/util/xtest/server/yell"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/douyu/jupiter/proto/testproto/v1"
	"google.golang.org/grpc"
)

var directClient testproto.GreeterServiceClient

var testconf = `
[jupiter.logger.jupiter]
	level = "debug"
	debug = true
[jupiter.registry.default]
    endpoints = ["http://localhost:2379"]
	timeout = "3s"
`

func init() {

	l, _ := startServer("127.0.0.1:9528", "srv1")
	time.Sleep(200 * time.Millisecond)

	startServer("127.0.0.1:9529", "srv1")

	time.Sleep(200 * time.Millisecond)
	err := conf.LoadFromReader(bytes.NewBufferString(testconf), toml.Unmarshal)
	if err != nil {
		panic(err)
	}

	cfg := DefaultConfig()
	cfg.Addr = l.Addr().String()

	conn := newGRPCClient(cfg)
	directClient = testproto.NewGreeterServiceClient(conn)
}

func startServer(addr, name string) (net.Listener, *grpc.Server) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		xlog.Panic("failed start server:" + err.Error())
	}

	xlog.Jupiter().Info("startServer", xlog.String("addr", addr), xlog.String("name", name))

	gserver := grpc.NewServer()
	grpcServer := &yell.FooServer{}
	grpcServer.SetName(name)

	testproto.RegisterGreeterServiceServer(gserver, grpcServer)
	go func() {
		if err := gserver.Serve(l); err != nil {
			panic("failed serve:" + err.Error())
		}
	}()

	reg := etcdv3.DefaultConfig().MustSingleton()
	reg.RegisterService(context.Background(), &server.ServiceInfo{
		Name:    name,
		Address: addr,
		Scheme:  "grpc",
	})
	return l, gserver
}
