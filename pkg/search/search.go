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
		Board:    board.New(fen),
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
	Board  *board.Board
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
func (c *Context) Search(depth int) (move.Move, evaluation.Rel, error) {
	// treat this function as the root call to Negamax
	// Negamax(depth, -Inf, Inf)
	alpha := evaluation.Rel(-evaluation.Inf)
	beta := evaluation.Rel(evaluation.Inf)

	moves := c.Board.GenerateMoves()

	switch {
	// king can be captured: illegal position
	case c.Board.IsInCheck(c.Board.SideToMove.Other()):
		return 0, evaluation.Inf, ErrIllegal // king is captured

	// no legal moves: position is mate
	case len(moves) == 0:
		return 0, evaluation.Mate, ErrMate

	case c.Board.IsDraw():
		return 0, evaluation.Draw, ErrDraw

	default:
		// keep track of the best move
		var bestMove move.Move
		score := evaluation.Rel(-evaluation.Inf)

		for _, m := range moves {
			c.Board.MakeMove(m)
			// one side's win is other side's loss
			// one move has been made so ply 1 from root
			curr := -c.Negamax(1, depth-1, -beta, -alpha)
			c.Board.UnmakeMove()

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
	if entry, hit := c.ttable.Get(c.Board.Hash); hit && entry.Depth >= depth {
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

	moves := c.Board.GenerateMoves()

	switch {
	// position is mate
	case len(moves) == 0:
		if c.Board.CheckN > 0 {
			// prefer the longer lines if getting mated, and vice versa
			return evaluation.Rel(-evaluation.Mate + plys)
		}

		return evaluation.Draw // stalemate

	case c.Board.IsDraw():
		return evaluation.Draw

	// depth 0 reached
	case depth == 0:
		return c.Quiescence(plys, alpha, beta)

	// keep searching
	default:
		score := evaluation.Rel(-evaluation.Inf)
		for i := 0; i < len(moves); i++ {
			c.orderMoves(moves, i)

			c.Board.MakeMove(moves[i])
			curr := -c.Negamax(plys+1, depth-1, -beta, -alpha)
			c.Board.UnmakeMove()

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
		c.ttable.Put(c.Board.Hash, transposition.TableEntry{
			Value: transposition.EvalFrom(score, plys),
			Depth: depth,
			Type:  entryType,
		})
		return score
	}
}

func (c *Context) orderMoves(moveList []move.Move, index int) {
	bestMove := evaluation.Move(-10000)
	bestIndex := -1
	for i, m := range moveList {
		if eval := evaluation.OfMove(c.Board, m); eval > bestMove {
			bestMove = eval
			bestIndex = i
		}
	}

	moveList[index], moveList[bestIndex] = moveList[bestIndex], moveList[index]
}
