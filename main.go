package main

import (
	"fmt"
	"os"

	"laptudirm.com/x/mess/internal/build"
	"laptudirm.com/x/mess/internal/engine"
)

func main() {
	// run engine
	if err := run(); err != nil {
		// exit with error
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// quiet exit
}

func run() error {
	// create new UCI client
	client := engine.NewClient()

	// engine header with name, version, and author
	fmt.Printf("Mess %s by Rak Laptudirm\n", build.Version)

	switch args := os.Args[1:]; {
	case len(args) == 0:
		// no command-line arguments: start repl
		return client.Start()

	default:
		// command-line arguments: evaluate arguments as an UCI command
		// since we are not in a repl don't run any commands in parallel
		return client.RunWith(args, false)
	}
}
