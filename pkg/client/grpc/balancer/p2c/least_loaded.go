package p2c

import (
	"context"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

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

	var scs = make([]*subConn, 0, len(readySCs))

	for _, sc := range readySCs {
		sc := &subConn{
			conn:     sc,
			inflight: 0,
		}

		scs = append(scs, sc)
	}
	rp := &p2cPicker{
		subConns: scs,
		rand:     rand.New(rand.NewSource(time.Now().Unix())),
	}
	return rp
}

type subConn struct {
	conn balancer.SubConn

	// grpc client inflight count
	inflight int64
}

type p2cPicker struct {
	// subConns is the snapshot of the roundrobin balancer when this picker was
	// created. The map is immutable. Each Get() will do a p2c
	// selection from it and return the selected SubConn.
	subConns []*subConn
	mu       sync.Mutex

	rand *rand.Rand
}

// Pick ...
func (p *p2cPicker) Pick(ctx context.Context, opts balancer.PickOptions) (balancer.SubConn, func(balancer.DoneInfo), error) {
	var sc, backsc *subConn

	if p == nil || len(p.subConns) <= 0 {
		return nil, nil, balancer.ErrNoSubConnAvailable
	} else if len(p.subConns) == 1 {
		sc = p.subConns[0]
	} else {

		// rand需要加锁
		p.mu.Lock()
		a := p.rand.Intn(len(p.subConns))
		b := p.rand.Intn(len(p.subConns) - 1)
		p.mu.Unlock()

		if b >= a {
			b = b + 1
		}
		sc, backsc = p.subConns[a], p.subConns[b]

		// 根据inflight选择更优节点
		if sc.inflight > backsc.inflight {
			sc, backsc = backsc, sc
		}
	}

	atomic.AddInt64(&sc.inflight, 1)

	return sc.conn, func(di balancer.DoneInfo) {
		atomic.AddInt64(&sc.inflight, -1)
	}, nil
}
