package registry

import (
	"context"
	"fmt"

	"github.com/douyu/jupiter/pkg/component"
	"github.com/douyu/jupiter/pkg/server"
)

var _ component.Component = &ServerComponent{}

type ServerComponent struct {
	server.Server
}

func (c ServerComponent) Start(stopCh <-chan struct{}) error {
	DefaultRegisterer.RegisterService(context.Background(), c.Info())
	var errCh = make(chan error)
	go func() {
		fmt.Println("before serve...")
		errCh <- c.Serve()
		fmt.Println("after server...")
	}()
	go func() {
		defer DefaultRegisterer.UnregisterService(context.Background(), c.Info())
		select {
		case <-stopCh:
			fmt.Println("stop...")
			c.GracefulStop(context.Background())
		case <-errCh:
			fmt.Println("err occur...")
		}
		fmt.Println("close err ch")
		close(errCh)
	}()
	return nil
}

func (c ServerComponent) ShouldBeLeader() bool {
	return false
}
