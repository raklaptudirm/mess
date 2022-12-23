package flag

import (
	"fmt"
)

// NewSchema initializes a new flag Schema.
func NewSchema() Schema {
	return Schema{
		flags: make(map[string]Flag),
	}
}

// Schema contains the flag schema for a command.
type Schema struct {
	// flags maps each flag's name to it's collection function
	flags map[string]Flag
}

// Parse parses the given argument list according to the given flag
// schema. It returns the values for each command and an error.
func (s Schema) Parse(args []string) (Values, error) {
	values := make(Values)

	// nil schema
	if s.flags == nil {
		return values, nil
	}

	for len(args) > 0 {
		name := args[0]

		collect, isFlag := s.flags[name]
		if !isFlag {
			return values, fmt.Errorf("parse flags: unknown flag %q", name)
		}

		// collect value from arguments
		value, newArgs, err := collect(args[1:])
		if err != nil {
			return values, err
		}

		args = newArgs
		values[name] = Value{
			Set:   true,
			Value: value,
		}
	}

	return values, nil
}

// Button adds a button flag with the given name to the schema. A button
// flag is a flag without any arguments, it is either set or not set. All
// the values of Button flags are equal to nil.
func (s Schema) Button(name string) {
	s.flags[name] = func(args []string) (any, []string, error) {
		return nil, args, nil
	}
}

// Single adds a single flag with the given name to the schema. A single
// flag is a flag with a single argument. Values of single flags are of
// type string.
func (s Schema) Single(name string) {
	s.flags[name] = func(args []string) (any, []string, error) {
		if len(args) == 0 {
			return nil, nil, argNumErr(name, 1, 0)
		}

		return args[0], args[1:], nil
	}
}

// Array adds an array flag with the given name and argument number to the
// schema. An array flag is a flag with a fixed number of arguments. Values
// of array flags are of type []string.
func (s Schema) Array(name string, argN int) {
	s.flags[name] = func(args []string) (any, []string, error) {
		value := make([]string, argN)
		if collected := copy(value, args); collected != argN {
			return nil, nil, argNumErr(name, argN, collected)
		}

		return value, args[argN:], nil
	}
}

// Variadic adds a variadic flag with the given name to the schema. A
// variadic flag is a flag which collects all the remaining arguments.
// Values of variadic flags are of type []string.
func (s Schema) Variadic(name string) {
	s.flags[name] = func(s []string) (any, []string, error) {
		return s, []string{}, nil
	}
}

// Flag represents a flag of an uci command. Flag is a collector function
// which collects it's arguments from the provided list, and return's it's
// value, the remaining arguments, and an error, if any.
type Flag func([]string) (any, []string, error)

// Values map's each flag's name to it's value in the current interaction.
type Values map[string]Value

// Value represents the value of a flag.
type Value struct {
	// Set stores whether or not this flag was set.
	Set bool

	// Value contains the value of the flag. It should be type casted to
	// it's proper type before use. See the documentation of the various
	// flag's for their value's data types.
	Value any
}

func argNumErr(flag string, expected, collected int) error {
	return fmt.Errorf("flag %s: expected %d args, collected %d args", flag, expected, collected)
}
