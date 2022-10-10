package evaluation

import (
	"laptudirm.com/x/mess/pkg/piece"
)

// Rel represents a relative centipawn evaluation where > 0 is better for
// the side to move, while < 0 is better for the other side.
type Rel int

// constants representing useful relative evaluations
const (
	// limits to differenciate between regular and mate in n evaluations
	WinInMaxPly  Rel = Mate - 2*10000
	LoseInMaxPly Rel = -WinInMaxPly
)

// Abs converts a Rel from the perspective of s to an absolute evaluation.
func (r Rel) Abs(s piece.Color) Abs {
	switch s {
	case piece.White:
		return Abs(r)
	case piece.Black:
		return Abs(-r)
	default:
		panic("bad color")
	}
}
