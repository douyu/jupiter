package gorm

import (
	"bytes"
	"testing"
	"time"

	"github.com/BurntSushi/toml"
	cfg "github.com/douyu/jupiter/pkg/conf"
	"github.com/stretchr/testify/assert"
)

type User struct {
	Birthday time.Time
	Age      int
	Name     string `gorm:"size:255"` // string默认长度为255, 使用这种tag重设。
}

func (u User) TableName() string {
	return "users"
}

func ExampleClient_New() {
	// 1. 从配置文件中创建client
	_ = StdNew("demo") // 读取jupiter.mysql.demo配置，创建实例

	// 2. 从配置文件按照自定义格式创建client
	_ = dial("demo", RawConfig("jupiter.pq.demo"))

	// 3. 自定义配置，创建client
	_ = dial("demo", &Config{
		Debug:           false,
		MaxIdleConns:    0,
		MaxOpenConns:    0,
		ConnMaxLifetime: 0,
		OnDialError:     "",
		raw:             nil,
	})
}

func TestStdNew(t *testing.T) {
	type args struct {
		name string
		opts []interface{}
	}
	tests := []struct {
		name    string
		args    args
		config  string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "std new",
			args: args{
				name: "demo",
				opts: []interface{}{},
			},
			wantErr: false,
			config: `
			[jupiter.mysql.demo]
				dsn="root:123456@tcp(mysql)/mysql?timeout=20s&readTimeout=20s"
				debug=true
				maxIdleConns=50
				connMaxLifeTime="20m"
				level="panic"
				slowThreshold="400ms"
				dialTimeout="1s"
			`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Nil(t, cfg.LoadFromReader(bytes.NewReader([]byte(tt.config)), toml.Unmarshal))
			defer cfg.Reset()
			if tt.wantErr {
				assert.Panics(t, func() { StdNew(tt.args.name, tt.args.opts...) })
			} else {
				assert.NotPanics(t, func() { StdNew(tt.args.name, tt.args.opts...) })
			}
		})

		t.Run(tt.name, func(t *testing.T) {
			assert.Nil(t, cfg.LoadFromReader(bytes.NewReader([]byte(tt.config)), toml.Unmarshal))
			defer cfg.Reset()
			assert.NotPanics(t, func() { Invoker(tt.args.name) })
		})
	}
}
