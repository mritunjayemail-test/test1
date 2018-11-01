package config

import (
	"fmt"

	"github.com/hashicorp/hcl2/gohcl"
	"github.com/hashicorp/hcl2/hcl"
	packerhcl "github.com/hashicorp/packer/packer/hcl"
)

// Loader loads a Root from a path
type Loader struct {
	FileLoader *packerhcl.Loader
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

	root := &Root{
		Files: hclFiles,
	}
	for _, file := range hclFiles {
		current := &Root{}
		diags = append(diags, gohcl.DecodeBody(file.Body, nil, current)...)
		root.Merge(current)
	}
	if diags.HasErrors() {
		return nil, diags
	}

	// relocate artifacts with a source
	// to their corresponding location.
	for i := 0; i < len(root.Artifacts); {
		if root.Artifacts[i].Source == nil {
			i++
			continue
		}
		// source is set
		// remove i from the array
		toRelocate := root.Artifacts[i]
		root.Artifacts = append(root.Artifacts[:i], root.Artifacts[i+1:]...) // remove

		// find source location
		location := findArtifactNamed(root.Artifacts, *toRelocate.Source)
		if location == nil {
			diags = append(diags, &hcl.Diagnostic{
				Severity: hcl.DiagError,
				Summary:  "Source not found",
				Detail:   fmt.Sprintf("cannot find source '%s' for %s", *toRelocate.Source, toRelocate.FullName()),
			})
			return nil, diags
		}
		toRelocate.Source = nil
		location.Artifacts = append(location.Artifacts, toRelocate)
	}

	return root, diags
}

func findArtifactNamed(artifacts []Artifact, name string) *Artifact {
	for i := range artifacts {
		if artifacts[i].FullName() == name {
			return &artifacts[i]
		}
		if art := findArtifactNamed(artifacts[i].Artifacts, name); art != nil {
			return art
		}
	}
	return nil
}
