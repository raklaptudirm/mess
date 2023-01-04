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

// quiescence search is a type of limited search which only evaluates 'quiet'
// positions, i.e. positions with no tactical moves like captures or promotions.
// This search is needed to avoid the horizon effect.
// https://www.chessprogramming.org/Quiescence_Search
func (search *Context) quiescence(plys int, alpha, beta eval.Eval) eval.Eval {
	// quick exit clauses
	switch {
	case search.shouldStop():
		return 0 // return value doesn't matter

	case search.Board.IsDraw():
		return search.draw()

	case plys >= MaxDepth:
		return search.score()
	}

	bestScore := search.score() // standing pat
	if bestScore >= beta {
		return bestScore // fail high
	}

	alpha = util.Max(alpha, bestScore)

	// generate tactical (captures and promotions) moves only
	moves := search.Board.GenerateMoves(true)

	// move ordering
	list := move.ScoreMoves(moves, eval.OfMove(search.Board, move.Null))
	for i := 0; i < list.Length; i++ {
		m := list.PickMove(i)

		// node amount updates are done here to prevent duplicates
		// when quiescence search is called from the negamax function.
		// In other words, node amount updates for a quiescence search
		// is done by the caller function, which in this case is the
		// quiescence search itself.
		search.stats.Nodes++

		search.Board.MakeMove(m)
		score := -search.quiescence(plys+1, -beta, -alpha)
		search.Board.UnmakeMove()

		// update score and bounds
		if score > bestScore {
			// better move found
			bestScore = score

			// check if move is new pv move
			if score > alpha {
				// new pv so alpha increases
				alpha = score

				if alpha >= beta {
					break // fail high
				}
			}
		}
	}

	return bestScore
}
