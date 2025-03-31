FROM ubuntu:24.10 AS builder
RUN apt-get update && apt-get install -y \
    build-essential \
    cmake \
    git \
    wget \
    autoconf \
    automake \
    libtool \
    pkg-config \
    curl \
    && rm -rf /var/lib/apt/lists/*
WORKDIR /opt
RUN git clone --recurse-submodules -b v1.71.0 https://github.com/grpc/grpc.git && \
    cd grpc && \
    mkdir -p cmake/build && cd cmake/build && \
    cmake -DgRPC_INSTALL=ON \
          -DgRPC_BUILD_TESTS=OFF \
          -DCMAKE_INSTALL_PREFIX=/usr/local \
          ../.. && \
    make -j$(nproc) && \
    make install
WORKDIR /app
COPY . /app
RUN mkdir -p build && cd build && \
    cmake -DCMAKE_PREFIX_PATH=/usr/local .. && \
    make -j$(nproc)

FROM gcr.io/distroless/cc
WORKDIR /app
COPY --from=builder /app/build/worker .
EXPOSE 5002
EXPOSE 15002
CMD ["/app/worker"]