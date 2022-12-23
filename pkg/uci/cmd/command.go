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

// RunWith runs commands from the command schema according
// to the given arguments.
func (l *Schema) RunWith(args []string) error {
	if len(args) == 0 {
		return nil
	}

	cmd, isCmd := l.commands[args[0]]
	if !isCmd {
		return fmt.Errorf("%s: command not found", args[0])
	}

	values, err := cmd.Flags.Parse(args[1:])
	if err != nil {
		return err
	}

	return cmd.Run(Interaction{
		stdout:  l.replyWriter,
		Command: cmd,

		Values: values,
	})
}

// Command represents the schema of a GUI to Engine command.
type Command struct {
	// name of the command
	// this is used as a token to identify if this command has been run
	Name string

	// Run is the actual work function for the command. It is provided
	// with an interaction which contains the relevant information
	// about the command interaction by the GUI.
	Run func(Interaction) error

	// Flags contains the flag schema of this command. The flags the
	// parsed from the provided args before the Run function is called.
	Flags flag.Schema
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
