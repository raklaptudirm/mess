cmake_minimum_required(VERSION 3.15)
project(mess_binary)

set(SOURCE "${PROJECT_SOURCE_DIR}/..")

set(MessMainFile "${SOURCE}/cmd/main.cpp")

file(GLOB MessHppFiles "${SOURCE}/pkg/chess/**/*.hpp")
file(GLOB MessCppFiles "${SOURCE}/lib/**/*.h")

set(SOURCES
    ${MessMainFile}
    ${MessHppFiles} ${MessCppFiles}
)

add_executable(${PROJECT_NAME} ${SOURCES})

target_link_libraries(${PROJECT_NAME} types)
