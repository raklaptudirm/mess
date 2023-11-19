// Copyright © 2023 Rak Laptudirm <rak@laptudirm.com>
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

package fen

import (
	"fmt"
	"strings"
)

// FromString returns a FEN parsed from the given string.
func FromString(fen string) String {
	return FromSlice(strings.Fields(fen))
}

// FromSlice returns a FEN parsed from the given slice.
func FromSlice(fen []string) String {
	if len(fen) == 4 {
		fen = append(fen, "0", "1")
	}
	return [6]string(fen)
}

// String represents a String position string.
type String [6]string

// String returns the string representation of the given fen string.
func (fen String) String() string {
	fenString := fmt.Sprint([6]string(fen))
	return fenString[1 : len(fenString)-1]
}
