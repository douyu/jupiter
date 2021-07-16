package protoc

import "testing"

func TestUcfirst(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestUcfirst",
			args: args{"helloService"},
			want: "HelloService",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Ucfirst(tt.args.str); got != tt.want {
				t.Errorf("Ucfirst() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestLcfirst(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "TestLcfirst",
			args: args{"HelloService"},
			want: "helloService",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Lcfirst(tt.args.str); got != tt.want {
				t.Errorf("Lcfirst() = %v, want %v", got, tt.want)
			}
		})
	}
}
