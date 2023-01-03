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
	"fmt"
	"time"

	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/board/move"
	"laptudirm.com/x/mess/pkg/search/eval"
)

// iterativeDeepening is the main search function. It implements an iterative
// deepening loop which call's the negamax search function for each iteration.
// It returns the principal variation and it's evaluation.
// https://www.chessprogramming.org/Iterative_Deepening
func (search *Context) iterativeDeepening() (move.Variation, eval.Eval) {
	var score eval.Eval
	var pv move.Variation

	start := time.Now()

	// iterative deepening loop, starting from 1, call negamax for each depth
	// until the depth limit is reached or time runs out. This allows us to
	// search to any depth depending on the allocated time. Previous iterations
	// also populate the transposition table with scores and pv moves which makes
	// iterative deepening to a depth faster that directly searching that depth.
	for search.depth = 1; search.depth <= search.limits.Depth; search.depth++ {

		// the new pv isn't directly stored into the pv variable since it will
		// pollute the correct pv if the next search is incomplete. Instead the
		// old pv is overwritten only if the search is found to be complete.
		var childPV move.Variation
		score = search.negamax(0, search.depth, -eval.Inf, eval.Inf, &childPV)

		if search.stopped {
			// don't use the new pv if search was stopped since the
			// search is probably unfinished

			// search.shouldStop is not used since the new pv is
			// only bad if the search was stopped in the middle
			// of the iteration, and not in here
			break
		}

		// search successfully completed, so update pv
		pv = childPV

		// print some info for the GUI
		searchTime := time.Since(start)
		fmt.Printf(
			"info depth %d score %s nodes %d nps %.f time %d pv %s\n",
			search.depth, score, search.nodes,
			float64(search.nodes)/util.Max(0.001, searchTime.Seconds()),
			searchTime.Milliseconds(), pv,
		)
	}

	return pv, score
}
