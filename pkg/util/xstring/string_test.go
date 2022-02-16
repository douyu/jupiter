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
	"testing"
)

func TestKickEmpty(t *testing.T) {
	type args struct {
		ss []string
	}
	tests := []struct {
		name string
		args args
		want Strings
	}{
		// TODO: Add test cases.
		{
			name: "testing",
			args: args{
				ss: []string{"", "1", "2", ""},
			},
			want: Strings{"1", "2"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := KickEmpty(tt.args.ss); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KickEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKick(t *testing.T) {
	type args struct {
		ss     []string
		remove func(string) bool
	}
	tests := []struct {
		name string
		args args
		want Strings
	}{
		// TODO: Add test cases.
		{
			name: "testing",
			args: args{
				ss: []string{"0", "1", "2", "3"},
				remove: func(item string) bool {
					return item == "1" || item == "2"
				},
			},
			want: Strings{"0", "3"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Kick(tt.args.ss, tt.args.remove); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("KickEmpty() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyBlank(t *testing.T) {
	type args struct {
		ss []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "1",
			args: args{
				ss: []string{"", "1", "2", ""},
			},
			want: true,
		},
		{
			name: "1",
			args: args{
				ss: []string{"1", "2"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AnyBlank(tt.args.ss); got != tt.want {
				t.Errorf("AnyBlank() = %v, want %v", got, tt.want)
			}
		})
	}
}
