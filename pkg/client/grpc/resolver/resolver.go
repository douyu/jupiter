package resolver

import (
	"context"
	"fmt"
	"net/url"
	"sync"

	"github.com/douyu/jupiter/pkg/registry"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/resolver"
)

// Register ...
func Register(name string, reg registry.Registry) {
	resolver.Register(&baseBuilder{
		name: name,
		reg:  reg,
	})
}

type baseBuilder struct {
	name string
	reg  registry.Registry
}

// Build ...
func (b *baseBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	return b.buildResolver(target, cc)
}

// Scheme ...
func (b baseBuilder) Scheme() string {
	return b.name
}

func (b *baseBuilder) buildResolver(target resolver.Target, cc resolver.ClientConn) (*baseResolver, error) {
	r := &baseResolver{
		target:    target,
		cc:        cc,
		builder:   b,
		addresses: make([]resolver.Address, 0),
	}
	go r.watchServiceNode(target)
	return r, nil
}

type baseResolver struct {
	target  resolver.Target
	cc      resolver.ClientConn
	builder *baseBuilder
	cancel  context.CancelFunc

	mtx      sync.RWMutex
	regInfos map[string]*url.Values // 注册信息
	cfgInfos map[string]*url.Values // 配置信息

	resolverAttributes *attributes.Attributes
	addresses          []resolver.Address
}

// ResolveNow ...
func (b *baseResolver) ResolveNow(options resolver.ResolveNowOptions) {}

// Close ...
func (b *baseResolver) Close() { b.cancel() }

func (b *baseResolver) watchServiceNode(target resolver.Target) {
	ctx, cancel := context.WithCancel(context.Background())
	b.cancel = cancel
	services, eventChan, err := b.builder.reg.WatchServices(ctx, target.Endpoint, "grpc")
	if err != nil {
		// todo(gorexlv): handle exception
		panic(err)
	}

	for _, service := range services {
		b.addresses = append(b.addresses, resolver.Address{
			Addr:       service.Address,
			ServerName: service.Name,
			Attributes: attributes.New(),
		})
	}

	for message := range eventChan {
		switch message.Event {
		case registry.EventUpdate:
			b.updateAddressList(message)
		case registry.EventDelete:
			b.deleteAddressList(message)
		default:
			panic("invalid event")
		}
	}
}

func (b *baseResolver) updateAddressList(message registry.EventMessage) {

}

func (b *baseResolver) deleteAddressList(message registry.EventMessage) {
	fmt.Printf("message = %+v\n", message)
}

func (b *baseResolver) updateClientConnState() {
	b.cc.UpdateState(resolver.State{
		Addresses:     b.addresses,
		ServiceConfig: nil,
		Attributes:    b.resolverAttributes, // resolver attributes
	})
}
