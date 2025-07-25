<h1 align="center">
    crie.
</h1>
<p align="center">
    Universal meta-linter using containerized execution
</p>
<p align="center">
    <img src="https://raw.githubusercontent.com/tyhal/crie/main/doc/demo.svg?sanitize=true" width="572" alt="crie cli demo">
</p>

## Features

#### This tool enables teams of developers to use static-analysis where they wouldn't:

* No more config chaos - one simple setup for all your tools
* Quality checks that just work, right out of the box
* Container-based isolation for consistent tool execution
* Container runtime flexibility - supports both Docker and Podman
* Drop in and get coding - minimal setup required
* Catch every blind spot with full coverage detection

## Install

Installing Crie

```bash
    go install github.com/tyhal/crie/cmd/crie@latest
```

Install either [Podman](https://podman.io/docs/installation) or [Docker](https://docs.docker.com/get-started/get-docker/)

## Usage

The suggested way to start running crie is to run `chk` at the top of your project and add `--continue` to see every error in the project (this will not change any code)

```bash
    crie chk --continue
```

## Docs

* [Autocompletion](doc/completion.md) - Setup tab completion for your shell
* [Example Project Settings](doc/example.yml) - If you want to override any defaults, add more languages, or ignore certain files you can follow this example for your own optional crie.yml

***

<div align="center">
    <a href="https://codecov.io/gh/tyhal/crie"> 
        <img alt="coverage" src="https://codecov.io/gh/tyhal/crie/graph/badge.svg?token=SSAG0W1TZB"/> 
    </a>
</div>
