package xorm

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"xorm.io/core"
)

type (
	Db = *xorm.Engine
)

func Open(driverName string, options *Config) ( *xorm.Engine , error ) {

	var err error

	inner, err := xorm.NewEngine(driverName, options.DSN)

	if err != nil {

		return nil ,err
	}

	inner.SetLogger(NewLogger())

	inner.ShowSQL(options.ShowSql)
	if options.Debug {

		inner.SetLogLevel(core.LOG_DEBUG)
	}else{

		inner.SetLogLevel(core.LOG_OFF)
	}



	inner.SetMaxIdleConns(options.MaxIdleConns)

	inner.SetMaxOpenConns(options.MaxOpenConns)

	return inner , nil

}