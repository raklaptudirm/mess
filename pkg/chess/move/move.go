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

// Package move declares types and constants pertaining to chess moves.
package move

import (
	"laptudirm.com/x/mess/pkg/chess/piece"
	"laptudirm.com/x/mess/pkg/chess/square"
)

// Move represents a chess move. It contains various metadata regarding
// the move including the source and target squares, the moving piece,
// the promoted piece and wether the move is a capture.
//
// Format: MSB -> LSB
// [20 isCapture bool 20] \
// [19 toPiece piece.Piece 16][15 fromPiece piece.Piece 12] \
// [11 target square.Square 6][05 source square.Square  00]
type Move uint

// MaxN is the maximum number of plys in a chess game.
const MaxN = 1024

// Null Move represents a "do nothing" move on the chessboard. It is
// represented by "0000", and is useful for returning errors and pruning.
const Null Move = 0

const (
	// bit offsets of infos
	sourceOffset    = 0
	targetOffset    = 6
	fromPieceOffset = 12
	toPieceOffset   = 16

	// bit masks of infos
	sourceMask    = 0x3f
	targetMask    = 0xfc0
	fromPieceMask = 0xf000
	toPieceMask   = 0xf0000
	captureMask   = 0x100000
)

// New creates a new Move value which is populated with the provided data.
func New(source, target square.Square, from piece.Piece, isCapture bool) Move {
	m := Move(source) << sourceOffset
	m |= Move(target) << targetOffset
	m |= Move(from) << fromPieceOffset
	m |= Move(from) << toPieceOffset
	if isCapture {
		m |= captureMask
	}
	return m
}

// String converts a move to it's long algebraic notation form.
// For example "e2e4", "e1g1"(castling), "d7d8q"(promotion), "0000"(null).
func (m Move) String() string {
	// null move is a special case
	if m == Null {
		return "0000"
	}

	s := m.Source().String() + m.Target().String()

	// add promotion indicator
	if m.IsPromotion() {
		s += m.ToPiece().String()
	}

	return s
}

// SetPromotion sets the promotion field of the move to the given piece.
func (m Move) SetPromotion(p piece.Piece) Move {
	m &^= toPieceMask
	m |= Move(p) << toPieceOffset
	return m
}

// Source returns the source square of the move.
func (m Move) Source() square.Square {
	return square.Square((m & sourceMask) >> sourceOffset)
}

// Target returns the target square of the move.
func (m Move) Target() square.Square {
	return square.Square((m & targetMask) >> targetOffset)
}

// FromPiece returns the piece that is being moved.
func (m Move) FromPiece() piece.Piece {
	return piece.Piece((m & fromPieceMask) >> fromPieceOffset)
}

// ToPiece returns the piece after moving. This is the same as FromPiece
// for normal moves, and is only useful in promotions, where it returns
// the promoted piece.
func (m Move) ToPiece() piece.Piece {
	return piece.Piece((m & toPieceMask) >> toPieceOffset)
}

// IsCapture checks whether the move is a capture.
func (m Move) IsCapture() bool {
	return m&captureMask != 0
}

// IsPromotion checks if the move is a promotion.
func (m Move) IsPromotion() bool {
	return m.FromPiece() != m.ToPiece()
}

// IsEnPassant checks if the move is en passant given the target square.
func (m Move) IsEnPassant(ep square.Square) bool {
	return m.Target() == ep && m.FromPiece().Type() == piece.Pawn
}

// IsQuiet checks if the move is a quiet move. A quiet move is a move
// which does not create huge material differences when played, unlike
// captures and promotions.
func (m Move) IsQuiet() bool {
	return !m.IsCapture() && !m.IsPromotion()
}

// IsReversible checks if the move is reversible. A move is termed as
// reversible if it is possible to "undo" the move, like moving a knight
// back. Captures and pawn moves are not reversible.
func (m Move) IsReversible() bool {
	return !m.IsCapture() && m.FromPiece().Type() != piece.Pawn
}
