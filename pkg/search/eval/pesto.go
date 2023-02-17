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
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
)

//go:generate go run laptudirm.com/x/mess/internal/generator/pesto

// phaseInc is the effect that each piece type has on the game phase.
var phaseInc = [piece.TypeN]Eval{0, 0, 1, 1, 2, 4, 0}

// OTSePUE (back-acronym of Efficiently Updatable PeSTO) is an efficiently
// updatable PeSTO evaluation function.
type OTSePUE struct {
	mg, eg [piece.ColorN]Eval // middle-game and end-game evaluations of both sides
	phase  Eval               // the game phase to lerp between middle and end game
}

// compile time check that OTSePUE implements eval.EfficientlyUpdatable
var _ EfficientlyUpdatable = (*OTSePUE)(nil)

// FillSquare adds the given piece to the given square of a chessboard.
func (pesto *OTSePUE) FillSquare(s square.Square, p piece.Piece) {
	// add the value of the new piece to
	// the middle and end game evaluations
	pesto.mg[p.Color()] += mgTable[p][s]
	pesto.eg[p.Color()] += egTable[p][s]

	// increase phase by the piece's weight
	pesto.phase += phaseInc[p.Type()]
}

// ClearSquare removes the given piece from the given square.
func (pesto *OTSePUE) ClearSquare(s square.Square, p piece.Piece) {
	// remove the value of the new piece to
	// the middle and end game evaluations
	pesto.mg[p.Color()] -= mgTable[p][s]
	pesto.eg[p.Color()] -= egTable[p][s]

	// decrease phase by the piece's weight
	pesto.phase -= phaseInc[p.Type()]
}

// Accumulate accumulates the efficiently updated variables into the
// evaluation of the position from the perspective of the given side.
func (pesto *OTSePUE) Accumulate(stm piece.Color) Eval {
	// find the middle and end game evaluations of the
	// position from the perspective of the given side
	mgScore := pesto.mg[stm] - pesto.mg[stm.Other()]
	egScore := pesto.eg[stm] - pesto.eg[stm.Other()]

	// calculate the effect that effect that the score
	// of each phase will have on the final evaluation
	// where (phase/24)*score is the value that the
	// phase will give to the final evaluation
	mgPhase := util.Min(pesto.phase, 24)
	egPhase := 24 - mgPhase

	// add the effective scores of each game phase to
	// find the final evaluation of the position
	return (mgScore*mgPhase + egScore*egPhase) / 24
}
