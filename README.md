<h1 align="center">
    crie.
</h1>
<p align="center">
    <img src="https://raw.githubusercontent.com/tyhal/crie/master/doc/demo.svg?sanitize=true" width="572" alt="crie cli demo">
</p>
<p align="center">
    Effectively format and lint code for a variety of languages
</p>

[![Docker Cloud Build Status](https://img.shields.io/docker/cloud/build/tyhal/crie.svg)](https://hub.docker.com/r/tyhal/crie)

## Features

#### This tool enables teams of developers to use static-analysis where they normally wouldn't:

-   Avoid remembering multiple run configurations
-   Avoid various install instructions

#### Quality of Life:

-   Git friendly - Check only changed files in the last few commits
-   Extendable for more languages
-   Fast and clean output
-   Batteries included but replaceable
-   Identify files lacking any static-analysis
-   Send logs via logstash to create analytics

## Install

Bundled all-together with Docker

```bash
    git clone https://github.com/tyhal/crie /tmp/crie; sudo /tmp/crie/script/crie install
```

The suggested way to start running crie is to run `chk` at the top of your project and add `--continue` to see every error in the project (this will not change any code)

```bash
    crie chk --continue
```
