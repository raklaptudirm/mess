// Copyright © 2022 Rak Laptudirm <raklaptudirm@gmail.com>
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
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/move"
	"laptudirm.com/x/mess/pkg/square"
)

// kingAttacksFrom generates an attack bitboard containing all the
// possible squares a king can move to from the given square.
func kingAttacksFrom(from square.Square) bitboard.Board {
	b := board{origin: from}

	// set all possible attack squares
	b.addAttack(1, 0)   // E
	b.addAttack(1, 1)   // SE
	b.addAttack(0, 1)   // S
	b.addAttack(-1, 0)  // W
	b.addAttack(0, -1)  // N
	b.addAttack(1, -1)  // NE
	b.addAttack(-1, 1)  // SW
	b.addAttack(-1, -1) // NW

	return b.board
}

// King acts as a wrapper method for the precalculated attack bitboards of
// a king from every position on the chessboard. It returns the attack
// bitboard for the provided square.
func King(s square.Square, friends, occupied bitboard.Board, cr move.CastlingRights) bitboard.Board {
	base := kingAttacks[s] &^ friends

	switch s {
	case square.E1:
		kingsideMask := bitboard.Board(0x6000000000000000)
		queensideMask := bitboard.Board(0xe00000000000000)

		if cr&move.CastleWhiteKingside != 0 && occupied&kingsideMask == 0 {
			base.Set(square.G1)
		}

		if cr&move.CastleWhiteQueenside != 0 && occupied&queensideMask == 0 {
			base.Set(square.C1)
		}
	case square.E8:
		kingsideMask := bitboard.Board(0x60)
		queensideMask := bitboard.Board(0xe)

		if cr&move.CastleBlackKingside != 0 && occupied&kingsideMask == 0 {
			base.Set(square.G8)
		}

		if cr&move.CastleBlackQueenside != 0 && occupied&queensideMask == 0 {
			base.Set(square.C8)
		}
	}

	return base
}
