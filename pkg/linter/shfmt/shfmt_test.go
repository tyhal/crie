package shfmt

// NOTE This mostly exists to just to be an easy boilerplate for testing other linter implementations

import (
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
	"math"
	"mvdan.cc/sh/v3/syntax"
	"sync"
	"testing"
)

func TestLint_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name       string
		yaml       string
		expectLang syntax.LangVariant
		wantErr    bool
	}{
		{
			name:       "lang bash",
			yaml:       `language: bash`,
			expectLang: syntax.LangBash,
			wantErr:    false,
		},
		{
			name:       "lang sh",
			yaml:       `language: sh`,
			expectLang: syntax.LangPOSIX,
			wantErr:    false,
		},
		{
			name:       "lang posix",
			yaml:       `language: posix`,
			expectLang: syntax.LangPOSIX,
			wantErr:    false,
		},
		{
			name:       "lang mksh",
			yaml:       `language: mksh`,
			expectLang: syntax.LangMirBSDKorn,
			wantErr:    false,
		},
		{
			name:    "err",
			yaml:    `language: none`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var sh Lint
			err := yaml.Unmarshal([]byte(tt.yaml), &sh)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectLang, sh.Language)
			}
		})
	}
}

func TestLint_Name(t *testing.T) {
	l := &Lint{}
	assert.Equal(t, "shfmt", l.Name())
}

func TestLint_WillRun(t *testing.T) {
	l := &Lint{}
	assert.NoError(t, l.WillRun())
}

func TestLint_Cleanup(t *testing.T) {
	l := &Lint{}
	var wg sync.WaitGroup
	wg.Add(1)
	l.Cleanup(&wg)
	wg.Wait()
}

func TestLint_MaxConcurrency(t *testing.T) {
	l := &Lint{}
	assert.Equal(t, math.MaxInt32, l.MaxConcurrency())
}

// TODO either mock or make dummy files
//func TestLint_Run(t *testing.T) {
//	l := &Lint{}
//	rep := make(chan linter.Report, 1)
//
//	l.Run("test.txt", rep)
//
//	report := <-rep
//	assert.Equal(t, "test.txt", report.File)
//	assert.NoError(t, report.Err)
//	assert.Nil(t, report.StdOut)
//	assert.Nil(t, report.StdErr)
//}
