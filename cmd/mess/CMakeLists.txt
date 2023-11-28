cmake_minimum_required(VERSION 3.15)
project(mess)

# ========================== #
# CMake Flags Initialization #
# ========================== #
set(CMAKE_CXX_STANDARD 20)
set(CMAKE_CXX_FLAGS "-march=native -flto")
set(CMAKE_CXX_FLAGS_DEBUG "-Wall -Wextra -g -ftime-report -fsanitize=undefined -fsanitize=address -fsanitize=leak")
set(CMAKE_CXX_FLAGS_RELEASE "-O3 -DNDEBUG")
set(CMAKE_CXX_FLAGS_RELWITHDEBINFO "-O3")

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

# ===================== #
# Setup Mess Executable #
# ===================== #

add_executable(${PROJECT_NAME} main.cpp)

# =================================== #
# Include and Link Required Libraries #
# =================================== #

set(SOURCE ${PROJECT_SOURCE_DIR}/../..)

add_subdirectory(${SOURCE}/lib/strutil lib/strutil)

add_subdirectory(${SOURCE}/pkg/chess pkg/chess)

target_link_libraries(${PROJECT_NAME}
    chess
)
