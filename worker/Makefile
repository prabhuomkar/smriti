.PHONY: build
build:
	@mkdir -p build
	@cmake -S . -B build -DCMAKE_BUILD_TYPE=Release -DBUILD_TESTS=ON
	@cd build && make -j 4

.PHONY: build-with-tests
build-with-tests:
	@mkdir -p build
	@cmake -S . -B build -DCMAKE_PREFIX_PATH=${MY_INSTALL_DIR} -DBUILD_TESTS=ON
	@cd build && make -j 4

.PHONY: run
run:
	@./build/worker

.PHONY: test
test: 
	@cd build/tests && \
		./worker_tests && \
		./worker_benchmarks

.PHONY: cover
cover: 
	@echo "TODO(omkar)"

.PHONY: cover-html
cover-html: cover
	@echo "TODO(omkar)"

.PHONY: lint
lint:
	@cpplint --recursive include/ src/ tests/

.PHONY: proto
proto:
	@protoc -I../protos --cpp_out=./generated/protos ../protos/*.proto
	@protoc -I../protos --grpc_out=generate_mock_code=true:./generated/protos --plugin=protoc-gen-grpc=`which grpc_cpp_plugin` ../protos/*.proto
