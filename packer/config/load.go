package config

import (
	"fmt"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	packerhcl "github.com/hashicorp/packer/packer/hcl"
	"github.com/imdario/mergo"
)

// Loader loads a Root from a path
type Loader struct {
	FileLoader *packerhcl.Loader
}

// Root of a packer configuration
type Root struct {
	Artifacts []Artifact `hcl:"artifact,block"`
}

// Merge two roots together
func (root *Root) Merge(toMerge *Root) {
	err := mergo.Merge(root, toMerge, mergo.WithAppendSlice)
	if err != nil {
		panic(fmt.Sprintf("merge: %v", err)) // TODO: remove me
	}
}

// Artifact represents a packer artifact
type Artifact struct {
	Type   string   `hcl:"type,label"`
	Name   string   `hcl:"name,label"`
	Source *string  `hcl:"source,attr"`
	Remain hcl.Body `hcl:",remain"` // remainin body will be used by artifact implementers
}

// DefaultLoader is the default loader object
var DefaultLoader = &Loader{
	FileLoader: packerhcl.DefaultLoader,
}

// Load calls Load from DefaultLoader
func Load(location string) (root *Root, diags hcl.Diagnostics) { return DefaultLoader.Load(location) }

// Load location and return a root packer configuration
func (loader *Loader) Load(location string) (*Root, hcl.Diagnostics) {
	hclFiles, diags := loader.FileLoader.Load(location)
	if diags.HasErrors() {
		return nil, diags
	}

	root := &Root{}
	for _, file := range hclFiles {
		current := &Root{}
		diags = append(diags, gohcl.DecodeBody(file.Body, nil, current)...)
		root.Merge(current)
	}

	return root, nil
}
