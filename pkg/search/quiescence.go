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
	"laptudirm.com/x/mess/pkg/search/tt"
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

	case search.board.DrawClock >= 100,
		search.board.IsInsufficientMaterial(),
		search.board.IsRepetition():
		return search.draw()

	case plys >= MaxDepth:
		return search.score()
	}

	// check for transposition table hits
	if entry, hit := search.tt.Probe(search.board.Hash); hit {
		search.stats.TTHits++

		// check if the tt entry can be used to exit the search early
		// on this node. If we have an exact value, we can safely
		// return it. If we have a new upper bound or lower bound,
		// check if it causes a beta cutoff.
		switch value := entry.Value.Eval(plys); {
		case entry.Type == tt.ExactEntry, // exact score
			entry.Type == tt.LowerBound && value >= beta,  // fail high
			entry.Type == tt.UpperBound && alpha >= value: // fail high
			// exit search early cause we have an exact
			// score or a beta cutoff from the tt entry
			return value
		}
	}

	bestScore := search.score() // standing pat
	if bestScore >= beta {
		return bestScore // fail high
	}

	alpha = util.Max(alpha, bestScore)

	// generate tactical (captures and promotions) moves only
	moves := search.board.GenerateMoves(true)

	// move ordering
	list := move.ScoreMoves(moves, eval.OfMove(eval.ModeEvalInfo{
		Board:   &search.board.Position,
		Killers: search.killers[plys],
	}))

	for i := 0; i < list.Length; i++ {
		m := list.PickMove(i)

		// node amount updates are done here to prevent duplicates
		// when quiescence search is called from the negamax function.
		// In other words, node amount updates for a quiescence search
		// is done by the caller function, which in this case is the
		// quiescence search itself.
		search.stats.Nodes++

		search.board.MakeMove(m)
		score := -search.quiescence(plys+1, -beta, -alpha)
		search.board.UnmakeMove()

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

	if !search.stopped {
		// update transposition table
		search.tt.Store(tt.Entry{
			Hash:  search.board.Hash,
			Value: tt.EvalFrom(bestScore, plys),
			Depth: 0,
			Type:  tt.ExactEntry,
		})
	}

	return bestScore
}
