FROM rust:1.95.0-alpine3.23 AS rust_layer
RUN apk --no-cache add llvm21-dev clang21-dev
#    clang21-static llvm21-dev

FROM rust_layer AS rubyfmt_layer
ADD https://github.com/fables-tales/rubyfmt.git#v0.13.0 /tmp/rubyfmt
WORKDIR /tmp/rubyfmt
RUN cargo build --release

FROM alpine:3.23.4

RUN adduser -D standards

# [ Golang ]
COPY --from=rubyfmt_layer /tmp/rubyfmt/release/rubyfmt-main /bin/revive

USER standards