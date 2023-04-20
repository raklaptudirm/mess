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

package attacks

//go:generate go run laptudirm.com/x/mess/internal/generator/attack

import (
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
)

// Of returns the attack set of the given piece on the given square
// and with the given blocker set. The blocker set is unused while
// calculating the attacks sets of non-sliding pieces.
func Of(p piece.Piece, s square.Square, blockers bitboard.Board) bitboard.Board {
	switch p.Type() {
	case piece.Pawn:
		return Pawn[p.Color()][s]
	case piece.Knight:
		return Knight[s]
	case piece.Bishop:
		return Bishop(s, blockers)
	case piece.Rook:
		return Rook(s, blockers)
	case piece.Queen:
		return Queen(s, blockers)
	case piece.King:
		return King[s]
	default:
		panic("attacks.Of: unknown piece type")
	}
}

// PawnPush gives the result after pushing every pawn in the given BB.
func PawnPush(pawns bitboard.Board, color piece.Color) bitboard.Board {
	return pawns.Up(color)
}

func Pawns(pawns bitboard.Board, color piece.Color) bitboard.Board {
	return PawnsLeft(pawns, color) | PawnsRight(pawns, color)
}

// PawnsLeft gives the result after every pawn captures left in the given BB.
func PawnsLeft(pawns bitboard.Board, color piece.Color) bitboard.Board {
	return pawns.Up(color).West()
}

// PawnsRight gives the result after every pawn captures right in the given BB.
func PawnsRight(pawns bitboard.Board, color piece.Color) bitboard.Board {
	return pawns.Up(color).East()
}

// Bishop returns the attack set for a bishop on the given square and with
// the given blocker set(occupied squares).
func Bishop(s square.Square, blockers bitboard.Board) bitboard.Board {
	return bishopTable.Probe(s, blockers)
}

// Rook returns the attack set for a rook on the given square and with
// the given blocker set(occupied squares).
func Rook(s square.Square, blockers bitboard.Board) bitboard.Board {
	return rookTable.Probe(s, blockers)
}

// Queen returns the attack set for a queen on the given square and with
// the given blocker set(occupied squares). It is calculated as the union
// of the attack sets of a bishop and a rook on the given square.
func Queen(s square.Square, occ bitboard.Board) bitboard.Board {
	return Rook(s, occ) | Bishop(s, occ)
}
