package xgf

import (
	"context"
	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/constant"
	"strconv"

	//"github.com/douyu/jupiter/pkg/ecode"
	//"github.com/douyu/jupiter/pkg/xlog"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"net"
)

type Server struct {
	*ghttp.Server
	config *Config
	listener net.Listener
}

func newServer(config *Config) *Server {
	//listener, err := net.Listen("tcp", config.Address())
	//
	//if err != nil {
	//	config.logger.Panic("new xgf server err", xlog.FieldErrKind(ecode.ErrKindListenErr), xlog.FieldErr(err))
	//}
	//config.Port = listener.Addr().(*net.TCPAddr).Port

	s := new(Server)
	s.Server = g.Server()

	s.SetPort(config.Port)

	s.config = config
	//s.listener = listener

	return s
}

func (s *Server)Serve() error  {
	s.Run()
	return nil
}

func (s *Server)Stop() error  {
	return s.Shutdown()
}

func (s *Server)GracefulStop(ctx context.Context) error  {
	return s.Stop()
}

func (s *Server)Info() *server.ServiceInfo  {
	return &server.ServiceInfo{
		Name:      pkg.Name(),
		Scheme:    "http",
		Address:   s.config.Host+":"+strconv.Itoa(s.config.Port),
		Weight:    0.0,
		Enable:    false,
		Healthy:   false,
		Metadata:  map[string]string{},
		Region:    "",
		Zone:      "",
		GroupName: "",
		Kind:      constant.ServiceProvider,
	}
}
