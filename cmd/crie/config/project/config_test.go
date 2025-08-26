package project

//func TestSaveConfiguration(t *testing.T) {
//	fmtLinter := &noop.LintNoop{}
//	chkLinter := &noop.LintNoop{}
//	regex := regexp.MustCompile(`\.go$`)
//
//	cli := &Config{
//		Config: language2.ConfigLanguages{
//			Ignore: []string{"\\*.tmp"},
//			Languages: map[string]language2.ConfigLanguage{
//				"go": {
//					Regex: &language2.ConfigRegex{Regexp: regex},
//					Fmt:   language2.ConfigLinter{Linter: fmtLinter},
//					Chk:   language2.ConfigLinter{Linter: chkLinter},
//				},
//			},
//		},
//	}
//
//	cli.SaveConfiguration()
//
//	expected := &linter.Language{
//		Regex: regex,
//		Fmt:   fmtLinter,
//		Chk:   chkLinter,
//	}
//
//	assert.Equal(t, expected, cli.Crie.Languages["go"])
//}
