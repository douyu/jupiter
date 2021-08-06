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

package conf

// GetOption ...
type (
	GetOption  func(o *GetOptions)
	GetOptions struct {
		TagName   string
		Namespace string
		Module    string
	}
)

var defaultGetOptions = GetOptions{
	TagName: "mapstructure",
}

// 设置Tag
func TagName(tag string) GetOption {
	return func(o *GetOptions) {
		o.TagName = tag
	}
}

func TagNameJSON() GetOption {
	return TagName("json")
}

func TagNameTOML() GetOption {
	return TagName("toml")
}

func TagNameYAML() GetOption {
	return TagName("yaml")
}

func BuildinModule(module string) GetOption {
	return func(o *GetOptions) {
		o.Namespace = "jupiter"
		o.Module = module
	}
}

func Namespace(namespace string) GetOption {
	return func(o *GetOptions) {
		o.Namespace = namespace
	}
}

func Module(module string) GetOption {
	return func(o *GetOptions) {
		o.Module = module
	}
}
