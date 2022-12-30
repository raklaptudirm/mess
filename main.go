package main

import (
	"os"

	"laptudirm.com/x/mess/internal/engine"
)

func main() {
	client := engine.NewClient()
	client.Println("Mess by Rak Laptudirm")
	if err := client.Start(); err != nil {
		client.Println(err)
		os.Exit(1)
	}
}
