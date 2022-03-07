package xxl

import (
	"bytes"
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/executor"
	"github.com/douyu/jupiter/pkg/executor/xxl/constants"
	"github.com/go-basic/ipv4"
	"github.com/stretchr/testify/assert"
)

func Test_NewExecutor(t *testing.T) {
	configStr := `
[xxl]
	[xxl.job]
		[xxl.job.executor]
			appname = "jupiter-xxl-job-demo"  # 执行器名称
			port = "59001"                    # 开启执行器的服务端口
			access_token = "jupiter-token"    # xxl-job需要的token信息
			address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名`
	// 测试配置加载
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	cc := conf.GetString("xxl.job.executor.appname")
	fmt.Println(cc)
	// 测试用例
	tests := []struct {
		name string
		opts []Option
		want Options
	}{
		{
			name: "useDefaultOptions",
			want: Options{
				ServerAddr:    "http://127.0.0.1:8080/xxl-job-admin",
				AccessToken:   "jupiter-token",
				ExecutorIp:    ipv4.LocalIP(),
				ExecutorPort:  "59001",
				RegistryKey:   "jupiter-xxl-job-demo",
				RegistryGroup: "EXECUTOR",
				LogDir:        constants.BasePath + "jupiter-xxl-job-demo" + "/jobhandler/",
				Switch:        true,
				Debug:         false,
			},
		},
		{
			name: "useOptions",
			opts: []Option{ExecutorHost("127.0.0.1"), ExecutorPort("7894"), AccessToken("test-token"),
				Debug(), RegistryKey("jupiter-demo"), RegistryGroup("CONTAINS"), Switch(false),
				ServerAddr("http://test/")},
			want: Options{
				ServerAddr:    "http://test/",
				AccessToken:   "test-token",
				ExecutorIp:    "127.0.0.1",
				ExecutorPort:  "7894",
				RegistryKey:   "jupiter-demo",
				RegistryGroup: "CONTAINS",
				LogDir:        constants.BasePath + "jupiter-demo/jobhandler/",
				Switch:        false,
				Debug:         true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := newJobExecutor(tt.opts...); !reflect.DeepEqual(got.opts, tt.want) {
				t.Errorf("newExecutor() = %v, want %v", got.opts, tt.want)
			}
		})
	}
}

func Test_Run(t *testing.T) {
	tests := []struct {
		name string
		opts []Option
		want Options
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := StdNewExecutor()
			executor.RegXJob(
				&TestIJob{},
			)
			assert.Nil(t, executor.Run())
		})
	}
}

// 测试IJob
type TestIJob struct{}

func (ins *TestIJob) GetJobName() string {
	return "test"
}

func (ins *TestIJob) Run(ctx context.Context, param *executor.RunReq) (msg string, err error) {
	fmt.Println("test")
	return "success", nil
}
