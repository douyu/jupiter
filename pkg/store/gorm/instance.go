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
	"errors"

	prome "github.com/douyu/jupiter/pkg/core/metric"
	"github.com/douyu/jupiter/pkg/util/xretry"
	"github.com/douyu/jupiter/pkg/xlog"
	"gorm.io/gorm"
)

// dial returns a new DB connection or an error.
func dial(name string, config *Config) *gorm.DB {
	if config.DSN == "" {
		xlog.Jupiter().Panic("empty dsn", xlog.FieldName(name))
	}

	dsn, err := parseDSN(config.DSN)
	if err != nil {
		xlog.Jupiter().Panic("parse dsn", xlog.FieldName(name), xlog.FieldErr(err))
	}

	var db *gorm.DB
	err = xretry.Do(config.Retry, config.RetryWaitTime, func() error {
		db, err = open(config)
		if err != nil {
			prome.LibHandleCounter.WithLabelValues(prome.TypeMySQL, name+".dial", dsn.Addr, err.Error()).Inc()
			return errors.New("dial nil" + err.Error())
		}

		if db == nil {
			prome.LibHandleCounter.WithLabelValues(prome.TypeMySQL, name+".dial", dsn.Addr, "nil db").Inc()
			return errors.New("db nil" + err.Error())
		}
		sql, err := db.DB()
		if err != nil {
			prome.LibHandleCounter.WithLabelValues(prome.TypeMySQL, name+".ping", dsn.Addr, err.Error()).Inc()
			return errors.New("db " + err.Error())
		}

		if err := sql.Ping(); err != nil {
			prome.LibHandleCounter.WithLabelValues(prome.TypeMySQL, name+".ping", dsn.Addr, err.Error()).Inc()
			return errors.New("ping " + err.Error())
		}

		return nil
	})
	if err != nil {
		if config.OnDialError == "panic" {
			xlog.Jupiter().Panic("dial mysql db", xlog.FieldErr(err), xlog.FieldAddr(dsn.Addr), xlog.FieldName(name))
		} else {
			xlog.Jupiter().Panic("dial mysql db", xlog.FieldAddr(dsn.Addr), xlog.FieldErr(err), xlog.FieldName(name))
		}
	}

	return db
}
