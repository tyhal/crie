package filelist

import (
	"os"
	"path/filepath"
)

// FromDir walks the provided directory and returns a list of files
// relative to dir. Directories themselves are not included.
func FromDir(dir string) ([]string, error) {
	var files []string

	// Create an initial file list
	err := filepath.Walk(dir, func(currPath string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			relPath, err := filepath.Rel(dir, currPath)
			if err != nil {
				return err
			}
			files = append(files, relPath)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}
