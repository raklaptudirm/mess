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

package eval

import (
	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
)

//go:generate go run laptudirm.com/x/mess/internal/generator/pesto

var phaseInc = [piece.TypeN]int{0, 0, 1, 1, 2, 4, 0}

// PeSTO is an evaluation.Func which uses PeSTO evaluation.
// https://www.chessprogramming.org/PeSTO%27s_Evaluation_Function
func PeSTO(b *board.Board) Eval {
	var mg [piece.ColorN]Eval
	var eg [piece.ColorN]Eval

	var gamePhase int

	for s := square.A8; s < square.N; s++ {
		p := b.Position[s]
		if p != piece.NoPiece {
			mg[p.Color()] += mgTable[p][s]
			eg[p.Color()] += egTable[p][s]

			gamePhase += phaseInc[p.Type()]
		}
	}

	// tapered evaluation

	mgScore := mg[b.SideToMove] - mg[b.SideToMove.Other()]
	egScore := eg[b.SideToMove] - eg[b.SideToMove.Other()]

	mgPhase := util.Min(gamePhase, 24)
	egPhase := 24 - mgPhase

	return (mgScore*Eval(mgPhase) + egScore*Eval(egPhase)) / 24
}
