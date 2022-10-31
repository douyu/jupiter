package xxl

import (
	"bytes"
	"fmt"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/executor/xxl/constants"
	"github.com/stretchr/testify/assert"
)

func Test_DefaultOptions(t *testing.T) {
	tests := []struct {
		name string
		want *Options
	}{
		{
			name: "DefaultOptions",
			want: &Options{
				ExecutorIp:    DefaultExecuteIp,
				ExecutorPort:  DefaultExecutorPort,
				RegistryKey:   DefaultRegistryKey,
				RegistryGroup: DefaultRegistryGroup,
				AccessToken:   DefaultAccessToken,
				ServerAddr:    "",
				Switch:        DefaultSwitch,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := DefaultOptions(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultOptions() failed")
			}
		})
	}
}

func Test_DefaultConfig(t *testing.T) {
	configStr := `
	[xxl]
		[xxl.job]
			[xxl.job.admin]
				appname = "jupiter-xxl-job-demo"  # 执行器名称
				port = "59001"                    # 开启执行器的服务端口
				access_token = "jupiter-token"    # xxl-job需要的token信息
				address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名`
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	tests := []struct {
		name string
		want *Options
	}{
		{
			name: "DefaultConfig",
			want: &Options{
				ExecutorIp:    DefaultExecuteIp,
				ExecutorPort:  "59001",
				RegistryKey:   "jupiter-xxl-job-demo",
				RegistryGroup: DefaultRegistryGroup,
				AccessToken:   "jupiter-token",
				ServerAddr:    "http://127.0.0.1:8080/xxl-job-admin",
				Switch:        DefaultSwitch,
				LogDir:        constants.BasePath + "jupiter-xxl-job-demo" + "/jobhandler/",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DefaultConfig()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_Build(t *testing.T) {
	options := DefaultOptions()
	options.ServerAddr = "127.0.0.1"
	tests := []struct {
		name string
		want *JobExecutor
	}{
		{
			name: "Build",
			want: &JobExecutor{
				opts: *options,
				regList: &taskList{
					data: make(map[string]*TaskWithPending),
				},
				runList: &taskList{
					data: make(map[string]*TaskWithPending),
				},
				address: fmt.Sprintf("%s:%s", options.ExecutorIp, options.ExecutorPort),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := options.Build()
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultOptions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_CheckOptions(t *testing.T) {
	options := DefaultOptions()
	options.ServerAddr = "127.0.0.1"
	t.Run("succuss", func(t *testing.T) {
		CheckOptions(options)
	})
}

func Test_ServerAddr(t *testing.T) {
	options := DefaultOptions()
	options.ServerAddr = "127.0.0.1"
	tests := []struct {
		args string
		name string
		want string
	}{
		{
			name: "Build",
			args: "172.10.1.1",
			want: "172.10.1.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			function := ServerAddr("172.10.1.1")
			function(options)
			if !reflect.DeepEqual(options.ServerAddr, tt.want) {
				t.Errorf("ServerAddr() = %v, want %v", options.ServerAddr, tt.want)
			}
		})
	}
}

func Test_AccessToken(t *testing.T) {
	options := DefaultOptions()
	tests := []struct {
		args string
		name string
		want string
	}{
		{
			name: "AccessToken",
			args: "test_token",
			want: "test_token",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			function := AccessToken(tt.args)
			function(options)
			if !reflect.DeepEqual(options.AccessToken, tt.want) {
				t.Errorf("AccessToken() = %v, want %v", options.AccessToken, tt.want)
			}
		})
	}
}

func Test_ExecutorHost(t *testing.T) {
	options := DefaultOptions()
	tests := []struct {
		args string
		name string
		want string
	}{
		{
			name: "ExecutorHost",
			args: "172.10.1.1",
			want: "172.10.1.1",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			function := ExecutorHost(tt.args)
			function(options)
			if !reflect.DeepEqual(options.ExecutorIp, tt.want) {
				t.Errorf("ExecutorHost() = %v, want %v", options.ExecutorIp, tt.want)
			}
		})
	}
}

func Test_ExecutorPort(t *testing.T) {
	options := DefaultOptions()
	tests := []struct {
		args string
		name string
		want string
	}{
		{
			name: "ExecutorPort",
			args: "8888",
			want: "8888",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			function := ExecutorPort(tt.args)
			function(options)
			if !reflect.DeepEqual(options.ExecutorPort, tt.want) {
				t.Errorf("ExecutorPort() = %v, want %v", options.ExecutorPort, tt.want)
			}
		})
	}
}

func Test_RegistryKey(t *testing.T) {
	options := DefaultOptions()
	tests := []struct {
		args string
		name string
		want string
	}{
		{
			name: "RegistryKey",
			args: "jupiter-application-demo",
			want: "jupiter-application-demo",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			function := RegistryKey(tt.args)
			function(options)
			if !reflect.DeepEqual(options.RegistryKey, tt.want) {
				t.Errorf("RegistryKey() = %v, want %v", options.RegistryKey, tt.want)
			}
		})
	}
}

func Test_RegistryGroup(t *testing.T) {
	options := DefaultOptions()
	tests := []struct {
		args string
		name string
		want string
	}{
		{
			name: "RegistryGroup",
			args: "CONTAINS",
			want: "CONTAINS",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			function := RegistryGroup(tt.args)
			function(options)
			if !reflect.DeepEqual(options.RegistryGroup, tt.want) {
				t.Errorf("RegistryGroup() = %v, want %v", options.RegistryGroup, tt.want)
			}
		})
	}
}

func Test_Switch(t *testing.T) {
	options := DefaultOptions()
	tests := []struct {
		args bool
		name string
		want bool
	}{
		{
			name: "Switch1",
			args: true,
			want: true,
		},
		{
			name: "Switch2",
			args: false,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			function := Switch(tt.args)
			function(options)
			if !reflect.DeepEqual(options.Switch, tt.want) {
				t.Errorf("Switch() = %v, want %v", options.Switch, tt.want)
			}
		})
	}
}

func Test_Debug(t *testing.T) {
	options := DefaultOptions()
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "Debug",
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			function := Debug()
			function(options)
			if !reflect.DeepEqual(options.Debug, tt.want) {
				t.Errorf("Debug() = %v, want %v", options.Debug, tt.want)
			}
		})
	}
}
