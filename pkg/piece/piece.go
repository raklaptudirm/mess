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
		return WhiteColor
	case "b":
		return BlackColor
	default:
		panic("new color: invalid color id")
	}
}

// Color represents the color of a Piece
type Color int

// various piece colors
const (
	EmptyColor Color = iota
	WhiteColor
	BlackColor
)

// String converts a Color to it's string representation.
func (c Color) String() string {
	switch c {
	case BlackColor:
		return "b"
	case WhiteColor:
		return "w"
	case EmptyColor:
		return "e"
	default:
		panic("new color: invalid color id")
	}
}

// New creates an instance of Piece from the given piece id.
func New(id string) Piece {
	switch id {
	case "K":
		return White + King
	case "Q":
		return White + Queen
	case "R":
		return White + Rook
	case "N":
		return White + Knight
	case "B":
		return White + Bishop
	case "P":
		return White + Pawn
	case "k":
		return Black + King
	case "q":
		return Black + Queen
	case "r":
		return Black + Rook
	case "n":
		return Black + Knight
	case "b":
		return Black + Bishop
	case "p":
		return Black + Pawn
	default:
		panic("new piece: invalid piece id")
	}
}

// Piece represents a chess piece.
type Piece int

// various chess pieces
const (
	Empty Piece = iota
	King
	Queen
	Rook
	Knight
	Bishop
	Pawn
)

// Number of pieces
const N = 7

var Promotions = []Piece{
	Queen, Rook, Bishop, Knight,
}

// colors of chess pieces for easy creation
// for example: piece.Black + piece.King
const (
	White Piece = 0
	Black Piece = 6
)

// String converts a Piece into it's string representation.
func (p Piece) String() string {
	pieces := [...]string{
		Empty:          " ",
		White + King:   "K",
		White + Queen:  "Q",
		White + Rook:   "R",
		White + Knight: "N",
		White + Bishop: "B",
		White + Pawn:   "P",
		Black + King:   "k",
		Black + Queen:  "q",
		Black + Rook:   "r",
		Black + Knight: "n",
		Black + Bishop: "b",
		Black + Pawn:   "p",
	}

	return pieces[p]
}

// Type returns the piece type of the given Piece.
func (p Piece) Type() Piece {
	// convert black pieces to white
	if p > Pawn {
		p -= Pawn
	}

	return p
}

// Color returns the piece color of the given Piece.
func (p Piece) Color() Color {
	switch {
	case p == Empty:
		return EmptyColor
	case p > Pawn:
		return BlackColor
	default:
		return WhiteColor
	}
}

// Is checks if the type of the given Piece matches the given type.
func (p Piece) Is(target Piece) bool {
	t := p.Type()
	return t == target
}

// IsColor checks if the color of the given Piece matches the given Color.
func (p Piece) IsColor(target Color) bool {
	c := p.Color()
	return c == target
}
