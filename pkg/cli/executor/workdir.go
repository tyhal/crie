package executor

import (
	"fmt"
	"os"
	"path/filepath"
)

func containerWorkdir(filePath string, chdir, noFileArg bool) (string, error) {
	wdContainer, err := GetWorkdirAsLinuxPath()
	if err != nil {
		return "", fmt.Errorf("getting containerWorkdir: %w", err)
	}

	if chdir {
		wdHost, err := os.Getwd()
		if err != nil {
			return "", fmt.Errorf("getting working directory: %w", err)
		}
		relPath, err := filepath.Rel(wdHost, filePath)
		if err != nil {
			return "", fmt.Errorf("getting relative path: %w", err)
		}
		relPath = ToLinuxPath(relPath)

		if noFileArg {
			wdContainer = filepath.Join(wdContainer, relPath)
		} else {
			wdContainer = filepath.Join(wdContainer, filepath.Dir(relPath))
		}
	}

	return wdContainer, nil
}

func hostWorkdir(filePath string, chdir, noFileArg bool) (string, error) {
	if !chdir {
		return os.Getwd()
	}

	if noFileArg {
		return filePath, nil
	}

	return filepath.Dir(filePath), nil
}
