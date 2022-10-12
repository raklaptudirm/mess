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

		if !b.IsInCheck(b.SideToMove.Other()) {
			newNodes := perft(b, depth-1)
			fmt.Printf("%s: %d\n", move, newNodes)
			nodes += newNodes
		}

		b.UnmakeMove()
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

		if !b.IsInCheck(b.SideToMove.Other()) {
			nodes += perft(b, depth-1)
		}

		b.UnmakeMove()
	}

	return nodes
}
