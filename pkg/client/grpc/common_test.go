package grpc

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/douyu/jupiter/pkg/util/xtest/server/yell"
	"github.com/douyu/jupiter/proto/testproto"
	"google.golang.org/grpc"
)

var directClient testproto.GreeterClient

func TestMain(m *testing.M) {
	l, s := startServer("127.0.0.1:0", "srv1")
	time.Sleep(200 * time.Millisecond)

	cfg := DefaultConfig()
	cfg.Addr = l.Addr().String()

	conn := newGRPCClient(cfg)
	directClient = testproto.NewGreeterClient(conn)
	m.Run()
	s.Stop()
	os.Exit(0)
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
