
FROM gcc:15.2.0
RUN apt-get update -y \
    && apt-get install -y \
        --no-install-recommends \
        cmake=3.31.6* \
    && rm -rf /var/lib/apt/lists/*
COPY . /src/
WORKDIR /build/
RUN cmake /src/ \
    && cmake --build .
CMD ["./example"]