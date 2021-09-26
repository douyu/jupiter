package governor

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/douyu/jupiter/pkg/component"
	"github.com/douyu/jupiter/pkg/xlog"
)

type Governor struct {
	component.BaseComponent
	mux  *http.ServeMux
	addr string
}

func New(addr string) *Governor {
	return &Governor{
		mux:  http.NewServeMux(),
		addr: addr,
	}
}

func (g *Governor) HandleFunc(pattern string, handler http.HandlerFunc) {
	g.mux.HandleFunc(pattern, handler)
}

func (g *Governor) PrintRoutes() {}

func (g *Governor) Start(stopCh <-chan struct{}) error {
	var listener, err = net.Listen("tcp4", g.addr)
	if err != nil {
		xlog.Panic("governor start error", xlog.FieldErr(err))
	}
	var server = &http.Server{Addr: g.addr, Handler: g.mux}
	var errCh = make(chan error)
	go func() {
		fmt.Println("start governor")
		errCh <- server.Serve(listener)
		fmt.Println("stop governor")
	}()
	go func() {
		select {
		case <-stopCh:
			fmt.Println("stop...")
			server.Shutdown(context.Background())
		case <-errCh:
			fmt.Println("err occur...")
		}
		fmt.Println("close err ch")
		close(errCh)
	}()
	return nil
}
