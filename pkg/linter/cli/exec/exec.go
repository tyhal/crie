package exec

import (
	"io"
	"os"
	"path/filepath"
)

// Executor is an abstraction to allow any cli tool to run anywhere
type Executor interface {
	Setup() error
	Exec(bin string, frontParams []string, filePath string, endParams []string, chdir bool, stdout io.Writer, stderr io.Writer) error
	Cleanup() error
}

// ToLinuxPath ensures windows paths can be mapped to linux container paths
func ToLinuxPath(path string) string {
	// Remove Windows drive letter if present
	if vol := filepath.VolumeName(path); vol != "" {
		path = path[len(vol):]
	}
	// Convert backslashes to slashes
	return filepath.ToSlash(path)
}

// GetWorkdirAsLinuxPath is useful for consistently mapping any host fs to what is typically in container environments (Linux)
func GetWorkdirAsLinuxPath() (string, error) {
	// Ensure we can mount our filesystem to the same path inside the container
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	dir, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}
	linuxDir := ToLinuxPath(dir)

	return linuxDir, nil
}
