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

// Package mailbox implements a 8x8 mailbox chessboard representation.
// https://www.chessprogramming.org/8x8_Board
package mailbox

import (
	"fmt"

	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
)

// Board represents a 8x8 chessboard consisting of pieces.
type Board [square.N]piece.Piece

// String converts a Board into it's human readable string representation.
func (b Board) String() string {
	// leading divider
	s := "+---+---+---+---+---+---+---+---+\n"

	for rank := square.Rank8; rank <= square.Rank1; rank++ {
		s += "| "

		for file := square.FileA; file <= square.FileH; file++ {
			square := square.New(file, rank)
			s += b[square].String() + " | "
		}

		s += rank.String()
		s += "\n+---+---+---+---+---+---+---+---+\n"
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
		currSquare := square.Square(i)

		if p == piece.NoPiece {
			// increase empty square count
			empty++
		} else {

			if empty > 0 {
				fen += fmt.Sprint(empty)
				empty = 0
			}

			fen += p.String()
		}

		// rank separators after last file
		if currSquare.File() == square.FileH {
			if empty > 0 {
				fen += fmt.Sprint(empty)
				empty = 0
			}

			// no trailing separator after last rank
			if currSquare.Rank() != square.Rank1 {
				fen += "/"
			}
		}
	}

	return fen
}
