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

if(APPLE)
    set(HOMEBREW_PREFIX "/opt/homebrew")
    include_directories(${HOMEBREW_PREFIX}/include)
    link_directories(${HOMEBREW_PREFIX}/lib)
    set(LIBOMP_ROOT ${HOMEBREW_PREFIX}/opt/libomp)
endif()

find_package(gRPC REQUIRED)
find_package(Protobuf REQUIRED)
find_package(ImageMagick COMPONENTS Magick++ REQUIRED)
find_package(PkgConfig REQUIRED)
find_package(OpenCV REQUIRED)
find_package(gflags REQUIRED)
pkg_check_modules(LibRaw REQUIRED libraw)

file(GLOB EXEC_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/src/*.cc)

file(GLOB GENERATED_SOURCES ${CMAKE_CURRENT_SOURCE_DIR}/generated/protos/*.cc)

file(GLOB_RECURSE PADDLEOCR_SOURCES ${paddleocr_SOURCE_DIR}/deploy/cpp_infer/src/*.cc)

add_executable(worker ${EXEC_SOURCES} ${GENERATED_SOURCES} ${PADDLEOCR_SOURCES})

set(TARGET_INCLUDE_DIRS 
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/generated
    ${CMAKE_CURRENT_SOURCE_DIR}/include
    ${paddleocr_SOURCE_DIR}/deploy/cpp_infer
    ${ONNXRUNTIME_INCLUDE_DIRS}
    ${POLYCLIPPING_INCLUDE_DIR}
    ${PADDLE_INFERENCE_INCLUDE_DIRS}
    ${IMAGEMAGICK_INCLUDE_DIRS}
    ${OpenCV_INCLUDE_DIRS}
    ${LibRaw_INCLUDE_DIRS})

set(TARGET_LINK_DIRS ${ONNXRUNTIME_LIB} ${PADDLE_INFERENCE_LIB})

set(TARGET_LINK_LIBS spdlog::spdlog
    cpr::cpr
    simdjson::simdjson
    gRPC::grpc++
    gRPC::grpc
    protobuf::libprotobuf
    ImageMagick::Magick++
    gflags
    polyclipping
    ${ONNXRUNTIME_LIBRARIES}
    ${PADDLE_INFERENCE_LIBRARIES}
    ${OpenCV_LIBRARIES}
    ${LibRaw_LIBRARIES})

if(APPLE)
    list(APPEND TARGET_INCLUDE_DIRS ${HOMEBREW_PREFIX}/opt/libomp/include)
    list(APPEND TARGET_LINK_DIRS ${HOMEBREW_PREFIX}/opt/libomp/lib)
    list(APPEND TARGET_LINK_LIBS omp)
endif()

target_include_directories(worker PRIVATE ${TARGET_INCLUDE_DIRS})

target_compile_definitions(worker PRIVATE
    DEFAULT_VERSION="${VERSION}"
    DEFAULT_GIT_SHA="${GIT_SHA}")

target_link_directories(worker PRIVATE ${TARGET_LINK_DIRS})
    
target_link_libraries(worker PRIVATE ${TARGET_LINK_LIBS})

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