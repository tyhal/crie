package orchestrator

import (
	"fmt"
	"os"
	"path/filepath"
)

// findModuleRoot walks up from filePath until it finds a directory containing marker.
func findModuleRoot(filePath, marker string) (string, error) {
	dir := filepath.Dir(filePath)
	for {
		if _, err := os.Stat(filepath.Join(dir, marker)); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", fmt.Errorf("%s not found above %s", marker, filePath)
		}
		dir = parent
	}
}

// groupByModule returns files grouped by their nearest ancestor directory containing marker.
// Results are cached per directory so files sharing a parent don't repeat the filesystem walk.
func groupByModule(files []string, marker string) map[string][]string {
	cache := make(map[string]string) // dir → module root ("" means not found)
	groups := make(map[string][]string)
	for _, f := range files {
		dir := filepath.Dir(f)
		root, seen := cache[dir]
		if !seen {
			root, _ = findModuleRoot(f, marker)
			cache[dir] = root
		}
		if root != "" {
			groups[root] = append(groups[root], f)
		}
	}
	return groups
}
