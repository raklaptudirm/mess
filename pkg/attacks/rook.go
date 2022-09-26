package attacks

import (
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/square"
)

func Rook(s square.Square, friends, occ bitboard.Board) bitboard.Board {
	fileMask := bitboard.Files[s.File()]
	fileAttacks := hyperbola(s, occ, fileMask)

	rankMask := bitboard.Ranks[s.Rank()]
	rankAttacks := hyperbola(s, occ, rankMask)

	return (fileAttacks | rankAttacks) &^ friends
}
