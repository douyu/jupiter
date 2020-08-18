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

import (
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/douyu/jupiter/pkg/util/xcast"
	"github.com/douyu/jupiter/pkg/util/xmap"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
)

// Configuration provides configuration for application.
type Configuration struct {
	mu       sync.RWMutex
	override map[string]interface{}
	keyDelim string

	keyMap    *sync.Map
	onChanges []func(*Configuration)

	watchers map[string][]func(*Configuration)
}

const (
	defaultKeyDelim = "."
)

// New constructs a new Configuration with provider.
func New() *Configuration {
	return &Configuration{
		override:  make(map[string]interface{}),
		keyDelim:  defaultKeyDelim,
		keyMap:    &sync.Map{},
		onChanges: make([]func(*Configuration), 0),
		watchers:  make(map[string][]func(*Configuration)),
	}
}

// SetKeyDelim set keyDelim of a defaultConfiguration instance.
func (c *Configuration) SetKeyDelim(delim string) {
	c.keyDelim = delim
}

// Sub returns new Configuration instance representing a sub tree of this instance.
func (c *Configuration) Sub(key string) *Configuration {
	return &Configuration{
		keyDelim: c.keyDelim,
		override: c.GetStringMap(key),
	}
}

// WriteConfig ...
func (c *Configuration) WriteConfig() error {
	// return c.provider.Write(c.override)
	return nil
}

// OnChange 注册change回调函数
func (c *Configuration) OnChange(fn func(*Configuration)) {
	c.onChanges = append(c.onChanges, fn)
}

// LoadFromDataSource ...
func (c *Configuration) LoadFromDataSource(ds DataSource, unmarshaller Unmarshaller) error {
	content, err := ds.ReadConfig()
	if err != nil {
		return err
	}

	if err := c.Load(content, unmarshaller); err != nil {
		return err
	}

	go func() {
		for range ds.IsConfigChanged() {
			if content, err := ds.ReadConfig(); err == nil {
				_ = c.Load(content, unmarshaller)
				for _, change := range c.onChanges {
					change(c)
				}
			}
		}
	}()

	return nil
}

// Load ...
func (c *Configuration) Load(content []byte, unmarshal Unmarshaller) error {
	configuration := make(map[string]interface{})
	if err := unmarshal(content, &configuration); err != nil {
		return err
	}
	return c.apply(configuration)
}

// Load loads configuration from provided data source.
func (c *Configuration) LoadFromReader(reader io.Reader, unmarshaller Unmarshaller) error {
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return err
	}
	return c.Load(content, unmarshaller)
}

func (c *Configuration) apply(conf map[string]interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	var changes = make(map[string]interface{})

	xmap.MergeStringMap(c.override, conf)
	for k, v := range c.traverse(c.keyDelim) {
		orig, ok := c.keyMap.Load(k)
		if ok && !reflect.DeepEqual(orig, v) {
			changes[k] = v
		}
		c.keyMap.Store(k, v)
	}

	if len(changes) > 0 {
		c.notifyChanges(changes)
	}

	return nil
}

func (c *Configuration) notifyChanges(changes map[string]interface{}) {
	var changedWatchPrefixMap = map[string]struct{}{}

	for watchPrefix := range c.watchers {
		for key := range changes {
			// 前缀匹配即可
			// todo 可能产生错误匹配
			if strings.HasPrefix(key, watchPrefix) {
				changedWatchPrefixMap[watchPrefix] = struct{}{}
			}
		}
	}

	for changedWatchPrefix := range changedWatchPrefixMap {
		for _, handle := range c.watchers[changedWatchPrefix] {
			go handle(c)
		}
	}
}

