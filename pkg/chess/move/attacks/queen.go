package attacks

import (
	"laptudirm.com/x/mess/pkg/chess/bitboard"
	"laptudirm.com/x/mess/pkg/chess/square"
)

func Queen(s square.Square, occ bitboard.Board) bitboard.Board {
	return Rook(s, occ) | Bishop(s, occ)
}
