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
	WhitePawn   Piece = Piece(Pawn)
	WhiteKnight Piece = Piece(Pawn) + 1
	WhiteBishop Piece = Piece(Pawn) + 2
	WhiteRook   Piece = Piece(Pawn) + 3
	WhiteQueen  Piece = Piece(Pawn) + 4
	WhiteKing   Piece = Piece(Pawn) + 5

	// black pieces
	BlackPawn   Piece = Piece(Pawn) + 8
	BlackKnight Piece = Piece(Pawn) + 9
	BlackBishop Piece = Piece(Pawn) + 10
	BlackRook   Piece = Piece(Pawn) + 11
	BlackQueen  Piece = Piece(Pawn) + 12
	BlackKing   Piece = Piece(Pawn) + 13
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
	switch {
	case p == NoPiece:
		return NoType
	default:
		return Type(p & typeMask)
	}
}

// Color returns the piece color of the given Piece.
func (p Piece) Color() Color {
	if p == NoPiece {
		panic("color of piece: can't find color of NoPiece")
	}

	return Color(p >> colorOffset)
}

// Is checks if the type of the given Piece matches the given type.
func (p Piece) Is(target Type) bool {
	t := p.Type()
	return t == target
}

// IsColor checks if the color of the given Piece matches the given Color.
func (p Piece) IsColor(target Color) bool {
	c := p.Color()
	return c == target
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
