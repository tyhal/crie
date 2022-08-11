	
FROM	haskell:9.4.1
COPY	src/	/src/
RUN	ghc -o hello /src/hello.hs
CMD	["./hello"]
