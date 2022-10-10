package evaluation

import (
	"fmt"

	"laptudirm.com/x/mess/pkg/piece"
)

// Abs represents an absolute centipawn evaluation where > 0 is better
// for white and < 0 is better for black white 0 is drawn.
type Abs int

// constants representing useful absolute evaluations
const (
	WhiteWon Abs = Mate
	BlackWon Abs = -WhiteWon

	// limits to differenciate between regular and mate in n evaluations
	WhiteWinInMaxPly Abs = Abs(WinInMaxPly)
	BlackWinInMaxPly Abs = -WhiteWinInMaxPly
)

// String returns the string representation of the given absolute evaluation.
func (a Abs) String() string {
	var str string
	negative := false

	if a < Draw {
		a = -a
		negative = true
	}

	switch {
	case a == Inf:
		str = "(king captured)"
	case a == WhiteWon:
		str = "(checkmate)"
	case a >= WhiteWinInMaxPly:
		plys := WhiteWon - a
		str = fmt.Sprintf("#%d", (plys/2)+(plys%2))
	default:
		str = fmt.Sprintf("%d.%d", a/100, a%100)
	}

	if negative {
		return "-" + str
	}

	return str
}

// Rel converts an Abs to a Rel from the perspective of s.
func (a Abs) Rel(s piece.Color) Rel {
	switch s {
	case piece.White:
		return Rel(a)
	case piece.Black:
		return Rel(-a)
	default:
		panic("bad color")
	}
}
