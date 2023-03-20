package xnet

import (
	"net"
	"testing"

	"github.com/samber/lo"
)

func TestAddress(t *testing.T) {
	host, _, _ := GetLocalMainIP()
	if host == "" {
		t.Fail()
	}

	type args struct {
		listener net.Listener
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "ipv4",
			args: args{
				listener: lo.Must(net.Listen("tcp", "127.0.0.1:48080")),
			},
			want: "127.0.0.1:48080",
		},
		{
			name: "0.0.0.0",
			args: args{
				listener: lo.Must(net.Listen("tcp", "0.0.0.0:48081")),
			},
			want: host + ":48081",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Address(tt.args.listener); got != tt.want {
				t.Errorf("Address() = %v, want %v", got, tt.want)
			}
		})
	}
}
