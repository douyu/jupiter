package xxl

import (
	"fmt"
	"strings"
	"time"

	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/executor/xxl/constants"
	"github.com/douyu/jupiter/pkg/xlog"
	"github.com/go-basic/ipv4"
)

type Options struct {
	ServerAddr    string        `json:"address" toml:"address"`                //调度中心地址
	AccessToken   string        `json:"access_token" toml:"access_token"`      //请求令牌
	Timeout       time.Duration `json:"timeout" toml:"timeout"`                //接口超时时间
	ExecutorIp    string        `json:"executor_ip" toml:"executor_ip"`        //本地(执行器)IP(可自行获取)
	ExecutorPort  string        `json:"port" toml:"port"`                      //本地(执行器)端口
	RegistryKey   string        `json:"appname" toml:"appname"`                //执行器名称
	RegistryGroup string        `json:"registry_group " toml:"registry_group"` //执行器组，默认EXECUTOR
	LogDir        string        `json:"log_dir" toml:"log_dir"`                //日志目录
	Switch        bool          `json:"switch" toml:"switch"`                  //开关
	Debug         bool          `json:"debug" toml:"debug"`                    //开关
}

var (
	DefaultExecutorPort  = "59000"
	DefaultAccessToken   = "jupiter-task-token"
	DefaultRegistryKey   = "jupiter-demo-app"
	DefaultRegistryGroup = "EXECUTOR"
	DefaultSwitch        = true
	DefaultExecuteIp     = ipv4.LocalIP()
)

func DefaultOptions() *Options {
	opt := &Options{
		ExecutorIp:    DefaultExecuteIp,
		ExecutorPort:  DefaultExecutorPort,
		RegistryKey:   DefaultRegistryKey,
		RegistryGroup: DefaultRegistryGroup,
		AccessToken:   DefaultAccessToken,
		ServerAddr:    "",
		Switch:        DefaultSwitch,
	}
	return opt
}

func DefaultConfig() *Options {
	// 加载框架默认Options
	options := DefaultOptions()
	// 加载用户自定义Options
	if err := conf.UnmarshalKey("xxl.job.executor", options, conf.TagName("toml")); err != nil {
		xlog.Jupiter().Panic("unmarshal config", xlog.FieldName("xxl.job.executor"), xlog.FieldErr(err))
	}
	if strings.TrimSpace(options.LogDir) == "" {
		options.LogDir = constants.BasePath + options.RegistryKey + "/jobhandler/"
	}
	return options
}

func (options *Options) Build() *JobExecutor {
	// 校验配置项，如果校验失败将会panic
	CheckOptions(options)
	// 创建执行器
	executor := &JobExecutor{
		opts: *options,
		regList: &taskList{
			data: make(map[string]*TaskWithPending),
		},
		runList: &taskList{
			data: make(map[string]*TaskWithPending),
		},
		address: fmt.Sprintf("%s:%s", options.ExecutorIp, options.ExecutorPort),
	}
	return executor
}

type Option func(o *Options)

func CheckOptions(opts *Options) {
	// 校验调度中心地址
	if opts.ServerAddr == "" {
		panic("invalid admin server")
	}
	// TODO:others check
}

// 设置调度中心地址
func ServerAddr(addr string) Option {
	return func(o *Options) {
		o.ServerAddr = addr
	}
}

// 请求令牌
func AccessToken(token string) Option {
	return func(o *Options) {
		o.AccessToken = token
	}
}

// 设置执行器IP
func ExecutorHost(ip string) Option {
	return func(o *Options) {
		o.ExecutorIp = ip
	}
}

// 设置执行器端口
func ExecutorPort(port string) Option {
	return func(o *Options) {
		o.ExecutorPort = port
	}
}

// 设置执行器标识
func RegistryKey(registryKey string) Option {
	return func(o *Options) {
		o.RegistryKey = registryKey
		o.LogDir = constants.BasePath + o.RegistryKey + "/jobhandler/"
	}
}

// 设置执行器组
func RegistryGroup(registryGroup string) Option {
	return func(o *Options) {
		o.RegistryGroup = registryGroup
	}
}

// 设置默认开关
func Switch(s bool) Option {
	return func(o *Options) {
		o.Switch = s
	}
}

// 本地调试
func Debug() Option {
	return func(o *Options) {
		o.Debug = true
	}
}
