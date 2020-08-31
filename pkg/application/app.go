package application

import (
	"os"
	"time"
)

func init() {
	hn, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	hostname = hn
}

// 通用状态信息
type RuntimeStats struct {
	IP       string `json:"ip"`
	Hostname string `json:"hostname"`
	Time     string `json:"time"` // 每次展示数据的时间
	Err      string `json:"err"`
}

func NewRuntimeStats() RuntimeStats {
	return RuntimeStats{
		IP:       EnvServerHost(),
		Hostname: Hostname(),
		Time:     time.Now().Format("2006-01-02 15:04:05"),
	}
}

var hostname string

// Hostname gets hostname.
func Hostname() string {
	return hostname
}
