package istats

import (
	"sync/atomic"
	"time"
)

// flowInfo 请求流量信息
type FlowInfoBase struct {
	Ctime               int64  `json:"ctime"`               // 创建时间
	HasFlow             bool   `json:"hasFlow"`             // 是否有请求流量
	FirstFlowTime       int64  `json:"firstFlowTime"`       // 首次请求时间
	LastFlowTime        int64  `json:"lastFlowTime"`        // 最后一次请求时间
	HasShadowFlow       bool   `json:"hasShadowFlow"`       // 是否有影子表请求流量
	FirstShadowFlowTime int64  `json:"firstShadowFlowTime"` // 首次请求影子表时间
	LastShadowFlowTime  int64  `json:"lastShadowFlowTime"`  // 最后一次请求影子表时间
	ShadowSwitch        string `json:"shadowSwitch"`        // 影子流量开关， on打开， off关闭， watch观察者模式（关闭且打印影子日志）
}

// UpdateFlow 流量状态更新
func (f *FlowInfoBase) UpdateFlow() {
	now := time.Now().Unix()
	if atomic.CompareAndSwapInt64(&f.FirstFlowTime, 0, now) {
		f.HasFlow = true
	}
	atomic.StoreInt64(&f.LastFlowTime, now)
}

// UpdateFlow 压测流量状态更新
func (f *FlowInfoBase) UpdateShadowFlow() {
	now := time.Now().Unix()
	if atomic.CompareAndSwapInt64(&f.FirstShadowFlowTime, 0, now) {
		f.HasShadowFlow = true
	}
	atomic.StoreInt64(&f.LastShadowFlowTime, now)
}

func NewFlowInfoBase(shadowSwitch string) FlowInfoBase {
	return FlowInfoBase{
		Ctime:        time.Now().Unix(),
		ShadowSwitch: shadowSwitch,
	}
}
