## crie chk

Run linters that only check code

### Synopsis

Check all code standards for coding conventions

```
crie chk [flags]
```

### Options

```
  -a, --continue            show all errors rather than stopping at the first
  -g, --git-diff            only check files changed in git
  -t, --git-target string   a target branch to compare against e.g 'remote/branch' or 'branch'
  -h, --help                help for chk
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

