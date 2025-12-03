package runner

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/tyhal/crie/pkg/linter/noop"
)

// disableLogging changes the logger output and returns a restore function.
// The caller is expected to defer the returned function.
func disableLogging() func() {
	originalOutput := logrus.StandardLogger().Out
	logrus.SetOutput(io.Discard)
	return func() {
		logrus.SetOutput(originalOutput)
	}
}

func Test_getName(t *testing.T) {
	assert.Empty(t, getName(nil))
	assert.Equal(t, "noop", getName(&noop.LintNoop{}))
}

//func TestRunConfiguration_GetLanguage(t *testing.T) {
//	config := &RunConfiguration{
//		Languages: map[string]*Language{
//			"test": {
//				Chk:       &noop.LintNoop{},
//				Fmt:       &noop.LintNoop{},
//				FileMatch: regexp.MustCompile(`\.test$`),
//			},
//		},
//	}
//
//	// Test existing language
//	lang, err := config.GetLanguage("test")
//	assert.NoError(t, err)
//	assert.NotNil(t, lang)
//
//	// Test non-existent language
//	_, err = config.GetLanguage("nonexistent")
//	assert.Error(t, err)
//	assert.Contains(t, err.Error(), "language 'nonexistent' not found")
//}

//func TestRunConfiguration_runLinter(t *testing.T) {
//	config := &RunConfiguration{
//		Languages: map[string]*Language{
//			"go": {
//				Chk:       &noop.LintNoop{},
//				FileMatch: regexp.MustCompile(`\.go$`),
//			},
//		},
//	}
//
//	fileList := []string{"test.go"}
//
//	var cleanupGroup sync.WaitGroup
//	err := config.runLinter(&cleanupGroup, "go", LintTypeChk, fileList)
//	assert.NoError(t, err)
//
//	// Wait for cleanup to complete
//	cleanupGroup.Wait()
//}

func TestRunConfiguration_runLinters(t *testing.T) {
	tests := []struct {
		name       string
		config     *RunConfiguration
		files      []string
		expectErr  bool
		errMessage string
	}{
		{
			name: "default runLinters - happy path",
			config: &RunConfiguration{
				Languages: map[string]*Language{
					"go": {
						Chk:       &noop.LintNoop{},
						FileMatch: regexp.MustCompile(`\.go$`),
					},
				},
			},
			files:     []string{"test.go"},
			expectErr: false,
		},
		{
			name: "runLinters with single valid language (go)",
			config: &RunConfiguration{
				Languages: map[string]*Language{
					"go": {
						Chk:       &noop.LintNoop{},
						FileMatch: regexp.MustCompile(`\.go$`),
					},
					"test": {
						Chk:       &noop.LintNoop{},
						FileMatch: regexp.MustCompile(`\.test$`),
					},
				},
				Options: Options{
					Only: "go",
				},
			},
			files:     []string{"test.go"},
			expectErr: false,
		},
		{
			name: "runLinters with nonexistent language in 'Only' option",
			config: &RunConfiguration{
				Languages: map[string]*Language{
					"go": {
						Chk:       &noop.LintNoop{},
						FileMatch: regexp.MustCompile(`\.go$`),
					},
				},
				Options: Options{
					Only: "nonexistent",
				},
			},
			files:      []string{"test.go"},
			expectErr:  true,
			errMessage: "language 'nonexistent' not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.runLinters(LintTypeChk, tt.files)

			if tt.expectErr {
				if assert.Error(t, err) && tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// charFromIndex returns a lowercase letter based on the given index (only uses 20 chars for round numbers)
func charFromIndex(i int) byte {
	const letters = "abcdefghijklmnopqrst"
	return letters[i%len(letters)]
}

// genFilenames generates a list of filenames based on the given count [ 0.a 1.b 2.c ... 19.t 20.a ]
func genFilenames(count int) []string {
	filenames := make([]string, count)
	for i := 0; i < count; i++ {
		filenames = append(filenames, fmt.Sprintf("%d.%c", i, charFromIndex(i)))
	}
	return filenames
}

func genLangs(count int) Languages {
	langs := make(Languages, count)
	for i := 0; i < count; i++ {
		langs[strconv.Itoa(i)] = &Language{
			// use a different linter object for each language
			Chk:       &noop.LintNoop{},
			FileMatch: regexp.MustCompile(fmt.Sprintf(`\.%c$`, charFromIndex(i))),
		}
	}
	return langs
}

func BenchmarkRunConfiguration_runLinters(b *testing.B) {
	defer disableLogging()()

	opts := Options{
		StrictLogging: true,
	}

	benchs := []struct {
		name   string
		config *RunConfiguration
		files  []string
	}{
		// comment-<FileCount>F-<LanguageCount>L
		{
			name: "noskip-20F-20L",
			config: &RunConfiguration{
				Options:   opts,
				Languages: genLangs(20),
			},
			files: genFilenames(20),
		},
		{
			name: "halfskip-100F-10L",
			config: &RunConfiguration{
				Options:   opts,
				Languages: genLangs(10),
			},
			files: genFilenames(100),
		},
		{
			name: "noskip-100F-20L",
			config: &RunConfiguration{
				Options:   opts,
				Languages: genLangs(20),
			},
			files: genFilenames(100),
		},
		{
			name: "halfskip-10F-100L",
			config: &RunConfiguration{
				Options:   opts,
				Languages: genLangs(100),
			},
			files: genFilenames(10),
		},
		{
			name: "noskip-20F-100L",
			config: &RunConfiguration{
				Options:   opts,
				Languages: genLangs(100),
			},
			files: genFilenames(20),
		},
		{
			name: "noskip-100F-100L",
			config: &RunConfiguration{
				Options:   opts,
				Languages: genLangs(100),
			},
			files: genFilenames(100),
		},
	}

	b.ResetTimer()
	for _, tt := range benchs {
		b.Run(tt.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				err := tt.config.runLinters(LintTypeChk, tt.files)
				if err != nil {
					b.Fatalf("unexpected error at iteration %d: %v", i, err)
				}
			}
		})
	}
}

func TestRunConfiguration_Run(t *testing.T) {
	tests := []struct {
		name       string
		config     *RunConfiguration
		expectErr  bool
		errMessage string
	}{
		{
			name: "no errors",
			config: &RunConfiguration{
				Languages: map[string]*Language{
					"go": {
						Chk:       noop.WithErr(nil, nil),
						FileMatch: regexp.MustCompile(`\.go$`),
					},
				},
			},
			expectErr: false,
		},
		{
			name: "linter startup error",
			config: &RunConfiguration{
				Languages: map[string]*Language{
					"go": {
						Chk:       noop.WithErr(errors.New("startup err"), nil),
						FileMatch: regexp.MustCompile(`\.go$`),
					},
				},
			},
			expectErr:  true,
			errMessage: "1 language(s) failed while chk'ing",
		},
		{
			name: "linter run error",
			config: &RunConfiguration{
				Languages: map[string]*Language{
					"go": {
						Chk:       noop.WithErr(nil, errors.New("run err")),
						FileMatch: regexp.MustCompile(`\.go$`),
					},
				},
			},
			expectErr:  true,
			errMessage: "1 language(s) failed while chk'ing",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			err := os.Chdir(dir)
			assert.NoError(t, err)

			file, err := os.Create(path.Join(dir, "test.go"))
			defer func() {
				_ = file.Close()
			}()

			err = tt.config.Run(LintTypeChk)

			if tt.expectErr {
				if assert.Error(t, err) && tt.errMessage != "" {
					assert.Contains(t, err.Error(), tt.errMessage)
				}
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
