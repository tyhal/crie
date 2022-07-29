	
FROM	haskell:9.2.4
COPY	src/	/src/
RUN	ghc -o hello /src/hello.hs
CMD	["./hello"]
