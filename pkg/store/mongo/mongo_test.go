package mongo

import (
	"reflect"
	"testing"

	"github.com/globalsign/mgo"
)

func TestStdNew(t *testing.T) {
	type args struct {
		name string
		opts []interface{}
	}
	tests := []struct {
		name string
		args args
		want *mgo.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StdNew(tt.args.name, tt.args.opts...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("StdNew() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		name   string
		config Config
	}
	tests := []struct {
		name string
		args args
		want *mgo.Session
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.name, tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
