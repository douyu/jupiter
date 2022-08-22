package sentinel

import (
	"git.dz11.com/vega/minerva/prome"
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/douyu/jupiter/pkg/hooks"
	"github.com/douyu/jupiter/pkg/xlog"
)

var (
	_logger = xlog.Jupiter().With(xlog.FieldMod("sentinel"))

	sentinelReqeust = prome.NewCounterVec("sentinel_request",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelSuccess = prome.NewCounterVec("sentinel_success",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelExceptionsThrown = prome.NewCounterVec("sentinel_exceptions_thrown",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelBlocked = prome.NewCounterVec("sentinel_blocked",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelRt = prome.NewHistogramVec("sentinel_rt",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelState = prome.NewGaugeVec("sentinel_state",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})

	WithError = base.WithError
)

const (
	language = "go"

	SENTINEL_DATASOURCE_ETCD  = "etcd"
	SENTINEL_DATASOURCE_FILES = "files"
)

type (
	SentinelEntry = base.SentinelEntry
	BlockError    = base.BlockError
	EntryContext  = base.EntryContext
	EntryOption   = sentinel.EntryOption
)

var (
	stdConfig Config
)

func init() {
	hooks.Register(hooks.Stage_AfterLoadConfig, func() {
		_ = build()
	})
}

// build 基于标准配置构建sentinel.
func build() error {
	stdConfig = StdConfig()

	return stdConfig.Build()
}

// Entry 执行熔断策略.
func Entry(resource string, opts ...EntryOption) (*SentinelEntry, *BlockError) {
	return stdConfig.Entry(resource, opts...)
}
