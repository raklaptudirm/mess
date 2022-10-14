package attacks

import (
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/square"
	"laptudirm.com/x/mess/pkg/util"
)

const MaxRookBlockerSets = 4096
const MaxBishopBlockerSets = 512

var RookMagics [square.N]Magic
var BishopMagics [square.N]Magic

var RookMoves [square.N][MaxRookBlockerSets]bitboard.Board
var BishopMoves [square.N][MaxBishopBlockerSets]bitboard.Board

var MagicSeeds = [8]uint64{255, 16645, 15100, 12281, 32803, 55013, 10316, 728}

type Magic struct {
	Number      uint64
	BlockerMask bitboard.Board
	Shift       byte
}

func generateRookMagics() {
	var rand util.PRNG

	for s := square.A8; s <= square.H1; s++ {
		magic := &RookMagics[s]

		magic.BlockerMask = rook(s, bitboard.Empty, true)
		bitCount := magic.BlockerMask.CountBits()
		magic.Shift = uint8(64 - bitCount)

		permutationsN := 1 << bitCount
		permutations := make([]bitboard.Board, permutationsN)
		blockers := bitboard.Empty

		for index := 0; blockers != bitboard.Empty || index == 0; index++ {
			permutations[index] = blockers
			blockers = (blockers - magic.BlockerMask) & magic.BlockerMask
		}

		rand.Seed(MagicSeeds[s.Rank()])

	searchingMagicNumber:
		for {
			candidate := rand.SparseUint64()
			magic.Number = candidate

			RookMoves[s] = [MaxRookBlockerSets]bitboard.Board{}

			for i := 0; i < permutationsN; i++ {
				blockers := permutations[i]
				index := (uint64(blockers) * candidate) >> magic.Shift
				attacks := rook(s, blockers, false)

				if RookMoves[s][index] != bitboard.Empty && RookMoves[s][index] != attacks {
					continue searchingMagicNumber
				}

				RookMoves[s][index] = attacks
			}

			break
		}
	}
}

func generateBishopMagics() {
	var rand util.PRNG

	for s := square.A8; s <= square.H1; s++ {
		magic := &BishopMagics[s]

		magic.BlockerMask = bishop(s, bitboard.Empty, true)
		bitCount := magic.BlockerMask.CountBits()
		magic.Shift = uint8(64 - bitCount)

		permutationsN := 1 << bitCount
		permutations := make([]bitboard.Board, permutationsN)
		blockers := bitboard.Empty
		index := 0

		for blockers != bitboard.Empty || index == 0 {
			permutations[index] = blockers
			index++
			blockers = (blockers - magic.BlockerMask) & magic.BlockerMask
		}

		rand.Seed(MagicSeeds[s.Rank()])

	searchingMagicNumber:
		for {
			candidate := rand.SparseUint64()
			magic.Number = candidate

			BishopMoves[s] = [MaxBishopBlockerSets]bitboard.Board{}

			for i := 0; i < permutationsN; i++ {
				blockers := permutations[i]
				index := (uint64(blockers) * candidate) >> magic.Shift
				attacks := bishop(s, blockers, false)

				if BishopMoves[s][index] != bitboard.Empty && BishopMoves[s][index] != attacks {
					continue searchingMagicNumber
				}

				BishopMoves[s][index] = attacks
			}

			break
		}
	}
}

func init() {
	generateRookMagics()
	generateBishopMagics()
}
