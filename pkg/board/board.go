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

	"laptudirm.com/x/mess/pkg/attacks"
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/board/mailbox"
	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
)

// Board represents the state of a chessboard at a given position.
type Board struct {
	// position data
	position  mailbox.Board      // 8x8 for fast lookup
	bitboards [13]bitboard.Board // bitboards for eval

	// useful bitboards
	friends bitboard.Board
	enemies bitboard.Board

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
	return fmt.Sprintf("%s\nFEN: %s\n", b.position, b.FEN())
}

// MakeMove plays a legal move on the Board.
func (b *Board) MakeMove(from, to square.Square) {
	if attackSet := b.MovesOf(from); !attackSet.IsSet(to) {
		// move not in attack board, illegal move
		panic(fmt.Sprintf("invalid move: piece can't move to given square\n%s", attackSet))
	}

	// half-move clock stuff
	switch {
	case b.position[from].Type() == piece.Pawn, b.position[to] != piece.Empty:
		// reset clock
		b.halfMoves = 0
	default:
		b.halfMoves++
	}

	// move piece in bitboard
	b.bitboards[b.position[from]].Unset(from)
	b.bitboards[b.position[from]].Set(to)
	// remove captured piece from bitboard
	b.bitboards[b.position[to]].Unset(to)

	// move piece in 8x8 board
	b.position[to] = b.position[from]
	b.position[from] = piece.Empty

	b.switchTurn()
	b.updateBitboards()
}

func (b *Board) switchTurn() {
	switch b.sideToMove {
	case piece.WhiteColor:
		b.sideToMove = piece.BlackColor
	case piece.BlackColor:
		b.sideToMove = piece.WhiteColor
		b.fullMoves++ // turn completed
	}
}

func (b *Board) updateBitboards() {
	b.friends = bitboard.Empty
	b.enemies = bitboard.Empty

	for p := piece.King + piece.White; p <= piece.Pawn+piece.Black; p++ {
		if p.Color() == b.sideToMove {
			b.friends |= b.bitboards[p]
		} else {
			b.enemies |= b.bitboards[p]
		}
	}
}

func (b *Board) MovesOf(index square.Square) bitboard.Board {
	p := b.position[index]
	if p.Color() != b.sideToMove {
		// other side has no moves
		return bitboard.Empty
	}

	var attackBoard bitboard.Board
	switch p.Type() {
	case piece.King:
		attackBoard = attacks.King(index)
	case piece.Knight:
		attackBoard = attacks.Knight(index)
	default:
		attackBoard = bitboard.Empty
	}

	// can't move to squares occupied by friendly pieces
	return attackBoard &^ b.friends
}
