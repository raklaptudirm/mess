package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"time"
)

func main() {
	// config
	//timeCon := "46.25+0.46"
	timeCon := "40+0.4s"
	gameNum := "5000"
	threads := "8"

	opponent := fmt.Sprintf("./testing/engines/%s", os.Args[1])

	// stage engines
	fmt.Print("info: staging engines... ")
	assert(run("make", "EXE=./testing/stage/mess"))
	fmt.Println("done.")

	// elo difference test
	assert(run(
		"cutechess-cli",
		"-repeat", "-recover", "-resign", "movecount=3", "score=400",
		"-draw", "movenumber=40", "movecount=8", "score=10", "-srand", strconv.Itoa(int(time.Now().Unix())),
		"-variant", "standard", "-concurrency", threads, "-games", gameNum,
		"-engine", "cmd=./testing/stage/mess", "proto=uci", "tc="+timeCon, "option.Hash=64", "name=mess", "stderr=testing/stderr.log",
		"-engine", "cmd="+opponent, "proto=uci", "tc="+timeCon, "option.Hash=64", "name=bench",
		"-openings", "file=testing/books/Openings.pgn", "format=pgn", "order=random", "plies=16", "-pgnout", "testing/pgns/games.pgn",
		//"-debug",
	))
}

func run(path string, args ...string) error {
	cmd := exec.Command(path, args...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func assert(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
