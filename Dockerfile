# syntax=docker/dockerfile:1.2

FROM golang:1.23-alpine3.21 as go_layer
RUN apk --no-cache add git wget
ENV CGO_ENABLED=0
COPY go.mod /crie/go.mod
COPY go.sum /crie/go.sum
WORKDIR /crie
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY cmd /crie/cmd
COPY internal /crie/internal
COPY pkg /crie/pkg
RUN --mount=type=cache,target=/root/.cache/go-build go build ./cmd/crie

# Alpine :ok_hand:
FROM alpine:3.18.2
RUN apk --no-cache add git wget ca-certificates \
    && update-ca-certificates

# hadolint ignore=DL3059
RUN adduser -D standards

# [ Run Scripts ]
COPY --from=go_layer /crie/crie /bin/crie

ENTRYPOINT ["/bin/crie"]

USER standards
