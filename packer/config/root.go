package config

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/imdario/mergo"
)

// Root of a packer configuration
type Root struct {
	Artifacts []Artifact `hcl:"artifact,block"`
}

// Merge toMerge into root
func (root *Root) Merge(toMerge *Root) {
	err := mergo.Merge(root, toMerge, mergo.WithAppendSlice)
	if err != nil {
		panic(fmt.Sprintf("merge: %v", err)) // TODO: remove me
	}
}

// Artifact represents a packer artifact
type Artifact struct {
	Type      string     `hcl:"type,label"`
	Name      string     `hcl:"name,label"`
	Source    *string    `hcl:"source,attr"`
	Artifacts []Artifact `hcl:"artifact,block"`
	Remain    hcl.Body   `hcl:",remain"` // remainin body will be used by artifact implementers
}

// FullName returns the full addressable name of this artifact
func (a *Artifact) FullName() string {
	name := strings.Join([]string{"artifact", a.Type, a.Name}, ".")
	return name
}
