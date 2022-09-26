package attacks

import (
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/square"
)

func Queen(s square.Square, friends, occ bitboard.Board) bitboard.Board {
	return Rook(s, friends, occ) | Bishop(s, friends, occ)
}
