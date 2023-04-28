package classical_test

import (
	"testing"

	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/search/eval/classical"
)

func BenchmarkAccumulate(b *testing.B) {
	evaluator := classical.EfficientlyUpdatable{}
	chessboard := board.New(board.EU(&evaluator))
	evaluator.Board = chessboard
	chessboard.UpdateWithFEN(board.StartFEN)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		evaluator.Accumulate(piece.White)
	}
}
