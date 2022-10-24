// Package search implements various functions used to search a
// position for the best move. The search functions can be configured to
// use any evaluation function during it's search.
package search

import (
	"errors"

	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/evaluation"
	"laptudirm.com/x/mess/pkg/move"
	"laptudirm.com/x/mess/pkg/search/transposition"
	"laptudirm.com/x/mess/pkg/util"
)

func NewContext(fen string) Context {
	return Context{
		board:    board.New(fen),
		evalFunc: evaluation.Of,
		ttable:   transposition.NewTable(40),
	}
}

// Context stores various options, state, and debug variables regarding a
// particular search.
type Context struct {
	// search options
	evalFunc evaluation.Func

	// search state
	board  *board.Board
	ttable *transposition.Table
}

// error values for illegal or mated positions
var (
	ErrMate    = errors.New("search move: position is mate")
	ErrDraw    = errors.New("search move: position is draw")
	ErrIllegal = errors.New("search move: position is illegal")
)

// Search function searches the given position for the best move and
// returns the position's absolute evaluation, best move, and any
// encountered error. It is very similar to the Negamax function except
// for the fact that it keeps track of the best move along with the
// evaluation.
func Search(fen string, depth int) (move.Move, evaluation.Rel, error) {
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

	case c.board.IsDraw():
		return 0, evaluation.Draw, ErrDraw

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

		return bestMove, score, nil
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
	if entry, hit := c.ttable.Get(c.board.Hash); hit && entry.Depth >= depth {
		value := entry.Value.Rel(plys)

		switch entry.Type {
		case transposition.ExactEntry:
			return value
		case transposition.LowerBound:
			alpha = util.Max(alpha, value)
		case transposition.UpperBound:
			beta = util.Min(beta, value)
		}

		if alpha >= beta {
			return value
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

	case c.board.IsDraw():
		return evaluation.Draw

	// depth 0 reached
	case depth == 0:
		return c.Quiescence(plys, alpha, beta)

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

		var entryType transposition.TableEntryType
		switch {
		case score <= originalAlpha:
			// if score <= alpha, it is a worse position for the max player than
			// a previously explored line, since the move's exact score is at
			// most score. Therefore, it is an upperbound on the exact score.
			entryType = transposition.UpperBound
		case score >= beta:
			// if score >= beta, it is a worse position for the min player than
			// a previously explored line, singe the move's exact score is at
			// least score. Therefore, it is a lowerbound on the exact score.
			entryType = transposition.LowerBound
		default:
			// if score is inside the bounds of alpha and beta, both the players
			// have been able to improve their position and it is an exact score.
			entryType = transposition.ExactEntry
		}

		// update transposition table
		c.ttable.Put(c.board.Hash, transposition.TableEntry{
			Value: transposition.EvalFrom(score, plys),
			Depth: depth,
			Type:  entryType,
		})
		return score
	}
}
