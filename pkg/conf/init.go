// Copyright 2021 rex lv
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

package conf

import (
	"encoding/json"
	"log"
	"net/url"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/douyu/jupiter/pkg/core/hooks"
	"github.com/douyu/jupiter/pkg/flag"
	"gopkg.in/yaml.v3"
)

const DefaultEnvPrefix = "APP_"

func init() {
	flag.Register(&flag.StringFlag{Name: "envPrefix", Usage: "--envPrefix=APP_", Default: DefaultEnvPrefix, Action: func(key string, fs *flag.FlagSet) {
		var envPrefix = fs.String(key)
		defaultConfiguration.LoadEnvironments(envPrefix)
	}})

	flag.Register(&flag.StringFlag{Name: "config", Usage: "--config=config.toml", Action: func(key string, fs *flag.FlagSet) {
		hooks.Do(hooks.Stage_BeforeLoadConfig)

		var configAddr = fs.String(key)
		log.Printf("read config: %s", configAddr)
		datasource, err := NewDataSource(configAddr)
		if err != nil {
			log.Fatalf("build datasource[%s] failed: %v", configAddr, err)
		}

		path := configAddr
		if uri, err := url.ParseRequestURI(configAddr); err == nil {
			path = uri.Path
		}

		unmarshaler := toml.Unmarshal
		switch filepath.Ext(path) {
		case ".toml":
			// default config type
		case ".yaml", ".yml":
			unmarshaler = yaml.Unmarshal
		case ".json":
			unmarshaler = json.Unmarshal
		default:
			log.Fatalf("unsupported config type: %s", filepath.Ext(configAddr))
		}

		if err := LoadFromDataSource(datasource, unmarshaler); err != nil {
			log.Fatalf("load config from datasource[%s] failed: %v", configAddr, err)
		}
		log.Printf("load config from datasource[%s] completely!", configAddr)

		hooks.Do(hooks.Stage_AfterLoadConfig)
	}})

	flag.Register(&flag.StringFlag{Name: "config-tag", Usage: "--config-tag=mapstructure", Default: "mapstructure", Action: func(key string, fs *flag.FlagSet) {
		defaultGetOptions.TagName = fs.String("config-tag")
	}})

	flag.Register(&flag.StringFlag{Name: "config-namespace", Usage: "--config-namespace=jupiter, 配置内建组件的默认命名空间, 默认是jupiter", Default: "jupiter", Action: func(key string, fs *flag.FlagSet) {
		defaultGetOptions.Namespace = fs.String("config-namespace")
	}})

	flag.Register(&flag.BoolFlag{Name: "watch", Usage: "--watch, watch config change event", Default: false, EnvVar: "JUPITER_CONFIG_WATCH", Action: func(key string, fs *flag.FlagSet) {
		log.Printf("load config watch: %v", fs.Bool(key))
	}})
}
