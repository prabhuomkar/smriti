file(GLOB TEST_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/*_test.cc)

file(GLOB EXEC_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/../src/*.cc)

file(GLOB GENERATED_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/../generated/protos/*.cc)

list(FILTER EXEC_SOURCES EXCLUDE REGEX ".*/main\\.cc$")

add_executable(worker_tests ${TEST_SOURCES} ${EXEC_SOURCES} ${GENERATED_SOURCES})

target_include_directories(worker_tests PRIVATE ${TARGET_INCLUDE_DIRS})

target_link_directories(worker_tests PRIVATE ${TARGET_LINK_DIRS})

target_link_libraries(worker_tests PRIVATE 
                      gtest_main
                      gmock_main
                      ${TARGET_LINK_LIBS})

include(GoogleTest)

gtest_discover_tests(worker_tests)

file(GLOB BENCHMARK_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/*_benchmark.cc)

add_executable(worker_benchmarks ${BENCHMARK_SOURCES} ${EXEC_SOURCES} ${GENERATED_SOURCES})

target_include_directories(worker_benchmarks PRIVATE ${TARGET_INCLUDE_DIRS})

target_link_directories(worker_benchmarks PRIVATE ${TARGET_LINK_DIRS})

target_link_libraries(worker_benchmarks PRIVATE 
                      benchmark_main
                      gmock_main
                      ${TARGET_LINK_LIBS})