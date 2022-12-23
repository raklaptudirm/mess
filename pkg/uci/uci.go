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

// Start makes the client start listening for commands.
func (c *Client) Start() error {
	reader := bufio.NewReader(c.stdin)

	for {
		// command strings are newline terminated
		prompt, err := reader.ReadString('\n')
		if err != nil {
			return err
		}

		switch err = c.commands.RunWith(strings.Fields(prompt)); err {
		case nil:
			// continue repl
		case errQuit:
			return nil
		default:
			c.Println(err)
		}
	}
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
