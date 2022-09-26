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
	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
)

func whitePawnMovesFrom(s square.Square) bitboard.Board {
	b := board{origin: s}
	b.addAttack(0, -1)
	return b.board
}

func blackPawnMovesFrom(s square.Square) bitboard.Board {
	b := board{origin: s}
	b.addAttack(0, 1)
	return b.board
}

func whitePawnAttacksFrom(s square.Square) bitboard.Board {
	b := board{origin: s}

	b.addAttack(1, -1)  // left
	b.addAttack(-1, -1) // right

	return b.board
}

func blackPawnAttacksFrom(s square.Square) bitboard.Board {
	b := board{origin: s}

	b.addAttack(1, 1)  // left
	b.addAttack(-1, 1) // right

	return b.board
}

func Pawn(s, ep square.Square, c piece.Color, friends, enemies bitboard.Board) bitboard.Board {
	var occupied = friends | enemies
	var attackSet bitboard.Board

	enemies.Set(ep)

	switch c {
	case piece.WhiteColor:
		attackSet = whitePawnMoves[s] &^ occupied  // 1 square ahead
		attackSet |= (attackSet >> 8) &^ occupied  // 2 squares ahead
		attackSet |= whitePawnAttacks[s] & enemies // diagonal attacks

	case piece.BlackColor:
		attackSet = blackPawnMoves[s] &^ occupied  // 1 square ahead
		attackSet |= (attackSet >> 8) &^ occupied  // 2 squares ahead
		attackSet |= blackPawnAttacks[s] & enemies // diagonal attacks

	default:
		panic("pawn attacks: invalid color")
	}

	return attackSet
}
