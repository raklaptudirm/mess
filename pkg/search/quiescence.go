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
	"laptudirm.com/x/mess/pkg/search/evaluation"
	"laptudirm.com/x/mess/internal/util"
)

// Quiescence search is a type of limited search which only evaluates 'quiet'
// positions, i.e. positions with no tactical moves like captures or promotions.
// This search is needed to avoid the horizon effect.
//
// https://www.chessprogramming.org/Quiescence_Search
//
func (c *Context) Quiescence(plys int, alpha, beta evaluation.Rel) evaluation.Rel {
	score := c.evalFunc(c.Board) // standing pat
	alpha = util.Max(alpha, score)
	if alpha >= beta {
		return score
	}

	// searching similar to Negamax, but only considering tactical moves

	moves := c.Board.GenerateMoves()

	switch {
	case len(moves) == 0:
		if c.Board.CheckN > 0 {
			// prefer the longer lines if getting mated, and vice versa
			return evaluation.Rel(-evaluation.Mate + plys)
		}

		return evaluation.Draw // stalemate

	case c.Board.IsDraw():
		return evaluation.Draw

	default:
		for _, m := range moves {
			var curr evaluation.Rel
			c.Board.MakeMove(m)

			if !m.IsCapture() && !m.IsPromotion() {
				c.Board.UnmakeMove()
				continue
			}

			curr = -c.Quiescence(plys+1, -beta, -alpha)

			c.Board.UnmakeMove()

			score = util.Max(score, curr)
			alpha = util.Max(alpha, score)

			if alpha >= beta {
				break
			}
		}
	}

	return score
}
