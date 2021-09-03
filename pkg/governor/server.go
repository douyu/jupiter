package governor

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/douyu/jupiter/pkg/component"
	"github.com/douyu/jupiter/pkg/xlog"
)

// Server ...
type Server struct {
	component.BaseComponent
	*http.Server
	listener net.Listener
	*Config
}

func newServer(config *Config) *Server {
	var listener, err = net.Listen("tcp4", config.Address())
	if err != nil {
		xlog.Panic("governor start error", xlog.FieldErr(err))
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

//Serve ..
func (s *Server) Start(stopCh <-chan struct{}) error {
	var errCh = make(chan error)
	go func() {
		fmt.Println("start governor")
		errCh <- s.Server.Serve(s.listener)
		fmt.Println("stop governor")
	}()
	go func() {
		select {
		case <-stopCh:
			fmt.Println("stop...")
			s.Shutdown(context.Background())
		case <-errCh:
			fmt.Println("err occur...")
		}
		fmt.Println("close err ch")
		close(errCh)
	}()
	return nil
}
