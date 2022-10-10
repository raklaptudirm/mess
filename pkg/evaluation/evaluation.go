package evaluation

import (
	"math"

	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
)

// Type Func represents a board evaluation function.
type Func func(*board.Board) Rel

// constants representing useful evaluations
const (
	Inf  = math.MaxInt32 / 2
	Mate = Inf - 1
	Draw = 0
)

// Of is a simple evaluation.Func which evaluates a position based on the
// material which each side has.
func Of(b *board.Board) Rel {
	var eval Abs

	var material = [...]Abs{
		piece.NoPiece:     0,
		piece.WhitePawn:   100,
		piece.WhiteKnight: 300,
		piece.WhiteBishop: 300,
		piece.WhiteRook:   500,
		piece.WhiteQueen:  900,
		piece.WhiteKing:   0,
		piece.BlackPawn:   -100,
		piece.BlackKnight: -300,
		piece.BlackBishop: -300,
		piece.BlackRook:   -500,
		piece.BlackQueen:  -900,
		piece.BlackKing:   0,
	}

	for s := square.A8; s <= square.H1; s++ {
		eval += material[b.Position[s]]
	}

	return eval.Rel(b.SideToMove)
}
