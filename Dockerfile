# syntax=docker/dockerfile:1.2
# ~~~ Languages ~~~

FROM tyhal/hadolint:0.0.2 as hadolint_layer
FROM hashicorp/terraform:1.1.3 as terraform_layer

FROM golang:1.17-alpine3.13 as go_layer
RUN apk --no-cache add git wget
ENV CGO_ENABLED=0

# TODO implement imp/golint.go
FROM go_layer as golint_layer
RUN go get -u golang.org/x/lint/golint

FROM go_layer as crie_layer
COPY go.mod /crie/go.mod
COPY go.sum /crie/go.sum
WORKDIR /crie
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY cmd /crie/cmd
COPY internal /crie/internal
COPY pkg /crie/pkg
RUN --mount=type=cache,target=/root/.cache/go-build go build ./cmd/crie

FROM alpine:3.13.0 as clang_layer
RUN apk --no-cache add clang

# ~~~           ~~~ ~~~~~~~~~~~~~~~~~ ~~~           ~~~
# ~~~~~~~~~~~~~~~~~ ~~~ TOP LAYER ~~~ ~~~~~~~~~~~~~~~~~
# ~~~           ~~~ ~~~~~~~~~~~~~~~~~ ~~~           ~~~

# Alpine :ok_hand:
FROM alpine:3.13.0
RUN apk --no-cache add git wget ca-certificates \
    && update-ca-certificates

# hadolint ignore=DL3059
RUN adduser -D standards

# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# [ NPM pkgs]
RUN apk add --no-cache npm && npm install -g jsonlint2 remark-cli remark-preset-lint-recommended standard

# [ OS pkgs] - We pull clang-format out specifically because we don't need the rest of clang
RUN apk --no-cache add gmp libxml2
COPY --from=clang_layer /usr/lib/libclang-cpp.so.10 /usr/lib/libclang-cpp.so.10
COPY --from=clang_layer /usr/lib/libLLVM-10.so /usr/lib/libLLVM-10.so
COPY --from=clang_layer /usr/bin/clang-format /bin/clang-format
RUN apk --no-cache add cppcheck asciidoctor

# [ Pips ]
ENV CRYPTOGRAPHY_DONT_BUILD_RUST=1
COPY requirements.txt /requirements.txt
ENV BUILD_LIBS="python3-dev build-base libffi-dev libressl-dev"
RUN apk add --no-cache python3 py3-pip $BUILD_LIBS \
    && pip3 --no-cache-dir install -r requirements.txt \
    && apk del --no-cache $BUILD_LIBS

COPY --from=hadolint_layer /bin/hadolint /bin/hadolint

# [ Bash ]
COPY --from=hadolint_layer /bin/shellcheck /bin/shellcheck

# [ Golang ]
COPY --from=go_layer /usr/local/go/bin/gofmt /bin/gofmt
COPY --from=golint_layer /go/bin/golint /bin/golint

# [ Terraform ]
COPY --from=terraform_layer /bin/terraform /bin/terraform

# [ Run Scripts ]
COPY --from=crie_layer /crie/crie /bin/crie

# [ Conf ]
COPY conf /etc/crie/
RUN chown -R standards:standards /etc/crie/
WORKDIR /l

# Give permission to non root to cache dirs
RUN mkdir /.standard-v14-cache /.ansible \
    && chmod -R o+rw /home /.standard-v14-cache /.ansible

ENTRYPOINT ["/bin/crie"]

USER standards
