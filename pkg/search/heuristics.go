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

package search

import (
	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/board/move"
	"laptudirm.com/x/mess/pkg/search/eval"
)

// storeKiller tries to store the given move from the given depth as one
// of the two killer moves.
func (search *Context) storeKiller(plys int, killer move.Move) {
	if !killer.IsCapture() && killer != search.killers[plys][0] {
		// different move in killer 1
		// move it to killer 2 position
		search.killers[plys][1] = search.killers[plys][0]
		search.killers[plys][0] = killer // new killer 1
	}
}

// updateHistory updates the history score of the given move with the given
// bonus. It also verifies that the move is a quiet move.
func (search *Context) updateHistory(m move.Move, bonus eval.Move) {
	if !m.IsCapture() {
		entry := &search.history[search.board.SideToMove][m.Source()][m.Target()]
		*entry += bonus - *entry*util.Abs(bonus)/32768
	}
}

// depthBonus returns the the history bonus for a particular depth.
func depthBonus(depth int) eval.Move {
	return eval.Move(util.Min(2000, depth*155))
}

// seeMargins returns the see pruning thresholds for the given depth.
func seeMargins(depth int) (quiet, noisy eval.Eval) {
	quiet = eval.Eval(-64 * depth)
	noisy = eval.Eval(-19 * depth * depth)
	return quiet, noisy
}
