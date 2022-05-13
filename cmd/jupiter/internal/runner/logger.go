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

package runner

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

var (
	rawColor = "raw"
	// TODO: support more colors
	colorMap = map[string]color.Attribute{
		"red":     color.FgRed,
		"green":   color.FgGreen,
		"yellow":  color.FgYellow,
		"blue":    color.FgBlue,
		"magenta": color.FgMagenta,
		"cyan":    color.FgCyan,
		"white":   color.FgWhite,
	}
)

type logFunc func(string, ...interface{})

type logger struct {
	config  *config
	colors  map[string]string
	loggers map[string]logFunc
}

func newLogger(cfg *config) *logger {
	if cfg == nil {
		return nil
	}

	colors := cfg.colorInfo()
	loggers := make(map[string]logFunc, len(colors))
	for name, nameColor := range colors {
		loggers[name] = newLogFunc(nameColor, cfg.Log)
	}
	loggers["default"] = defaultLogger()
	return &logger{
		config:  cfg,
		colors:  colors,
		loggers: loggers,
	}
}

func newLogFunc(colorname string, cfg cfgLog) logFunc {
	return func(msg string, v ...interface{}) {
		// There are some escape sequences to format color in terminal, so cannot
		// just trim new line from right.
		msg = strings.Replace(msg, "\n", "", -1)
		msg = strings.TrimSpace(msg)
		if len(msg) == 0 {
			return
		}
		// TODO: filter msg by regex
		msg = msg + "\n"
		if cfg.AddTime {
			t := time.Now().Format("15:04:05")
			msg = fmt.Sprintf("[%s] %s", t, msg)
		}
		if colorname == rawColor {
			_, _ = fmt.Fprintf(os.Stdout, msg, v...)
		} else {
			_, _ = color.New(getColor(colorname)).Fprintf(color.Output, msg, v...)
		}
	}
}

func getColor(name string) color.Attribute {
	if v, ok := colorMap[name]; ok {
		return v
	}
	return color.FgWhite
}

func (l *logger) main() logFunc {
	return l.getLogger("main")
}

func (l *logger) build() logFunc {
	return l.getLogger("build")
}

func (l *logger) runner() logFunc {
	return l.getLogger("runner")
}

func (l *logger) watcher() logFunc {
	return l.getLogger("watcher")
}

func rawLogger() logFunc {
	return newLogFunc("raw", defaultConfig().Log)
}

func defaultLogger() logFunc {
	return newLogFunc("white", defaultConfig().Log)
}

func (l *logger) getLogger(name string) logFunc {
	v, ok := l.loggers[name]
	if !ok {
		return rawLogger()
	}
	return v
}
