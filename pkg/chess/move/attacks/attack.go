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

//go:generate go run laptudirm.com/x/mess/pkg/generator/attack

import (
	"laptudirm.com/x/mess/pkg/chess/bitboard"
	"laptudirm.com/x/mess/pkg/chess/piece"
	"laptudirm.com/x/mess/pkg/chess/square"
)

func PawnPush(pawns bitboard.Board, color piece.Color) bitboard.Board {
	switch color {
	case piece.White:
		return pawns.North()
	case piece.Black:
		return pawns.South()
	default:
		panic("bad color")
	}
}

func PawnsLeft(pawns bitboard.Board, color piece.Color) bitboard.Board {
	switch color {
	case piece.White:
		return pawns.North().West()
	case piece.Black:
		return pawns.South().West()
	default:
		panic("bad color")
	}
}

func PawnsRight(pawns bitboard.Board, color piece.Color) bitboard.Board {
	switch color {
	case piece.White:
		return pawns.North().East()
	case piece.Black:
		return pawns.South().East()
	default:
		panic("bad color")
	}
}

func Bishop(s square.Square, blockers bitboard.Board) bitboard.Board {
	return bishopTable.Probe(s, blockers)
}

func Rook(s square.Square, blockers bitboard.Board) bitboard.Board {
	return rookTable.Probe(s, blockers)
}

func Queen(s square.Square, occ bitboard.Board) bitboard.Board {
	return Rook(s, occ) | Bishop(s, occ)
}
