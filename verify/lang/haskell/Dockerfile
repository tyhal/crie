
FROM haskell:9.12.2
COPY src/ /src/
RUN ghc -o hello /src/hello.hs
CMD ["./hello"]