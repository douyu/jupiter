package governor

import (
	"context"
	"net"
	"net/http"

	"github.com/douyu/jupiter/pkg/core/constant"
	"github.com/douyu/jupiter/pkg/server"
	"github.com/douyu/jupiter/pkg/xlog"
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
		xlog.Jupiter().Panic("governor start error", xlog.FieldErr(err))
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

// Serve ..
func (s *Server) Serve() error {
	err := s.Server.Serve(s.listener)
	if err == http.ErrServerClosed {
		return nil
	}
	return err

}

// Stop ..
func (s *Server) Stop() error {
	return s.Server.Close()
}

// GracefulStop ..
func (s *Server) GracefulStop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}

// Healthz
// TODO(roamerlv):
func (s *Server) Healthz() bool {
	return true
}

// Info ..
func (s *Server) Info() *server.ServiceInfo {
	serviceAddr := s.listener.Addr().String()
	if s.Config.ServiceAddress != "" {
		serviceAddr = s.Config.ServiceAddress
	}

	info := server.ApplyOptions(
		server.WithScheme("govern"),
		server.WithAddress(serviceAddr),
		server.WithKind(constant.ServiceGovernor),
	)
	// info.Name = info.Name + "." + ModName
	return &info
}
