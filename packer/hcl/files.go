package hcl

import (
	"io/ioutil"
	"path/filepath"
	"strings"
)

const extensionSeparator = "."

// listDirectory returns a list of filenames from dir that end with
// loader.Extensions()
func (loader *Loader) listDirectory(dirname string) ([]string, error) {
	files, err := ioutil.ReadDir(dirname)
	if err != nil {
		return nil, err
	}

	var res []string
	for _, f := range files {
		res = append(res, filepath.Join(dirname, f.Name()))
	}
	return res, nil
}

// extension of filename
func extension(filename string) (extension string) {
	parts := strings.Split(filename, extensionSeparator)
	if len(parts) == 0 {
		return ""
	}
	return parts[len(parts)-1]
}

func (loader *Loader) isJSON(extension string) bool {
	for _, validExtension := range loader.JSONExtensions {
		if strings.EqualFold(extension, validExtension) {
			return true
		}
	}
	return false
}

func (loader *Loader) isHCL(extension string) bool {
	for _, validExtension := range loader.HCLExtensions {
		if strings.EqualFold(extension, validExtension) {
			return true
		}
	}
	return false
}
