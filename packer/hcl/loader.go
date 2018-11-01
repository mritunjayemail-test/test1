package hcl

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hclparse"
)

// Loader deals with hcl loading
type Loader struct {
	// Allowed extension for json and hcl; defaults to ".json" and ".hcl"
	JSONExtensions, HCLExtensions []string
}

// DefaultLoader is the default configuration for a loader
var DefaultLoader = &Loader{
	JSONExtensions: []string{"json"},
	HCLExtensions:  []string{"hcl"},
}

// Load content of location using DefaultLoader.
func Load(location string) (files map[string]*hcl.File, diags hcl.Diagnostics) {
	return DefaultLoader.Load(location)
}

// Load content of location.
// location can be a directory or a file.
// When location is a directory, files with JSONExtensions/HCLExtensions
// will be listed and loaded if the extension matches.
// When location is a file, said file will then be loaded.
// When no file matching recognized extensions were found,
// an error is returned.
func (loader *Loader) Load(location string) (map[string]*hcl.File, hcl.Diagnostics) {
	fi, err := os.Stat(location)
	if err != nil {
		return nil, hcl.Diagnostics{
			{
				Severity: hcl.DiagError,
				Summary:  fmt.Sprintf("Failed to stat '%s'", location),
				Detail:   fmt.Sprintf("%v", err),
			},
		}
	}
	var files []string
	switch mode := fi.Mode(); {
	case mode.IsDir():
		files, err = loader.listDirectory(location)
		if err != nil {
			return nil, hcl.Diagnostics{
				{
					Severity: hcl.DiagError,
					Summary:  "Failed to list location",
					Detail:   fmt.Sprintf("The location %q could not be list; %v", location, err),
				},
			}
		}
	case mode.IsRegular():
		files = []string{location}
	}

	hclfiles, hcldiags := loader.load(files)
	if len(hclfiles) == 0 && len(hcldiags) == 0 {
		return hclfiles, hcl.Diagnostics{
			{
				Severity: hcl.DiagWarning,
				Summary:  "No recognized file type found",
				Detail:   fmt.Sprintf("Recognized file extensions: %+v", *loader),
			},
		}
	}
	return hclfiles, hcldiags
}

// load files into obj using absolute path
func (loader *Loader) load(files []string) (hclfiles map[string]*hcl.File, diags hcl.Diagnostics) {
	parser := hclparse.NewParser()
	hclfiles = map[string]*hcl.File{}
	for _, fileName := range files {
		var moreFiles *hcl.File
		var moreDiags hcl.Diagnostics

		switch ex := extension(fileName); {
		case loader.isJSON(ex):
			moreFiles, moreDiags = parser.ParseJSONFile(fileName)
		case loader.isHCL(ex):
			moreFiles, moreDiags = parser.ParseHCLFile(fileName)
		default:
			continue // file type not recognized
		}
		diags = append(diags, moreDiags...)
		if diags.HasErrors() {
			continue // for now, it's probably better to get all errors from beginning
		}
		hclfiles[fileName] = moreFiles
	}
	return
}
