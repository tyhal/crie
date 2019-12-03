# Crie enforcement

[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/tyhal/crie.svg)](https://hub.docker.com/r/tyhal/crie)

Crie is an effective way to format and lint code for a variety of languages

-   Alpine based image to reduce download size
-   Extendable for more languages
-   Opinionated configuration to avoid wasting time discussing styles

## Install options

#### Docker Based Install

Bundled all-together with Docker 

```bash
    git clone https://github.com/tyhal/crie /tmp/crie; sudo /tmp/crie/script/crie install
```

#### Local Binary

Local binary requiring all linters to be installed

```bash
    go get -u github.com/tyhal/crie
```

## Usage

```bash
    crie chk
```

```bash
    crie fmt
```
