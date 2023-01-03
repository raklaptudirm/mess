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

package main

import (
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/board/square"
)

func whitePawnAttacksFrom(s square.Square) bitboard.Board {
	pawnUp := bitboard.Squares[s].North()
	return pawnUp.East() | pawnUp.West()
}

func blackPawnAttacksFrom(s square.Square) bitboard.Board {
	pawnUp := bitboard.Squares[s].South()
	return pawnUp.East() | pawnUp.West()
}

// knightAttacksFrom generates an attack bitboard containing all the
// possible squares a knight can move to from the given square.
func knightAttacksFrom(from square.Square) bitboard.Board {
	knight := bitboard.Squares[from]

	knightNorth := knight.North().North()
	knightSouth := knight.South().South()

	knightEast := knight.East().East()
	knightWest := knight.West().West()

	knightAttacks := knightNorth.East() | knightNorth.West()
	knightAttacks |= knightSouth.East() | knightSouth.West()

	knightAttacks |= knightEast.North() | knightEast.South()
	knightAttacks |= knightWest.North() | knightWest.South()

	return knightAttacks
}

// kingAttacksFrom generates an attack bitboard containing all the
// possible squares a king can move to from the given square.
func kingAttacksFrom(from square.Square) bitboard.Board {
	king := bitboard.Squares[from]

	kingNorth := king.North()
	kingSouth := king.South()
	kingEast := king.East()
	kingWest := king.West()

	kingAttacks := kingNorth | kingSouth | kingEast | kingWest

	kingAttacks |= kingNorth.East() | kingNorth.West()
	kingAttacks |= kingSouth.East() | kingSouth.West()

	return kingAttacks
}
