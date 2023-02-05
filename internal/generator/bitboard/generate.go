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

package main

import (
	_ "embed"

	"laptudirm.com/x/mess/internal/generator"
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/board/square"
)

type bitboardStruct struct {
	Between [square.N][square.N]bitboard.Board
}

//go:embed .gotemplate
var template string

func main() {
	var b bitboardStruct

	// initialize Between
	for s1 := square.A8; s1 <= square.H1; s1++ {
		for s2 := square.A8; s2 <= square.H1; s2++ {
			sqs := bitboard.Square(s1) | bitboard.Square(s2)
			var mask bitboard.Board

			switch {
			case s1.File() == s2.File():
				mask = bitboard.Files[s1.File()]
			case s1.Rank() == s2.Rank():
				mask = bitboard.Ranks[s1.Rank()]
			case s1.Diagonal() == s2.Diagonal():
				mask = bitboard.Diagonals[s1.Diagonal()]
			case s1.AntiDiagonal() == s2.AntiDiagonal():
				mask = bitboard.AntiDiagonals[s1.AntiDiagonal()]
			default:
				// the squares don't have their file, rank, diagonal, or
				// anti-diagonal in common, so path is Empty (default).
				continue
			}

			b.Between[s1][s2] = bitboard.Hyperbola(s1, sqs, mask) & bitboard.Hyperbola(s2, sqs, mask)
		}
	}

	generator.Generate("bitboards", template, b)
}
