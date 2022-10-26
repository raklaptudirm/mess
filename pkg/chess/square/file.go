// Copyright © 2022 Rak Laptudirm <rak@laptudirm.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package square

// File represents a file on the chessboard.
type File int

// Constants representing every file on the chessboard.
const (
	FileA File = iota
	FileB
	FileC
	FileD
	FileE
	FileF
	FileG
	FileH
)

// Number of files
const FileN = 8

// String converts a File into it's string representation.
func (f File) String() string {
	files := [...]string{
		FileA: "a",
		FileB: "b",
		FileC: "c",
		FileD: "d",
		FileE: "e",
		FileF: "f",
		FileG: "g",
		FileH: "h",
	}

	return files[f]
}

// fileFrom creates an instance of a File from the given file id.
func fileFrom(id string) File {
	switch id {
	case "a":
		return FileA
	case "b":
		return FileB
	case "c":
		return FileC
	case "d":
		return FileD
	case "e":
		return FileE
	case "f":
		return FileF
	case "g":
		return FileG
	case "h":
		return FileH
	default:
		panic("new file: invalid file id")
	}
}
