FROM alpine:3.21.3 as clang_layer
RUN apk --no-cache add clang

FROM alpine:3.21.3

# hadolint ignore=DL3059
RUN adduser -D standards

# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# [ OS pkgs] - We pull clang-format out specifically because we don't need the rest of clang
RUN apk --no-cache add gmp libxml2
COPY --from=clang_layer /usr/lib/libclang-cpp.so.10 /usr/lib/libclang-cpp.so.10
COPY --from=clang_layer /usr/lib/libLLVM-10.so /usr/lib/libLLVM-10.so
COPY --from=clang_layer /usr/bin/clang-format /bin/clang-format
RUN apk --no-cache add cppcheck asciidoctor