	
FROM	haskell:8.0.2
COPY	src/	/src/
RUN	ghc -o hello /src/hello.hs
CMD	["./hello"]
