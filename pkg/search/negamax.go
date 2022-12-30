// Copyright © 2022 Rak Laptudirm <rak@laptudirm.com>
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

// negamax is a simplified version of the minmax searching algorithm, which
// uses a single function for both the maximizing and minimizing players.
// This can be achieved because chess is a zero-sum game and one player's
// advantage is the other's disadvantage.
// https://www.chessprogramming.org/Negamax
//
// This function also implements alpha-beta pruning to reduce the amount of
// nodes that need to be searched, due to the fact that a single refutation
// is enough to mark a position as worse compared to an already found one.
// https://www.chessprogramming.org/Alpha-Beta
func (search *Context) negamax(plys, depth int, alpha, beta eval.Eval, pv *move.Variation) eval.Eval {
	search.nodes++

	// quick exit clauses
	switch {
	case search.shouldStop():
		// some search limit has been breached
		// the return value doesn't matter since this search's result
		// will be trashed and the previous iteration's pv will be used
		return 0

	case search.Board.IsDraw():
		// position is draw due to 50-move rule or threefold-repetition
		return search.draw()

	case depth <= 0, plys >= MaxDepth:
		// depth 0 reached, drop to quiescence search to prevent
		// the horizon effect from making the evaluation bad
		return search.quiescence(plys, alpha, beta)
	}

	// node properties
	isPVNode := beta-alpha != 1 // beta = alpha + 1 during PVS

	// generate all moves
	moves := search.Board.GenerateMoves(false)
	if len(moves) == 0 {
		// no legal moves, so some type of mate

		if search.Board.IsInCheck(search.Board.SideToMove) {
			return eval.MatedIn(plys) // checkmate
		}

		return eval.Draw // stalemate
	}

	// keep track of the original value of alpha for determining whether
	// the score will act as an upper bound entry in the transposition table
	originalAlpha := alpha

	// keep track of best move and score
	bestMove := move.Null
	bestEval := -eval.Inf

	// check for transposition table hits
	if entry, hit := search.tt.Probe(search.Board.Hash); hit {
		// use pv move for move ordering in any case
		bestMove = entry.Move

		// only use entry if current node is not a pv node and
		// entry depth is >= current depth (not worse quality)
		if !isPVNode && entry.Depth >= depth {
			search.ttHits++
			value := entry.Value.Eval(plys)

			switch entry.Type {
			case tt.ExactEntry:
				return value
			case tt.LowerBound:
				alpha = util.Max(alpha, value)
			case tt.UpperBound:
				beta = util.Min(beta, value)
			}

			if alpha >= beta {
				return value // fail high
			}
		}
	}

	// move ordering; score the generated moves
	list := move.ScoreMoves(moves, eval.OfMove(search.Board, bestMove))
	for i := 0; i < list.Length; i++ {
		var childPV move.Variation

		move := list.PickMove(i)

		search.Board.MakeMove(move)

		// Principal Variation Search

		var eval eval.Eval

		if !isPVNode || i > 0 {
			// null window search for non-pv nodes
			eval = -search.negamax(plys+1, depth-1, -alpha-1, -alpha, &childPV)
		}

		if isPVNode && ((eval > alpha && eval < beta) || i == 0) {
			// full window search for pv nodes
			eval = -search.negamax(plys+1, depth-1, -beta, -alpha, &childPV)
		}

		search.Board.UnmakeMove()

		// update score and bounds
		if eval > bestEval {
			// better move found
			bestMove = move
			bestEval = eval

			// check if move is new pv move
			if eval > alpha {
				// new pv so alpha increases
				alpha = eval

				// update parent pv
				pv.Update(move, childPV)

				if alpha >= beta {
					break // fail high
				}
			}
		}
	}

	// if search is stopped, score may be of a bad quality and
	// thus can pollute the transposition table for future searches
	if !search.stopped {
		var entryType tt.EntryType
		switch {
		case bestEval <= originalAlpha:
			// if score <= alpha, it is a worse position for the max player than
			// a previously explored line, since the move's exact score is at
			// most score. Therefore, it is an upperbound on the exact score.
			entryType = tt.UpperBound
		case bestEval >= beta:
			// if score >= beta, it is a worse position for the min player than
			// a previously explored line, singe the move's exact score is at
			// least score. Therefore, it is a lowerbound on the exact score.
			entryType = tt.LowerBound
		default:
			// if score is inside the bounds of alpha and beta, both the players
			// have been able to improve their position and it is an exact score.
			entryType = tt.ExactEntry
		}

		// update transposition table
		search.tt.Store(tt.Entry{
			Hash:  search.Board.Hash,
			Value: tt.EvalFrom(bestEval, plys),
			Move:  bestMove,
			Depth: depth,
			Type:  entryType,
		})
	}

	return bestEval
}