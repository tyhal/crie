FROM golang:1.25-alpine3.22 AS go_layer
RUN apk --no-cache add git wget
ENV CGO_ENABLED=0

# TODO implement each as a pkg/linter/*
FROM go_layer AS golint_layer
# hadolint ignore=DL3062
RUN go install github.com/errata-ai/vale@latest \
    && go install github.com/mgechev/revive@latest \
    && go install golang.org/x/lint/golint@latest \
    && go install github.com/yoheimuta/protolint/cmd/protolint@latest


FROM alpine:3.22.2

RUN adduser -D standards

# [ Golang ]
COPY --from=go_layer /usr/local/go/bin/gofmt /bin/gofmt
COPY --from=golint_layer /go/bin/revive /bin/revive
COPY --from=golint_layer /go/bin/golint /bin/golint
COPY --from=golint_layer /go/bin/vale /bin/vale
COPY --from=golint_layer /go/bin/protolint /bin/protolint

ENV VALELINT_CONFIG_PATH=/etc/vale/.vale.ini
COPY --chmod=644 conf/markdown/.vale.ini $VALELINT_CONFIG_PATH
ENV REVIVE_CONFIG_PATH=/etc/revive/revive.toml
COPY --chmod=644 conf/revive/revive.toml $REVIVE_CONFIG_PATH
RUN chmod -R 755 /etc/vale /etc/revive

# hadolint ignore=DL3059
RUN vale sync

USER standards