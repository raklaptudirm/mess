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

// Package package board contains the main board representation used by the
// mess chess engine. It may be used as a library when developing other
// engines. It also contains various sub-packages related to the board
// representation.
package board

import (
	"fmt"

	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/board/mailbox"
	"laptudirm.com/x/mess/pkg/board/move"
	"laptudirm.com/x/mess/pkg/board/move/attacks"
	"laptudirm.com/x/mess/pkg/board/move/castling"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
	"laptudirm.com/x/mess/pkg/board/zobrist"
)

// Board represents the state of a chessboard at a given position. It
// contains two representations of the chess board: a 8x8 mailbox which is
// used to easily look up what piece occupies a given square, and a
// bitboard representation, used for various bitwise calculations like
// calculating the attack sets of pieces.
//
// Board contains the additional state information of a chessboard like the
// half-move clock, en passant square, number of plys, etc.
//
// Various pre-calculated utility information like check masks and pin
// masks are stored in Board to prevent the need for expensive calculation.
type Board struct {
	// main position data:
	// these are the basic position data for a chessboard

	// zobrist hash of current position
	Hash zobrist.Key

	// 8x8 mailbox board representation
	Position mailbox.Board

	// bitboard board representation
	PieceBBs [piece.TypeN]bitboard.Board
	ColorBBs [piece.ColorN]bitboard.Board

	// other necessary information
	SideToMove      piece.Color
	EnPassantTarget square.Square
	CastlingRights  castling.Rights

	// move counters
	Plys      int
	FullMoves int
	DrawClock int

	// game history
	History [move.MaxN]BoardState
}

// BoardState contains the irreversible position data of a given board
// state. This is used to rollback to a previous position in UnmakeMove.
type BoardState struct {
	// move information
	Move          move.Move   // move made on this BoardState
	CapturedPiece piece.Piece // piece captured by playing Move

	// irreversible information
	CastlingRights  castling.Rights
	EnPassantTarget square.Square
	DrawClock       int

	// zobrist key is reversible but is stored for repetition detection
	Hash zobrist.Key
}

// String converts a Board into a human readable string.
func (b Board) String() string {
	return fmt.Sprintf("%s\nFen: %s\nKey: %X\n", b.Position, b.FEN(), b.Hash)
}

// IsDraw checks if the given position is a draw either by the 50 move rule
// or by a repetition. Threefold repetition is not calculated as it is just
// simpler to evaluate any repetition as a draw.
func (b *Board) IsDraw() bool {
	return b.DrawClock >= 100 || b.IsThreefoldRepetition()
}

// IsRepetition checks if the current position has occurred in the game
// before. This is done by probing the game history till the last
// irreversible move, which is pawn pushes or a capture.
func (b *Board) IsRepetition() bool {
	// probe till game start or last irreversible move, whichever is closer
	depth := util.Max(0, b.Plys-b.DrawClock)

	for i := b.Plys - 2; i >= depth; i -= 2 {
		if b.History[i].Hash == b.Hash {
			return true
		}
	}

	return false
}

func (b *Board) IsThreefoldRepetition() bool {
	// probe till game start or last irreversible move, whichever is closer
	depth := util.Max(0, b.Plys-b.DrawClock)

	repetitions := 1 // current position is a repetition
	for i := b.Plys - 2; i >= depth; i -= 2 {
		if b.History[i].Hash == b.Hash {
			repetitions++
			if repetitions >= 3 {
				return true
			}
		}
	}

	return false
}

// ClearSquare removes the piece occupying the given square and updates the
// dependent position information accordingly.
func (b *Board) ClearSquare(s square.Square) {
	p := b.Position[s]

	if p == piece.NoPiece {
		return
	}

	b.ColorBBs[p.Color()].Unset(s)

	// remove piece from other records
	b.PieceBBs[p.Type()].Unset(s)       // piece bitboard
	b.Position[s] = piece.NoPiece       // mailbox board
	b.Hash ^= zobrist.PieceSquare[p][s] // zobrist hash
}

// FillSquare fills the given square with the given piece. Callers should
// make sure that the provided square is unoccupied, otherwise the
// incrementally updating the position will give wrong results.
func (b *Board) FillSquare(s square.Square, p piece.Piece) {
	c := p.Color()
	t := p.Type()

	b.ColorBBs[c].Set(s)

	b.PieceBBs[t].Set(s)                // piece bitboard
	b.Position[s] = p                   // mailbox board
	b.Hash ^= zobrist.PieceSquare[p][s] // zobrist hash
}

// IsInCheck checks if the side with the given color is in check.
func (b *Board) IsInCheck(c piece.Color) bool {
	return b.IsAttacked(b.KingBB(c).FirstOne(), c.Other())
}

// IsAttacked checks if the given squares is attacked by pieces of the
// given color.
func (b *Board) IsAttacked(s square.Square, them piece.Color) bool {
	if attacks.Pawn[them.Other()][s]&b.PawnsBB(them) != bitboard.Empty {
		return true
	}

	if attacks.Knight[s]&b.KnightsBB(them) != bitboard.Empty {
		return true
	}

	if attacks.King[s]&b.KingBB(them) != bitboard.Empty {
		return true
	}

	blockers := b.ColorBBs[piece.White] | b.ColorBBs[piece.Black]
	queens := b.QueensBB(them)

	if attacks.Bishop(s, blockers)&(b.BishopsBB(them)|queens) != bitboard.Empty {
		return true
	}

	return attacks.Rook(s, blockers)&(b.RooksBB(them)|queens) != bitboard.Empty
}

// PawnsBB returns a bitboard of all the pawns of the given color.
func (b *Board) PawnsBB(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Pawn] & b.ColorBBs[c]
}

// KnightsBB returns a bitboard of all the knights of the given color.
func (b *Board) KnightsBB(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Knight] & b.ColorBBs[c]
}

// BishopsBB returns a bitboard of all the bishops of the given color.
func (b *Board) BishopsBB(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Bishop] & b.ColorBBs[c]
}

// RooksBB returns a bitboard of all the rooks of the given color.
func (b *Board) RooksBB(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Rook] & b.ColorBBs[c]
}

// QueensBB returns a bitboard of all the queens of the given color.
func (b *Board) QueensBB(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Queen] & b.ColorBBs[c]
}

// KingBB returns a bitboard containing the king of the given color.
func (b *Board) KingBB(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.King] & b.ColorBBs[c]
}
