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
	"laptudirm.com/x/mess/pkg/board/move"
	"laptudirm.com/x/mess/pkg/search/eval"
)

// iterativeDeepening is the main search function. It implements an iterative
// deepening loop which call's the negamax search function for each iteration.
// It returns the principal variation and it's evaluation.
// https://www.chessprogramming.org/Iterative_Deepening
func (search *Context) iterativeDeepening() (move.Variation, eval.Eval) {

	// iterative deepening loop, starting from 1, call negamax for each depth
	// until the depth limit is reached or time runs out. This allows us to
	// search to any depth depending on the allocated time. Previous iterations
	// also populate the transposition table with scores and pv moves which makes
	// iterative deepening to a depth faster that directly searching that depth.
	for search.stats.Depth = 1; search.stats.Depth <= search.limits.Depth; search.stats.Depth++ {

		// the new pv isn't directly stored into the pv variable since it will
		// pollute the correct pv if the next search is incomplete. Instead the
		// old pv is overwritten only if the search is found to be complete.
		score, pv := search.aspirationWindow(search.stats.Depth, search.pvScore)

		if search.stopped {
			// don't use the new pv if search was stopped since the
			// search is probably unfinished

			// search.shouldStop is not used since the new pv is
			// only bad if the search was stopped in the middle
			// of the iteration, and not in here
			break
		}

		// search successfully completed, so update pv
		search.pv = pv
		search.pvScore = score

		// print some info for the GUI
		search.reporter(search.GenerateReport())

		if search.time.OptimisticExpired() {
			break
		}
	}

	if search.stats.Depth < search.limits.Depth && search.limits.Infinite {
		for search.limits.Infinite && !search.shouldStop() {
			// if in infinite mode, wait for stop
		}
	}

	return search.pv, search.pvScore
}
