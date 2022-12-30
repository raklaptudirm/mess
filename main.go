package main

import (
	"os"

	"laptudirm.com/x/mess/internal/build"
	"laptudirm.com/x/mess/internal/engine"
)

func main() {
	client := engine.NewClient()
	client.Printf("Mess %s by Rak Laptudirm\n", build.Version)
	if err := client.Start(); err != nil {
		client.Println(err)
		os.Exit(1)
	}
}
