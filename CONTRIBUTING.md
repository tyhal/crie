# Contributing

Adding new languages is easy, just add a new entry to:
[main language configuration file](internal/config/language/defaults.yml)

You should add to the verify folder for the new language.

## Building and testing

```bash
    # install deps
    script/bootstrap
    
    # run tests
    script/test
```

## View the docs

[CLI Docs](doc/cli/crie.md)

Code Docs

```
godoc -http=:6060
```

[View here](http://localhost:6060/pkg/github.com/tyhal/crie/crie/#pg-overview)

## Project structure reference

```
crie/
├── cmd/                    # Application entry points
│   ├── crie/              # Main CLI binary
│   └── docgen/            # Documentation generator
│
├── internal/              # Private application code
│   ├── cli/              # Cobra commands & CLI logic
│   ├── config/           # Configuration handling
│   │   ├── language/     # Language definitions & parsing
│   │   └── project/      # Project settings & parsing
│   └── runner/           # Core execution engine
│
├── pkg/                   # Public, reusable packages
│   └── linter/           # Linter interface & implementations
│
├── doc/                   # Documentation (to be renamed from doc/)
│   ├── cli/              # Auto-generated command docs
│   ├── config/           # Configuration examples
│   └── demo/             # Demo assets
│
├── res/                   # Static resources
│   ├── completion/       # Shell completion scripts
│   └── schema/           # JSON schemas
│
├── script/                # Build & development scripts
│
└── .github/              # CI/CD workflows
```
