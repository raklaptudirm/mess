cmake_minimum_required(VERSION 3.15)
project(chess)

set(chess_VERSION_MAJOR 1)
set(chess_VERSION_MINOR 0)
set(chess_VERSION_PATCH 0)
set(chess_VERSION "${chess_VERSION_MAJOR}.${chess_VERSION_MINOR}.${chess_VERSION_PATCH}")

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
