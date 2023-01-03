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

package bitboard

import (
	"math/bits"

	"laptudirm.com/x/mess/pkg/board/square"
)

// Hyperbola implements hyperbola quintessence given a from square,
// occupancy, and occupancy mask on the given bitboard.Board.
// https://www.chessprogramming.org/Hyperbola_Quintessence
func Hyperbola(s square.Square, occ, mask Board) Board {
	r := Squares[s]
	o := occ & mask // masked occupancy
	return ((o - 2*r) ^ reverse(reverse(o)-2*reverse(r))) & mask
}

// reverse is a simple function to reduce the verbosity of the code.
// It is inlined by the go compiler during compilation into a binary.
func reverse(b Board) Board {
	return Board(bits.Reverse64(uint64(b)))
}
