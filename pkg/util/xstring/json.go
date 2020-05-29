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

package xstring

import (
	"reflect"
	"strings"

	jsoniter "github.com/json-iterator/go"

	"github.com/modern-go/reflect2"
)

var _jsonPrettyAPI = jsoniter.Config{
	IndentionStep:                 4,
	MarshalFloatWith6Digits:       false,
	EscapeHTML:                    true,
	SortMapKeys:                   false,
	UseNumber:                     false,
	DisallowUnknownFields:         false,
	TagKey:                        "",
	OnlyTaggedField:               false,
	ValidateJsonRawMessage:        false,
	ObjectFieldMustBeSimpleString: false,
	CaseSensitive:                 false,
}.Froze()

var _jsonAPI = jsoniter.Config{
	SortMapKeys:            true,
	UseNumber:              true,
	CaseSensitive:          true,
	EscapeHTML:             true,
	ValidateJsonRawMessage: true,
}.Froze()

// OmitDefaultAPI ...
var OmitDefaultAPI = jsoniter.Config{
	SortMapKeys:            true,
	UseNumber:              true,
	CaseSensitive:          true,
	EscapeHTML:             true,
	ValidateJsonRawMessage: true,
}.Froze()

func init() {
	OmitDefaultAPI.RegisterExtension(new(emitDefaultExtension))
}

// Json ...
func Json(obj interface{}) string {
	aa, _ := _jsonAPI.Marshal(obj)
	return string(aa)
}

// JsonBytes ...
func JsonBytes(obj interface{}) []byte {
	aa, _ := _jsonAPI.Marshal(obj)
	return aa
}

// PrettyJson ...
func PrettyJson(obj interface{}) string {
	aa, _ := _jsonPrettyAPI.MarshalIndent(obj, "", "    ")
	return string(aa)
}

// PrettyJSONBytes ...
func PrettyJSONBytes(obj interface{}) []byte {
	aa, _ := _jsonPrettyAPI.MarshalIndent(obj, "", "    ")
	return aa
}

type emitDefaultExtension struct {
	jsoniter.DummyExtension
}

// UpdateStructDescriptor ...
func (ed emitDefaultExtension) UpdateStructDescriptor(structDescriptor *jsoniter.StructDescriptor) {
	for _, field := range structDescriptor.Fields {
		var hasOmitEmpty bool
		tagParts := strings.Split(field.Field.Tag().Get("json"), ",")
		for _, tagPart := range tagParts[1:] {
			if tagPart == "omitempty" {
				hasOmitEmpty = true
				break
			}
		}
		if hasOmitEmpty {
			oldField := field.Field
			field.Field = &myfield{oldField}
		}
	}
}

type myfield struct{ reflect2.StructField }

// Tag 不得不用这么骚的操作
func (mf *myfield) Tag() reflect.StructTag {
	return reflect.StructTag(strings.Replace(string(mf.StructField.Tag()), ",omitempty", "", -1))
}
