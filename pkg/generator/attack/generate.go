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
	"os"
	"text/template"

	"laptudirm.com/x/mess/pkg/chess/bitboard"
	"laptudirm.com/x/mess/pkg/chess/move/attacks/magic"
	"laptudirm.com/x/mess/pkg/chess/piece"
	"laptudirm.com/x/mess/pkg/chess/square"
)

type attackStruct struct {
	King   [square.N]bitboard.Board
	Knight [square.N]bitboard.Board
	Pawn   [piece.NColor][square.N]bitboard.Board

	Rook   magic.Table
	Bishop magic.Table
}

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
	{
		a.Rook = magic.Table{
			MaxMaskN: 4096, MoveFunc: rook,
		}

		a.Bishop = magic.Table{
			MaxMaskN: 512, MoveFunc: bishop,
		}

		a.Rook.Populate()
		a.Bishop.Populate()

		// the MoveFunc property is unnecessary after table population is
		// complete, so they are excluded from the generated file
		a.Rook.MoveFunc = nil
		a.Bishop.MoveFunc = nil
	}

	f, err := os.Create("../../move/attacks/gen_attacks.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	attackTemplate.Execute(f, a)
}

var attackTemplate = template.Must(template.New("array").Parse(`
// Code generated by go generate; DO NOT EDIT OR COMMIT THIS FILE
// The source code for the generator can be found at generator/attack

package attacks

import (
	"laptudirm.com/x/mess/pkg/chess/bitboard"
	"laptudirm.com/x/mess/pkg/chess/move/attacks/magic"
)

// non-sliding pieces; standard tables
var King = {{ printf "%#v" .King }}
var Knight = {{ printf "%#v" .Knight }}
var Pawn = {{ printf "%#v" .Pawn }}

// sliding pieces; magic tables
// these tables are un-exported since utility functions will be declared
// to make probing these tables simpler in the attacks package
var rookTable = {{ printf "%#v" .Rook }}
var bishopTable = {{ printf "%#v" .Bishop }}
`))
