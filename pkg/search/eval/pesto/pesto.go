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

package pesto

import (
	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
	"laptudirm.com/x/mess/pkg/search/eval"
)

//go:generate go run laptudirm.com/x/mess/internal/generator/pesto

// EfficientlyUpdatable (back-acronym of Efficiently Updatable PeSTO) is an efficiently
// updatable PeSTO evaluation function.
type EfficientlyUpdatable struct {
	Board *board.Board

	score [piece.ColorN]Score // middle-game and end-game evaluations of both sides
	phase eval.Eval           // the game phase to lerp between middle and end game

	// pawnN records the number of pawns for each color on each file
	pawnN [piece.ColorN][square.FileN]int
}

// compile time check that OTSePUE implements eval.EfficientlyUpdatable
var _ eval.EfficientlyUpdatable = (*EfficientlyUpdatable)(nil)

// FillSquare adds the given piece to the given square of a chessboard.
func (pesto *EfficientlyUpdatable) FillSquare(s square.Square, p piece.Piece) {
	pType := p.Type()
	color := p.Color()

	// add piece's contribution to the evaluation
	pesto.score[color] += table[p][s]
	pesto.phase += phaseInc[pType]

	if pType == piece.Pawn {
		pesto.pawnN[color][s.File()]++
	}
}

// ClearSquare removes the given piece from the given square.
func (pesto *EfficientlyUpdatable) ClearSquare(s square.Square, p piece.Piece) {
	pType := p.Type()
	color := p.Color()

	// remove piece's contribution from the evaluation
	pesto.score[color] -= table[p][s]
	pesto.phase -= phaseInc[pType]

	if pType == piece.Pawn {
		pesto.pawnN[color][s.File()]--
	}
}

// Accumulate accumulates the efficiently updated variables into the
// evaluation of the position from the perspective of the given side.
func (pesto *EfficientlyUpdatable) Accumulate(stm piece.Color) eval.Eval {
	xtm := stm.Other()

	// create copy of stm and xstm scores
	stmScore := pesto.score[stm]
	xtmScore := pesto.score[xtm]

	// penalty for having stacked pawns
	for file := square.FileA; file <= square.FileH; file++ {
		stmScore -= stackedPawnPenalty[pesto.pawnN[stm][file]]
		xtmScore -= stackedPawnPenalty[pesto.pawnN[xtm][file]]
	}

	// score from side to move's perspective
	score := stmScore - xtmScore

	// linearly interpolate between the end game and middle game
	// evaluations using phase/startposPhase as the contribution
	// of the middle game to the final evaluation
	phase := util.Min(pesto.phase, startposPhase)
	return util.Lerp(score.EG(), score.MG(), phase, startposPhase)
}
