cmake_minimum_required(VERSION 3.15)
project(Mess)

# ========================== #
# CMake Flags Initialization #
# ========================== #
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_FLAGS "-march=native")
set(CMAKE_CXX_FLAGS_DEBUG "-Wall -Wextra -g -ftime-report -fsanitize=undefined -fsanitize=address -fsanitize=leak")
set(CMAKE_CXX_FLAGS_RELEASE "-O3 -DNDEBUG")

# ======================== #
# Measure Compilation Time #
# ======================== #
set_property(GLOBAL PROPERTY RULE_LAUNCH_COMPILE "${CMAKE_COMMAND} -E time")

# =========================== #
# Setup CMake Package Manager #
# =========================== #

set(CPM_DOWNLOAD_VERSION 0.38.0)

if(CPM_SOURCE_CACHE)
    set(CPM_DOWNLOAD_LOCATION "${CPM_SOURCE_CACHE}/cpm/CPM_${CPM_DOWNLOAD_VERSION}.cmake")
elseif(DEFINED ENV{CPM_SOURCE_CACHE})
    set(CPM_DOWNLOAD_LOCATION "$ENV{CPM_SOURCE_CACHE}/cpm/CPM_${CPM_DOWNLOAD_VERSION}.cmake")
else()
    set(CPM_DOWNLOAD_LOCATION "${CMAKE_BINARY_DIR}/cmake/CPM_${CPM_DOWNLOAD_VERSION}.cmake")
endif()

if(NOT (EXISTS ${CPM_DOWNLOAD_LOCATION}))
    message(STATUS "Downloading CPM.cmake to ${CPM_DOWNLOAD_LOCATION}")
    file(DOWNLOAD
            https://github.com/TheLartians/CPM.cmake/releases/download/v${CPM_DOWNLOAD_VERSION}/CPM.cmake
            ${CPM_DOWNLOAD_LOCATION}
        )
endif()

include(${CPM_DOWNLOAD_LOCATION})

# ===================== #
# Setup Mess Executable #
# ===================== #
file(GLOB MessPackages
    "${CMAKE_CURRENT_SOURCE_DIR}/src/pkg/chess/*"
    "${CMAKE_CURRENT_SOURCE_DIR}/src/pkg/util/*"
)

file(GLOB MessLibraries
    "${CMAKE_CURRENT_SOURCE_DIR}/src/lib/*"
)

add_executable(${PROJECT_NAME} "${CMAKE_CURRENT_SOURCE_DIR}/src/main.cpp" ${MessPackages} ${MessLibraries} src/pkg/uci/uci.hpp)
