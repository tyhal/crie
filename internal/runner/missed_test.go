package runner

import (
	"bytes"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_printCoverageStats(t *testing.T) {
	buf := &bytes.Buffer{}
	err := printCoverageStats(buf, map[string]int{".bash": 1, ".go": 2})
	assert.NoError(t, err)
	exp := `
┌───────────┬───────┐
│ EXTENSION │ COUNT │
├───────────┼───────┤
│ .bash     │ 1     │
│ .go       │ 2     │
└───────────┴───────┘
`
	assert.Equal(t, exp, "\n"+buf.String())
}

func Test_noCoverageStats(t *testing.T) {
	l := RunConfiguration{
		NamedMatches: NamedMatches{
			"shell": LinterMatch{
				FileMatch: regexp.MustCompile("\\.sh$"),
			},
		},
	}
	stats := l.noCoverageStats([]string{
		"test.bash",
		"something.go",
		"shell.sh",
	})
	assert.Equal(t, map[string]int{".bash": 1, ".go": 1}, stats)
}
