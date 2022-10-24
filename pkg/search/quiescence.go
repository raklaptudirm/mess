package search

import (
	"laptudirm.com/x/mess/pkg/evaluation"
	"laptudirm.com/x/mess/pkg/util"
)

// Quiescence search is a type of limited search which only evaluates 'quiet'
// positions, i.e. positions with no tactical moves like captures or promotions.
// This search is needed to avoid the horizon effect.
//
// https://www.chessprogramming.org/Quiescence_Search
//
func (c *Context) Quiescence(plys int, alpha, beta evaluation.Rel) evaluation.Rel {
	score := c.evalFunc(c.board) // standing pat
	alpha = util.Max(alpha, score)
	if alpha >= beta {
		return score
	}

	// searching similar to Negamax, but only considering tactical moves

	moves := c.board.GenerateMoves()

	switch {
	case len(moves) == 0:
		if c.board.CheckN > 0 {
			// prefer the longer lines if getting mated, and vice versa
			return evaluation.Rel(-evaluation.Mate + plys)
		}

		return evaluation.Draw // stalemate

	case c.board.IsDraw():
		return evaluation.Draw

	default:
		for _, m := range moves {
			var curr evaluation.Rel
			c.board.MakeMove(m)

			if !m.IsCapture() && !m.IsPromotion() {
				c.board.UnmakeMove()
				continue
			}

			curr = -c.Quiescence(plys+1, -beta, -alpha)

			c.board.UnmakeMove()

			score = util.Max(score, curr)
			alpha = util.Max(alpha, score)

			if alpha >= beta {
				break
			}
		}
	}

	return score
}
