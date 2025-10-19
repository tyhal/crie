## crie lnt

Runs both fmt and then chk

### Synopsis

Runs both format and then check

```
crie lnt [flags]
```

### Options

```
  -e, --continue            show all errors rather than stopping at the first
  -g, --git-diff            only use files from the current commit to (git-target)
  -t, --git-target string   the branch to compare against to find changed files (default "origin/main")
  -h, --help                help for lnt
      --only crie ls        run with only one language (see crie ls for available options)
  -p, --passes              show files that passed
```

### Options inherited from parent commands

```
  -c, --conf string        project configuration file (default "crie.yml")
  -j, --json               turn on json output
  -l, --lang-conf string   language configuration file (default "crie.lang.yml")
  -q, --quiet              only prints critical errors (suppresses verbose)
  -v, --verbose            turn on verbose printing for reports
```

### SEE ALSO

* [crie](crie.md)	 - crie is a formatter and linter for many languages.

