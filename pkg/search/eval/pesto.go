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

var phaseInc = [piece.TypeN]Eval{0, 0, 1, 1, 2, 4, 0}

type OTSePUE struct {
	mg, eg [piece.ColorN]Eval
	phase  Eval
}

var _ EfficientlyUpdatable = (*OTSePUE)(nil)

func (pesto *OTSePUE) FillSquare(s square.Square, p piece.Piece) {
	pesto.mg[p.Color()] += mgTable[p][s]
	pesto.eg[p.Color()] += egTable[p][s]

	pesto.phase += phaseInc[p.Type()]
}

func (pesto *OTSePUE) ClearSquare(s square.Square, p piece.Piece) {
	pesto.mg[p.Color()] -= mgTable[p][s]
	pesto.eg[p.Color()] -= egTable[p][s]

	pesto.phase -= phaseInc[p.Type()]
}

func (pesto *OTSePUE) Accumulate(stm piece.Color) Eval {
	mgScore := pesto.mg[stm] - pesto.mg[stm.Other()]
	egScore := pesto.eg[stm] - pesto.eg[stm.Other()]

	mgPhase := util.Min(pesto.phase, 24)
	egPhase := 24 - mgPhase

	return (mgScore*mgPhase + egScore*egPhase) / 24
}
