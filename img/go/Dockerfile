FROM golang:1.25-alpine3.22 AS go_layer
RUN apk --no-cache add git wget
ENV CGO_ENABLED=0

# TODO implement pkg/linter/golint
FROM go_layer AS golint_layer
# hadolint ignore=DL3062
RUN go install golang.org/x/lint/golint@latest
# hadolint ignore=DL3059
RUN go install github.com/errata-ai/vale@v3.12.0
# hadolint ignore=DL3059
RUN go install github.com/yoheimuta/protolint/cmd/protolint@v0.56.4


FROM alpine:3.22.2

RUN adduser -D standards

# [ Golang ]
COPY --from=go_layer /usr/local/go/bin/gofmt /bin/gofmt
COPY --from=golint_layer /go/bin/golint /bin/golint
COPY --from=golint_layer /go/bin/vale /bin/vale
COPY --from=golint_layer /go/bin/protolint /bin/protolint

ENV VALELINT_CONFIG_PATH=/etc/vale/.vale.ini
COPY --chmod=644 conf/markdown/.vale.ini $VALELINT_CONFIG_PATH
RUN chmod -R 755 /etc/vale

# hadolint ignore=DL3059
RUN vale sync

USER standards


