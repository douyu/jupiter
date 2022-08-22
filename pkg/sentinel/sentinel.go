package sentinel

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/douyu/jupiter/pkg/hooks"
	"github.com/douyu/jupiter/pkg/metric"
	"github.com/douyu/jupiter/pkg/xlog"
)

var (
	_logger = xlog.Jupiter().With(xlog.FieldMod("sentinel"))

	sentinelReqeust = metric.NewCounterVec("sentinel_request",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelSuccess = metric.NewCounterVec("sentinel_success",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelExceptionsThrown = metric.NewCounterVec("sentinel_exceptions_thrown",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelBlocked = metric.NewCounterVec("sentinel_blocked",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelRt = metric.NewHistogramVec("sentinel_rt",
		[]string{"resource", "language", "appName", "aid", "region", "zone", "iid", "mode"})
	sentinelState = metric.NewGaugeVec("sentinel_state",
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
