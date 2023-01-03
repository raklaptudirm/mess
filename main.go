// Copyright © 2023 Rak Laptudirm <rak@laptudirm.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
