# Effective Standards Development


# ~~~ Languages ~~~

FROM alpine:3.9.4 as haskell_layer
RUN apk add --no-cache git ghc=8.4.3-r0 xz wget build-base make ca-certificates \
        && update-ca-certificates
ENV HADOVER=tags/v1.15.0
RUN git clone --recursive https://github.com/lukasmartinelli/hadolint.git
WORKDIR /hadolint
RUN git checkout "$HADOVER"
ENV STACK_DIR stack-1.9.1-linux-x86_64-static
RUN wget --no-check-certificate -O  /stack.tar.gz https://github.com/commercialhaskell/stack/releases/download/v1.9.1/$STACK_DIR.tar.gz
RUN tar -xvzf /stack.tar.gz
RUN chmod +x $STACK_DIR/stack
RUN $STACK_DIR/stack update
RUN $STACK_DIR/stack --system-ghc install
RUN $STACK_DIR/stack --system-ghc install ShellCheck

FROM golang:1.12.7-alpine3.9 as go_layer
RUN apk --no-cache add git wget
ENV CGO_ENABLED=0

FROM go_layer as golint_layer
RUN go get -u golang.org/x/lint/golint

FROM go_layer as vale_layer
RUN go get -u github.com/errata-ai/vale

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
COPY crie.go /crie/crie.go
RUN go build 

FROM alpine:3.9.4 as clang_layer
RUN apk --no-cache add clang

FROM alpine:3.9.4 as terraform_layer
RUN apk --no-cache add git wget zip
ENV TERRA_VER 0.11.13
RUN wget "https://releases.hashicorp.com/terraform/$TERRA_VER/terraform_${TERRA_VER}_$(uname -s | tr '[:upper:]' '[:lower:]')_amd64.zip"
RUN unzip "terraform_${TERRA_VER}_$(uname -s | tr '[:upper:]' '[:lower:]')_amd64.zip"
RUN pwd

# ~~~           ~~~ ~~~~~~~~~~~~~~~~~ ~~~           ~~~
# ~~~~~~~~~~~~~~~~~ ~~~ TOP LAYER ~~~ ~~~~~~~~~~~~~~~~~
# ~~~           ~~~ ~~~~~~~~~~~~~~~~~ ~~~           ~~~

# Alpine :ok_hand:
FROM alpine:3.9.4
RUN apk --no-cache add git wget ca-certificates \
    && update-ca-certificates

RUN adduser -D standards

# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# [ Pips ]
COPY requirements.txt /requirements.txt
ENV BUILD_LIBS="python3-dev build-base libffi-dev libressl-dev"
RUN apk add --no-cache python3 $BUILD_LIBS \
    && pip3 install -r requirements.txt \
    && apk del --no-cache $BUILD_LIBS

# [ Docker ]
RUN apk --no-cache add gmp
COPY --from=haskell_layer /root/.local/bin/hadolint /bin/hadolint
#COPY --from=go_layer /go/bin/dockfmt /bin/dockfmt

# [ Bash ]
COPY --from=shfmt_layer /go/bin/shfmt /bin/shfmt
COPY --from=haskell_layer /root/.local/bin/shellcheck /bin/shellcheck

# [ Javascript ]
RUN apk add --no-cache nodejs-npm \
    && npm install -g standard

# [ Golang ]
COPY --from=go_layer /usr/local/go/bin/gofmt /bin/gofmt
COPY --from=golint_layer /go/bin/golint /bin/golint

# [ Markdown + AsciiDoctor ]
RUN apk add --no-cache nodejs-npm asciidoctor \
    && npm install -g remark-cli remark-preset-lint-recommended
COPY --from=vale_layer /go/bin/vale /bin/vale

# [ JSON ]
RUN apk add --no-cache nodejs-npm \
    && npm install -g jsonlint2

# [ CPP ]
RUN apk add --no-cache cppcheck
COPY --from=clang_layer /usr/bin/clang-format /bin/clang-format

# [ Terraform ]
COPY --from=terraform_layer /terraform /bin/terraform

# [ Run Scripts ]
COPY --from=crie_layer /crie/crie /bin/crie

# [ Conf ]
COPY conf /conf/
RUN chown -R standards:standards /conf
ENV PATH /node_modules/.bin:$PATH
WORKDIR /check

# Give permission to non root to cache dirs
RUN mkdir /.standard-v14-cache /.ansible
RUN chmod -R o+rw /.standard-v14-cache /.ansible

ENTRYPOINT ["/bin/crie"]

USER standards
