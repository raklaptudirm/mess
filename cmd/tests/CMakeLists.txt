cmake_minimum_required(VERSION 3.15)
project(tests)

set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_FLAGS "-march=native -O3")

enable_testing()

set(SOURCE ${PROJECT_SOURCE_DIR}/../..)

configure_file(${SOURCE}/tests/common/common.hpp.in ${PROJECT_BINARY_DIR}/common.hpp)

file(GLOB TEST_FILES ${SOURCE}/tests/**/*.cpp)
add_executable(${PROJECT_NAME} ${TEST_FILES})

add_subdirectory(${SOURCE}/pkg/chess   pkg/chess)
add_subdirectory(${SOURCE}/lib/catch2  lib/catch2)
add_subdirectory(${SOURCE}/lib/strutil lib/strutil)

target_link_libraries(${PROJECT_NAME}
    chess strutil
    Catch2::Catch2WithMain
)

target_include_directories(${PROJECT_NAME}
    PUBLIC
        ${CMAKE_BINARY_DIR}
)
