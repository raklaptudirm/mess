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

package uci

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"laptudirm.com/x/mess/pkg/uci/cmd"
)

// NewClient creates a new uci.Client which is listening to the stdin for
// commands and with the default isready and quit commands added.
func NewClient() Client {
	client := Client{
		// communication streams
		stdin:  os.Stdin,
		stdout: os.Stdout,
	}

	client.commands = cmd.NewSchema(client.stdout)

	// add default commands
	client.AddCommand(cmdQuit)
	client.AddCommand(cmdIsReady)

	return client
}

// Client represents an UCI client.
type Client struct {
	stdin  io.Reader // GUI to Engine commands
	stdout io.Writer // Engine to GUI commands

	commands cmd.Schema // commands schema
}

// AddCommand adds the given command to the client's schema.
func (c *Client) AddCommand(cmd cmd.Command) {
	c.commands.Add(cmd)
}

// Start starts a repl listening for UCI commands which match the client's
// Schema. It listen's on the Client's stdin.
func (c *Client) Start() error {
	reader := bufio.NewReader(c.stdin)

	// read-eval-print loop
	for {
		// read prompt form client's stdin
		prompt, err := reader.ReadString('\n')
		if err != nil {
			// read errors are probably fatal
			return err
		}

		// parse arguments from prompt
		args := strings.Fields(prompt)

		// since we are in a repl run commands in parallel if needed
		switch err := c.RunWith(args, true); err {
		case nil:
			// no error: continue repl

		case errQuit:
			// errQuit is returned by quit command to stop the repl
			// so honour the request and return, stopping the repl
			return nil

		default:
			// non-nil error: print and continue
			c.Println(err)
		}
	}
}

// Run is a simple utility function which runs the provided arguments as a
// command without parallelization.
func (c *Client) Run(args ...string) error {
	return c.RunWith(args, false)
}

// RunWith finds a command whose name matches the first element of the args
// array, and runs it with the remaining args. It returns any error sent
// by the command. It honours the cmd.Parallel property if parallelize is
// set to true.
func (c *Client) RunWith(args []string, parallelize bool) error {
	// separate command name and arguments
	name, args := args[0], args[1:]

	// get uci command
	cmd, found := c.commands.Get(name)
	if !found {
		// command with given name not found
		return fmt.Errorf("%s: command not found", name)
	}

	// run command with given arguments
	return cmd.RunWith(args, parallelize, c.commands)
}

// Print acts as fmt.Print on the client's stdout.
func (c *Client) Print(a ...any) (int, error) {
	return fmt.Fprint(c.stdout, a...)
}

// Printf acts as fmt.Printf on the client's stdout.
func (c *Client) Printf(format string, a ...any) (int, error) {
	return fmt.Fprintf(c.stdout, format, a...)
}

// Println acts as fmt.Println on the client's stdout.
func (c *Client) Println(a ...any) (int, error) {
	return fmt.Fprintln(c.stdout, a...)
}
