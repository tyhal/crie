```
┌─────────────┬────────────┬──────────────┬─────────────────────────────────────┐
│  LANGUAGE   │  CHECKER   │  FORMATTER   │          ASSOCIATED FILES           │
├─────────────┼────────────┼──────────────┼─────────────────────────────────────┤
│ python      │ pylint     │ black        │ \.py$                               │
│ cppheaders  │            │ clang-format │ \.h$|\.hpp$                         │
│ sh          │ shellcheck │ shfmt        │ \.sh$|/script/[^.]*$|^script/[^.]*$ │
│ proto       │ protolint  │ protolint    │ \.proto$                            │
│ bash        │ shellcheck │ shfmt        │ \.bash$                             │
│ javascript  │ standard   │ standard     │ \.js$|\.jsx$                        │
│ docker      │ hadolint   │              │ Dockerfile$                         │
│ json        │ jsonlint   │ jsonlint     │ \.json$|\.JSON$                     │
│ cmake       │ cmakelint  │              │ CMakeLists.txt$|\.cmake$            │
│ golang      │ golint     │ gofmt        │ \.go$                               │
│ markdown    │ vale       │ remark       │ \.md$                               │
│ terraform   │ tflint     │ terraform    │ \.tf$                               │
│ cpp         │ cppcheck   │ clang-format │ \.cc$|\.cpp$                        │
│ c           │ cppcheck   │ clang-format │ \.c$                                │
│ yml         │ yamllint   │ yamlfmt      │ \.yml$|\.yaml$                      │
│ asciidoctor │ vale       │              │ \.adoc$                             │
└─────────────┴────────────┴──────────────┴─────────────────────────────────────┘
```
