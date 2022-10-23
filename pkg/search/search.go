// Package search implements various functions used to search a
// position for the best move. The search functions can be configured to
// use any evaluation function during it's search.
package search

import (
	"errors"

	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/evaluation"
	"laptudirm.com/x/mess/pkg/move"
	"laptudirm.com/x/mess/pkg/util"
)

func NewContext(fen string) Context {
	return Context{
		board:    board.New(fen),
		evalFunc: evaluation.Of,
		ttable:   make(transpositionTable),
	}
}

// Context stores various options, state, and debug variables regarding a
// particular search.
type Context struct {
	// search options
	evalFunc evaluation.Func

	// search state
	board  *board.Board
	ttable transpositionTable
}

// error values for illegal or mated positions
var (
	ErrMate    = errors.New("search move: position is mate")
	ErrIllegal = errors.New("search move: position is illegal")
)

// Search function searches the given position for the best move and
// returns the position's absolute evaluation, best move, and any
// encountered error. It is very similar to the Negamax function except
// for the fact that it keeps track of the best move along with the
// evaluation.
func Search(fen string, depth int) (move.Move, evaluation.Abs, error) {
	c := NewContext(fen)

	// treat this function as the root call to Negamax
	// Negamax(depth, -Inf, Inf)
	alpha := evaluation.Rel(-evaluation.Inf)
	beta := evaluation.Rel(evaluation.Inf)

	moves := c.board.GenerateMoves()

	switch {
	// king can be captured: illegal position
	case c.board.IsInCheck(c.board.SideToMove.Other()):
		return 0, evaluation.Inf, ErrIllegal // king is captured

	// no legal moves: position is mate
	case len(moves) == 0:
		return 0, evaluation.Mate, ErrMate

	default:
		// keep track of the best move
		var bestMove move.Move
		score := evaluation.Rel(-evaluation.Inf)

		for _, m := range moves {
			c.board.MakeMove(m)
			// one side's win is other side's loss
			// one move has been made so ply 1 from root
			curr := -c.Negamax(1, depth-1, -beta, -alpha)
			c.board.UnmakeMove()

			if curr > score {
				// better move found
				score = curr
				bestMove = m
			}

			alpha = util.Max(alpha, score)
		}

		return bestMove, score.Abs(c.board.SideToMove), nil
	}
}

// Negamax determines the evaluation of a particular position after a
// particular depth using the Negamax search algorithm.
//
// This function also implements alpha-beta pruning for lossless faster
// evaluation, and calls quiescence search to prevent the horizon effect.
// This function also uses a transposition table to prevent redoing work.
//
// https://www.chessprogramming.org/Negamax
// https://www.chessprogramming.org/Alpha-Beta
// https://www.chessprogramming.org/Quiescence_Search
// https://www.chessprogramming.org/Transposition_Table
//
func (c *Context) Negamax(plys, depth int, alpha, beta evaluation.Rel) evaluation.Rel {
	// keep track of the original value of alpha for determining whether
	// the score will act as an upper bound entry in the transposition table
	originalAlpha := alpha

	// check for transposition table hits
	if entry, hit := c.ttable.Get(c.board.Hash, plys, depth); hit {
		switch entry.eType {
		case exact:
			return entry.value
		case lowerBound:
			alpha = util.Max(alpha, entry.value)
		case upperBound:
			beta = util.Min(beta, entry.value)
		}

		if alpha >= beta {
			return entry.value
		}
	}

	// search moves

	moves := c.board.GenerateMoves()

	switch {
	// position is mate
	case len(moves) == 0:
		if c.board.CheckN > 0 {
			// prefer the longer lines if getting mated, and vice versa
			return evaluation.Rel(-evaluation.Mate + plys)
		}

		return evaluation.Draw // stalemate

	// depth 0 reached
	case depth == 0:
		return evaluation.Of(c.board)

	// keep searching
	default:
		score := evaluation.Rel(-evaluation.Inf)
		for _, m := range moves {
			c.board.MakeMove(m)
			curr := -c.Negamax(plys+1, depth-1, -beta, -alpha)
			c.board.UnmakeMove()

			// update score and bounds

			score = util.Max(score, curr)
			alpha = util.Max(alpha, score)

			if alpha >= beta {
				break
			}
		}

		// update transposition table
		c.ttable.Put(c.board.Hash, plys, depth, score, originalAlpha, beta)
		return score
	}
}
