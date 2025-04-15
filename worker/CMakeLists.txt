cmake_minimum_required(VERSION 3.30)

set(CMAKE_CXX_STANDARD 17)

project(worker)

option(BUILD_TESTS "Build Tests and Benchmarks" OFF) 

if(NOT DEFINED VERSION)
    set(VERSION "dev")
endif()

execute_process(
    COMMAND git rev-parse HEAD
    WORKING_DIRECTORY ${CMAKE_SOURCE_DIR}
    OUTPUT_VARIABLE GIT_SHA
    ERROR_QUIET
    OUTPUT_STRIP_TRAILING_WHITESPACE
)
if(NOT GIT_SHA)
    set(GIT_SHA "-")
endif()

include(cmake/Dependencies.cmake)

file(GLOB EXEC_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/src/*.cc)

add_executable(worker ${EXEC_SOURCES})

target_include_directories(worker PRIVATE 
    ${CMAKE_CURRENT_SOURCE_DIR}/third_party/simdjson/include
    ${CMAKE_CURRENT_SOURCE_DIR}/third_party/simdjson/src
    ${CMAKE_CURRENT_SOURCE_DIR}/third_party/cpr/include
    ${CMAKE_CURRENT_SOURCE_DIR}/include)

target_compile_definitions(worker PRIVATE
    DEFAULT_VERSION="${VERSION}"
    DEFAULT_GIT_SHA="${GIT_SHA}")

target_link_libraries(worker PRIVATE
    spdlog::spdlog
    cpr::cpr
    simdjson::simdjson)

set_target_properties(worker PROPERTIES 
    CXX_STANDARD 17 
    CXX_EXTENSIONS OFF)

if(BUILD_TESTS)
    enable_testing()

    set(BENCHMARK_ENABLE_TESTING OFF)
    set(GTEST_CREATE_TESTS OFF)
    set(INSTALL_GTEST OFF)

    include(FetchContent)
    FetchContent_MakeAvailable(googletest benchmark)

    enable_testing()

    add_subdirectory(tests)
endif()