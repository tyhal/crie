FROM rust:1.87.0-alpine3.22

RUN adduser -D standards

# rustfmt and clippy ship with the stable toolchain; be explicit so the layer
# is self-documenting and survives slim/minimal base variants
RUN rustup component add rustfmt

USER standards