package juno

import (
	"errors"
	"testing"

	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/assert"
)

type mockCC struct {
	res *resty.Response
	err error
}

func newMockCC(res *resty.Response, err error) *mockCC {
	return &mockCC{
		res: res,
		err: err,
	}
}

func (m *mockCC) Get(url string) (*resty.Response, error) {
	return m.res, m.err
}

func TestNewDataSource(t *testing.T) {
	type args struct {
		path  string
		watch bool
		cc    client
	}
	tests := []struct {
		name string
		args args
		want *resty.Response
		err  string
	}{
		{
			name: "test1",
			args: args{
				path:  "juno://",
				watch: true,
				cc:    newMockCC(nil, errors.New("error test")),
			},
			want: nil,
			err:  "juno agent address is empty",
		},
		{
			name: "test2",
			args: args{
				path:  "juno://localhost:8080",
				watch: true,
				cc:    newMockCC(nil, errors.New("error test")),
			},
			want: nil,
			err:  "juno env is empty",
		},
		{
			name: "test3",
			args: args{
				path:  "juno://localhost:8080/path/to/config?env=dev",
				watch: true,
				cc:    newMockCC(nil, errors.New("error test")),
			},
			want: nil,
			err:  "error test",
		},
		{
			name: "test4",
			args: args{
				path:  "juno://localhost:8080/path/to/config?env=dev",
				watch: true,
				cc:    newMockCC(&resty.Response{}, nil),
			},
			want: &resty.Response{},
			err:  "",
		},
		{
			name: "test5",
			args: args{
				path:  "::ss",
				watch: true,
				cc:    newMockCC(nil, errors.New("error test")),
			},
			want: nil,
			err:  "parse \"::ss\": missing protocol scheme",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDataSource(tt.args.path, tt.args.watch)
			got.cc = tt.args.cc

			content, err := got.ReadConfig()
			if err != nil {
				assert.EqualError(t, err, tt.err)
			}

			if tt.want != nil {
				assert.Equal(t, tt.want.Body(), content)
			}

			got.IsConfigChanged()
			got.Close()
		})
	}
}
