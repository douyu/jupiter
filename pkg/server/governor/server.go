package governor

import (
	"context"
	"github.com/douyu/jupiter/pkg/constant"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/xlog"
	"net"
	"net/http"
)

// Server ...
type Server struct {
	*http.Server
	listener net.Listener
	*Config
}

func newServer(config *Config) *Server {
	var listener, err = net.Listen("tcp4", config.Address())
	if err != nil {
		xlog.Panic("start governor", xlog.FieldErr(err))
	}

	return &Server{
		Server: &http.Server{
			Addr:    config.Address(),
			Handler: DefaultServeMux,
		},
		listener: listener,
		Config:   config,
	}
}

func (s *Server) Serve() error {
	s.logger.Info("start governor", xlog.FieldAddr("http://"+s.listener.Addr().String()))
	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err

}

func (s *Server) Stop() error {
	return s.Server.Close()
}

func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

func (s *Server) Info() *server.ServiceInfo {
	info := server.ApplyOptions(
		server.WithScheme("http"),
		server.WithAddress(s.listener.Addr().String()),
		server.WithKind(constant.ServiceGovernor),
	)
	return &info
}
