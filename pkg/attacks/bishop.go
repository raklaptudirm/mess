package attacks

import (
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/square"
)

func bishop(s square.Square, occ bitboard.Board, isMask bool) bitboard.Board {
	diagonalMask := bitboard.Diagonals[s.Diagonal()]
	diagonalAttack := hyperbola(s, occ, diagonalMask)

	antiDiagonalMask := bitboard.AntiDiagonals[s.AntiDiagonal()]
	antiDiagonalAttack := hyperbola(s, occ, antiDiagonalMask)

	attacks := diagonalAttack | antiDiagonalAttack
	if isMask {
		attacks &^= bitboard.Rank1 | bitboard.Rank8 | bitboard.FileA | bitboard.FileH
	}

	return attacks
}

func Bishop(s square.Square, blockers bitboard.Board) bitboard.Board {
	magic := BishopMagics[s]
	blockers &= magic.BlockerMask
	return BishopMoves[s][(uint64(blockers)*magic.Number) >> magic.Shift]
}
