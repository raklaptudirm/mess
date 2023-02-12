package util

import (
	"os"
	"os/exec"
	"strings"
)

// RunNormal runs the given command with the standard input and output.
func RunNormal(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// RunQuiet runs the given command with output silenced.
func RunQuiet(args ...string) error {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

// RunWithOutput runs the given command and returns the output.
func RunWithOutput(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stderr = os.Stderr   // print any errors
	out, err := cmd.Output() // copy the stdout

	return strings.TrimSuffix(string(out), "\n"), err
}
