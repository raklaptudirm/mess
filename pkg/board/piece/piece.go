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

// Package piece implements representations of all the chess pieces and
// colors, and related utility functions.
//
// The King, Queen, Rook, Knight, Bishop, and Pawn are represented by the
// K, Q, R, N, B, and P strings respectively, with uppercase for white and
// lower case for black.
//
// The strings w, and b are used for representing the White and Black
// colors respectively.
package piece

// New creates a new Piece with the given type and color.
func New(t Type, c Color) Piece {
	return Piece(c<<colorOffset) | Piece(t)
}

// NewFromString creates an instance of Piece from the given piece id.
func NewFromString(id string) Piece {
	switch id {
	case "K":
		return WhiteKing
	case "Q":
		return WhiteQueen
	case "R":
		return WhiteRook
	case "N":
		return WhiteKnight
	case "B":
		return WhiteBishop
	case "P":
		return WhitePawn
	case "k":
		return BlackKing
	case "q":
		return BlackQueen
	case "r":
		return BlackRook
	case "n":
		return BlackKnight
	case "b":
		return BlackBishop
	case "p":
		return BlackPawn
	default:
		panic("new piece: invalid piece id")
	}
}

// Piece represents a colored chess piece.
// Format: MSB [color 1 bit][type 3 bits] LSB
type Piece uint8

// constants representing colored chess pieces
const (
	NoPiece Piece = 0

	// white pieces
	WhitePawn   Piece = Piece(White)<<3 | Piece(Pawn)
	WhiteKnight Piece = Piece(White)<<3 | Piece(Knight)
	WhiteBishop Piece = Piece(White)<<3 | Piece(Bishop)
	WhiteRook   Piece = Piece(White)<<3 | Piece(Rook)
	WhiteQueen  Piece = Piece(White)<<3 | Piece(Queen)
	WhiteKing   Piece = Piece(White)<<3 | Piece(King)

	// black pieces
	BlackPawn   Piece = Piece(Black)<<3 | Piece(Pawn)
	BlackKnight Piece = Piece(Black)<<3 | Piece(Knight)
	BlackBishop Piece = Piece(Black)<<3 | Piece(Bishop)
	BlackRook   Piece = Piece(Black)<<3 | Piece(Rook)
	BlackQueen  Piece = Piece(Black)<<3 | Piece(Queen)
	BlackKing   Piece = Piece(Black)<<3 | Piece(King)
)

// N is the number of chess piece-color combinations there are. Ideally it
// should be 6x2 = 12, but the number is bloated due to separating the bit
// offsets of piece type and color to make getting them easier.
const N = 16

// constants representing field offsets in Piece
const (
	colorOffset = 3
	typeMask    = (1 << colorOffset) - 1
)

// String converts a Piece into it's string representation. THe pieces are
// represented using their standard alphabets, with white pieces having
// upper case letters and black pieces having lower case ones.
func (p Piece) String() string {
	const pieceToStr = " PNBRQK  pnbrqk"
	return string(pieceToStr[p])
}

// Type returns the piece type of the given Piece.
func (p Piece) Type() Type {
	return Type(p & typeMask)
}

// Color returns the piece color of the given Piece.
func (p Piece) Color() Color {
	return Color(p >> colorOffset)
}

// Is checks if the type of the given Piece matches the given type.
func (p Piece) Is(target Type) bool {
	return p.Type() == target
}

// IsColor checks if the color of the given Piece matches the given Color.
func (p Piece) IsColor(target Color) bool {
	return p.Color() == target
}

// Type represents the type/kind of chess piece.
type Type uint8

// constants representing chess piece types
const (
	NoType Type = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King
)

// TypeN is the number of chess piece types, including NoType.
const TypeN = 7

func (t Type) String() string {
	const typeToStr = " pnbrqk"
	return string(typeToStr[t])
}
