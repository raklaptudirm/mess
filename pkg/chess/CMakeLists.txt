cmake_minimum_required(VERSION 3.15)
project(chess VERSION 1.0.0)

configure_file(${PROJECT_SOURCE_DIR}/chess/version.hpp.in ${PROJECT_SOURCE_DIR}/chess/version.hpp)

add_library(${PROJECT_NAME} INTERFACE)

target_include_directories(${PROJECT_NAME}
    INTERFACE
        ${CMAKE_CURRENT_SOURCE_DIR}
)

target_link_libraries(${PROJECT_NAME}
    INTERFACE
        strutil
)
