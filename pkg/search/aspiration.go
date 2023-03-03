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

// aspirationWindow implements aspiration windows, which are a way to
// reduce the search space in an alpha-beta search. The technique is to
// use a guess of the expected value (usually from the last iteration in
// iterative deepening), and use a window around this as the alpha-beta
// bounds. Because the window is narrower, more beta cutoffs are achieved,
// and the search takes a shorter time. The drawback is that if the true
// score is outside this window, then a costly re-search must be made.
func (search *Context) aspirationWindow(depth int, prevEval eval.Eval) (eval.Eval, move.Variation) {
	// default values for alpha and beta
	alpha := eval.Eval(-eval.Inf)
	beta := eval.Eval(eval.Inf)

	initialDepth := depth

	// aspiration window size
	var windowSize eval.Eval = 50

	// only do aspiration search at greater than depth 5
	if depth >= 5 {
		// reduce search window
		alpha = prevEval - windowSize
		beta = prevEval + windowSize
	}

	for {
		if search.shouldStop() {
			// some search limit has been breached
			// the return value doesn't matter since this search's result
			// will be trashed and the previous iteration's pv will be used
			return 0, move.Variation{}
		}

		var pv move.Variation
		result := search.negamax(0, depth, alpha, beta, &pv)

		switch {
		// result <= alpha: search failed low
		case result <= alpha:
			beta = (alpha + beta) / 2
			alpha = util.Max(alpha-windowSize, -eval.Inf)

			// reset reduced depth
			depth = initialDepth

		// result >= beta: search failed high
		case result >= beta:
			beta = util.Min(beta+windowSize, eval.Inf)

			// unless we are mating, reduce depth
			if util.Abs(result) <= eval.Inf/2 {
				depth--
			}

		// exact score is inside bounds
		default:
			// return exact score
			return result, pv
		}

		// score out of bounds, research needed

		// increase window size
		windowSize += windowSize / 2
	}
}
