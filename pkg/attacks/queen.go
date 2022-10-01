package attacks

import (
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/square"
)

func Queen(s square.Square, occ bitboard.Board) bitboard.Board {
	return Rook(s, occ) | Bishop(s, occ)
}
