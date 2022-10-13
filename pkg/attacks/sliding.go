package attacks

import (
	"math/bits"

	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/square"
)

func hyperbola(s square.Square, occ, mask bitboard.Board) bitboard.Board {
	r := bitboard.Squares[s]
	o := occ & mask // masked occupancy
	return ((o - 2*r) ^ reverse(reverse(o)-2*reverse(r))) & mask
}

func reverse(b bitboard.Board) bitboard.Board {
	return bitboard.Board(bits.Reverse64(uint64(b)))
}
