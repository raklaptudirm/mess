cmake_minimum_required(VERSION 3.15)
project(Mess)

# ========================== #
# CMake Flags Initialization #
# ========================== #
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_FLAGS "-march=native -flto")
set(CMAKE_CXX_FLAGS_DEBUG "-Wall -Wextra -g -ftime-report -fsanitize=undefined -fsanitize=address -fsanitize=leak")
set(CMAKE_CXX_FLAGS_RELEASE "-O3 -DNDEBUG")

# ============================ #
# Profile Guided Optimizations #
# ============================ #
if(PGO MATCHES "True")
    message(STATUS "PGO: Enabled")
    if(EXISTS ${CMAKE_BINARY_DIR}/pgo.profdata)
        message(STATUS "PGO: Using existing profile data at: ${CMAKE_BINARY_DIR}/pgo.profdata")
        set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fprofile-use=${CMAKE_BINARY_DIR}/pgo.profdata")
    else()
        message(STATUS "PGO: No existing profile data found. Compiling with instrumentation.")
        set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -fprofile-instr-generate=${CMAKE_BINARY_DIR}/pgo.profraw")
    endif()
endif()

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

#set(SOURCE "${CMAKE_CURRENT_SOURCE_DIR}/src")

#set(MessMainFile "${SOURCE}/main.cpp")

#file(GLOB MessHppFiles "${SOURCE}/pkg/chess/**/*.hpp")
#file(GLOB MessCppFiles "${SOURCE}/pkg/chess/**/*.cpp")

#set(SOURCES
#    ${MessMainFile}
#    ${MessHppFiles} ${MessCppFiles}
##)

#add_executable(${PROJECT_NAME} ${SOURCES})

#target_link_libraries(${PROJECT_NAME} types)

add_subdirectory(src/pkg/types)
add_subdirectory(src/cmd)
