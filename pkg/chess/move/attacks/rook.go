package attacks

import (
	"laptudirm.com/x/mess/pkg/chess/bitboard"
	"laptudirm.com/x/mess/pkg/chess/square"
)

func rook(s square.Square, occ bitboard.Board, isMask bool) bitboard.Board {
	fileMask := bitboard.Files[s.File()]
	fileAttacks := hyperbola(s, occ, fileMask)

	rankMask := bitboard.Ranks[s.Rank()]
	rankAttacks := hyperbola(s, occ, rankMask)

	if isMask {
		fileAttacks &^= bitboard.Rank1 | bitboard.Rank8
		rankAttacks &^= bitboard.FileA | bitboard.FileH
	}

	return fileAttacks | rankAttacks
}

func Rook(s square.Square, blockers bitboard.Board) bitboard.Board {
	magic := RookMagics[s]
	blockers &= magic.BlockerMask
	return RookMoves[s][(uint64(blockers)*magic.Number)>>magic.Shift]
}
