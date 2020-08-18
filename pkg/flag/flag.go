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

package flag

import (
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

var (
	flagset *FlagSet
)

func init() {
	// procName := filepath.Base(os.Args[0])
	// nfs := flag.NewFlagSet(procName, flag.ExitOnError)
	flagset = &FlagSet{
		FlagSet:  flag.CommandLine,
		flags:    defaultFlags,
		actions:  make(map[string]func(string, *FlagSet)),
		environs: make(map[string]string),
	}
	testing.Init()
}

// Flag ...
type (
	// Flag defines application flag.
	Flag interface {
		Apply(*FlagSet)
	}

	// FlagSet wraps a set of Flags.
	FlagSet struct {
		*flag.FlagSet
		flags    []Flag
		actions  map[string]func(string, *FlagSet)
		environs map[string]string
	}
)

// Register ...
func Register(fs ...Flag) {
	flagset.Register(fs...)
}

// Register ...
func (fs *FlagSet) Register(flags ...Flag) {
	fs.flags = append(fs.flags, flags...)
}

// With adds flags to the flagset.
func With(fs ...Flag) { flagset.With(fs...) }

// With adds flags to provided flagset.
func (fs *FlagSet) With(flags ...Flag) {
	fs.flags = append(fs.flags, flags...)
}

// Parse parses the flagset.
func Parse() error {
	return flagset.Parse()
}

// Lookup lookup flag value by name
// priority: flag > default > env
func (fs *FlagSet) Lookup(name string) *flag.Flag {
	flag := fs.FlagSet.Lookup(name)
	if flag != nil {
		if flag.Value.String() == "" {
			if env, ok := fs.environs[name]; ok {
				flag.Value.Set(env)
			}
		}
		if flag.Value.String() == "" {
			flag.Value.Set(flag.DefValue)
		}
	}
	return flag
}

// Parse parses provided flagset.
func (fs *FlagSet) Parse() error {
	if fs.Parsed() {
		return nil
	}
	for _, f := range fs.flags {
		f.Apply(fs)
	}

	if err := fs.FlagSet.Parse(os.Args[1:]); err != nil {
		return err
	}

	fs.FlagSet.Visit(func(f *flag.Flag) {
		if action, ok := fs.actions[f.Name]; ok && action != nil {
			action(f.Name, fs)
		}
		if env, ok := fs.environs[f.Name]; ok {
			fs.environs[f.Name] = env
		}
	})

	return nil
}

// BoolE parses bool flag of the flagset with error returned.
func BoolE(name string) (bool, error) { return flagset.BoolE(name) }

// BoolE parses bool flag of provided flagset with error returned.
func (fs *FlagSet) BoolE(name string) (bool, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return strconv.ParseBool(flag.Value.String())
	}

	return false, fmt.Errorf("undefined flag name: %s", name)
}

// Bool parses bool flag of the flagset.
func Bool(name string) bool { return flagset.Bool(name) }

// Bool parses bool flag of provided flagset.
func (fs *FlagSet) Bool(name string) bool {
	ret, _ := fs.BoolE(name)
	return ret
}

// StringE parses string flag of the flagset with error returned.
func StringE(name string) (string, error) { return flagset.StringE(name) }

// StringE parses string flag of provided flagset with error returned.
func (fs *FlagSet) StringE(name string) (string, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return flag.Value.String(), nil
	}

	return "", fmt.Errorf("undefined flag name: %s", name)
}

// String parses string flag of the flagset.
func String(name string) string { return flagset.String(name) }

// String parses string flag of provided flagset.
func (fs *FlagSet) String(name string) string {
	ret, _ := fs.StringE(name)
	return ret
}

// IntE parses int flag of the flagset with error returned.
func IntE(name string) (int64, error) { return flagset.IntE(name) }

// IntE parses int flag of provided flagset with error returned.
func (fs *FlagSet) IntE(name string) (int64, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return strconv.ParseInt(flag.Value.String(), 10, 64)
	}

	return 0, fmt.Errorf("undefined flag name: %s", name)
}

// Int parses int flag of the flagset.
func Int(name string) int64 { return flagset.Int(name) }

// Int parses int flag of provided flagset.
func (fs *FlagSet) Int(name string) int64 {
	ret, _ := fs.IntE(name)
	return ret
}

