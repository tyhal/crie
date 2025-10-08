package exec

import (
	"fmt"
	"strings"
	"testing"
)

func TestToLinuxPath(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"unix_path", "/usr/local/bin", "/usr/local/bin"},
		{"backslashes_only", "dir\\sub\\file.txt", "dir/sub/file.txt"},
		{"windows_drive_backslashes", "C:\\Users\\alice\\work", "/Users/alice/work"},
		{"windows_drive_mixed", "D:/Projects/go/src", "/Projects/go/src"},
		{"mixed_separators", "some\\mixed/path\\here", "some/mixed/path/here"},
		{"relative_path", "..\\..\\bin\\file.txt", "../../bin/file.txt"},
		{"root", "/", "/"},
		{"empty", "", "."},
		{"empty", ".", "."},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := ToLinuxPath(tc.in)
			want := tc.want
			fmt.Println(tc.in)
			if got != want {
				t.Fatalf("ToLinuxPath(%q) = %s; want %s", tc.in, got, want)
			}
			// Ensure result uses forward slashes only
			if strings.Contains(got, "\\") {
				t.Fatalf("expected no backslashes in result, got %q", got)
			}
		})
	}
}
