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
	"laptudirm.com/x/mess/pkg/board/mailbox"
	"laptudirm.com/x/mess/pkg/castling"
	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
	"laptudirm.com/x/mess/pkg/zobrist"
)

// Board represents the state of a chessboard at a given position.
type Board struct {
	// position data
	hash      zobrist.Key
	position  mailbox.Board           // 8x8 for fast lookup
	bitboards [piece.N]bitboard.Board // bitboards for eval

	// useful bitboards
	friends bitboard.Board
	enemies bitboard.Board

	sideToMove      piece.Color
	enPassantTarget square.Square
	castlingRights  castling.Rights

	// move counters
	halfMoves int
	fullMoves int
}

// String converts a Board into a human readable string.
func (b Board) String() string {
	return fmt.Sprintf("%s\nFen: %s\nKey: %X\n", b.position, b.FEN(), b.hash)
}

func (b *Board) ClearSquare(s square.Square) {
	p := b.position[s]

	// the piece can only be in one of the bitboards, so
	// a conditional is unnecessary and both can be unset
	b.friends.Unset(s) // friends bitboard
	b.enemies.Unset(s) // enemies bitboard

	// remove piece from other records
	b.bitboards[p].Unset(s)             // piece bitboard
	b.position[s] = piece.NoPiece       // mailbox board
	b.hash ^= zobrist.PieceSquare[p][s] // zobrist hash
}

func (b *Board) FillSquare(s square.Square, p piece.Piece) {
	if p.Color() == b.sideToMove {
		b.friends.Set(s) // friends bitboard
	} else {
		b.enemies.Set(s) // enemies bitboard
	}

	b.bitboards[p].Set(s)               // piece bitboard
	b.position[s] = p                   // mailbox board
	b.hash ^= zobrist.PieceSquare[p][s] // zobrist hash
}
