FROM ubuntu:26.04 AS builder
ENV DEBIAN_FRONTEND=noninteractive
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        build-essential \
        ca-certificates \
        cmake \
        curl \
        git \
        libcurl4-openssl-dev \
        libgflags-dev \
        libgrpc++-dev \
        libmagick++-dev \
        libopencv-dev \
        libprotobuf-dev \
        libraw-dev \
        libyaml-cpp-dev \
        ninja-build \
        patch \
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
      -DBUILD_TESTS=OFF && \
    cmake --build build --target worker -j"$(nproc)"

FROM ubuntu:26.04 AS runtime
ENV DEBIAN_FRONTEND=noninteractive
ENV LD_LIBRARY_PATH=/opt/worker/lib:/opt/worker/lib/onnxruntime:${LD_LIBRARY_PATH}
RUN apt-get update && \
    apt-get install -y --no-install-recommends \
        ca-certificates \
        libimage-exiftool-perl \
        imagemagick \
        libcurl4 \
        libgflags2.2 \
        libgrpc++1.51t64 \
        libgrpc29t64 \
        libgomp1 \
        libheif1 \
        libheif-plugin-libde265 \
        libde265-0 \
        libmagick++-7.q16-5 \
        libmagickcore-7.q16-10 \
        libmagickcore-7.q16-10-extra \
        libmagickwand-7.q16-10 \
        libopencv-core410 \
        libopencv-dnn410 \
        libopencv-imgcodecs410 \
        libopencv-imgproc410 \
        libprotobuf32t64 \
        libraw23t64 \
        libyaml-cpp0.8 && \
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/build/worker .
COPY --from=builder /app/third_party/onnxruntime/lib /opt/worker/lib/onnxruntime
EXPOSE 5002
EXPOSE 15002
ENTRYPOINT ["/app/worker"]