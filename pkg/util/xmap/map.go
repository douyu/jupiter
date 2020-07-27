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

package xmap

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/douyu/jupiter/pkg/util/xcast"
	"github.com/mitchellh/mapstructure"
)

// Unmarshaller ...
type Unmarshaller = func([]byte, interface{}) error

// KeySpliter ...
var KeySpliter = "."

// FlatMap ...
type FlatMap struct {
	data   map[string]interface{}
	mu     sync.RWMutex
	keyMap sync.Map
}

// NewFlatMap ...
func NewFlatMap() *FlatMap {
	return &FlatMap{
		data: make(map[string]interface{}),
	}
}

// Load ...
func (flat *FlatMap) Load(content []byte, unmarshal Unmarshaller) error {
	data := make(map[string]interface{})
	if err := unmarshal(content, &data); err != nil {
		return err
	}
	return flat.apply(data)
}

func (flat *FlatMap) apply(data map[string]interface{}) error {
	flat.mu.Lock()
	defer flat.mu.Unlock()

	MergeStringMap(flat.data, data)
	var changes = make(map[string]interface{})
	for k, v := range flat.traverse(KeySpliter) {
		orig, ok := flat.keyMap.Load(k)
		if ok && !reflect.DeepEqual(orig, v) {
			changes[k] = v
		}
		flat.keyMap.Store(k, v)
	}

	return nil
}

// Set ...
func (flat *FlatMap) Set(key string, val interface{}) error {
	paths := strings.Split(key, KeySpliter)
	lastKey := paths[len(paths)-1]
	m := deepSearch(flat.data, paths[:len(paths)-1])
	m[lastKey] = val
	return flat.apply(m)
}

// Get returns the value associated with the key
func (flat *FlatMap) Get(key string) interface{} {
	return flat.find(key)
}

// GetString returns the value associated with the key as a string.
func (flat *FlatMap) GetString(key string) string {
	return xcast.ToString(flat.Get(key))
}

// GetBool returns the value associated with the key as a boolean.
func (flat *FlatMap) GetBool(key string) bool {
	return xcast.ToBool(flat.Get(key))
}

// GetInt returns the value associated with the key as an integer.
func (flat *FlatMap) GetInt(key string) int {
	return xcast.ToInt(flat.Get(key))
}

// GetInt64 returns the value associated with the key as an integer.
func (flat *FlatMap) GetInt64(key string) int64 {
	return xcast.ToInt64(flat.Get(key))
}

// GetFloat64 returns the value associated with the key as a float64.
func (flat *FlatMap) GetFloat64(key string) float64 {
	return xcast.ToFloat64(flat.Get(key))
}

// GetTime returns the value associated with the key as time.
func (flat *FlatMap) GetTime(key string) time.Time {
	return xcast.ToTime(flat.Get(key))
}

// GetDuration returns the value associated with the key as a duration.
func (flat *FlatMap) GetDuration(key string) time.Duration {
	return xcast.ToDuration(flat.Get(key))
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (flat *FlatMap) GetStringSlice(key string) []string {
	return xcast.ToStringSlice(flat.Get(key))
}

// GetSlice returns the value associated with the key as a slice of strings.
func (flat *FlatMap) GetSlice(key string) []interface{} {
	return xcast.ToSlice(flat.Get(key))
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (flat *FlatMap) GetStringMap(key string) map[string]interface{} {
	return xcast.ToStringMap(flat.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (flat *FlatMap) GetStringMapString(key string) map[string]string {
	return xcast.ToStringMapString(flat.Get(key))
}

// GetSliceStringMap returns the value associated with the slice of maps.
func (flat *FlatMap) GetSliceStringMap(key string) []map[string]interface{} {
	return xcast.ToSliceStringMap(flat.Get(key))
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (flat *FlatMap) GetStringMapStringSlice(key string) map[string][]string {
	return xcast.ToStringMapStringSlice(flat.Get(key))
}

// UnmarshalKey takes a single key and unmarshal it into a Struct.
func (flat *FlatMap) UnmarshalKey(key string, rawVal interface{}, tagName string) error {
	config := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     rawVal,
		TagName:    tagName,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}
	if key == "" {
		flat.mu.RLock()
		defer flat.mu.RUnlock()
		return decoder.Decode(flat.data)
	}

	value := flat.Get(key)
	if value == nil {
		return fmt.Errorf("invalid key %s, maybe not exist in config", key)
	}

	return decoder.Decode(value)
}

// Reset ...
func (flat *FlatMap) Reset() {
	flat.mu.Lock()
	defer flat.mu.Unlock()

	flat.data = make(map[string]interface{})
	// erase map
	flat.keyMap.Range(func(key interface{}, value interface{}) bool {
		flat.keyMap.Delete(key)
		return true
	})
}

func (flat *FlatMap) find(key string) interface{} {
	dd, ok := flat.keyMap.Load(key)
	if ok {
		return dd
	}

	paths := strings.Split(key, KeySpliter)
	flat.mu.RLock()
	defer flat.mu.RUnlock()
	m := DeepSearchInMap(flat.data, paths[:len(paths)-1]...)
	dd = m[paths[len(paths)-1]]
	flat.keyMap.Store(key, dd)
	return dd
}

func lookup(prefix string, target map[string]interface{}, data map[string]interface{}, sep string) {
	for k, v := range target {
		pp := fmt.Sprintf("%s%s%s", prefix, sep, k)
		if prefix == "" {
			pp = fmt.Sprintf("%s", k)
		}
		if dd, err := xcast.ToStringMapE(v); err == nil {
			lookup(pp, dd, data, sep)
		} else {
			data[pp] = v
		}
	}
}

func (flat *FlatMap) traverse(sep string) map[string]interface{} {
	data := make(map[string]interface{})
	lookup("", flat.data, data, sep)
	return data
}

func deepSearch(m map[string]interface{}, path []string) map[string]interface{} {
	for _, k := range path {
		m2, ok := m[k]
		if !ok {
			m3 := make(map[string]interface{})
			m[k] = m3
			m = m3
			continue
		}
		m3, ok := m2.(map[string]interface{})
		if !ok {
			m3 = make(map[string]interface{})
			m[k] = m3
		}
		m = m3
	}
	return m
}
