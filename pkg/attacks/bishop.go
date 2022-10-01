package attacks

import (
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/square"
)

func Bishop(s square.Square, occ bitboard.Board) bitboard.Board {
	diagonalMask := bitboard.Diagonals[s.Diagonal()]
	diagonalAttack := hyperbola(s, occ, diagonalMask)

	antiDiagonalMask := bitboard.AntiDiagonals[s.AntiDiagonal()]
	antiDiagonalAttack := hyperbola(s, occ, antiDiagonalMask)

	return (diagonalAttack | antiDiagonalAttack)
}
