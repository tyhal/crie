# Effective Standards Development


# ~~~ Languages ~~~

FROM alpine:3.9 as haskelldeps
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

FROM golang:1.11-alpine3.8 as godeps
RUN apk --no-cache add git wget
RUN apk --no-cache add clang clang-dev musl-dev git gcc
ENV CXX=clang++
ENV CC=clang
ARG GOARCH=amd64
ENV GOARCH=$GOARCH
RUN go get -u github.com/golang/lint/golint
RUN go get -u github.com/errata-ai/vale
RUN go get -u mvdan.cc/sh/cmd/shfmt
#RUN go get -u github.com/jessfraz/dockfmt
ENV CRIE github.com/tyhal/crie
COPY crie /go/src/$CRIE/crie
COPY api /go/src/$CRIE/api
RUN go get $CRIE/crie
RUN go build $CRIE/crie

FROM alpine:3.9 as zipdeps
RUN apk --no-cache add git wget zip
ENV TERRA_VER 0.11.8
RUN wget "https://releases.hashicorp.com/terraform/$TERRA_VER/terraform_${TERRA_VER}_$(uname -s | tr '[:upper:]' '[:lower:]')_amd64.zip"
RUN unzip "terraform_${TERRA_VER}_$(uname -s | tr '[:upper:]' '[:lower:]')_amd64.zip"
RUN pwd

# ~~~           ~~~ ~~~~~~~~~~~~~~~~~ ~~~           ~~~
# ~~~~~~~~~~~~~~~~~ ~~~ TOP LAYER ~~~ ~~~~~~~~~~~~~~~~~
# ~~~           ~~~ ~~~~~~~~~~~~~~~~~ ~~~           ~~~

# Alpine :ok_hand:
FROM alpine:3.9
RUN apk --no-cache add git wget ca-certificates \
    && update-ca-certificates

RUN adduser -D standards

# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# [ Python ]
RUN apk add --no-cache python3 python3-dev build-base \
    && pip3 install pylint==2.2.2 autopep8==1.4.3 \
    && apk del --no-cache python3-dev build-base

# [ Docker ]
RUN apk --no-cache add gmp
COPY --from=haskelldeps /root/.local/bin/hadolint /bin/hadolint
#COPY --from=godeps /go/bin/dockfmt /bin/dockfmt

# [ Bash ]
COPY --from=godeps /go/bin/shfmt /bin/shfmt
COPY --from=haskelldeps /root/.local/bin/shellcheck /bin/shellcheck

# [ YML ]
RUN apk add --no-cache python3 python3-dev build-base \
    && pip3 install yamllint==1.13.0 \
    && apk del --no-cache python3-dev build-base

# [ Javascript ]
RUN apk add --no-cache nodejs-npm \
    && npm install -g standard

# [ Golang ]
COPY --from=godeps /usr/local/go/bin/gofmt /bin/gofmt
COPY --from=godeps /go/bin/golint /bin/golint

# [ Markdown + AsciiDoctor ]
RUN apk add --no-cache nodejs-npm asciidoctor \
    && npm install -g remark-cli remark-preset-lint-recommended
COPY --from=godeps /go/bin/vale /bin/vale

# [ JSON ]
RUN apk add --no-cache nodejs-npm \
    && npm install -g jsonlint2

# [ CPP ]
RUN apk add --no-cache cppcheck
COPY --from=godeps /usr/bin/clang-format /bin/clang-format

# [ Doxygen ]
RUN apk add --no-cache doxygen

# [ Cmake ]
RUN pip3 install cmakelint==1.3.4.1




# [ Docker Compose ]
RUN pip3 install docker-compose==1.23.2

# [ Terraform ]
COPY --from=zipdeps /terraform /bin/terraform

# [ Run Scripts ]
COPY --from=godeps /go/bin/crie /bin/crie

# [ Conf ]
COPY conf /conf/
RUN chown -R standards:standards /conf
ENV PATH /node_modules/.bin:$PATH
WORKDIR /check

# BUG https://priapus.atlassian.net/browse/BUG-337
RUN mkdir /.standard-v12-cache
RUN chmod -R o+rw /.standard-v12-cache

ENTRYPOINT ["crie"]

USER standards
