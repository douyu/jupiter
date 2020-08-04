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
	"github.com/douyu/jupiter/pkg"
	"github.com/douyu/jupiter/pkg/conf"
	jsoniter "github.com/json-iterator/go"
	"net/http"
	"os"
)

func init() {
	HandleFunc("/configs", func(w http.ResponseWriter, r *http.Request) {
		encoder := json.NewEncoder(w)
		if r.URL.Query().Get("pretty") == "true" {
			encoder.SetIndent("", "    ")
		}
		encoder.Encode(conf.Traverse("."))
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
