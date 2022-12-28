// Copyright © 2022 Rak Laptudirm <rak@laptudirm.com>
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

package cmd

import (
	"fmt"
	"io"

	"laptudirm.com/x/mess/pkg/uci/flag"
)

// NewSchema initializes a new command schema.
func NewSchema(replyWriter io.Writer) Schema {
	return Schema{
		replyWriter: replyWriter,
		commands:    make(map[string]Command),
	}
}

// Schema contains a command schema for a client.
type Schema struct {
	replyWriter io.Writer
	commands    map[string]Command
}

// Add adds the given command to the Schema.
func (l *Schema) Add(c Command) {
	l.commands[c.Name] = c
}

func (l *Schema) Get(name string) (Command, bool) {
	cmd, found := l.commands[name]
	return cmd, found
}

// Command represents the schema of a GUI to Engine command.
type Command struct {
	// name of the command
	// this is used as a token to identify if this command has been run
	Name string

	// If Parallel is true, the listener will not wait for the command
	// to finish before accepting new commands.
	Parallel bool

	// Run is the actual work function for the command. It is provided
	// with an interaction which contains the relevant information
	// about the command interaction by the GUI.
	Run func(Interaction) error

	// Flags contains the flag schema of this command. The flags the
	// parsed from the provided args before the Run function is called.
	Flags flag.Schema
}

func (c Command) RunWith(args []string, schema Schema) error {
	values, err := c.Flags.Parse(args)
	if err != nil {
		return err
	}

	return c.Run(Interaction{
		stdout:  schema.replyWriter,
		Command: c,

		Values: values,
	})
}

// Interaction encapsulates relevant information about a Command sent to
// the Engine by the GUI.
type Interaction struct {
	stdout io.Writer

	Command // parent Command

	// values provided for the command's flags
	Values flag.Values
}

// Reply writes to the GUI's input. It is similar to fmt.Println.
func (i *Interaction) Reply(a ...any) (int, error) {
	return fmt.Fprintln(i.stdout, a...)
}

// Replyf writes to the GUI's input. It is similar to fmt.Printf with
// a newline terminator.
func (i *Interaction) Replyf(format string, a ...any) (int, error) {
	return fmt.Fprintf(i.stdout, format+"\n", a...)
}
