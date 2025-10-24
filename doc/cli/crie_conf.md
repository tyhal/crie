## crie conf

Print configuration settings

### Synopsis

Print what crie has parsed from flags, env, the project file, and then defaults

```
crie conf [flags]
```

### Options

```
  -e, --continue            show all errors rather than stopping at the first
  -g, --git-diff            only check files changed in git
  -t, --git-target string   a target branch to compare against e.g 'remote/branch' or 'branch'
  -h, --help                help for conf
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

