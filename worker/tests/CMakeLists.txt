file(GLOB TEST_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/*_test.cc)

file(GLOB EXEC_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/../src/*.cc)

file(GLOB GENERATED_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/../generated/protos/*.cc)

list(FILTER EXEC_SOURCES EXCLUDE REGEX ".*/main\\.cc$")

add_executable(worker_tests ${TEST_SOURCES} ${EXEC_SOURCES} ${GENERATED_SOURCES})

target_include_directories(worker_tests
  PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}/../generated
    ${CMAKE_CURRENT_SOURCE_DIR}/../third_party/simdjson/include
    ${CMAKE_CURRENT_SOURCE_DIR}/../third_party/simdjson/src
    ${CMAKE_CURRENT_SOURCE_DIR}/../third_party/cpr/include
    ${CMAKE_CURRENT_SOURCE_DIR}/../include
)

target_sources(worker_tests PRIVATE ${EXEC_SOURCES})

target_link_libraries(worker_tests
                      gtest_main
                      gmock_main
                      spdlog::spdlog
                      cpr::cpr
                      simdjson::simdjson
                      gRPC::grpc++
                      gRPC::grpc
                      protobuf::libprotobuf)

include(GoogleTest)

gtest_discover_tests(worker_tests)

file(GLOB BENCHMARK_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/*_benchmark.cc)

add_executable(worker_benchmarks ${BENCHMARK_SOURCES})

target_include_directories(worker_benchmarks
  PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}/../generated
    ${CMAKE_CURRENT_SOURCE_DIR}/../third_party/simdjson/include
    ${CMAKE_CURRENT_SOURCE_DIR}/../third_party/simdjson/src
    ${CMAKE_CURRENT_SOURCE_DIR}/../third_party/cpr/include
    ${CMAKE_CURRENT_SOURCE_DIR}/../include
)

target_sources(worker_benchmarks PRIVATE ${EXEC_SOURCES} ${EXEC_SOURCES} ${GENERATED_SOURCES})

target_link_libraries(worker_benchmarks
                      benchmark_main
                      gmock_main
                      spdlog::spdlog
                      cpr::cpr
                      simdjson::simdjson
                      gRPC::grpc++
                      gRPC::grpc
                      protobuf::libprotobuf)