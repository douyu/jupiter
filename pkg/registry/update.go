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

package registry

import (
	"github.com/douyu/jupiter/pkg/server"
)

type Operation uint8

const (
	// Add indicates an Endpoint is added.
	Add Operation = iota
	// Delete indicates an existing address is deleted.
	Delete
)

// Update defines a name resolution update. Notice that it is not valid having both
// empty string Addr and nil Metadata in an Update.
type Update struct {
	// Op indicates the operation of the update.
	Op Operation
	// Addr is the updated address. It is empty string if there is no address update.
	Addr string
	// Metadata is the updated metadata. It is nil if there is no metadata update.
	// Metadata is not required for a custom naming implementation.
	MetadataX *server.ServiceInfo
}
