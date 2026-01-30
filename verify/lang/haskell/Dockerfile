
FROM haskell:9.14.1
COPY src/ /src/
RUN ghc -o hello /src/hello.hs
CMD ["./hello"]