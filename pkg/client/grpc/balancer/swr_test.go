package balancer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/attributes"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
)

type mockSubConn struct { num int }
func (m mockSubConn) UpdateAddresses(addresses []resolver.Address) { }
func (m mockSubConn) Connect() { }

func Test_weightPicker(t *testing.T) {
	readySCs := map[balancer.SubConn]base.SubConnInfo{}
	for i:=0;i<10;i++ {
		var subConn = &mockSubConn{num:i}
		readySCs[subConn]= base.SubConnInfo{
			Address: resolver.Address{
				Addr: fmt.Sprintf("127.0.0.1:909%d", i),
				ServerName: "server_name",
				Attributes: attributes.New("meta", &Config{
					Group:  "red",
					Weight: 10,
				}),
			},
		}
	}

	t.Run("round_robin", func(t *testing.T) {
		picker := newWeightPicker(readySCs)
		counter := map[balancer.SubConn]int{}
		for i:=0;i<10;i++ {
			result, err := picker.Pick(balancer.PickInfo{})
			assert.Nil(t, err)
			assert.Contains(t, readySCs, result.SubConn)
			if _, ok := counter[result.SubConn]; !ok {
				counter[result.SubConn]=0
			}
			counter[result.SubConn] += 1
		}
		assert.Equal(t, len(counter), 10)
	})

	t.Run("no group provided", func(t *testing.T) {
	})

	t.Run("no group registered", func(t *testing.T) {

	})
}