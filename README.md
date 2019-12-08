# Crie enforcement

[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/tyhal/crie.svg)](https://hub.docker.com/r/tyhal/crie)

Crie is an effective way to format and lint code for a variety of languages

-   Alpine based image to reduce download size
-   Extendable for more languages
-   Batteries included but replaceable - default configuration (/imp) is separated from core library

## Install options

#### Docker Based Install (Recommended)

Bundled all-together with Docker

```bash
    git clone https://github.com/tyhal/crie /tmp/crie; sudo /tmp/crie/script/crie install
```

#### Local Binary

Local binary requiring all linters to be installed

```bash
    # Ensure your path contains $GOPATH/bin
    git clone https://github.com/tyhal/crie /tmp/crie; cd /tmp/crie
    go install

    # Additional tools used to lint
    go get -u mvdan.cc/sh/cmd/shfmt golang.org/x/lint/golint
    pip3 install -r requirements.txt

    sudo npm install -g jsonlint2 remark-cli remark-preset-lint-recommended standard
    sudo apt install cppcheck shellcheck clang-format

    # TODO install help:
    # hadolint
```

## Usage

```bash
    crie chk
```

```bash
    crie fmt
```
