FROM alpine:3.22.2 AS clang_layer
RUN apk --no-cache add clang20-extra-tools
# RUN ldd /usr/bin/clang-format && exit 1

FROM alpine:3.22.2

# hadolint ignore=DL3059
RUN adduser -D standards

# ~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~~

# [ OS pkgs] - We pull clang-format out specifically because we don't need the rest of clang
# hadolint ignore=DL3059
RUN apk --no-cache add gmp libxml2 libffi zstd-libs
COPY --from=clang_layer /usr/lib/llvm20/lib/libclang-cpp.so.20.1 /usr/lib/llvm20/lib/libclang-cpp.so.20.1
COPY --from=clang_layer /usr/lib/llvm20/lib/libLLVM.so.20.1 /usr/lib/llvm20/lib/libLLVM.so.20.1
COPY --from=clang_layer /usr/bin/clang-format /bin/clang-format
RUN apk --no-cache add cppcheck asciidoctor

USER standards