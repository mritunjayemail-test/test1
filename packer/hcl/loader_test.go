package hcl

import (
	"path/filepath"
	"reflect"
	"testing"
)

func testdir(name string) string {
	return filepath.Join("./testdir", name)
}

func TestLoad(t *testing.T) {
	type args struct {
		location string
	}
	tests := []struct {
		name      string
		args      args
		wantFiles int
		wantDiags int
	}{
		{"error", args{"error"}, 0, 1},
		{"0", args{testdir("0")}, 0, 1},
		{"1", args{testdir("1")}, 1, 0},
		{"2", args{testdir("2")}, 2, 0},
		{"file", args{testdir("1/valid.json")}, 1, 0},
		{"file", args{testdir("1/valid_json")}, 0, 1},
		{"3", args{testdir("3")}, 0, 3},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFiles, gotDiags := Load(tt.args.location)
			if !reflect.DeepEqual(len(gotFiles), tt.wantFiles) {
				t.Errorf("Load() gotFiles = %v, want %v", gotFiles, tt.wantFiles)
			}
			if !reflect.DeepEqual(len(gotDiags), tt.wantDiags) {
				t.Errorf("Load() gotDiags = %v, want %v", gotDiags, tt.wantDiags)
			}
		})
	}
}
