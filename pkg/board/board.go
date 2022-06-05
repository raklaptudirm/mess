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

// Package board implements a complete chess board along with valid move
// generation and other related utilities.
package board

import (
	"fmt"

	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
)

// Board represents the state of a chessboard at a given position.
type Board struct {
	// position data
	position  [64]piece.Piece       // 8x8 for fast lookup
	bitboards [13]bitboard.Bitboard // bitboards for eval

	sideToMove piece.Color

	enPassantTarget square.Square

	// castling rights
	blackCastleKingside  bool
	blackCastleQueenside bool
	whiteCastleKingside  bool
	whiteCastleQueenside bool

	// move counters
	halfMoves int
	fullMoves int
}

// String converts a Board into a human readable string.
func (b Board) String() string {
	divider := "+---+---+---+---+---+---+---+---+\n"

	s := divider

	for rank := 0; rank < 8; rank++ {
		s += "| "

		for file := 0; file < 8; file++ {
			square := square.Square(rank*8 + file)
			s += b.position[square].String() + " | "
		}

		s += fmt.Sprintln(8 - rank)
		s += divider
	}

	s += "  a   b   c   d   e   f   g   h\n\n"
	s += fmt.Sprintf("FEN: %s\n", b.FEN())

	return s
}

// MakeMove plays a legal move on the Board.
func (b *Board) MakeMove(from, to square.Square) {

	// basic legality checks
	switch {
	case b.position[from] == piece.Empty:
		panic("invalid move: empty from square")

	case b.position[from].Color() != b.sideToMove:
		panic("invalid move: from square occupied by enemy piece")

	case b.position[to] == piece.Empty:
		break

	case b.position[to].Color() == b.sideToMove:
		panic("invalid move: to square occupied by friendly piece")
	}

	// move piece in 8x8 board
	b.position[to] = b.position[from]
	b.position[from] = piece.Empty

	// move piece in bitboard
	b.bitboards[b.position[to]].Unset(from)
	b.bitboards[b.position[to]].Set(to)

	// switch turn
	switch b.sideToMove {
	case piece.WhiteColor:
		b.sideToMove = piece.BlackColor
	case piece.BlackColor:
		b.sideToMove = piece.WhiteColor
		b.fullMoves++ // turn completed
	}
}
