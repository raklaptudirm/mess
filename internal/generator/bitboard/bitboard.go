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

package main

import (
	"os"
	"text/template"

	"laptudirm.com/x/mess/pkg/chess/bitboard"
	"laptudirm.com/x/mess/pkg/chess/square"
)

type bitboardStruct struct {
	Squares [square.N]bitboard.Board
	Between [square.N][square.N]bitboard.Board
}

func main() {
	var b bitboardStruct

	// initialize Squares
	for s := square.A8; s <= square.H1; s++ {
		b.Squares[s] = 1 << s
	}

	// initialize Between
	for s1 := square.A8; s1 <= square.H1; s1++ {
		for s2 := square.A8; s2 <= square.H1; s2++ {
			sqs := b.Squares[s1] | b.Squares[s2]
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

	f, err := os.Create("gen_bitboard.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	bitboardTemplate.Execute(f, b)
}

var bitboardTemplate = template.Must(template.New("bitboard").Parse(`
// Code generated by go generate; DO NOT EDIT OR COMMIT THIS FILE
// The source code for the generator can be found at generator/attack

package bitboard

// Between contains bitboards which have the path between two squares set.
// The definition of path is only valid for squares which lie on the same
// file, rank, diagonal, or anti-diagonal. For all other square
// combinations, the path is Empty.
var Between = [64][64]Board{
	{{ range .Between }}
		{
			{{ range . }}
				{{ printf "%d" . }},
			{{ end }}
		},
	{{ end }}
}

// Squares contains bitboards for each square where only that square's bit
// is set. Squares[square] is equivalent to Board(1 << square).
var Squares = [64]Board{
	{{ range .Squares }}
		{{ printf "%d" .}},
	{{ end }}
}
`))