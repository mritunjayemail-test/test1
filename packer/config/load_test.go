package config

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testdir(name string) string { return filepath.Join("./testdir", name) }
func ptrstr(name string) *string { return &name }

func TestLoad(t *testing.T) {
	type args struct {
		location string
	}
	tests := []struct {
		name            string
		args            args
		wantRoot        *Root
		wantDiagnostics int
	}{
		// {"not-found", args{"404"}, nil, 1},
		{"google-simple", args{testdir("google-simple")}, &Root{
			Artifacts: []Artifact{
				{"googlecompute", "ubuntu-1804-lts", nil, nil},
				{"googlecompute", "ubuntu-1804-lts-consul", ptrstr("artifact.googlecompute.ubuntu-1804-lts"), nil},
				{"compress", "ubuntu-1804-lts-consul.gz", ptrstr("artifact.googlecompute.ubuntu-1804-lts-consul"), nil},
				{"googlecompute", "ubuntu-1804-lts-vault", ptrstr("artifact.googlecompute.ubuntu-1804-lts"), nil},
				{"compress", "ubuntu-1804-lts-vault.gz", ptrstr("artifact.googlecompute.ubuntu-1804-lts-vault"), nil},
			},
		}, 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRoot, gotDiags := Load(tt.args.location)
			for i := range gotRoot.Artifacts {
				r := &gotRoot.Artifacts[i]
				r.Remain = nil // fix it to nil for now
			}
			if diff := cmp.Diff(tt.wantRoot, gotRoot); diff != "" {
				t.Errorf("Load() -want +got: \n%s", diff)
			}
			if !reflect.DeepEqual(len(gotDiags), tt.wantDiagnostics) {
				t.Errorf("Load() gotDiags = %v, want %v", gotDiags, tt.wantDiagnostics)
			}
		})
	}
}
