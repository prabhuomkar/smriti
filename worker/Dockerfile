FROM ubuntu:26.04 AS builder
ENV DEBIAN_FRONTEND=noninteractive
ARG VERSION=dev
ARG GITSHA=-
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        cmake \
        git \
        libcurl4-openssl-dev \
        libgrpc++-dev \
        libmagick++-dev \
        libopencv-dev \
        libprotobuf-dev \
        libraw-dev \
        ninja-build \
        pkg-config \
        protobuf-compiler \
        protobuf-compiler-grpc && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY CMakeLists.txt ./
COPY cmake ./cmake
COPY include ./include
COPY src ./src
COPY --from=shared_protos . ./protos
RUN mkdir -p generated/protos && \
    protoc \
      --proto_path=protos \
      --cpp_out=generated/protos \
      --grpc_out=generated/protos \
      --plugin=protoc-gen-grpc=/usr/bin/grpc_cpp_plugin \
      protos/*.proto
RUN cmake -S . -B build \
      -G Ninja \
      -DCMAKE_POLICY_VERSION_MINIMUM=3.5 \
      -DCMAKE_BUILD_TYPE=Release \
      -DVERSION="${VERSION}" \
      -DGIT_SHA="${GITSHA}" \
      -DBUILD_TESTS=OFF && \
    cmake --build build --target worker -j"$(nproc)"
RUN strip --strip-unneeded /app/build/worker && \
    find /app/third_party/onnxruntime/lib -maxdepth 1 -type f -name 'libonnxruntime.so*' -exec strip --strip-unneeded {} + || true

FROM debian:trixie-slim AS runtime
ENV DEBIAN_FRONTEND=noninteractive
ENV LD_LIBRARY_PATH=/opt/worker/lib:/opt/worker/lib/onnxruntime:${LD_LIBRARY_PATH}
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        libimage-exiftool-perl \
        libcurl4 \
        libgrpc++1.51t64 \
        libgomp1 \
        libheif1 \
        libheif-plugin-libde265 \
        libmagick++-7.q16-5 \
        libopencv-core410 \
        libopencv-dnn410 \
        libopencv-imgcodecs410 \
        libopencv-imgproc410 \
        libprotobuf32t64 \
        libraw23t64 && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/*
WORKDIR /app
COPY --from=builder /app/build/worker .
COPY --from=builder /app/third_party/onnxruntime/lib/libonnxruntime.so* /opt/worker/lib/onnxruntime/
EXPOSE 5002
ENTRYPOINT ["/app/worker"]