package config

import (
	"path/filepath"
	"reflect"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func testdir(name string) string { return filepath.Join("./testdir", name) }

func nillifyRemains(artifacts []Artifact) {
	for i := range artifacts {
		r := &artifacts[i]
		r.Remain = nil // fix it to nil for now
		for i := range r.Provisioners {
			r := &r.Provisioners[i]
			r.Remain = nil
		}
		nillifyRemains(r.Artifacts)
	}
}

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
				{0, nil, nil, nil,
					"googlecompute", "ubuntu-1804-lts", nil,
					[]Provisioner{{"shell", nil}},
					[]Artifact{
						{
							0, nil, nil, nil,
							"googlecompute", "ubuntu-1804-lts-consul", nil,
							[]Provisioner{{"shell", nil}},
							[]Artifact{{0, nil, nil, nil, "compress", "ubuntu-1804-lts-consul.gz", nil, nil, nil, nil}},
							nil,
						},
						{
							0, nil, nil, nil,
							"googlecompute", "ubuntu-1804-lts-vault", nil,
							[]Provisioner{{"shell", nil}},
							[]Artifact{{0, nil, nil, nil, "compress", "ubuntu-1804-lts-vault.gz", nil, nil, nil, nil}},
							nil,
						},
					},
					nil},
			},
		}, 0},
		{"google-simple", args{testdir("google-source-not-found")}, nil, 1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRoot, gotDiags := Load(tt.args.location)
			if gotRoot != nil {
				nillifyRemains(gotRoot.Artifacts)
				gotRoot.Files = nil
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
