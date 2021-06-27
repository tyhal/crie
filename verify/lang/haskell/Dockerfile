	
FROM	haskell:8.10.4
COPY	src/	/src/
RUN	ghc -o hello /src/hello.hs
CMD	["./hello"]
