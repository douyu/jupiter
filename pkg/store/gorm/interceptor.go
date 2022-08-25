package gorm

import (
	"errors"
	"time"

	prome "github.com/douyu/jupiter/pkg/metric"
	"github.com/douyu/jupiter/pkg/xlog"
	"gorm.io/gorm"
)

type Handler func(*gorm.DB)
type Interceptor func(dsn *DSN, op string, options *Config, next Handler) Handler

var errSlowCommand = errors.New("mysql slow command")

func metricInterceptor() Interceptor {
	return func(dsn *DSN, op string, options *Config, next Handler) Handler {
		return func(scope *gorm.DB) {
			beg := time.Now()
			next(scope)
			cost := time.Since(beg)

			// error metric
			if scope.Error != nil {
				prome.LibHandleCounter.WithLabelValues(prome.TypeMySQL, dsn.DBName+"."+scope.Name(), dsn.Addr, getStatement(scope.Error.Error())).Inc()

				if scope.Error != gorm.ErrRecordNotFound {
					xlog.Jupiter().Error("mysql err", xlog.FieldErr(scope.Error), xlog.FieldName(dsn.DBName+"."+scope.Name()), xlog.FieldMethod(op))
				} else {
					xlog.Jupiter().Warn("record not found", xlog.FieldErr(scope.Error), xlog.FieldName(dsn.DBName+"."+scope.Name()), xlog.FieldMethod(op))
				}
			} else {
				prome.LibHandleCounter.WithLabelValues(prome.TypeMySQL, dsn.DBName+"."+scope.Name(), dsn.Addr, "OK").Inc()
			}

			prome.LibHandleHistogram.WithLabelValues(prome.TypeMySQL, dsn.DBName+"."+scope.Name(), dsn.Addr).Observe(cost.Seconds())

			if options.SlowThreshold > time.Duration(0) && options.SlowThreshold < cost {
				xlog.Jupiter().Error(
					"slow",
					xlog.FieldErr(errSlowCommand),
					xlog.FieldMethod(op),
					xlog.FieldExtMessage(logSQL(scope.Statement.SQL.String(), scope.Statement.Vars, options.DetailSQL)),
					xlog.FieldAddr(dsn.Addr),
					xlog.FieldName(dsn.DBName+"."+scope.Name()),
					xlog.FieldCost(cost),
				)
			}
		}
	}
}

func logSQL(sql string, args []interface{}, containArgs bool) string {
	if containArgs {
		return bindSQL(sql, args)
	}
	return sql
}

func traceInterceptor() Interceptor {
	return func(dsn *DSN, op string, options *Config, next Handler) Handler {
		return func(scope *gorm.DB) {
			// if ctx := scope.Statement.Context; ctx != nil {
			// 	tracer := xtrace.NewTracer(trace.SpanKindClient)

			// 	span, _ := trace.Start(
			// 		ctx,
			// 		"GORM",
			// 		trace.TagComponent("mysql"),
			// 		trace.TagSpanKind("client"),
			// 	)
			// 	defer span.Finish()

			// 	// 延迟执行 scope.CombinedConditionSql() 避免sqlVar被重复追加
			// 	next(scope)

			// 	span.SetTag("sql.inner", dsn.DBName)
			// 	span.SetTag("sql.addr", dsn.Addr)
			// 	span.SetTag("span.kind", "client")
			// 	span.SetTag("peer.service", "mysql")
			// 	span.SetTag("db.instance", dsn.DBName)
			// 	span.SetTag("peer.address", dsn.Addr)
			// 	span.SetTag("peer.statement", logSQL(scope.Statement.SQL.String(), scope.Statement.Vars, options.DetailSQL))

			// 	return
			// }

			next(scope)
		}
	}
}
