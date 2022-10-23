package board_test

import (
	"testing"

	"laptudirm.com/x/mess/pkg/board"
)

func TestFEN(t *testing.T) {
	tests := []string{
		"rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		"rnbqkbnr/pppp1ppp/8/4p3/4P3/8/PPPP1PPP/RNBQKBNR w KQkq - 0 2",
		"r1bqk1nr/pppp1ppp/2n5/2b1p3/2B1P3/5N2/PPPP1PPP/RNBQ1RK1 b kq - 5 4",
		"rnbq1rk1/ppp1bppp/4pn2/3p2B1/2PP4/2N2N2/PP2PPPP/R2QKB1R w KQ - 6 6",
		"rnbqkbnr/ppp2ppp/8/2Ppp3/8/8/PP1PPPPP/RNBQKBNR w KQkq d6 0 3",
		"rnbqkbnr/pp1ppppp/8/8/2pPP3/5N2/PPP2PPP/RNBQKB1R b KQkq d3 0 3",
		"rn3rk1/pbp1qpp1/1p5p/3p4/3P4/3BPN2/PP3PPP/R2Q1RK1 b - - 3 12",
	}

	for n, test := range tests {
		t.Run(test, func(t *testing.T) {
			b := board.New(test)
			newFEN := b.FEN()
			if test != newFEN {
				t.Errorf("test %d: wrong fen\n%s\n%s\n", n, test, newFEN)
			}
		})
	}
}
