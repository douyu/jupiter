package common

import (
	"log"
	"path/filepath"
	"testing"
)

func TestGetModPath(t *testing.T) {
	type args struct {
		projectPath string
	}
	tests := []struct {
		name        string
		args        args
		wantModPath string
	}{
		{
			name:        "getModPath",
			args:        args{projectPath: "./"},
			wantModPath: "github.com/douyu/jupiter/tools/jupiter/",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			projectPath, err := filepath.Abs(tt.args.projectPath)
			if err != nil {
				log.Fatalf("get absolute path failed:%v\n", err)
			}
			if gotModPath := GetModPath(projectPath); gotModPath != tt.wantModPath {
				t.Errorf("GetModPath() = %v, want %v", gotModPath, tt.wantModPath)
			}
		})
	}
}
