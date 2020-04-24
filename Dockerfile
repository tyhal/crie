# Effective Standards Development


# ~~~ Languages ~~~

FROM tyhal/hadolint-hadolint:v1.15.0 as hadolint_layer

FROM golang:1.12.14-alpine3.9 as go_layer
RUN apk --no-cache add git wget
ENV CGO_ENABLED=0

# TODO implement imp/golint.go
FROM go_layer as golint_layer
RUN go get -u golang.org/x/lint/golint

# TODO implement imp/shfmt.go
FROM go_layer as shfmt_layer
RUN go get -u mvdan.cc/sh/cmd/shfmt
#RUN go get -u github.com/jessfraz/dockfmt

FROM go_layer as crie_layer
COPY go.mod /crie/go.mod
COPY go.sum /crie/go.sum
WORKDIR /crie
RUN go mod download
COPY cli /crie/cli
COPY api /crie/api
COPY imp /crie/imp
COPY crie.go /crie/crie.go
RUN go build

FROM alpine:3.11.6 as clang_layer
RUN apk --no-cache add clang

FROM alpine:3.11.6 as terraform_layer
RUN apk --no-cache add git wget zip
ENV TERRA_VER 0.12.17
RUN wget "https://releases.hashicorp.com/terraform/$TERRA_VER/terraform_${TERRA_VER}_$(uname -s | tr '[:upper:]' '[:lower:]')_amd64.zip"
RUN unzip "terraform_${TERRA_VER}_$(uname -s | tr '[:upper:]' '[:lower:]')_amd64.zip"

# ~~~           ~~~ ~~~~~~~~~~~~~~~~~ ~~~           ~~~
# ~~~~~~~~~~~~~~~~~ ~~~ TOP LAYER ~~~ ~~~~~~~~~~~~~~~~~
# ~~~           ~~~ ~~~~~~~~~~~~~~~~~ ~~~           ~~~

# Alpine :ok_hand:
FROM alpine:3.11.6
RUN apk --no-cache add git wget ca-certificates \
    && update-ca-certificates

RUN adduser -D standards

# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# [ NPM pkgs]
RUN apk add --no-cache nodejs-npm && npm install -g jsonlint2 remark-cli remark-preset-lint-recommended standard

# [ OS pkgs]
RUN apk --no-cache add gmp libxml2
 COPY --from=clang_layer /usr/lib/libLLVM-8.so /usr/lib/libLLVM-8.so
 COPY --from=clang_layer /usr/bin/clang-format /bin/clang-format
RUN apk --no-cache add cppcheck asciidoctor

# [ Pips ]
COPY requirements.txt /requirements.txt
ENV BUILD_LIBS="python3-dev build-base libffi-dev libressl-dev"
RUN apk add --no-cache python3 $BUILD_LIBS \
    && pip3 install -r requirements.txt \
    && apk del --no-cache $BUILD_LIBS

COPY --from=hadolint_layer /bin/hadolint /bin/hadolint

# [ Bash ]
COPY --from=shfmt_layer /go/bin/shfmt /bin/shfmt
COPY --from=hadolint_layer /bin/shellcheck /bin/shellcheck

# [ Golang ]
COPY --from=go_layer /usr/local/go/bin/gofmt /bin/gofmt
COPY --from=golint_layer /go/bin/golint /bin/golint

# [ Terraform ]
COPY --from=terraform_layer /terraform /bin/terraform

# [ Run Scripts ]
COPY --from=crie_layer /crie/crie /bin/crie

# [ Conf ]
COPY conf /etc/crie/
RUN chown -R standards:standards /etc/crie/
WORKDIR /l

# Give permission to non root to cache dirs
RUN mkdir /.standard-v14-cache /.ansible
RUN chmod -R o+rw /home /.standard-v14-cache /.ansible

ENTRYPOINT ["/bin/crie"]

USER standards
