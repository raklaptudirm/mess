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

// Package mailbox implements a 8x8 mailbox chessboard representation.
// https://www.chessprogramming.org/8x8_Board
package mailbox

import (
	"fmt"

	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
)

// Board represents a 8x8 chessboard consisting of pieces.
type Board [8 * 8]piece.Piece

// String converts a Board into it's human readable string representation.
func (b Board) String() string {

	s := "+---+---+---+---+---+---+---+---+\n"

	for rank := 0; rank < 8; rank++ {
		s += "| "

		for file := 0; file < 8; file++ {
			square := square.Square(rank*8 + file)
			s += b[square].String() + " | "
		}

		s += fmt.Sprintln(8 - rank)
		s += "+---+---+---+---+---+---+---+---+\n"
	}

	s += "  a   b   c   d   e   f   g   h\n"
	return s
}

// FEN generates the position part of a fen string representing the current
// Board position. It can be used together with other information about the
// position to generate a complete fen string.
func (b *Board) FEN() string {
	var fen string

	empty := 0
	for i, p := range b {
		if p == piece.Empty {
			// increase empty square count
			empty++
		} else {

			if empty > 0 {
				fen += fmt.Sprint(empty)
				empty = 0
			}

			fen += p.String()
		}

		// rank separators
		if (i+1)%8 == 0 {
			if empty > 0 {
				fen += fmt.Sprint(empty)
				empty = 0
			}

			// no trailing separator
			if i < 63 {
				fen += "/"
			}
		}
	}

	return fen
}
