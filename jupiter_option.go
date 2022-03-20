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

package jupiter

import (
	"github.com/douyu/jupiter/pkg/application"
)

type Option = application.Option

type Disable = application.Disable

const (
	DisableParserFlag      Disable = application.DisableParserFlag
	DisableLoadConfig      Disable = application.DisableLoadConfig
	DisableDefaultGovernor Disable = application.DisableDefaultGovernor
)

var WithConfigParser = application.WithConfigParser
var WithDisable = application.WithDisable
