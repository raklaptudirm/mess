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

// phaseInc is the effect that each piece type has on the game phase.
var phaseInc = [piece.TypeN]eval.Eval{0, 0, 1, 1, 2, 4, 0}

// OTSePUE (back-acronym of Efficiently Updatable PeSTO) is an efficiently
// updatable PeSTO evaluation function.
type OTSePUE struct {
	Board *board.Board

	score [piece.ColorN]Score // middle-game and end-game evaluations of both sides
	phase eval.Eval           // the game phase to lerp between middle and end game

	PawnN [piece.ColorN][square.FileN]int
}

// compile time check that OTSePUE implements eval.EfficientlyUpdatable
var _ eval.EfficientlyUpdatable = (*OTSePUE)(nil)

// FillSquare adds the given piece to the given square of a chessboard.
func (pesto *OTSePUE) FillSquare(s square.Square, p piece.Piece) {
	pType := p.Type()
	color := p.Color()

	// add the value of the new piece to
	// the middle and end game evaluations
	pesto.score[color] += table[p][s]

	// increase phase by the piece's weight
	pesto.phase += phaseInc[pType]

	if pType == piece.Pawn {
		pesto.PawnN[color][s.File()]++
	}
}

// ClearSquare removes the given piece from the given square.
func (pesto *OTSePUE) ClearSquare(s square.Square, p piece.Piece) {
	pType := p.Type()
	color := p.Color()

	// remove the value of the new piece to
	// the middle and end game evaluations
	pesto.score[color] -= table[p][s]

	// decrease phase by the piece's weight
	pesto.phase -= phaseInc[pType]

	if pType == piece.Pawn {
		pesto.PawnN[color][s.File()]--
	}
}

// Accumulate accumulates the efficiently updated variables into the
// evaluation of the position from the perspective of the given side.
func (pesto *OTSePUE) Accumulate(stm piece.Color) eval.Eval {
	xstm := stm.Other()

	score := pesto.score[stm] - pesto.score[xstm]

	// stacked pawn penalties
	for file := square.FileA; file <= square.FileH; file++ {
		score -= stackedPawnPenalty[pesto.PawnN[stm][file]]
		score += stackedPawnPenalty[pesto.PawnN[xstm][file]]
	}

	// calculate the effect that effect that the score
	// of each phase will have on the final evaluation
	// where (phase/24)*score is the value that the
	// phase will give to the final evaluation
	mgPhase := util.Min(pesto.phase, 24)
	egPhase := 24 - mgPhase

	// add the effective scores of each game phase to
	// find the final evaluation of the position
	return (score.MG()*mgPhase + score.EG()*egPhase) / 24
}

func S(mg, eg eval.Eval) Score {
	return Score(uint64(eg)<<32) + Score(mg)
}

type Score int64

func (score Score) MG() eval.Eval {
	return eval.Eval(int32(uint32(uint64(score))))
}

func (score Score) EG() eval.Eval {
	return eval.Eval(int32(uint32(uint64(score+(1<<32)) >> 32)))
}
