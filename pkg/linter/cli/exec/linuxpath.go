package exec

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

// ToLinuxPath ensures windows paths can be mapped to linux container paths
func ToLinuxPath(p string) string {
	p = strings.ReplaceAll(p, `\`, `/`)
	// Remove Windows drive letter (e.g., "C:")
	if len(p) >= 2 && p[1] == ':' {
		p = p[2:]
	}
	p = path.Clean(p)
	return p
}

// GetWorkdirAsLinuxPath is useful for consistently mapping any host fs to what is typically in container environments (Linux)
func GetWorkdirAsLinuxPath() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	// ensure absolute path for consistency
	dir, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}
	return ToLinuxPath(dir), nil
}
