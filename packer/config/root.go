package config

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/imdario/mergo"
)

// Root of a packer configuration tree
//
type Root struct {
	Artifacts []Artifact `hcl:"artifact,block"`
}

// Merge toMerge into root
// slices will be appended together.
func (root *Root) Merge(toMerge *Root) {
	err := mergo.Merge(root, toMerge, mergo.WithAppendSlice)
	if err != nil {
		panic(fmt.Sprintf("merge: %v", err)) // TODO: remove me
	}
}

// Operation defines different lifecycle
// name/identifier for different lifecycle
// steps of packer.
//
// New Operations might be added in the future.
type Operation int

var (
	// Validation indicates that
	// a validation operation is in progress
	Validation Operation
	// Build indicates that
	// a build operation is in progress
	Build Operation = 200
)

// Artifact represents the configuration for
// a basic packer artifact run step.
// It should be viewed as a config.Artifact
// and not as a resulting artifact.
//
// All of it - except the Operation field -
// can be set in config file.
//
// Pointer mark fields as not required
type Artifact struct {
	// Operation is set at runtime and
	// tells what operation is being run.
	Operation

	Debug        *bool         `hcl:"debug"`
	Force        *bool         `hcl:"force"`
	OnError      *string       `hcl:"on_error"`
	Type         string        `hcl:"type,label"`
	Name         string        `hcl:"name,label"`
	Source       *string       `hcl:"source,attr"`
	Provisioners []Provisioner `hcl:"provisioner,block"`
	Artifacts    []Artifact    `hcl:"artifact,block"` // children
	Remain       hcl.Body      `hcl:",remain"`        // remaining body will be used by artifact implementer
}

// FullName returns the full addressable name of this artifact
func (a *Artifact) FullName() string {
	name := strings.Join([]string{"artifact", a.Type, a.Name}, ".")
	return name
}

// Provisioner represents a basic packer provisioner
type Provisioner struct {
	Type   string   `hcl:"type,label"`
	Remain hcl.Body `hcl:",remain"` // remaining body will be used by implementers
}
