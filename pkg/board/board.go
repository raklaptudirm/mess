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
	"laptudirm.com/x/mess/pkg/castling"
	"laptudirm.com/x/mess/pkg/move"
	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
	"laptudirm.com/x/mess/pkg/zobrist"
)

// Board represents the state of a chessboard at a given position.
type Board struct {
	// position data
	Hash     zobrist.Key
	Position mailbox.Board // 8x8 for fast lookup
	PieceBBs [piece.NType]bitboard.Board
	ColorBBs [piece.NColor]bitboard.Board

	Kings [piece.NColor]square.Square

	SideToMove      piece.Color
	EnPassantTarget square.Square
	CastlingRights  castling.Rights

	CheckN    int
	CheckMask bitboard.Board

	// move counters
	Plys      int
	FullMoves int
	DrawClock int

	// game data
	History [256]Undo
}

type Undo struct {
	Move            move.Move
	CastlingRights  castling.Rights
	CapturedPiece   piece.Piece
	EnPassantTarget square.Square
	DrawClock       int
	Hash            zobrist.Key
}

// String converts a Board into a human readable string.
func (b Board) String() string {
	return fmt.Sprintf("%s\nFen: %s\nKey: %X\n", b.Position, b.FEN(), b.Hash)
}

func (b *Board) Occupied() bitboard.Board {
	return b.ColorBBs[piece.White] | b.ColorBBs[piece.Black]
}

func (b *Board) ClearSquare(s square.Square) {
	p := b.Position[s]

	b.ColorBBs[p.Color()].Unset(s)

	// remove piece from other records
	b.PieceBBs[p.Type()].Unset(s)       // piece bitboard
	b.Position[s] = piece.NoPiece       // mailbox board
	b.Hash ^= zobrist.PieceSquare[p][s] // zobrist hash
}

func (b *Board) FillSquare(s square.Square, p piece.Piece) {
	c := p.Color()
	t := p.Type()

	b.ColorBBs[c].Set(s)

	if t == piece.King {
		b.Kings[c] = s
	}

	b.PieceBBs[t].Set(s)                // piece bitboard
	b.Position[s] = p                   // mailbox board
	b.Hash ^= zobrist.PieceSquare[p][s] // zobrist hash
}

func (b *Board) IsInCheck(c piece.Color) bool {
	return b.IsAttacked(b.Kings[c], c.Other())
}

func (b *Board) IsAttacked(s square.Square, them piece.Color) bool {
	occ := b.Occupied()

	if attacks.Pawn[them.Other()][s]&b.Pawns(them) != bitboard.Empty {
		return true
	}

	if attacks.Knight[s]&b.Knights(them) != bitboard.Empty {
		return true
	}

	if attacks.King[s]&b.King(them) != bitboard.Empty {
		return true
	}

	queens := b.Queens(them)

	if attacks.Bishop(s, occ)&(b.Bishops(them)|queens) != bitboard.Empty {
		return true
	}

	return attacks.Rook(s, occ)&(b.Rooks(them)|queens) != bitboard.Empty
}

func (b *Board) Pawns(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Pawn] & b.ColorBBs[c]
}

func (b *Board) Knights(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Knight] & b.ColorBBs[c]
}

func (b *Board) Bishops(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Bishop] & b.ColorBBs[c]
}

func (b *Board) Rooks(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Rook] & b.ColorBBs[c]
}

func (b *Board) Queens(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Queen] & b.ColorBBs[c]
}

func (b *Board) King(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.King] & b.ColorBBs[c]
}

func (b *Board) CalculateCheckmask() {
	occ := b.Occupied()

	us := b.SideToMove
	them := us.Other()

	b.CheckN = 0
	b.CheckMask = bitboard.Empty

	kingSq := b.Kings[us]

	pawns := b.Pawns(them) & attacks.Pawn[us][kingSq]
	knights := b.Knights(them) & attacks.Knight[kingSq]
	bishops := (b.Bishops(them) & b.Queens(them)) & attacks.Bishop(kingSq, occ)
	rooks := (b.Rooks(them) & b.Queens(them)) & attacks.Rook(kingSq, occ)

	switch {
	case pawns != bitboard.Empty:
		b.CheckMask |= pawns
		b.CheckN++

	case knights != bitboard.Empty:
		b.CheckMask |= knights
		b.CheckN++
	}

	if bishops != bitboard.Empty {
		bishopSq := bishops.FirstOne()
		b.CheckMask |= attacks.Between[kingSq][bishopSq] | bitboard.Squares[bishopSq]
		b.CheckN++
	}

	if b.CheckN < 2 && rooks != bitboard.Empty {
		if b.CheckN == 0 && rooks.Count() > 1 {
			b.CheckN++
		} else {
			rookSq := rooks.FirstOne()
			b.CheckMask |= attacks.Between[kingSq][rookSq] | bitboard.Squares[rookSq]
			b.CheckN++
		}
	}

	if b.CheckN == 0 {
		b.CheckMask = bitboard.Universe
	}
}