// UintE parses int flag of the flagset with error returned.
func UintE(name string) (uint64, error) { return flagset.UintE(name) }

// UintE parses int flag of provided flagset with error returned.
func (fs *FlagSet) UintE(name string) (uint64, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return strconv.ParseUint(flag.Value.String(), 10, 64)
	}

	return 0, fmt.Errorf("undefined flag name: %s", name)
}

// Uint parses int flag of the flagset.
func Uint(name string) uint64 { return flagset.Uint(name) }

// Uint parses int flag of provided flagset.
func (fs *FlagSet) Uint(name string) uint64 {
	ret, _ := fs.UintE(name)
	return ret
}

// Float64E parses int flag of the flagset with error returned.
func Float64E(name string) (float64, error) { return flagset.Float64E(name) }

// Float64E parses int flag of provided flagset with error returned.
func (fs *FlagSet) Float64E(name string) (float64, error) {
	flag := fs.Lookup(name)
	if flag != nil {
		return strconv.ParseFloat(flag.Value.String(), 64)
	}

	return 0.0, fmt.Errorf("undefined flag name: %s", name)
}

// Float64 parses int flag of the flagset.
func Float64(name string) float64 { return flagset.Float64(name) }

// Float64 parses int flag of provided flagset.
func (fs *FlagSet) Float64(name string) float64 {
	ret, _ := fs.Float64E(name)
	return ret
}

// BoolFlag is a bool flag implements of Flag interface.
type BoolFlag struct {
	Name     string
	Usage    string
	EnvVar   string
	Default  bool
	Variable *bool
	Action   func(string, *FlagSet)
}

// Apply implements of Flag Apply function.
func (f *BoolFlag) Apply(set *FlagSet) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)
		if f.Variable != nil {
			set.FlagSet.BoolVar(f.Variable, field, f.Default, f.Usage)
		}

		set.FlagSet.Bool(field, f.Default, f.Usage)
		set.actions[field] = f.Action
		set.environs[field] = os.Getenv(f.EnvVar)
	}
}

// StringFlag is a string flag implements of Flag interface.
type StringFlag struct {
	Name     string
	Usage    string
	EnvVar   string
	Default  string
	Variable *string
	Action   func(string, *FlagSet)
}

// Apply implements of Flag Apply function.
func (f *StringFlag) Apply(set *FlagSet) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)
		if f.Variable != nil {
			set.FlagSet.StringVar(f.Variable, field, f.Default, f.Usage)
		}
		set.FlagSet.String(field, f.Default, f.Usage)
		set.actions[field] = f.Action
		set.environs[field] = os.Getenv(f.EnvVar)
	}
}

// IntFlag is an int flag implements of Flag interface.
type IntFlag struct {
	Name     string
	Usage    string
	Default  int
	Variable *int
	Action   func(string, *FlagSet)
}

// Apply implements of Flag Apply function.
func (f *IntFlag) Apply(set *FlagSet) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)
		if f.Variable != nil {
			set.FlagSet.IntVar(f.Variable, field, f.Default, f.Usage)
		}
		set.FlagSet.Int(field, f.Default, f.Usage)
		set.actions[field] = f.Action
	}
}

// UintFlag is an uint flag implements of Flag interface.
type UintFlag struct {
	Name     string
	Usage    string
	Default  uint
	Variable *uint
	Action   func(string, *FlagSet)
}

// Apply implements of Flag Apply function.
func (f *UintFlag) Apply(set *FlagSet) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)
		if f.Variable != nil {
			set.FlagSet.UintVar(f.Variable, field, f.Default, f.Usage)
		}
		set.FlagSet.Uint(field, f.Default, f.Usage)
		set.actions[field] = f.Action
	}
}

// Float64Flag is a float flag implements of Flag interface.
type Float64Flag struct {
	Name     string
	Usage    string
	Default  float64
	Variable *float64
	Action   func(string, *FlagSet)
}

// Apply implements of Flag Apply function.
func (f *Float64Flag) Apply(set *FlagSet) {
	for _, field := range strings.Split(f.Name, ",") {
		field = strings.TrimSpace(field)
		if f.Variable != nil {
			set.FlagSet.Float64Var(f.Variable, field, f.Default, f.Usage)
		}
		set.FlagSet.Float64(field, f.Default, f.Usage)
		set.actions[field] = f.Action
	}
}
