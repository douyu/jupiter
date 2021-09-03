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

package governor

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/pprof"
	"os"
	"runtime/debug"

	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/conf"
	"github.com/douyu/jupiter/pkg/util/xstring"
	"github.com/douyu/jupiter/pkg/xlog"
	jsoniter "github.com/json-iterator/go"
)

var (
	// DefaultServeMux ...
	DefaultServeMux  = http.NewServeMux()
	routes           = []string{}
	DefaultConfigKey = "jupiter.server.govern"
)

func init() {
	conf.OnLoaded(func(c *conf.Configuration) {
		log.Print("hook config, init runtime(governor)")
		if c.Get(DefaultConfigKey) == nil {
			return
		}
		var config = DefaultConfig()
		if err := c.UnmarshalKey(DefaultConfigKey, &config); err != nil {
			config.logger.Panic("govern server parse config panic",
				xlog.FieldErr(err), xlog.FieldKey(DefaultConfigKey),
				xlog.FieldValueAny(config),
			)
		}
	})
	// 获取全部治理路由
	HandleFunc("/routes", func(resp http.ResponseWriter, req *http.Request) {
		json.NewEncoder(resp).Encode(routes)
	})

	HandleFunc("/debug/pprof/", pprof.Index)
	HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	HandleFunc("/debug/pprof/profile", pprof.Profile)
	HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	HandleFunc("/debug/pprof/trace", pprof.Trace)

	if info, ok := debug.ReadBuildInfo(); ok {
		HandleFunc("/modInfo", func(w http.ResponseWriter, r *http.Request) {
			encoder := json.NewEncoder(w)
			if r.URL.Query().Get("pretty") == "true" {
				encoder.SetIndent("", "    ")
			}
			_ = encoder.Encode(info)
		})
	}
	HandleFunc("/configs", func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		if r.URL.Query().Get("pretty") == "true" {
			encoder.SetIndent("", "    ")
		}
		encoder.Encode(conf.Traverse("."))
	})

	HandleFunc("/debug/config", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write(xstring.PrettyJSONBytes(conf.Traverse(".")))
	})

	HandleFunc("/debug/env", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_ = jsoniter.NewEncoder(w).Encode(os.Environ())
	})

	HandleFunc("/build/info", func(w http.ResponseWriter, r *http.Request) {
		serverStats := map[string]string{
			"name":           pkg.Name(),
			"appID":          pkg.AppID(),
			"appMode":        pkg.AppMode(),
			"appVersion":     pkg.AppVersion(),
			"jupiterVersion": pkg.JupiterVersion(),
			"buildUser":      pkg.BuildUser(),
			"buildHost":      pkg.BuildHost(),
			"buildTime":      pkg.BuildTime(),
			"startTime":      pkg.StartTime(),
			"hostName":       pkg.HostName(),
			"goVersion":      pkg.GoVersion(),
		}
		_ = jsoniter.NewEncoder(w).Encode(serverStats)
	})
}

// HandleFunc ...
func HandleFunc(pattern string, handler http.HandlerFunc) {
	// todo: 增加安全管控
	DefaultServeMux.HandleFunc(pattern, handler)
	routes = append(routes, pattern)
}
