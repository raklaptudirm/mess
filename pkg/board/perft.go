package board

import "fmt"

func Perft(b *Board, depth int) int {
	if depth == 0 {
		return 1
	}

	var nodes int
	moves := b.GenerateMoves()

	for _, move := range moves {
		b.MakeMove(move)
		newNodes := perft(b, depth-1)
		fmt.Printf("%s: %d\n", move, newNodes)
		nodes += newNodes
		b.UnmakeMove()
	}

	return nodes
}

func perft(b *Board, depth int) int {

	switch depth {
	case 0:
		return 1
	case 1:
		return len(b.GenerateMoves())
	default:
		var nodes int
		moves := b.GenerateMoves()

		for _, move := range moves {
			b.MakeMove(move)
			nodes += perft(b, depth-1)
			b.UnmakeMove()
		}

		return nodes
	}
}
