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

package jupiter

import (
	"github.com/douyu/jupiter/pkg/core/application"
	"github.com/douyu/jupiter/pkg/core/hooks"
)

var (
	//StageAfterStop after app stop
	StageAfterStop uint32 = uint32(hooks.Stage_AfterStop)
	//StageBeforeStop before app stop
	StageBeforeStop = uint32(hooks.Stage_BeforeStop)
)

// Application is the framework's instance, it contains the servers, workers, client and configuration settings.
// Create an instance of Application, by using &Application{}
type Application = application.Application

var New = application.New
var DefaultApp = application.DefaultApp
