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

# ==================== #
# Base Executable Name #
# ==================== #
EXE = ./bin/mess

# ==================== #
# Make a Native Binary #
# ==================== #

# Development Binary: Detailed Versioning
dev-binary:
	go run ./scripts/build -- EXE=${EXE} -- dev-build

# Release Binary: Tagged Versioning
release-binary:
	go run ./scripts/build -- EXE=${EXE} -- release-build

# ============================ #
# Make Cross-Platform Binaries #
# ============================ #
release-binaries:
	go run ./scripts/build -- GOOS=linux   GOARCH=arm   EXE=${EXE}-linux-arm     -- release-build
	go run ./scripts/build -- GOOS=linux   GOARCH=arm64 EXE=${EXE}-linux-arm64   -- release-build
	go run ./scripts/build -- GOOS=linux   GOARCH=amd64 EXE=${EXE}-linux-amd64   -- release-build
	go run ./scripts/build -- GOOS=darwin  GOARCH=amd64 EXE=${EXE}-darwin-amd64  -- release-build
	go run ./scripts/build -- GOOS=darwin  GOARCH=arm64 EXE=${EXE}-darwin-arm64  -- release-build
	go run ./scripts/build -- GOOS=windows GOARCH=amd64 EXE=${EXE}-windows-amd64 -- release-build
	go run ./scripts/build -- GOOS=windows GOARCH=386   EXE=${EXE}-windows-386   -- release-build

# ========================= #
# Make Generated Code Files #
# ========================= #
code-files:
	go generate ./...
