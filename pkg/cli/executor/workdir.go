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
		targetPath := ToLinuxPath(filePath)
		if noFileArg {
			wdContainer = filepath.Join(wdContainer, targetPath)
		} else {
			wdContainer = filepath.Join(wdContainer, filepath.Dir(targetPath))
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
