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

func PawnAll(s, ep square.Square, c piece.Color, occupied, enemies bitboard.Board) bitboard.Board {
	enemies.Set(ep)

	attackSet := PawnMoves[c][s] &^ occupied // 1 square ahead

	switch c {
	case piece.White:
		if s.Rank() == square.Rank2 {
			attackSet |= (attackSet >> 8) &^ occupied // 2 squares ahead
		}

	case piece.Black:
		if s.Rank() == square.Rank7 {
			attackSet |= (attackSet << 8) &^ occupied // 2 squares ahead
		}

	default:
		panic("pawn attacks: invalid color")
	}

	attackSet |= Pawn[c][s] & enemies // diagonal attacks

	return attackSet
}

func PawnPush(pawns, blockers bitboard.Board, color piece.Color) (bitboard.Board, bitboard.Board) {
	switch color {
	case piece.White:
		singlePush := pawns.North() &^ blockers
		doublePush := (singlePush.North() & bitboard.Rank4) &^ blockers
		return singlePush, doublePush
	case piece.Black:
		singlePush := pawns.South() &^ blockers
		doublePush := (singlePush.South() & bitboard.Rank5) &^ blockers
		return singlePush, doublePush
	default:
		panic("bad color")
	}
}

func PawnLeft(pawns, enemies bitboard.Board, color piece.Color) bitboard.Board {
	switch color {
	case piece.White:
		return pawns.North().West() & enemies
	case piece.Black:
		return pawns.South().West() & enemies
	default:
		panic("bad color")
	}
}

func PawnRight(pawns, enemies bitboard.Board, color piece.Color) bitboard.Board {
	switch color {
	case piece.White:
		return pawns.North().East() & enemies
	case piece.Black:
		return pawns.South().East() & enemies
	default:
		panic("bad color")
	}
}
