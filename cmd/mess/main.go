package main

import (
	"fmt"

	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/square"
)

func main() {
	b := board.New("rnbqkbnr/pppppppp/8/8/3Q7/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	// b := board.New("rnbqkbnr/pppppppp/8/8/2K5/8/PPPPPPPP/RNBQ1BNR w kq - 0 1")
	fmt.Println()
	fmt.Println(b)
	fmt.Println(b.MovesOf(square.D4))
}
