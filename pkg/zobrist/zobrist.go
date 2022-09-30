package zobrist

import (
	"laptudirm.com/x/mess/pkg/castling"
	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
)

type Key uint64

var PieceSquare [piece.N][square.N]Key
var EnPassant [square.FileN]Key
var Castling [castling.N]Key
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
	for r := castling.None; r <= castling.All; r++ {
		Castling[r] = Key(rng.Uint64())
	}

	// black to move number
	SideToMove = Key(rng.Uint64())
}
