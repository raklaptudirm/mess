package zobrist

import (
	"laptudirm.com/x/mess/pkg/move"
	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
)

type Key uint64

var PieceSquare [2 * piece.N][square.N]Key
var EnPassant [square.FileN]Key
var Castling [16]Key
var SideToMove Key

func init() {
	var rng PRNG
	rng.Seed(1070372)

	// piece square numbers
	for p := 0; p < piece.N; p++ {
		for s := square.A8; s <= square.H1; s++ {
			PieceSquare[p][s] = Key(rng.Uint64())
		}
	}

	// en passant file numbers
	for f := square.FileA; f <= square.FileH; f++ {
		EnPassant[f] = Key(rng.Uint64())
	}

	// castling right numbers
	for r := move.CastleNone; r <= move.CastleAll; r++ {
		Castling[r] = Key(rng.Uint64())
	}

	// black to move number
	SideToMove = Key(rng.Uint64())
}
