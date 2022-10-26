package castling

import (
	"laptudirm.com/x/mess/pkg/chess/piece"
	"laptudirm.com/x/mess/pkg/chess/square"
)

const Sides = 4

var Rooks = [square.N]struct {
	From, To square.Square
	RookType piece.Piece
}{
	square.G1: {
		From:     square.H1,
		To:       square.F1,
		RookType: piece.WhiteRook,
	},
	square.C1: {
		From:     square.A1,
		To:       square.D1,
		RookType: piece.WhiteRook,
	},
	square.G8: {
		From:     square.H8,
		To:       square.F8,
		RookType: piece.BlackRook,
	},
	square.C8: {
		From:     square.A8,
		To:       square.D8,
		RookType: piece.BlackRook,
	},
}
