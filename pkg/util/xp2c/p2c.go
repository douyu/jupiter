package xp2c

import "google.golang.org/grpc/balancer"

type P2c interface {
	// Next returns next selected item.
	Next() (interface{}, func(balancer.DoneInfo))
	// Add a item.
	Add(interface{})
}
