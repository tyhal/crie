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

## Useful for pinning versions in a Dockerfile

    RUN dpkg -l
    RUN pip freeze
