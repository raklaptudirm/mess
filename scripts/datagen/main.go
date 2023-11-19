package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/board/move"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/formats/fen"
	"laptudirm.com/x/mess/pkg/search"
	"laptudirm.com/x/mess/pkg/search/eval"
)

func main() {
	if err := Main(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Main() error {
	// Command-Line Flags:
	openings := flag.String("openings", "", "shuffled opening book containing a list of fens")
	offset := flag.Int("opening-offset", 0, "offset from which the books should be read")
	output := flag.String("output", "data.legacy", "output file for the generated fens and other data")
	games := flag.Int("games", 100_000, "number of games to generate data for (actual might be less)")
	threads := flag.Int("threads", 1, "number of threads to use for data generation work")
	winAdjudicateEval := flag.Uint("win-adjudicate-eval", uint(eval.Mate), "search score for which game will be adjudicated as a win")
	nodes := flag.Int("nodes", 10_000, "node limit for searches on a single fen")
	depth := flag.Int("depth", 0x0009, "depth limit for searches on a single fen")

	// Parse the CLI Flags.
	flag.Parse()

	// Create a new data generator.
	g, err := NewGenerator(*openings, *output, *offset, *games, *threads, *nodes, *depth, eval.Eval(*winAdjudicateEval))
	if err != nil {
		return err
	}

	// Generate Data.
	g.GenerateData()
	return nil
}

func NewGenerator(from, to string, offset, games, threads, nodes, depth int, winThreshold eval.Eval) (*Generator, error) {
	// Open the opening source epd file.
	i, err := os.Open(from)
	if err != nil {
		return nil, err
	}

	// Open the data target file.
	o, err := os.OpenFile(to, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Get most effective hash size for the node count.
	hashSize := util.Clamp((nodes*15)/(1024*1024), 1, 256)
	log.Printf("hash size used is %d mb\n", hashSize)

	// Create a new instance of a Generator.
	return &Generator{
		Input:  bufio.NewScanner(i),
		Output: bufio.NewWriterSize(o, 2000*100),

		Offset: offset,

		Openings: make(chan string),
		Data:     make(chan DataPoint),
		Deaths:   make(chan int),

		Games:   games,
		Threads: threads,

		Nodes: nodes,
		Depth: depth,

		HashSize: hashSize,

		WinThreshold: winThreshold,
	}, nil
}

type Generator struct {
	// Input and Output files.
	Input  *bufio.Scanner
	Output *bufio.Writer

	// Opening Offset in Input.
	Offset int

	// Sync channels.
	Openings chan string
	Data     chan DataPoint
	Deaths   chan int

	// Number of games done.
	Done int

	// Generator Configuration.

	Games   int // Total number of games to play.
	Threads int // Total number of threads to use.

	// Threshold for win adjudication.
	WinThreshold eval.Eval

	Nodes int // Node limit for searches.
	Depth int // Depth limit for searches.

	HashSize int // Most efficient hash size.
}

func (generator *Generator) GenerateData() {
	log.Printf("starting %d workers\n", generator.Threads)
	for i := 1; i <= generator.Threads; i++ {
		go generator.StartWorker(i)
	}

	log.Printf("playing %d games\n", generator.Games)
	go generator.ScheduleWork()

	start := time.Now()
	datapoints := 0
	deaths := 0

	for {
		select {
		case data := <-generator.Data:
			_, _ = generator.Output.WriteString(data.String())
			datapoints++

			if datapoints&4095 == 0 {
				delta := int(time.Since(start).Seconds()) + 1
				log.Printf(
					"%10d fens [%4d fens/second] %8d games [%2d games/second] [%3d fens/game]\n",
					datapoints, datapoints/delta, generator.Done, generator.Done/delta, datapoints/generator.Done,
				)
			}

		case <-generator.Deaths:
			if deaths++; deaths == generator.Threads {
				close(generator.Deaths)
				close(generator.Data)

				_ = generator.Output.Flush()

				log.Println("all workers are done")
				return
			}
		}
	}
}

func (generator *Generator) ScheduleWork() {
	for i, openings := 0, 0; openings < generator.Games && generator.Input.Scan(); i++ {
		if i >= generator.Offset {
			openings++
			generator.Openings <- generator.Input.Text()
		}
	}

	close(generator.Openings)
}

func (generator *Generator) StartWorker(id int) {
	data := make([]DataPoint, 0)

	limits := search.Limits{
		Depth:    generator.Depth,
		Nodes:    generator.Nodes,
		Infinite: true,
	}

	worker := search.NewContext(func(report search.Report) {}, generator.HashSize)

	for opening := range generator.Openings {
		worker.UpdatePosition(fen.FromString(opening))

		board := worker.Board()

		data = data[:0]
		var result = float32(0.5)

		for {
			if board.DrawClock >= 100 ||
				board.IsRepetition() {
				break
			}

			pv, score, _ := worker.Search(limits)

			score = util.Ternary(board.SideToMove == piece.White, score, -score)
			bestMove := pv.Move(0)

			if bestMove == move.Null || util.Abs(score) >= generator.WinThreshold {
				result = util.Ternary[float32](score > eval.Draw, 1.0, 0.0)
				break
			}

			// Position Filters: Tactical Positions
			if !bestMove.IsQuiet() && !board.IsInCheck(board.SideToMove) {
				goto nextMove
			}

			data = append(data, DataPoint{
				FEN:   board.FEN().String(),
				Score: score,
			})

		nextMove:
			worker.MakeMove(bestMove)
		}

		for i := 0; i < len(data); i++ {
			data[i].Result = result
			generator.Data <- data[i]
		}

		generator.Done++
	}

	generator.Deaths <- id
}

type DataPoint struct {
	FEN    string
	Score  eval.Eval
	Result float32
}

func (data *DataPoint) String() string {
	return fmt.Sprintf("%s | %d | %.1f\n", data.FEN, data.Score, data.Result)
}
