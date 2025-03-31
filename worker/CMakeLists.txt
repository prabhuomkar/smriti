cmake_minimum_required(VERSION 3.30)
project(worker)
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
add_subdirectory(third_party/spdlog)
add_executable(worker
    src/main.cc
    src/metadata.cc
    src/preview_thumbnail.cc
    src/places.cc)
target_include_directories(worker PRIVATE include)
target_compile_definitions(worker PRIVATE
    DEFAULT_VERSION="${VERSION}"
    DEFAULT_GIT_SHA="${GIT_SHA}"
)
target_link_libraries(worker PRIVATE spdlog::spdlog)
set_target_properties(worker PROPERTIES 
    CXX_STANDARD 17 
    CXX_EXTENSIONS OFF)