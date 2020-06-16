package gorm

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/douyu/jupiter/pkg/trace"
	"github.com/douyu/jupiter/pkg/util/xcolor"
	"github.com/douyu/jupiter/pkg/xlog"
)

// Handler ...
type Handler func(*Scope)

// Interceptor ...
type Interceptor func(*DSN, string, *Config) func(next Handler) Handler

func debugInterceptor(dsn *DSN, op string, options *Config) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(scope *Scope) {
			fmt.Printf("%-50s[%s] => %s\n", xcolor.Green(dsn.Addr+"/"+dsn.DBName), time.Now().Format("04:05.000"), xcolor.Green("Send: "+logSQL(scope.SQL, scope.SQLVars, true)))
			next(scope)
			if scope.HasError() {
				fmt.Printf("%-50s[%s] => %s\n", xcolor.Red(dsn.Addr+"/"+dsn.DBName), time.Now().Format("04:05.000"), xcolor.Red("Erro: "+scope.DB().Error.Error()))
			} else {
				fmt.Printf("%-50s[%s] => %s\n", xcolor.Green(dsn.Addr+"/"+dsn.DBName), time.Now().Format("04:05.000"), xcolor.Green("Affected: "+strconv.Itoa(int(scope.DB().RowsAffected))))
			}
		}
	}
}

func metricInterceptor(dsn *DSN, op string, options *Config) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(scope *Scope) {
			beg := time.Now()
			next(scope)
			cost := time.Since(beg)

			// error metric
			if scope.HasError() {
				// todo sql语句，需要转换成脱密状态才能记录到日志
				if scope.DB().Error != ErrRecordNotFound {
					options.logger.Error("mysql err", xlog.FieldErr(scope.DB().Error), xlog.FieldName(dsn.DBName+"."+scope.TableName()), xlog.FieldMethod(op))
				} else {
					options.logger.Warn("record not found", xlog.FieldErr(scope.DB().Error), xlog.FieldName(dsn.DBName+"."+scope.TableName()), xlog.FieldMethod(op))
				}
			}

			if options.SlowThreshold > time.Duration(0) && options.SlowThreshold < cost {
				options.logger.Error(
					"slow",
					xlog.FieldErr(errSlowCommand),
					xlog.FieldMethod(op),
					xlog.FieldExtMessage(logSQL(scope.SQL, scope.SQLVars, options.DetailSQL)),
					xlog.FieldAddr(dsn.Addr),
					xlog.FieldName(dsn.DBName+"."+scope.TableName()),
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

func traceInterceptor(dsn *DSN, op string, options *Config) func(Handler) Handler {
	return func(next Handler) Handler {
		return func(scope *Scope) {
			if val, ok := scope.Get("_context"); ok {
				if ctx, ok := val.(context.Context); ok {
					span, _ := trace.StartSpanFromContext(
						ctx,
						op,
						trace.TagComponent("mysql"),
						trace.TagSpanKind("client"),
					)
					defer span.Finish()

					// 延迟执行 scope.CombinedConditionSql() 避免sqlVar被重复追加
					next(scope)

					span.SetTag("sql.inner", dsn.DBName)
					span.SetTag("sql.addr", dsn.Addr)
					span.SetTag("span.kind", "client")
					span.SetTag("peer.service", "mysql")
					span.LogFields(trace.String("sql.query", logSQL(scope.SQL, scope.SQLVars, options.DetailSQL)))
					return
				}
			}

			next(scope)
		}
	}
}
