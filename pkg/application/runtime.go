package application

import (
	"os"
	"time"

	"github.com/douyu/jupiter/pkg/flag"
)

var (
	appHost  = os.Getenv("JUPITER_APP_HOST")
	hostname string
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

// Hostname gets hostname.
func Hostname() string {
	return hostname
}

// EnvServerHost gets JUPITER_APP_HOST.
func EnvServerHost() string {
	host := flag.String("host")
	if host != "" {
		return host
	}

	if appHost == "" {
		return "127.0.0.1"
	}
	return appHost
}
