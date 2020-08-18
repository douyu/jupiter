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

package constant

const (
	// EnvKeySentinelLogDir ...
	EnvKeySentinelLogDir = "SENTINEL_LOG_DIR"
	// EnvKeySentinelAppName ...
	EnvKeySentinelAppName = "SENTINEL_APP_NAME"
)

const (
	// EnvAppName ...
	EnvAppName = "APP_NAME"
	// EnvDeployment ...
	EnvDeployment = "APP_DEPLOYMENT"

	EnvAppLogDir   = "APP_LOG_DIR"
	EnvAppMode     = "APP_MODE"
	EnvAppRegion   = "APP_REGION"
	EnvAppZone     = "APP_ZONE"
	EnvAppHost     = "APP_HOST"
	EnvAppInstance = "APP_INSTANCE" // application unique instance id.
)

const (
	// DefaultDeployment ...
	DefaultDeployment = ""
	// DefaultRegion ...
	DefaultRegion = ""
	// DefaultZone ...
	DefaultZone = ""
)

const (
	// KeyBalanceGroup ...
	KeyBalanceGroup = "__group"

	// DefaultBalanceGroup ...
	DefaultBalanceGroup = "default"
)
