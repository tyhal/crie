# CLI Auto complete

you can generate autocomplete scripts for crie using `crie completion`

### Zsh

Here is **an** example:

```shell
mkdir -p ~/.zsh/completion
crie completion zsh > ~/.zsh/completion/_crie
```

Added to your `~/.zshrc`

```zsh
fpath=(~/.zsh/completion $fpath)
autoload -U compinit
compinit
```

Other shells are also available:

***

### Bash

```shell
crie completion bash
```

### Fish

```shell
crie completion fish
```

### Powershell

```shell
crie completion powershell
```
