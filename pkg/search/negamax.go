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
	"laptudirm.com/x/mess/pkg/board/bitboard"
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
	search.stats.Nodes++

	// node properties
	isCheck := search.Board.IsInCheck(search.Board.SideToMove)
	isPVNode := beta-alpha != 1 // beta = alpha + 1 during PVS
	isNullMove := search.Board.Plys > 0 && search.Board.History[search.Board.Plys-1].Move == move.Null

	// Check Extension: If position is in check, extend search depth so
	// that we don't push anything important over the horizon. This also
	// prevents search from going into qsearch while in check.
	if isCheck {
		depth++
	}

	// quick exit clauses
	switch {
	case search.shouldStop():
		// some search limit has been breached
		// the return value doesn't matter since this search's result
		// will be trashed and the previous iteration's pv will be used
		return 0

	case search.Board.DrawClock >= 100,
		plys == 0 && search.Board.IsThreefoldRepetition(),
		plys != 0 && search.Board.IsRepetition():
		// position is draw due to 50-move rule or threefold-repetition
		return search.draw()

	case plys >= MaxDepth:
		// maximum search depth reached, return static evaluation
		return search.score()

	case depth <= 0:
		// depth 0 reached, drop to quiescence search to prevent
		// the horizon effect from making the evaluation bad
		return search.quiescence(plys, alpha, beta)
	}

	// generate all moves
	moves := search.Board.GenerateMoves(false)
	if len(moves) == 0 {
		// position is mated; checkmate if king is in check
		return util.Ternary(isCheck, eval.MatedIn(plys), eval.Draw)
	}

	// keep track of the original value of alpha for determining whether
	// the score will act as an upper bound entry in the transposition table
	originalAlpha := alpha

	// keep track of best move and score
	bestMove := move.Null
	bestScore := -eval.Inf

	// non-static position evaluation used by
	// some heuristics and pruning techniques
	var posEval eval.Eval

	// check for transposition table hits
	if entry, hit := search.tt.Probe(search.Board.Hash); hit {
		// use pv move for move ordering in any case
		bestMove = entry.Move

		// use tt score as position eval when available
		posEval = entry.Value.Eval(plys)

		// only use entry if current node is not a pv node and
		// entry depth is >= current depth (not worse quality)
		if !isPVNode && !isNullMove && entry.Depth >= depth {
			search.stats.TTHits++
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
	} else {
		// when tt score is not available, use the position's
		// static evaluation as the position evaluation
		posEval = search.score()
	}

	// Internal Iterative Reduction (IIR): If a hash move is not found by
	// probing the transposition table, do a shallower search, as our move
	// ordering won't be as effective.
	if depth >= 4 && bestMove == move.Null {
		depth--
	}

	if !isPVNode && !isCheck {
		// Reverse Futility Pruning (RFP): The position is so far above
		// beta that we can expect the node to fail high and thus we can
		// safely prune this branch.
		if depth <= 5 && posEval >= beta && posEval-eval.Eval(75*depth) >= beta && posEval < eval.WinInMaxPly {
			return posEval
		}

		// Razoring: If static evaluation is really bad, drop into qsearch
		// and if qsearch score is <= alpha, don't spend any more time
		// searching this node which will probably fail low.
		if depth <= 3 && posEval+eval.Eval(200*depth) <= alpha {
			if score := search.quiescence(plys, alpha, beta); score <= alpha {
				return score
			}
		}

		// Null Move Pruning (NMP): Based on the Null Move Observation(given a
		// free move, the side to move can almost always improve their position)
		// NMP reduces the search tree by giving the opponent a free move in a
		// position where the position evaluation is enough to cause a beta cutoff.
		// If the score is still high enough to cause a beta cutoff after a
		// null move, the branch can be safely pruned.
		//
		// However, this method fails in Zugzwang positions, were it is better to
		// do nothing than to move. Therefore, NMP is not used in endgame positions
		// containing only pawns, where zugzwang positions occur most frequently.
		if !isNullMove && depth >= 3 && posEval >= beta &&
			search.Board.NonPawnMaterial(search.Board.SideToMove) != bitboard.Empty {

			reduction := 5 + util.Min(4, depth/5) + util.Min(3, int((posEval-beta)/214))

			search.Board.MakeMove(move.Null)
			score := -search.negamax(plys+1, depth-reduction, -beta, -beta+1, &move.Variation{})
			search.Board.UnmakeMove()

			if score >= beta {
				if score >= eval.WinInMaxPly {
					// don't return mate evaluations
					// from the null move search
					return beta
				}

				return score
			}
		}
	}

	lmrDepth := util.Ternary(isPVNode, 4, 2)
	historyBonus := depthBonus(depth)

	// move ordering; score the generated moves
	list := move.ScoreMoves(moves, eval.OfMove(eval.ModeEvalInfo{
		Board:   &search.Board.Position,
		PVMove:  bestMove,
		Killers: search.killers[plys],
		History: &search.history[search.Board.SideToMove],
	}))

	for i := 0; i < list.Length; i++ {
		var childPV move.Variation

		move := list.PickMove(i)

		search.Board.MakeMove(move)

		var score eval.Eval

		switch {
		// Late Move Reduction (LMR): Assuming that our move ordering is
		// good, later moves are less likely to raise alpha. LMR is used to
		// quickly prove that a move will be worse than alpha by searching
		// it at a lower(reduced) depth.
		case depth >= 3 && !isCheck && i > lmrDepth:
			rDepth := reductions[depth][i+1]
			rDepth = util.Clamp(depth-rDepth, 1, depth+1)

			// reduced depth search
			score = -search.negamax(plys+1, rDepth, -alpha-1, -alpha, &childPV)
			if score <= alpha {
				break
			}

			// lmr failed: do a full depth research
			fallthrough

		case !isPVNode || i > 0:
			// full depth search if lmr failed or for a non-PV node
			score = -search.negamax(plys+1, depth-1, -alpha-1, -alpha, &childPV)
		}

		// Principal Variation Search (PVS): Search PV nodes with a full
		// window as they are expected to be the best move in a position
		// while non PV nodes are searched with a null window to prove
		// that they are worse compared to the PV.
		if isPVNode && ((score > alpha && score < beta) || i == 0) {
			// full window search for pv nodes
			score = -search.negamax(plys+1, depth-1, -beta, -alpha, &childPV)
		}

		search.Board.UnmakeMove()

		// update score and bounds
		if score > bestScore {
			// better move found
			bestMove = move
			bestScore = score

			// check if move is new pv move
			if score > alpha {
				// new pv so alpha increases
				alpha = score

				// update parent pv
				pv.Update(move, childPV)

				if alpha >= beta {
					// move ordering heuristics
					search.storeKiller(plys, move)           // killer move
					search.updateHistory(move, historyBonus) // history bonus

					break // fail high
				}
			}
		}

		// beta wasn't raised, so give move a history penalty
		search.updateHistory(move, -historyBonus)
	}

	// if search is stopped, score may be of a bad quality and
	// thus can pollute the transposition table for future searches
	if !search.stopped {
		var entryType tt.EntryType
		switch {
		case bestScore <= originalAlpha:
			// if score <= alpha, it is a worse position for the max player than
			// a previously explored line, since the move's exact score is at
			// most score. Therefore, it is an upperbound on the exact score.
			entryType = tt.UpperBound
		case bestScore >= beta:
			// if score >= beta, it is a worse position for the min player than
			// a previously explored line, since the move's exact score is at
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
			Value: tt.EvalFrom(bestScore, plys),
			Move:  bestMove,
			Depth: depth,
			Type:  entryType,
		})
	}

	return bestScore
}
