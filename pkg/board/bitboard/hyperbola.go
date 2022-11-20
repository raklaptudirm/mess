package bitboard

import (
	"math/bits"

	"laptudirm.com/x/mess/pkg/board/square"
)

func Hyperbola(s square.Square, occ, mask Board) Board {
	r := Squares[s]
	o := occ & mask // masked occupancy
	return ((o - 2*r) ^ reverse(reverse(o)-2*reverse(r))) & mask
}

func reverse(b Board) Board {
	return Board(bits.Reverse64(uint64(b)))
}
