# syntax=docker/dockerfile:1.17

FROM golang:1.24-alpine3.21 AS go_layer
RUN apk --no-cache add git wget gcc gpgme-dev musl-dev btrfs-progs-dev
ENV CGO_ENABLED=1
COPY go.mod /crie/go.mod
COPY go.sum /crie/go.sum
WORKDIR /crie
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download
COPY cmd /crie/cmd
COPY pkg /crie/pkg
RUN --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg/mod \
    go build ./cmd/crie

# Alpine :ok_hand:
FROM alpine:3.22
RUN apk --no-cache add git wget ca-certificates \
    && update-ca-certificates


RUN adduser -D standards

# [ Run Scripts ]
COPY --from=go_layer /crie/crie /bin/crie

ENTRYPOINT ["/bin/crie"]

USER standards