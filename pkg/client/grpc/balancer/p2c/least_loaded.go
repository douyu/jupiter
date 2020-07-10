package p2c

import (
	"context"

	"github.com/douyu/jupiter/pkg/util/xp2c"
	"github.com/douyu/jupiter/pkg/util/xp2c/leastloaded"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/grpclog"
	_ "google.golang.org/grpc/health"
	"google.golang.org/grpc/resolver"
)

// Name is the name of p2c with least loaded balancer.
const (
	Name = "p2c_least_loaded"
)

// newBuilder creates a new balance builder.
func newBuilder() balancer.Builder {
	return base.NewBalancerBuilderWithConfig(Name, &p2cPickerBuilder{}, base.Config{HealthCheck: true})
}

func init() {
	balancer.Register(newBuilder())
}

type p2cPickerBuilder struct{}

func (*p2cPickerBuilder) Build(readySCs map[resolver.Address]balancer.SubConn) balancer.Picker {
	grpclog.Infof("p2cPickerBuilder: newPicker called with readySCs: %v", readySCs)
	if len(readySCs) == 0 {
		return base.NewErrPicker(balancer.ErrNoSubConnAvailable)
	}

	var p2c = leastloaded.New()

	for _, sc := range readySCs {
		p2c.Add(sc)
	}

	rp := &p2cPicker{
		p2c: p2c,
	}
	return rp
}

type p2cPicker struct {
	p2c xp2c.P2c
}

// Pick ...
func (p *p2cPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {

	item, done := p.p2c.Next()
	if item == nil {
		return nil, nil, balancer.ErrNoSubConnAvailable
	}

	return item.(balancer.SubConn), done, nil
}
