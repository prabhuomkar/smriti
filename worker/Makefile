GRPC_VERSION ?= v1.71.0
MY_INSTALL_DIR ?= ${HOME}/.local

.PHONY: install
install:
	@if [ ! -d "grpc" ]; then \
		git clone --recurse-submodules -b ${GRPC_VERSION} --depth 1 --shallow-submodules https://github.com/grpc/grpc; \
	else \
		echo "grpc already downloaded"; \
		cd grpc && git checkout ${GRPC_VERSION}; \
		git submodule update --init --recursive; \
	fi
	@cd grpc && \
		mkdir -p cmake/build && \
		cd cmake/build && \
		cmake -DgRPC_INSTALL=ON \
			-DgRPC_BUILD_TESTS=OFF \
			-DCMAKE_CXX_STANDARD=17 \
			-DCMAKE_INSTALL_PREFIX=${MY_INSTALL_DIR} \
      		../.. && \
		make -j 4 && \
		make install && \
		cd ../..
	@rm -rf grpc

.PHONY: build
build:
	@mkdir -p build
	@cmake -DCMAKE_PREFIX_PATH=${MY_INSTALL_DIR} . -B build
	@cd build && make -j 4

.PHONY: run
run:
	@./build/worker

.PHONY: test-install
test-install:
	@git submodule update --init --recursive

.PHONY: test
test: 
	@./build/worker_tests
	@./build/worker_benchmarks

.PHONY: cover
cover: 
	@echo "TODO(omkar)"

.PHONY: cover-html
cover-html: cover
	@echo "TODO(omkar)"

.PHONY: lint
lint:
	@cpplint --recursive include/ src/ tests/
	@clang-tidy src/*.cc -- -std=c++17 -I./include -I/usr/include -I/usr/local/include/ -I/opt/homebrew/include/

.PHONY: proto
proto:
	@protoc -I../ --cpp_out=./generated ../protos/*.proto
