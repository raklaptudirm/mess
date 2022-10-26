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

package attacks

import (
	"math/bits"

	"laptudirm.com/x/mess/pkg/chess/bitboard"
	"laptudirm.com/x/mess/pkg/chess/square"
)

func hyperbola(s square.Square, occ, mask bitboard.Board) bitboard.Board {
	r := bitboard.Squares[s]
	o := occ & mask // masked occupancy
	return ((o - 2*r) ^ reverse(reverse(o)-2*reverse(r))) & mask
}

func reverse(b bitboard.Board) bitboard.Board {
	return bitboard.Board(bits.Reverse64(uint64(b)))
}
