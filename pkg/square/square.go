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

// Package square declares constants representing every square on a
// chessboard, and related utility functions.
//
// Squares are represented using the algebraic notation.
// https://www.chessprogramming.org/Algebraic_Chess_Notation
// The null square is represented using the "-" symbol.
package square

import "fmt"

// New creates a new instance of a Square from the given identifier.
func New(id string) Square {
	switch {
	case id == "-":
		return None
	case len(id) != 2:
		panic("new square: invalid square id")
	}

	// ascii code to square index
	return From(fileFrom(string(id[0])), rankFrom(string(id[1])))
}

// From creates a new instance of a Square from the given file and rank.
func From(file File, rank Rank) Square {
	return Square(int(rank*8) + int(file))
}

// Square represents a square on a chessboard.
type Square int

const None Square = -1

// constants representing various squares.
const (
	A8 Square = iota
	B8
	C8
	D8
	E8
	F8
	G8
	H8

	A7
	B7
	C7
	D7
	E7
	F7
	G7
	H7

	A6
	B6
	C6
	D6
	E6
	F6
	G6
	H6

	A5
	B5
	C5
	D5
	E5
	F5
	G5
	H5

	A4
	B4
	C4
	D4
	E4
	F4
	G4
	H4

	A3
	B3
	C3
	D3
	E3
	F3
	G3
	H3

	A2
	B2
	C2
	D2
	E2
	F2
	G2
	H2

	A1
	B1
	C1
	D1
	E1
	F1
	G1
	H1
)

// String converts a square into it's algebraic string representation.
func (s Square) String() string {
	if s == None {
		return "-"
	}

	// <file><rank>
	return fmt.Sprintf("%s%s", s.File(), s.Rank())
}

// File returns the file of the given square.
func (s Square) File() File {
	return File(s % 8)
}

// Rank returns the rank of the given square.
func (s Square) Rank() Rank {
	return Rank(s / 8)
}

func (s Square) Diagonal() Diagonal {
	return 14 - Diagonal(s.Rank()) - Diagonal(s.File())
}

func (s Square) AntiDiagonal() AntiDiagonal {
	return 7 - AntiDiagonal(s.Rank()) + AntiDiagonal(s.File())
}
