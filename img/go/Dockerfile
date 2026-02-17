FROM golang:1.26-alpine3.22 AS go_layer
RUN apk --no-cache add git wget
ENV CGO_ENABLED=0

FROM go_layer AS revive_layer
RUN go install github.com/mgechev/revive@v1.14.0
FROM go_layer AS protolint_layer
RUN go install github.com/yoheimuta/protolint/cmd/protolint@v0.56.4

FROM go_layer AS run_layer

RUN adduser -D standards

# [ Golang ]
COPY --from=revive_layer /go/bin/revive /bin/revive
COPY --from=protolint_layer /go/bin/protolint /bin/protolint

ENV REVIVE_CONFIG_PATH=/etc/revive/revive.toml
COPY --chmod=644 conf/revive/revive.toml $REVIVE_CONFIG_PATH
RUN chmod -R 755 /etc/revive

USER standards