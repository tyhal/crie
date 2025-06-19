package cli

import (
	"os"
	"path/filepath"
)

func getwdContainer() (string, error) {
	// Ensure we can mount our filesystem to the same path inside the container
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}
	linuxDir := toLinuxPath(dir)

	return linuxDir, nil
}
