```
┌─────────────┬────────────┬──────────────┬─────────────────────────────────────┐
│  LANGUAGE   │  CHECKER   │  FORMATTER   │          ASSOCIATED FILES           │
├─────────────┼────────────┼──────────────┼─────────────────────────────────────┤
│ asciidoctor │ vale       │              │ \.adoc$                             │
│ bash        │ shellcheck │ shfmt        │ \.bash$                             │
│ c           │ cppcheck   │ clang-format │ \.c$                                │
│ cmake       │ cmakelint  │              │ CMakeLists.txt$|\.cmake$            │
│ cpp         │ cppcheck   │ clang-format │ \.cc$|\.cpp$                        │
│ cppheaders  │            │ clang-format │ \.h$|\.hpp$                         │
│ docker      │ hadolint   │              │ (?i).*(Contain|Docker)file.*        │
│ golang      │ revive     │ gofmt        │ \.go$                               │
│ javascript  │ standard   │ standard     │ \.js$|\.jsx$                        │
│ json        │ jsonlint   │ jsonlint     │ \.json$|\.JSON$                     │
│ markdown    │ vale       │ remark       │ \.md$                               │
│ proto       │ protolint  │ protolint    │ \.proto$                            │
│ python      │ pylint     │ black        │ \.py$                               │
│ sh          │ shellcheck │ shfmt        │ \.sh$|/script/[^.]*$|^script/[^.]*$ │
│ terraform   │ tflint     │ terraform    │ \.tf$                               │
│ yml         │ yamllint   │ yamlfmt      │ \.ya?ml$                            │
└─────────────┴────────────┴──────────────┴─────────────────────────────────────┘
```
