package xxl

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/executor"
	"github.com/douyu/jupiter/pkg/executor/xxl/constants"
	"github.com/go-basic/ipv4"
	"github.com/stretchr/testify/assert"
)

func Test_StdNewExecutor(t *testing.T) {
	configStr := `
	[xxl]
		[xxl.job]
			[xxl.job.admin]
				appname = "jupiter-xxl-job-demo"  # 执行器名称
				port = "59001"                    # 开启执行器的服务端口
				access_token = "jupiter-token"    # xxl-job需要的token信息
				address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名`
	// 测试配置加载
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "get",
			args: "task1",
			want: "127.0.0.1:59001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StdNewExecutor(ExecutorHost("127.0.0.1")); !reflect.DeepEqual(got.GetAddress(), tt.want) {
				t.Errorf("StdNewExecutor() = %v, want %v", got.GetAddress(), tt.want)
			}
		})
	}
}

func Test_NewExecutor(t *testing.T) {
	configStr := `
[xxl]
	[xxl.job]
		[xxl.job.admin]
			appname = "jupiter-xxl-job-demo"  # 执行器名称
			port = "59001"                    # 开启执行器的服务端口
			access_token = "jupiter-token"    # xxl-job需要的token信息
			address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名`
	// 测试配置加载
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
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
	configStr := `
	[xxl]
		[xxl.job]
			[xxl.job.admin]
				appname = "jupiter-xxl-job-demo"  # 执行器名称
				port = "59001"                    # 开启执行器的服务端口
				access_token = "jupiter-token"    # xxl-job需要的token信息
				address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名`
	// 测试配置加载
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	tests := []struct {
		name string
	}{
		{
			name: "run",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := newJobExecutor()
			executor.RegXJob(
				&TestIJob{},
			)
			assert.Nil(t, executor.Run())
		})
	}
}

func Test_RegXJob(t *testing.T) {
	configStr := `
	[xxl]
		[xxl.job]
			[xxl.job.admin]
				appname = "jupiter-xxl-job-demo"  # 执行器名称
				port = "59001"                    # 开启执行器的服务端口
				access_token = "jupiter-token"    # xxl-job需要的token信息
				address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名`
	// 测试配置加载
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	tests := []struct {
		name string
	}{
		{
			name: "RegXJob",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := newJobExecutor()
			executor.RegXJob(
				&TestIJob{},
			)
		})
	}
}

// TODO
func Test_runTaskHandler(t *testing.T) {
}

func Test_registryRemove(t *testing.T) {
	configStr := `
	[xxl]
		[xxl.job]
			[xxl.job.admin]
				appname = "jupiter-xxl-job-demo"  # 执行器名称
				port = "59001"                    # 开启执行器的服务端口
				access_token = "jupiter-token"    # xxl-job需要的token信息
				address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名`
	// 测试配置加载
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	tests := []struct {
		name string
	}{
		{
			name: "registryRemove",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := newJobExecutor()
			executor.RegXJob(
				&TestIJob{},
			)
			executor.registryRemove()
		})
	}
}

func Test_GracefulStop(t *testing.T) {
	configStr := `
	[xxl]
		[xxl.job]
			[xxl.job.admin]
				appname = "jupiter-xxl-job-demo"  # 执行器名称
				port = "59001"                    # 开启执行器的服务端口
				access_token = "jupiter-token"    # xxl-job需要的token信息
				address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名`
	// 测试配置加载
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	tests := []struct {
		name string
	}{
		{
			name: "GracefulStop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := newJobExecutor()
			executor.RegXJob(
				&TestIJob{},
			)
			executor.GracefulStop()
		})
	}
}

func Test_Stop(t *testing.T) {
	configStr := `
	[xxl]
		[xxl.job]
			[xxl.job.admin]
				appname = "jupiter-xxl-job-demo"  # 执行器名称
				port = "59001"                    # 开启执行器的服务端口
				access_token = "jupiter-token"    # xxl-job需要的token信息
				address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名`
	// 测试配置加载
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	tests := []struct {
		name string
	}{
		{
			name: "Stop",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := newJobExecutor()
			executor.RegXJob(
				&TestIJob{},
			)
			executor.Stop()
		})
	}
}

func Test_callback(t *testing.T) {
	configStr := `
	[xxl]
		[xxl.job]
			[xxl.job.admin]
				appname = "jupiter-xxl-job-demo"  # 执行器名称
				port = "59001"                    # 开启执行器的服务端口
				access_token = "jupiter-token"    # xxl-job需要的token信息
				address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名`
	// 测试配置加载
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	job := &TestXJob2{}
	ctx, cancel := context.WithCancel(context.Background())
	task := &Task{
		Id:   1,
		Name: job.GetJobName(),
		fn:   job.Run,
		Param: &executor.RunReq{
			JobID:          123,
			ExecutorParams: "",
		},
		ctx:    ctx,
		cancel: cancel,
	}
	tests := []struct {
		name string
	}{
		{
			name: "callback",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := newJobExecutor()
			executor.RegXJob(
				&TestIJob{},
			)
			executor.callback(task, 0, "执行成功", false)
			executor.callback(task, 500, "执行失败", false)
			executor.callback(task, 0, "执行成功", true)
		})
	}
}

