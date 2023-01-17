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

// Package option implements functionality for using and supporting UCI
// options.
package option

import (
	"fmt"
	"strconv"
	"strings"
)

// NewSchema returns a new option schema.
func NewSchema() Schema {
	return Schema{
		options: make(map[string]Option),
	}
}

// Schema represents the schema of the options supported by an UCI client.
// It contains a map which maps an option's name to it's type.
type Schema struct {
	options map[string]Option
}

// AddOptions adds an option with the given name and type to the schema.
func (schema *Schema) AddOption(name string, option Option) {
	schema.options[name] = option
}

// SetDefaults sets the default values for all the options in the schema.
func (schema *Schema) SetDefaults() error {
	for _, option := range schema.options {
		if err := option.Initialize(); err != nil {
			return err
		}
	}

	return nil
}

// SetOption sets the given option to the given value.
func (schema *Schema) SetOption(name string, value []string) error {
	option, found := schema.options[name]
	if !found {
		return fmt.Errorf("set option: %q is not a valid option", name)
	}

	return option.Store(value)
}

// String converts the given Schema to an UCI compatible string
// representation. This should be printed when the `uci` command is
// received by a client.
func (s *Schema) String() string {
	var str string

	for name, option := range s.options {
		str += fmt.Sprintf("option name %s type %s\n", name, option.Type())
	}

	return str
}

// Option is the interface implemented by all the different types of
// options.
type Option interface {
	// type string
	Type() string

	// storage funcs
	Store([]string) error // normal storage
	Initialize() error    // default storage
}

// Check represents an UCI option of type check.
// UCI Specification: a checkbox that can either be true or false
type Check struct {
	Default bool

	// user defined storage function
	Storage func(bool) error
}

// compile time check that *Check implements Option.
var _ Option = (*Check)(nil)

// Type returns the type string of the given check option.
// Format: check default <default value>
func (option *Check) Type() string {
	return fmt.Sprintf("check default %v", option.Default)
}

// Store implements the storage function for a check option.
func (option *Check) Store(value []string) error {
	// expect 1 argument
	if len(value) != 1 {
		return fmt.Errorf("option check: expected %d values, received %d values", 1, len(value))
	}

	// argument should represent a boolean
	boolean, err := strconv.ParseBool(value[0])
	if err != nil {
		return err
	}

	// call the user defined storage function
	return option.Storage(boolean)
}

// Initialize stores the default value for the check option.
func (option *Check) Initialize() error {
	// call the user defined storage function on the default value
	return option.Storage(option.Default)
}

// Spin represents an UCI option of the type spin.
// UCI Specification: a spin wheel that can be an integer in a certain
// range defined by min and max
type Spin struct {
	Default  int
	Max, Min int

	// user defined storage function
	Storage func(int) error
}

// compile time check that *Spin represents Option.
var _ Option = (*Spin)(nil)

// Type returns the type string of the given string option.
// Format: spin default <default value> min <min value> max <max value>
func (option *Spin) Type() string {
	return fmt.Sprintf("spin default %v min %d max %d", option.Default, option.Min, option.Max)
}

// Store implements the storage function for a spin option.
func (option *Spin) Store(value []string) error {
	// expect 1 argument
	if len(value) != 1 {
		return fmt.Errorf("option spin: expected %d values, received %d values", 1, len(value))
	}

	// argument should represent an integer
	integer, err := strconv.Atoi(value[0])
	if err != nil {
		return err
	}

	// integer should be inside the provided bounds
	if integer < option.Min || integer > option.Max {
		return fmt.Errorf("option spin: value out of bounds [%d, %d]", option.Min, option.Max)
	}

	// call user defined storage function
	return option.Storage(integer)
}

// Initialize stores the default value for the spin option.
func (option *Spin) Initialize() error {
	// call the user defined storage function on the default value
	return option.Storage(option.Default)
}

// Button represents an UCI option of type button.
// UCI Specification: a button that can be pressed to send a command to the
// engine
type Button struct {
	// user defined ping function
	Ping func() error
}

// compile type check that *Button implements Option.
var _ Option = (*Button)(nil)

// Type returns the type string of the given button option.
// Format: button
func (option *Button) Type() string {
	return "button"
}

// Store implements the storage function for a button option.
func (option *Button) Store(value []string) error {
	// expect no arguments
	if len(value) > 0 {
		return fmt.Errorf("option button: expected %d values, received %d values", 0, len(value))
	}

	// ping the client
	return option.Ping()
}

// Initialize is a dummy function defined so that *Button can implement the
// Option interface. Buttons don't have any default values.
func (Option *Button) Initialize() error {
	// do nothing
	return nil
}

// String represents an UCI option of type string.
type String struct {
	Default string

	// user defined storage function
	Storage func(string) error
}

// compile time check that *String implements Option.
var _ Option = (*String)(nil)

// Type returns the type string of the given string option.
func (option *String) Type() string {
	return fmt.Sprintf("string default %s", option.Default)
}

// Store implements the storage function for a string option.
func (option *String) Store(value []string) error {
	return option.Storage(strings.Join(value, " "))
}

// Initialize stores the default value for the string option.
func (option *String) Initialize() error {
	return option.Storage(option.Default)
}
