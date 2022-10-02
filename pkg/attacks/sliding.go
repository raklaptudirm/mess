package attacks

import (
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/square"
)

func hyperbola(s square.Square, occ, mask bitboard.Board) bitboard.Board {
	var r bitboard.Board
	r.Set(s)

	o := occ & mask // masked occupancy
	return ((o - 2*r) ^ reverse(reverse(o)-2*reverse(r))) & mask
}

func reverse(b bitboard.Board) bitboard.Board {
	b = (b&0x5555555555555555)<<1 | ((b >> 1) & 0x5555555555555555)
	b = (b&0x3333333333333333)<<2 | ((b >> 2) & 0x3333333333333333)
	b = (b&0x0f0f0f0f0f0f0f0f)<<4 | ((b >> 4) & 0x0f0f0f0f0f0f0f0f)
	b = (b&0x00ff00ff00ff00ff)<<8 | ((b >> 8) & 0x00ff00ff00ff00ff)

	return (b << 48) | ((b & 0xffff0000) << 16) | ((b >> 16) & 0xffff0000) | (b >> 48)
}
