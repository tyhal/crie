	
FROM	haskell:9.2.3
COPY	src/	/src/
RUN	ghc -o hello /src/hello.hs
CMD	["./hello"]
