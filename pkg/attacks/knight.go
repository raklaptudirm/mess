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
	"laptudirm.com/x/mess/pkg/square"
)

// knightAttacksFrom generates an attack bitboard containing all the
// possible squares a knight can move to from the given square.
func knightAttacksFrom(from square.Square) bitboard.Board {
	b := board{origin: from}

	// set all possible attack squares
	b.addAttack(2, 1)   // soEaEa
	b.addAttack(1, 2)   // soSoEa
	b.addAttack(1, -2)  // noNoEa
	b.addAttack(2, -1)  // noEaEa
	b.addAttack(-1, 2)  // soSoWe
	b.addAttack(-2, 1)  // soWeWe
	b.addAttack(-2, -1) // noWeWe
	b.addAttack(-1, -2) // noNoWe

	return b.board
}

// Knight acts as a wrapper method on the precalculated attacks bitboards
// of knights from every square on the board. It returns the attack
// bitboard for the provided square.
func Knight(s square.Square, friends bitboard.Board) bitboard.Board {
	return knightAttacks[s] &^ friends
}
