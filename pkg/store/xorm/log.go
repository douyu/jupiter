// Copyright 2020 Douyu
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

package xorm

import (
	"fmt"
	"github.com/douyu/jupiter/pkg/xlog"
	"go.uber.org/zap"
	"log"
	"xorm.io/core"
)

type Logger struct{
	DEBUG   *log.Logger
	ERR     *log.Logger
	INFO    *log.Logger
	WARN    *log.Logger
	level   core.LogLevel
	logger  *xlog.Logger
	core.ILogger
	showSQL bool
}

func NewLogger () *Logger{

	return &Logger{
		level: core.LOG_DEBUG,
		logger: xlog.JupiterLogger,
	}
}

func (l *Logger) Debug(v ...interface{}) {

	if l.level <= core.LOG_INFO {

		l.logger.Debug("mysql" , xlog.FieldMod("xorm"),zap.String("msg", fmt.Sprintf("%s",v...)))
	}
	return
}

func (l *Logger) Debugf(format string, v ...interface{}) {

	if l.level <= core.LOG_INFO {

		l.logger.Debug("mysql" , xlog.FieldMod("xorm"),zap.String("msg", fmt.Sprintf(format,v...)))
	}
	return
}

func (l *Logger) Error(v ...interface{}) {

	if l.level <= core.LOG_INFO {

		l.logger.Error("mysql" , xlog.FieldMod("xorm"),zap.String("msg", fmt.Sprintf("%s",v...)))
	}
	return
}
func (l *Logger) Errorf(format string, v ...interface{}) {

	if l.level <= core.LOG_INFO {

		l.logger.Error("mysql" , xlog.FieldMod("xorm"),zap.String("msg", fmt.Sprintf(format,v...)))
	}
	return
}

func (l *Logger) Info(v ...interface{}) {

	if l.level <= core.LOG_INFO {

		l.logger.Info("mysql" , xlog.FieldMod("xorm"),zap.String("msg", fmt.Sprintf("%s",v...)))
	}
	return
}
func (l *Logger) Infof(format string,v ...interface{}) {

	if l.level <= core.LOG_INFO {

		l.logger.Info("mysql" , xlog.FieldMod("xorm"),zap.String("msg", fmt.Sprintf(format,v...)))
	}
	return
}

func (l *Logger) Warn(v ...interface{}) {

	if l.level <= core.LOG_INFO {

		l.logger.Warn("mysql" , xlog.FieldMod("xorm"),zap.String("msg", fmt.Sprintf("%s",v...)))
	}
	return
}
func (l *Logger) Warnf(format string,v ...interface{}) {

	if l.level <= core.LOG_INFO {

		l.logger.Warn("mysql" , xlog.FieldMod("xorm"),zap.String("msg", fmt.Sprintf(format,v...)))
	}
	return
}

func (l *Logger) Level() core.LogLevel{

	return l.level
}

func (l *Logger) SetLevel(level core.LogLevel){

	l.level = level
}

func (l *Logger) ShowSQL(show ...bool){

	if len(show) == 0 {
		l.showSQL = true
		return
	}
	l.showSQL = show[0]
}

func (l *Logger) IsShowSQL() bool {

	return l.showSQL
}

func createField (msg string ) []xlog.Field {

	var fields = make([]xlog.Field, 0, 8)

	fields = append(fields, zap.String("msg", msg))

	return fields
}