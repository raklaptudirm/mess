// Package search implements various functions used to search a
// position for the best move. The search functions can be configured to
// use any evaluation function during it's search.
package search

import (
	"errors"
	"fmt"

	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/evaluation"
	"laptudirm.com/x/mess/pkg/move"
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

	// king can be captured; illegal position
	if c.board.IsInCheck(c.board.SideToMove.Other()) {
		return move.Move{}, evaluation.Inf, ErrIllegal // king is captured
	}

	// keep track of the best move
	var bestMove move.Move

	moves := c.board.GenerateMoves()
	score := evaluation.Rel(-evaluation.Inf)
	for _, m := range moves {
		c.board.MakeMove(m)
		// one side's win is other side's loss
		// one move has been made so ply 1 from root
		curr := -c.Negamax(1, depth-1, -beta, -alpha)
		c.board.UnmakeMove(m)

		if curr != -evaluation.Inf {
			fmt.Printf("%s %s\n", m, curr.Abs(c.board.SideToMove))
		}

		if curr > score {
			// better move found
			score = curr
			bestMove = m
		}

		alpha = max(alpha, score)
	}

	// position is mate; no legal moves
	if score == -evaluation.Inf {
		return move.Move{}, score.Abs(c.board.SideToMove), ErrMate
	}

	return bestMove, score.Abs(c.board.SideToMove), nil
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
			alpha = max(alpha, entry.value)
		case upperBound:
			beta = min(beta, entry.value)
		}

		if alpha >= beta {
			return entry.value
		}
	}

	// depth == 0 or terminal node

	if c.board.IsInCheck(c.board.SideToMove.Other()) {
		return evaluation.Inf // king is captured
	}

	if depth == 0 {
		return evaluation.Of(c.board)
	}

	// search moves

	moves := c.board.GenerateMoves()
	score := evaluation.Rel(-evaluation.Inf)
	for _, m := range moves {
		c.board.MakeMove(m)
		curr := -c.Negamax(plys+1, depth-1, -beta, -alpha)
		c.board.UnmakeMove(m)

		// update score and bounds

		score = max(score, curr)
		alpha = max(alpha, score)

		if alpha >= beta {
			break
		}
	}

	// check for mate
	if score == -evaluation.Inf {
		score = evaluation.Draw // stalemate
		if c.board.IsInCheck(c.board.SideToMove) {
			// prefer the longer lines if getting mated, and vice versa
			score = evaluation.Rel(-evaluation.Mate + plys)
		}
	}

	// update transposition table
	c.ttable.Put(c.board.Hash, plys, depth, score, originalAlpha, beta)
	return score
}
