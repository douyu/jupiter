package gorm

import (
	"context"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type (
	DB    = gorm.DB
	Model = gorm.Model
)

var (
	// ErrRecordNotFound record not found error.
	ErrRecordNotFound = gorm.ErrRecordNotFound
)

// WithContext ...
func WithContext(ctx context.Context, db *gorm.DB) *gorm.DB {
	return db.WithContext(ctx)
}

func open(options *Config) (*gorm.DB, error) {
	inner, err := gorm.Open(mysql.Open(options.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// inner.(options.Debug)
	// 设置默认连接配置
	db, err := inner.DB()
	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(options.MaxIdleConns)
	db.SetMaxOpenConns(options.MaxOpenConns)

	if options.ConnMaxLifetime != 0 {
		db.SetConnMaxLifetime(options.ConnMaxLifetime)
	}

	// 开启 debug
	if options.Debug {
		inner = inner.Debug()
	}

	registerInterceptor(inner, options, metricInterceptor(), traceInterceptor())

	return inner, err
}

// Open ...
// Deprecated
func Open(options *Config) (*gorm.DB, error) {
	return open(options)
}

// 收敛status，避免prometheus日志太多
func getStatement(err string) string {
	if !strings.HasPrefix(err, "Errord") {
		return "Unknown"
	}
	slice := strings.Split(err, ":")
	if len(slice) < 2 {
		return "Unknown"
	}

	// 收敛错误
	return slice[0]
}

type processor interface {
	Get(name string) func(*gorm.DB)
	Replace(name string, handler func(*gorm.DB)) error
}

func registerInterceptor(db *gorm.DB, options *Config, interceptors ...Interceptor) {
	dsn, err := ParseDSN(options.DSN)
	if err != nil {
		panic(err)
	}

	var processors = []struct {
		Name      string
		Processor processor
	}{
		{"gorm:create", db.Callback().Create()},
		{"gorm:query", db.Callback().Query()},
		{"gorm:delete", db.Callback().Delete()},
		{"gorm:update", db.Callback().Update()},
		{"gorm:row", db.Callback().Row()},
		{"gorm:raw", db.Callback().Raw()},
	}

	for _, interceptor := range interceptors {
		for _, processor := range processors {
			handler := processor.Processor.Get(processor.Name)
			handler = interceptor(dsn, processor.Name, options, handler)

			if err := processor.Processor.Replace(processor.Name, handler); err != nil {
				panic(err)
			}
		}
	}
}
