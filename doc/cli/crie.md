## crie

crie is a formatter and linter for many languages.

### Synopsis


	crie brings together a vast collection of formatters and linters
	to create a handy tool that can sanity check any codebase.

### Examples

```

check all files changes since the target branch 
	$ crie chk --git-diff --git-target=origin/main

format all python files
	$ crie fmt --only python

```

### Options

```
  -c, --conf string        project configuration file (default "crie.yml")
  -d, --dir string         the directory to run crie in
  -h, --help               help for crie
  -j, --json               turn on json output
  -l, --lang-conf string   language configuration file (default "crie.lang.yml")
  -q, --quiet              only prints critical errors
  -v, --verbose            turn on verbose printing for reports
```

### SEE ALSO

* [crie chk](crie_chk.md)	 - Run linters that only check code
* [crie conf](crie_conf.md)	 - Print configuration settings
* [crie fmt](crie_fmt.md)	 - Run formatters
* [crie init](crie_init.md)	 - Create optional config files
* [crie lnt](crie_lnt.md)	 - Runs both fmt and then chk
* [crie ls](crie_ls.md)	 - Show languages
* [crie non](crie_non.md)	 - Show what isn't supported for this project
* [crie schema](crie_schema.md)	 - Print JSON schemas for config files

