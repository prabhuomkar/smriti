file(GLOB TEST_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/*_test.cc)

file(GLOB EXEC_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/../src/*.cc)

list(FILTER EXEC_SOURCES EXCLUDE REGEX ".*/main\\.cc$")

add_executable(worker_tests ${TEST_SOURCES})

target_include_directories(worker_tests
  PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}/../include
    ${CMAKE_CURRENT_SOURCE_DIR}/../src
)

target_sources(worker_tests PRIVATE ${EXEC_SOURCES})

target_link_libraries(worker_tests
                      gtest_main
                      gmock_main
                      spdlog::spdlog
                      cpr::cpr
                      simdjson::simdjson)

include(GoogleTest)

gtest_discover_tests(worker_tests)

file(GLOB BENCHMARK_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/*_benchmark.cc)

add_executable(worker_benchmarks ${BENCHMARK_SOURCES})

target_include_directories(worker_benchmarks
  PRIVATE
    ${CMAKE_CURRENT_SOURCE_DIR}/../include
    ${CMAKE_CURRENT_SOURCE_DIR}/../src
)

target_sources(worker_benchmarks PRIVATE ${EXEC_SOURCES})

target_link_libraries(worker_benchmarks
                      benchmark_main
                      gmock_main
                      spdlog::spdlog
                      cpr::cpr
                      simdjson::simdjson)