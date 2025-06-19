FROM golang:1.24-alpine3.21 AS go_layer
RUN apk --no-cache add git wget
ENV CGO_ENABLED=0

# TODO implement imp/golint.go
FROM go_layer AS golint_layer
RUN go install golang.org/x/lint/golint@latest
# hadolint ignore=DL3059
RUN go install github.com/errata-ai/vale@latest
# hadolint ignore=DL3059
RUN go install github.com/yoheimuta/protolint/cmd/protolint@latest


FROM alpine:3.13.0

RUN adduser -D standards

# [ Golang ]
COPY --from=go_layer /usr/local/go/bin/gofmt /bin/gofmt
COPY --from=golint_layer /go/bin/golint /bin/golint
COPY --from=golint_layer /go/bin/vale /bin/vale
COPY --from=golint_layer /go/bin/protolint /bin/protolint

ENV VALELINT_CONFIG_FILE=/home/standards/.config/.vale.ini
COPY --chown=standards:standards conf/markdown/.vale.ini $VALELINT_CONFIG_FILE

USER standards