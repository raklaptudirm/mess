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

// Package square declares constants representing every square on a
// chessboard, and related utility functions.
//
// Squares are represented using the algebraic notation.
// https://www.chessprogramming.org/Algebraic_Chess_Notation
// The null square is represented using the "-" symbol.
package square

// NewFromString creates a new instance of a Square from the given identifier.
func NewFromString(id string) Square {
	switch {
	case id == "-":
		return None
	case len(id) != 2:
		panic("new square: invalid square id")
	}

	// ascii code to square index
	return New(FileFrom(string(id[0])), RankFrom(string(id[1])))
}

// New creates a new instance of a Square from the given file and rank.
func New(file File, rank Rank) Square {
	return Square(int(rank)<<3 | int(file))
}

// Square represents a square on a chessboard.
type Square int8

// constants representing various squares.
const (
	None Square = -1

	A8, B8, C8, D8, E8, F8, G8, H8 Square = +0, +1, +2, +3, +4, +5, +6, +7
	A7, B7, C7, D7, E7, F7, G7, H7 Square = +8, +9, 10, 11, 12, 13, 14, 15
	A6, B6, C6, D6, E6, F6, G6, H6 Square = 16, 17, 18, 19, 20, 21, 22, 23
	A5, B5, C5, D5, E5, F5, G5, H5 Square = 24, 25, 26, 27, 28, 29, 30, 31
	A4, B4, C4, D4, E4, F4, G4, H4 Square = 32, 33, 34, 35, 36, 37, 38, 39
	A3, B3, C3, D3, E3, F3, G3, H3 Square = 40, 41, 42, 43, 44, 45, 46, 47
	A2, B2, C2, D2, E2, F2, G2, H2 Square = 48, 49, 50, 51, 52, 53, 54, 55
	A1, B1, C1, D1, E1, F1, G1, H1 Square = 56, 57, 58, 59, 60, 61, 62, 63
)

// Number of squares
const N = 64

// String converts a square into it's algebraic string representation.
func (s Square) String() string {
	if s == None {
		return "-"
	}

	// <file><rank>
	return s.File().String() + s.Rank().String()
}

// File returns the file of the given square.
func (s Square) File() File {
	return File(s) & 7
}

// Rank returns the rank of the given square.
func (s Square) Rank() Rank {
	return Rank(s) >> 3
}

// Diagonal returns the NE-SW diagonal of the given square.
func (s Square) Diagonal() Diagonal {
	return 14 - Diagonal(s.Rank()) - Diagonal(s.File())
}

// AntiDiagonal returns the NW-SE anti-diagonal of the given square.
func (s Square) AntiDiagonal() AntiDiagonal {
	return 7 - AntiDiagonal(s.Rank()) + AntiDiagonal(s.File())
}
