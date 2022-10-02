package board

import "fmt"

func Perft(fen string, depth int) int {
	if depth == 0 {
		return 1
	}

	b := New(fen)

	var nodes int
	moves := b.GenerateMoves()

	for _, move := range moves {
		b.MakeMove(move)

		if !b.IsAttacked(b.kings[b.sideToMove.Other()], b.sideToMove) {
			newNodes := perft(b, depth-1)
			fmt.Printf("%s: %d\n", move, newNodes)
			nodes += newNodes
		}

		b.Unmove(move)
	}

	return nodes
}

func perft(b *Board, depth int) int {
	if depth == 0 {
		return 1
	}

	var nodes int
	moves := b.GenerateMoves()

	for _, move := range moves {
		b.MakeMove(move)

		if !b.IsAttacked(b.kings[b.sideToMove.Other()], b.sideToMove) {
			nodes += perft(b, depth-1)
		}

		b.Unmove(move)
	}

	return nodes
}
