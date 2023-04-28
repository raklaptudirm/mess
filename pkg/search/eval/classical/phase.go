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

package classical

import (
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/search/eval"
)

// game phase increment of each piece type
const (
	pawnPhaseInc   eval.Eval = 0
	knightPhaseInc eval.Eval = 1
	bishopPhaseInc eval.Eval = 1
	rookPhaseInc   eval.Eval = 2
	queenPhaseInc  eval.Eval = 4
)

// phaseInc maps each piece type to it's phase increment.
var phaseInc = [piece.TypeN]eval.Eval{
	piece.Pawn:   pawnPhaseInc,
	piece.Knight: knightPhaseInc,
	piece.Bishop: bishopPhaseInc,
	piece.Rook:   rookPhaseInc,
	piece.Queen:  queenPhaseInc,
}

// MaxPhase is the phase of the starting position.
const MaxPhase = 16*pawnPhaseInc +
	4*knightPhaseInc + 4*bishopPhaseInc +
	4*rookPhaseInc + 2*queenPhaseInc
