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
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
	"laptudirm.com/x/mess/pkg/search/eval"
	"laptudirm.com/x/mess/pkg/search/pesto"
)

type pestoStruct struct {
	Table       [piece.N][square.N]pesto.Score
	StackedPawn [7]pesto.Score
}

//go:embed .gotemplate
var template string

func main() {
	var pesto pestoStruct

	for i := 2; i < 6; i++ {
		pesto.StackedPawn[i] = pesto.S(
			eval.Eval(15*(i-1)),
			eval.Eval(20*(i-1)),
		)
	}

	// initialize PESTO tables
	for s := square.A8; s < square.N; s++ {
		for p := piece.Pawn; p <= piece.King; p++ {
			white := piece.New(p, piece.White)
			black := piece.New(p, piece.Black)

			pesto.Table[white][s] = pesto.S(
				mgPieceValues[p]+mgPieceTable[p][s],
				egPieceValues[p]+egPieceTable[p][s],
			)

			pesto.Table[black][s] = pesto.S(
				mgPieceValues[p]+mgPieceTable[p][s^56],
				egPieceValues[p]+egPieceTable[p][s^56],
			)
		}
	}

	generator.Generate("tables", template, pesto)
}
