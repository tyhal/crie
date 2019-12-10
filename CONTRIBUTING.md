# Contributing

## Building and testing

```bash
    sudo -E script/bootstrap
    script/test
```

## View the docs

Host docs with :

    godoc -http=:6060

[View here](http://localhost:6060/pkg/github.com/tyhal/crie/crie/#pg-overview)

## Local Binary

Local binary requiring all linters to be installed

```bash
    # Ensure your path contains $GOPATH/bin
    go install

    # Additional tools used to lint
    go get -u mvdan.cc/sh/cmd/shfmt golang.org/x/lint/golint
    pip3 install -r requirements.txt

    sudo npm install -g jsonlint2 remark-cli remark-preset-lint-recommended standard
    sudo apt install cppcheck shellcheck clang-format

    # TODO install help:
    # hadolint
```
