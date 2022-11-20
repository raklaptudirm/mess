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

// This is a generator package used to generate go files containing data
// pertaining to attack bitboards of chess pieces.
package main

import (
	_ "embed"

	"laptudirm.com/x/mess/internal/generator"
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/board/move/attacks/magic"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
)

type attackStruct struct {
	King   [square.N]bitboard.Board
	Knight [square.N]bitboard.Board
	Pawn   [piece.ColorN][square.N]bitboard.Board

	Rook   magic.Table
	Bishop magic.Table
}

//go:embed .gotemplate
var template string

func main() {
	var a attackStruct

	// initialize standard lookup tables for non-sliding pieces
	for s := square.A8; s <= square.H1; s++ {
		// compute attack bitboards for current square
		a.King[s] = kingAttacksFrom(s)
		a.Knight[s] = knightAttacksFrom(s)
		a.Pawn[piece.White][s] = whitePawnAttacksFrom(s)
		a.Pawn[piece.Black][s] = blackPawnAttacksFrom(s)
	}

	// initialize magic lookup tables for sliding pieces
	a.Rook = *magic.NewTable(4096, rook)
	a.Bishop = *magic.NewTable(512, bishop)

	generator.Generate("tables", template, a)
}
