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

// Package castling provides various types and definitions which are useful
// when dealing with castling moves in a board representation.
package castling

import "laptudirm.com/x/mess/pkg/chess/square"

// Rights represents the current castling rights of the position.
// [Black Queen-side][Black King-side][White Queen-side][White King-side]
type Rights byte

// NewRights creates a new castling.Rights from the given string. It
// checks if the identifier for each possible castling is in the string
// in the proper order.
//
//     White King-side:  K
//     White Queen-side: Q
//     Black King-side:  k
//     Black Queen-side: q
//
// The string "-" represents castling.NoCasl.
func NewRights(r string) Rights {
	var rights Rights

	if r == "-" {
		return NoCasl
	}

	if r != "" && r[0] == 'K' {
		r = r[1:]
		rights |= WhiteK
	}

	if r != "" && r[0] == 'Q' {
		r = r[1:]
		rights |= WhiteQ
	}

	if r != "" && r[0] == 'k' {
		r = r[1:]
		rights |= BlackK
	}

	if r != "" && r[0] == 'q' {
		rights |= BlackQ
	}

	return rights
}

// Constants representing various castling rights.
const (
	WhiteK Rights = 1 << 0 // white king-side
	WhiteQ Rights = 1 << 1 // white queen-side
	BlackK Rights = 1 << 2 // black king-side
	BlackQ Rights = 1 << 3 // black queen-side

	NoCasl Rights = 0 // no castling possible

	WhiteA Rights = WhiteK | WhiteQ // only white can castle
	BlackA Rights = BlackK | BlackQ // only black can castle

	Kingside  Rights = WhiteK | BlackK // only king-side castling
	Queenside Rights = WhiteQ | BlackQ // only queen-side castling

	All Rights = WhiteA | BlackA // all castling possible
)

// N is the number of possible unique castling rights.
const N = 1 << 4 // 4 possible castling sides

// RightUpdates is a map of each chessboard square to the rights that
// need to be removed if a piece moves from or to that square. For example,
// if a piece moves from or two from the square A1, either the white rook
// has moved or it has been captured, so white can no longer castle
// queen-side. Squares which are not occupied by a king or a rook do not
// effect the castling rights. Squares occupied by a rook remove the
// castling rights of the rook's side and color. Squares occupied by a king
// remove it's color's castling rights.
var RightUpdates = [square.N]Rights{
	BlackQ, NoCasl, NoCasl, NoCasl, BlackA, NoCasl, NoCasl, BlackK,
	NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl,
	NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl,
	NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl,
	NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl,
	NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl,
	NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl,
	WhiteQ, NoCasl, NoCasl, NoCasl, WhiteA, NoCasl, NoCasl, WhiteK,
}

// String converts the given castling.Rights to a readable string.
func (c Rights) String() string {
	var str string

	if c&WhiteK != 0 {
		str += "K"
	}

	if c&WhiteQ != 0 {
		str += "Q"
	}

	if c&BlackK != 0 {
		str += "k"
	}

	if c&BlackQ != 0 {
		str += "q"
	}

	if str == "" {
		str = "-"
	}

	return str
}
