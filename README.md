<h1 align="center">
    crie.
</h1>
<p align="center">
    <img src="https://raw.githubusercontent.com/tyhal/crie/main/doc/demo.svg?sanitize=true" width="572" alt="crie cli demo">
</p>
<p align="center">
    Effectively format and lint code for a variety of languages
</p>

## Features

#### This tool enables teams of developers to use static-analysis where they normally wouldn't:

This tool is ideal for teams who want:

    Immediate code quality checks without extensive setup
    Consistent linting across multiple languages
    Container-based isolation for tools
    A solution that "just works" out of the box

## Install

Getting the tool

```bash
    go install github.com/tyhal/crie/cmd/crie@latest
```

The suggested way to start running crie is to run `chk` at the top of your project and add `--continue` to see every error in the project (this will not change any code)

```bash
    crie chk --continue
```
