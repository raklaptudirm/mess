package main

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/notnil/chess"
	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/board/move"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
	"laptudirm.com/x/mess/pkg/search"
)

func main() {
	engine := search.NewContext(func(r search.Report) {}, 256)
	limits := search.Limits{
		Infinite: true,
		MoveTime: math.MaxInt32,
		Depth:    7,
	}

	fenCount := 0

	start := time.Now()

	err := filepath.WalkDir("./data", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if !strings.HasSuffix(path, ".pgn") {
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			return err
		}
		defer f.Close()

		scanner := chess.NewScanner(f)
		games := []*chess.Game{}
		// while there's more to read in the file
		for scanner.Scan() {
			// scan the next game
			games = append(games, scanner.Next())
		}

		fmt.Fprintln(os.Stderr, len(games))

		for _, game := range games {
			var result string
			switch game.GetTagPair("Result").Value {
			case "1-0":
				result = "[1.0]"
			case "0-1":
				result = "[0.0]"
			case "1/2-1/2":
				result = "[0.5]"
			}

			moves := game.Moves()
			chessboard := board.New(board.FEN(board.StartFEN))
			for i, gameMove := range moves {
				if i == len(moves)-1 {
					break
				}

				source := square.Square(gameMove.S1())
				source = square.New(square.File(source%8), 7-square.Rank(source/8))
				target := square.Square(gameMove.S2())
				target = square.New(square.File(target%8), 7-square.Rank(target/8))

				boardMove := chessboard.NewMove(source, target)

				switch gameMove.Promo() {
				case chess.Knight:
					boardMove = boardMove.SetPromotion(piece.New(piece.Knight, chessboard.SideToMove))
				case chess.Bishop:
					boardMove = boardMove.SetPromotion(piece.New(piece.Bishop, chessboard.SideToMove))
				case chess.Rook:
					boardMove = boardMove.SetPromotion(piece.New(piece.Rook, chessboard.SideToMove))
				case chess.Queen:
					boardMove = boardMove.SetPromotion(piece.New(piece.Queen, chessboard.SideToMove))
				}

				chessboard.MakeMove(boardMove)

				if chessboard.IsInCheck(chessboard.SideToMove) {
					continue
				}

				fenString := chessboard.FEN()

				engine.UpdatePosition(fenString)
				variation, _, _ := engine.Search(limits)

				bestMove := variation.Move(0)

				if bestMove.IsCapture() || bestMove.IsPromotion() {
					continue
				}

				if bestMove == move.Null {
					fmt.Fprintln(os.Stderr, "Ahhhhh!")
				}

				// print out FEN for each move in the game
				fmt.Println(result, fenString)
				fenCount++
			}

			fmt.Fprintf(os.Stderr, "data-gen: %d fens generated (%d fens/s)\n", fenCount, fenCount/(int(time.Since(start).Seconds())+1))
		}

		return nil
	})

	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
