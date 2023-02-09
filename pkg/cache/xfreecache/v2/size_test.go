package xfreecache

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseSize(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    Size
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "12byte",
			args: args{s: "12byte"},
			want: 12 * Byte,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "24B",
			args: args{s: "24B"},
			want: 24 * Byte,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "24kB",
			args: args{s: "24kB"},
			want: 24 * KB,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "36MB",
			args: args{s: "36MB"},
			want: 36 * MB,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "22GB",
			args: args{s: "22GB"},
			want: 22 * GB,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "2TB",
			args: args{s: "2TB"},
			want: 2 * TB,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "2gb256mb24kb12b",
			args: args{s: "2gb256mb24kb12b"},
			want: 2*GB + 256*MB + 24*KB + 12*Byte,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return err == nil
			},
		},
		{
			name: "2kbb",
			args: args{s: "2kbb"},
			want: 0,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				fmt.Println(err)
				return err != nil
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSize(tt.args.s)
			if !tt.wantErr(t, err, fmt.Sprintf("ParseSize(%v)", tt.args.s)) {
				return
			}
			assert.Equalf(t, tt.want, got, "ParseSize(%v)", tt.args.s)
		})
	}
}
