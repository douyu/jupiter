package grpc

import (
	"github.com/douyu/jupiter/pkg/util/xtest/proto/testproto"
	"github.com/douyu/jupiter/pkg/util/xtest/server/yell"
	"google.golang.org/grpc"
	"net"
	"testing"
	"time"
)

var directClient testproto.GreeterClient

func TestMain(m *testing.M) {
	l, s := startServer("127.0.0.1:0", "srv1")
	time.Sleep(200 * time.Millisecond)

	cfg := DefaultConfig()
	cfg.Address = l.Addr().String()

	conn := newGRPCClient(cfg)
	directClient = testproto.NewGreeterClient(conn)
	m.Run()
	s.Stop()
}

func startServer(addr, name string) (net.Listener, *grpc.Server) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		panic("failed start server:" + err.Error())
	}
	server := grpc.NewServer()
	grpcServer := &yell.FooServer{}
	grpcServer.SetName(name)
	testproto.RegisterGreeterServer(server, grpcServer)
	go func() {
		if err := server.Serve(l); err != nil {
			panic("failed serve:" + err.Error())
		}
	}()
	return l, server
}
