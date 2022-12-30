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
// Every vertical line of squares is called a file.
//
//	a b c d e f g h
//	a b c d e f g h
//	a b c d e f g h
//	a b c d e f g h
//	a b c d e f g h
//	a b c d e f g h
//	a b c d e f g h
//	a b c d e f g h
type File int8

// constants representing various files
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

// FileN is the number of files.
const FileN = 8

// String converts a File into it's string representation.
func (f File) String() string {
	const fileToStr = "abcdefgh"
	return string(fileToStr[f])
}

// FileFrom creates an instance of a File from the given file id.
func FileFrom(id string) File {
	return File(id[0] - 'a')
}
