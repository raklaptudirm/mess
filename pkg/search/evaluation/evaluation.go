package evaluation

import (
	"math"

	"laptudirm.com/x/mess/pkg/chess"
	"laptudirm.com/x/mess/pkg/chess/piece"
	"laptudirm.com/x/mess/pkg/chess/square"
)

// Type Func represents a board evaluation function.
type Func func(*chess.Board) Rel

// constants representing useful evaluations
const (
	Inf  = math.MaxInt32 / 2
	Mate = Inf - 1
	Draw = 0
)

// Of is a simple evaluation.Func which evaluates a position based on the
// material which each side has.
func Of(b *chess.Board) Rel {
	var eval Abs

	for s := square.A8; s <= square.H1; s++ {
		p := b.Position[s]
		if p == piece.NoPiece {
			continue
		}

		eval += material[p] + squareBonuses[p][s].Abs(p.Color())
	}

	return eval.Rel(b.SideToMove)
}
