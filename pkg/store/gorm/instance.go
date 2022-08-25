package gorm

import (
	"errors"

	prome "github.com/douyu/jupiter/pkg/metric"
	"github.com/douyu/jupiter/pkg/util/xretry"
	"github.com/douyu/jupiter/pkg/xlog"
	"gorm.io/gorm"
)

// dial returns a new DB connection or an error.
func dial(name string, config *Config) *gorm.DB {
	if config.DSN == "" {
		xlog.Jupiter().Panic("empty dsn", xlog.FieldName(name))
	}

	d, err := ParseDSN(config.DSN)
	if err != nil {
		xlog.Jupiter().Panic("parse dsn", xlog.FieldName(name), xlog.FieldErr(err))
	}

	var db *gorm.DB
	err = xretry.Do(config.Retry, config.RetryWaitTime, func() error {
		db, err = open(config)
		if err != nil {
			prome.LibHandleCounter.WithLabelValues(prome.TypeMySQL, name+".dial", d.Addr, err.Error()).Inc()
			return errors.New("dial nil" + err.Error())
		}

		if db == nil {
			prome.LibHandleCounter.WithLabelValues(prome.TypeMySQL, name+".dial", d.Addr, "nil db").Inc()
			return errors.New("db nil" + err.Error())
		}
		sql, err := db.DB()
		if err != nil {
			prome.LibHandleCounter.WithLabelValues(prome.TypeMySQL, name+".ping", d.Addr, err.Error()).Inc()
			return errors.New("db " + err.Error())
		}

		if err := sql.Ping(); err != nil {
			prome.LibHandleCounter.WithLabelValues(prome.TypeMySQL, name+".ping", d.Addr, err.Error()).Inc()
			return errors.New("ping " + err.Error())
		}

		return nil
	})
	if err != nil {
		if config.OnDialError == "panic" {
			xlog.Jupiter().Panic("dial mysql db", xlog.FieldErr(err), xlog.FieldAddr(d.Addr), xlog.FieldName(name))
		} else {
			xlog.Jupiter().Panic("dial mysql db", xlog.FieldAddr(d.Addr), xlog.FieldErr(err), xlog.FieldName(name))
		}
	}

	return db
}
