# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #
# Copyright © 2023 Rak Laptudirm <rak@laptudirm.com>                        #
#                                                                           #
# Licensed under the Apache License, Version 2.0 (the "License");           #
# you may not use this file except in compliance with the License.          #
# You may obtain a copy of the License at                                   #
# http://www.apache.org/licenses/LICENSE-2.0                                #
#                                                                           #
# Unless required by applicable law or agreed to in writing, software       #
# distributed under the License is distributed on an "AS IS" BASIS,         #
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.  #
# See the License for the specific language governing permissions and       #
# limitations under the License.                                            #
# # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # # #

# ====================== #
# Engine Executable Name #
# ====================== #
EXE = mess

# ======================= #
# Code Building Directory #
# ======================= #
BUILD_DIR = cmake-build

# =========================================================== #
# Compiler Selection: If CC and CXX variables are not set,    #
# clang and clang++ compilers are used by default to provide  #
# binaries with the best performance, or CC and CXX are used. #
# =========================================================== #
ifeq ($(CC),)
    CC  = clang
    CXX = clang++
endif

# ===================================================== #
# OS Specific Tasks: Setup the correct copy command and #
# executable file extension according to the target os. #
# ===================================================== #
ifeq ($(OS),Windows_NT)
    CP = powershell cp
    EXTENSION = .exe
else
    CP = cp
    EXTENSION =
endif

# ======================= #
# Default Makefile Target #
# ======================= #
# TODO: set default to testing
.PHONY: default
default: release

# ============================ #
# All Make targets are Phonies #
# ============================ #
.PHONY: debug release testing   # Build Targets
.PHONY: cmake-setup cmake-build # CMake Utilities
.PHONY: clean help              # Other Utilities

# ======================================== #
# Targets for various binary Build Configs #
# ======================================== #

# TARGET=Debug
debug:
	@make cmake-build TARGET=Debug

# TARGET=Testing
testing:
	@make cmake-build TARGET=Testing

# Target=Release
release:
	@make cmake-build TARGET=Release

# =============================================================================== #
# CMake Setup: Setup Mess's build system with the Ninja Generator, which allows   #
# incremental compilations to the source code. This is useful for developers, but #
# probably won't cause any change for users building release binaries with it.    #
# =============================================================================== #
cmake-setup:
	cmake -B $(BUILD_DIR)/$(TARGET) -DCMAKE_BUILD_TYPE=$(TARGET)   \
	-DCMAKE_C_COMPILER=$(CC) -DCMAKE_CXX_COMPILER=$(CXX) \
	-G Ninja

cmake-build: cmake-setup
	cmake --build $(BUILD_DIR)/$(TARGET) --config $(TARGET) && \
	$(CP) $(BUILD_DIR)/$(TARGET)/Mess$(EXTENSION) $(EXE)$(EXTENSION)

# ===================================================================== #
# Utility Commands which are not build targets but useful for the user. #
# ===================================================================== #

clean:
	@rm -rf cmake-build
	@echo "make: removed cached build information"

help:
	@echo "USAGE:"
	@echo "    for builds: make [ clean ] [ TARGET ] { VARIABLE=VALUE }"
	@echo "    for others: make UTILITY"
	@echo
	@echo "TARGETS: (* = default)"
	@echo "    debug           create a mess binary with all  debugging info"
	@echo "    testing         create a mess binary with some debugging info"
	@echo "  * release         create a mess binary with no   debugging info"
	@echo
	@echo "UTILITIES:"
	@echo "    clean           remove all the cached target build information"
	@echo "    help            print this help message for ease of use"
	@echo
	@echo "VARIABLES: (VARIABLE = DEFAULT_VALUE)"
	@echo "    EXE = mess      the path to the built executable"
	@echo "    C   = clang     the C   compiler used for the build"
	@echo "    CXX = clang++   the C++ compiler used for the build"
	@echo
