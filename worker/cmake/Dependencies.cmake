include(FetchContent)
set(FETCHCONTENT_BASE_DIR ${CMAKE_SOURCE_DIR}/third_party)
function(fetch_dependency name repo tag)
    FetchContent_Declare(
        ${name}
        GIT_REPOSITORY ${repo}
        GIT_TAG        ${tag}
    )
    FetchContent_MakeAvailable(${name})
endfunction()

fetch_dependency(cpr https://github.com/libcpr/cpr.git 1.11.2)
fetch_dependency(simdjson https://github.com/simdjson/simdjson.git v3.12.3)
fetch_dependency(spdlog https://github.com/gabime/spdlog.git v1.15.2)

if(BUILD_TESTS)
    fetch_dependency(googletest https://github.com/google/googletest.git v1.16.0)
    fetch_dependency(benchmark https://github.com/google/benchmark.git v1.9.2)
endif()