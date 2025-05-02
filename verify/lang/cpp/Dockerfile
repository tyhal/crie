
FROM	gcc:15.1.0
RUN	apt-get update -y \
	&& apt-get install -y \
	--no-install-recommends \
	cmake=3.18.4* \
	&& rm -rf /var/lib/apt/lists/*
COPY	.	/src/
WORKDIR	/build/
RUN	cmake /src/ \
	&& cmake --build .
CMD	["./example"]
