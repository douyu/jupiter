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

package xtrace

import (
	"go.opentelemetry.io/otel/propagation"
	"google.golang.org/grpc/metadata"
)

// assert that MetadataReaderWriter implements the TextMapCarrier interface
var _ propagation.TextMapCarrier = (*MetadataReaderWriter)(nil)

// MetadataReaderWriter ...
type MetadataReaderWriter metadata.MD

func (m MetadataReaderWriter) Get(key string) string {
	values := metadata.MD(m).Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (m MetadataReaderWriter) Set(key, value string) {
	metadata.MD(m).Set(key, value)
}

func (m MetadataReaderWriter) Keys() []string {
	keys := make([]string, 0, len(m))
	for k := range metadata.MD(m) {
		keys = append(keys, k)
	}
	return keys
}