func Test_post(t *testing.T) {
	configStr := `
	[xxl]
		[xxl.job]
			[xxl.job.admin]
				appname = "jupiter-xxl-job-demo"  # 执行器名称
				port = "59001"                    # 开启执行器的服务端口
				access_token = "jupiter-token"    # xxl-job需要的token信息
				address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名`
	// 测试配置加载
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	tests := []struct {
		name string
	}{
		{
			name: "post",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := newJobExecutor()
			executor.RegXJob(
				&TestIJob{},
			)
			executor.post("register", "register")
		})
	}
}
func Test_handlerWithAuth(t *testing.T) {
	configStr := `
	[xxl]
		[xxl.job]
			[xxl.job.admin]
				appname = "jupiter-xxl-job-demo"  # 执行器名称
				port = "59001"                    # 开启执行器的服务端口
				access_token = "jupiter-token"    # xxl-job需要的token信息
				address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名`
	// 测试配置加载
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	executor := newJobExecutor(ExecutorHost("127.0.0.1"))
	req_succuss := &http.Request{
		Header: http.Header{},
	}
	req_succuss.Header.Add("XXL-JOB-ACCESS-TOKEN", "jupiter-token")
	req_fail := &http.Request{
		Header: http.Header{},
	}
	req_fail.Header.Add("XXL-JOB-ACCESS-TOKEN", "jupiter-token_fail")
	var init = func(http.ResponseWriter, *http.Request) {
	}

	tests := []struct {
		name string
		args HttpHandler
	}{
		{
			name: "verifyAuth nil",
			args: init,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := executor.handlerWithAuth(tt.args)
			got(&ResponseWriter{}, req_succuss)
			got(&ResponseWriter{}, req_fail)
		})
	}
}

type ResponseWriter struct{}

func (ins *ResponseWriter) Header() http.Header {
	return make(map[string][]string)
}

func (ins *ResponseWriter) Write([]byte) (int, error) {
	return 1, nil
}
func (ins *ResponseWriter) WriteHeader(statusCode int) {
}

func Test_verifyAuth(t *testing.T) {
	configStr := `
	[xxl]
		[xxl.job]
			[xxl.job.admin]
				appname = "jupiter-xxl-job-demo"  # 执行器名称
				port = "59001"                    # 开启执行器的服务端口
				access_token = "jupiter-token"    # xxl-job需要的token信息
				address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名`
	// 测试配置加载
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	executor := newJobExecutor(ExecutorHost("127.0.0.1"))
	req_succuss := &http.Request{
		Header: http.Header{},
	}
	req_succuss.Header.Add("XXL-JOB-ACCESS-TOKEN", "jupiter-token")
	req_fail := &http.Request{
		Header: http.Header{},
	}
	req_fail.Header.Add("XXL-JOB-ACCESS-TOKEN", "jupiter-token_fail")
	tests := []struct {
		name string
		args *http.Request
		want bool
	}{
		{
			name: "verifyAuth nil",
			args: nil,
			want: false,
		},
		{
			name: "verifyAuth success",
			args: req_succuss,
			want: true,
		},
		{
			name: "verifyAuth fail",
			args: req_fail,
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := executor.verifyAuth(tt.args); got != tt.want {
				t.Errorf("verifyAuth() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_GetAddress(t *testing.T) {
	configStr := `
	[xxl]
		[xxl.job]
			[xxl.job.admin]
				appname = "jupiter-xxl-job-demo"  # 执行器名称
				port = "59001"                    # 开启执行器的服务端口
				access_token = "jupiter-token"    # xxl-job需要的token信息
				address = "http://127.0.0.1:8080/xxl-job-admin"  # 注意换成XXL调度中心对应环境的域名`
	// 测试配置加载
	assert.Nil(t, conf.LoadFromReader(bytes.NewBufferString(configStr), toml.Unmarshal))
	tests := []struct {
		name string
		args string
		want string
	}{
		{
			name: "get",
			args: "127.0.0.1",
			want: "127.0.0.1:59001",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			executor := StdNewExecutor(ExecutorHost(tt.args))
			got := executor.GetAddress()
			if got != tt.want {
				t.Errorf("GetAddress() = %v, want %v", got, tt.want)
			}
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
