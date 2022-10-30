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
	"laptudirm.com/x/mess/pkg/chess/bitboard"
	"laptudirm.com/x/mess/pkg/chess/square"
)

func bishop(s square.Square, occ bitboard.Board, isMask bool) bitboard.Board {
	diagonalMask := bitboard.Diagonals[s.Diagonal()]
	diagonalAttack := hyperbola(s, occ, diagonalMask)

	antiDiagonalMask := bitboard.AntiDiagonals[s.AntiDiagonal()]
	antiDiagonalAttack := hyperbola(s, occ, antiDiagonalMask)

	attacks := diagonalAttack | antiDiagonalAttack
	if isMask {
		attacks &^= bitboard.Rank1 | bitboard.Rank8 | bitboard.FileA | bitboard.FileH
	}

	return attacks
}

func Bishop(s square.Square, blockers bitboard.Board) bitboard.Board {
	return BishopTable.Probe(s, blockers)
}
