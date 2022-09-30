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

// NewColor creates an instance of color from the given id.
func NewColor(id string) Color {
	switch id {
	case "w":
		return White
	case "b":
		return Black
	default:
		panic("new color: invalid color id")
	}
}

// Color represents the color of a Piece
type Color int

// various piece colors
const (
	White Color = iota
	Black

	NColor = 2
)

func (c Color) Other() Color {
	return c ^ Black
}

// String converts a Color to it's string representation.
func (c Color) String() string {
	switch c {
	case Black:
		return "b"
	case White:
		return "w"
	default:
		panic("new color: invalid color id")
	}
}

func New(t Type, c Color) Piece {
	return Piece(c << 3) + Piece(t)
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

type Type int

// various chess pieces
const (
	NoType Type = iota
	Pawn
	Knight
	Bishop
	Rook
	Queen
	King

	NType = 7
)

func (t Type) String() string {
	return Piece(t | 8).String()
}

// Piece represents a chess piece.
type Piece int

const (
	NoPiece Piece = 0

	WhitePawn   Piece = Piece(Pawn)
	WhiteKnight Piece = Piece(Pawn) + 1
	WhiteBishop Piece = Piece(Pawn) + 2
	WhiteRook   Piece = Piece(Pawn) + 3
	WhiteQueen  Piece = Piece(Pawn) + 4
	WhiteKing   Piece = Piece(Pawn) + 5

	BlackPawn   Piece = Piece(Pawn) + 8
	BlackKnight Piece = Piece(Pawn) + 9
	BlackBishop Piece = Piece(Pawn) + 10
	BlackRook   Piece = Piece(Pawn) + 11
	BlackQueen  Piece = Piece(Pawn) + 12
	BlackKing   Piece = Piece(Pawn) + 13

	N = 16
)

var Promotions = []Type{
	Queen, Rook, Bishop, Knight,
}

// String converts a Piece into it's string representation.
func (p Piece) String() string {
	pieces := [...]string{
		NoPiece:     " ",
		WhitePawn:   "P",
		WhiteKnight: "N",
		WhiteBishop: "B",
		WhiteRook:   "R",
		WhiteQueen:  "Q",
		WhiteKing:   "K",
		BlackPawn:   "p",
		BlackKnight: "n",
		BlackBishop: "b",
		BlackRook:   "r",
		BlackQueen:  "q",
		BlackKing:   "k",
	}

	return pieces[p]
}

// Type returns the piece type of the given Piece.
func (p Piece) Type() Type {
	switch {
	case p == NoPiece:
		return NoType
	default:
		return Type(p & 7)
	}
}

// Color returns the piece color of the given Piece.
func (p Piece) Color() Color {
	if p == NoPiece {
		panic("color of piece: can't find color of NoPiece")
	}

	return Color(p >> 3)
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
