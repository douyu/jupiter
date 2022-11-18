// Copyright 2022 Douyu
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gorm

import (
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

func open(options *Config) (*gorm.DB, error) {
	inner, err := gorm.Open(mysql.Open(options.DSN), &options.gormConfig)
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

	registerInterceptor(inner, options,
		metricInterceptor(),
		traceInterceptor(),
		sentinelInterceptor(),
	)

	return inner, err
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
	dsn, err := parseDSN(options.DSN)
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
