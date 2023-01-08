# Copyright © 2023 Rak Laptudirm <rak@laptudirm.com>
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
# http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Executable Name
EXE = mess

# Default Build Target
build:
	go build -o ${EXE} .

# Cross Build Target
build-all:
	GOOS=linux   GOARCH=arm   go build -o ${EXE}-linux-arm
	GOOS=linux   GOARCH=arm64 go build -o ${EXE}-linux-arm64
	GOOS=linux   GOARCH=amd64 go build -o ${EXE}-linux-amd64
	GOOS=darwin  GOARCH=amd64 go build -o ${EXE}-darwin-amd64
	GOOS=darwin  GOARCH=arm64 go build -o ${EXE}-darwin-arm64
	GOOS=windows GOARCH=amd64 go build -o ${EXE}-windows-amd64.exe
	GOOS=windows GOARCH=386   go build -o ${EXE}-windows-386.exe	