// Set ...
func (c *Configuration) Set(key string, val interface{}) error {
	paths := strings.Split(key, c.keyDelim)
	lastKey := paths[len(paths)-1]
	m := deepSearch(c.override, paths[:len(paths)-1])
	m[lastKey] = val
	return c.apply(m)
	// c.keyMap.Store(key, val)
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

// Get returns the value associated with the key
func (c *Configuration) Get(key string) interface{} {
	return c.find(key)
}

// GetString returns the value associated with the key as a string with default defaultConfiguration.
func GetString(key string) string {
	return defaultConfiguration.GetString(key)
}

// GetString returns the value associated with the key as a string.
func (c *Configuration) GetString(key string) string {
	return xcast.ToString(c.Get(key))
}

// GetBool returns the value associated with the key as a boolean with default defaultConfiguration.
func GetBool(key string) bool {
	return defaultConfiguration.GetBool(key)
}

// GetBool returns the value associated with the key as a boolean.
func (c *Configuration) GetBool(key string) bool {
	return xcast.ToBool(c.Get(key))
}

// GetInt returns the value associated with the key as an integer with default defaultConfiguration.
func GetInt(key string) int {
	return defaultConfiguration.GetInt(key)
}

// GetInt returns the value associated with the key as an integer.
func (c *Configuration) GetInt(key string) int {
	return xcast.ToInt(c.Get(key))
}

// GetInt64 returns the value associated with the key as an integer with default defaultConfiguration.
func GetInt64(key string) int64 {
	return defaultConfiguration.GetInt64(key)
}

// GetInt64 returns the value associated with the key as an integer.
func (c *Configuration) GetInt64(key string) int64 {
	return xcast.ToInt64(c.Get(key))
}

// GetFloat64 returns the value associated with the key as a float64 with default defaultConfiguration.
func GetFloat64(key string) float64 {
	return defaultConfiguration.GetFloat64(key)
}

// GetFloat64 returns the value associated with the key as a float64.
func (c *Configuration) GetFloat64(key string) float64 {
	return xcast.ToFloat64(c.Get(key))
}

// GetTime returns the value associated with the key as time with default defaultConfiguration.
func GetTime(key string) time.Time {
	return defaultConfiguration.GetTime(key)
}

// GetTime returns the value associated with the key as time.
func (c *Configuration) GetTime(key string) time.Time {
	return xcast.ToTime(c.Get(key))
}

// GetDuration returns the value associated with the key as a duration with default defaultConfiguration.
func GetDuration(key string) time.Duration {
	return defaultConfiguration.GetDuration(key)
}

// GetDuration returns the value associated with the key as a duration.
func (c *Configuration) GetDuration(key string) time.Duration {
	return xcast.ToDuration(c.Get(key))
}

// GetStringSlice returns the value associated with the key as a slice of strings with default defaultConfiguration.
func GetStringSlice(key string) []string {
	return defaultConfiguration.GetStringSlice(key)
}

// GetStringSlice returns the value associated with the key as a slice of strings.
func (c *Configuration) GetStringSlice(key string) []string {
	return xcast.ToStringSlice(c.Get(key))
}

// GetSlice returns the value associated with the key as a slice of strings with default defaultConfiguration.
func GetSlice(key string) []interface{} {
	return defaultConfiguration.GetSlice(key)
}

// GetSlice returns the value associated with the key as a slice of strings.
func (c *Configuration) GetSlice(key string) []interface{} {
	return xcast.ToSlice(c.Get(key))
}

// GetStringMap returns the value associated with the key as a map of interfaces with default defaultConfiguration.
func GetStringMap(key string) map[string]interface{} {
	return defaultConfiguration.GetStringMap(key)
}

// GetStringMap returns the value associated with the key as a map of interfaces.
func (c *Configuration) GetStringMap(key string) map[string]interface{} {
	return xcast.ToStringMap(c.Get(key))
}

// GetStringMapString returns the value associated with the key as a map of strings with default defaultConfiguration.
func GetStringMapString(key string) map[string]string {
	return defaultConfiguration.GetStringMapString(key)
}

// GetStringMapString returns the value associated with the key as a map of strings.
func (c *Configuration) GetStringMapString(key string) map[string]string {
	return xcast.ToStringMapString(c.Get(key))
}

// GetSliceStringMap returns the value associated with the slice of maps.
func (c *Configuration) GetSliceStringMap(key string) []map[string]interface{} {
	return xcast.ToSliceStringMap(c.Get(key))
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings with default defaultConfiguration.
func GetStringMapStringSlice(key string) map[string][]string {
	return defaultConfiguration.GetStringMapStringSlice(key)
}

// GetStringMapStringSlice returns the value associated with the key as a map to a slice of strings.
func (c *Configuration) GetStringMapStringSlice(key string) map[string][]string {
	return xcast.ToStringMapStringSlice(c.Get(key))
}

// UnmarshalWithExpect unmarshal key, returns expect if failed
func UnmarshalWithExpect(key string, expect interface{}) interface{} {
	return defaultConfiguration.UnmarshalWithExpect(key, expect)
}

// UnmarshalWithExpect unmarshal key, returns expect if failed
func (c *Configuration) UnmarshalWithExpect(key string, expect interface{}) interface{} {
	err := c.UnmarshalKey(key, expect)
	if err != nil {
		return expect
	}
	return expect
}

// UnmarshalKey takes a single key and unmarshal it into a Struct with default defaultConfiguration.
func UnmarshalKey(key string, rawVal interface{}, opts ...GetOption) error {
	return defaultConfiguration.UnmarshalKey(key, rawVal, opts...)
}

// ErrInvalidKey ...
var ErrInvalidKey = errors.New("invalid key, maybe not exist in config")

// UnmarshalKey takes a single key and unmarshal it into a Struct.
func (c *Configuration) UnmarshalKey(key string, rawVal interface{}, opts ...GetOption) error {
	var options = defaultGetOptions
	for _, opt := range opts {
		opt(&options)
	}

	config := mapstructure.DecoderConfig{
		DecodeHook: mapstructure.StringToTimeDurationHookFunc(),
		Result:     rawVal,
		TagName:    options.TagName,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return err
	}
	if key == "" {
		c.mu.RLock()
		defer c.mu.RUnlock()
		return decoder.Decode(c.override)
	}

	value := c.Get(key)
	if value == nil {
		return errors.Wrap(ErrInvalidKey, key)
	}

	return decoder.Decode(value)
}

func (c *Configuration) find(key string) interface{} {
	dd, ok := c.keyMap.Load(key)
	if ok {
		return dd
	}

	paths := strings.Split(key, c.keyDelim)
	c.mu.RLock()
	defer c.mu.RUnlock()
	m := xmap.DeepSearchInMap(c.override, paths[:len(paths)-1]...)
	dd = m[paths[len(paths)-1]]
	c.keyMap.Store(key, dd)
	return dd
}

func lookup(prefix string, target map[string]interface{}, data map[string]interface{}, sep string) {
	for k, v := range target {
		pp := fmt.Sprintf("%s%s%s", prefix, sep, k)
		if prefix == "" {
			pp = k
		}
		if dd, err := xcast.ToStringMapE(v); err == nil {
			lookup(pp, dd, data, sep)
		} else {
			data[pp] = v
		}
	}
}

func (c *Configuration) traverse(sep string) map[string]interface{} {
	data := make(map[string]interface{})
	lookup("", c.override, data, sep)
	return data
}
