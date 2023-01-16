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

package cmd

import (
	"errors"

	"laptudirm.com/x/mess/internal/engine/context"
	"laptudirm.com/x/mess/pkg/uci/cmd"
	"laptudirm.com/x/mess/pkg/uci/flag"
)

// UCI command setoption
//
// This is sent to the engine when the user wants to change the internal
// parameters of the engine. For the button type no value is needed.
//
// One string will be sent for each parameter and this will only be sent
// when the engine is waiting. The name and value of the option in id
// should not be case sensitive and can include spaces.
//
// The substrings value and name should be avoided in id and x to allow
// unambiguous parsing, for example do not use name = draw value.
func NewSetOption(engine *context.Engine) cmd.Command {
	flags := flag.NewSchema()

	flags.Single("name")
	flags.Variadic("value")

	return cmd.Command{
		Name: "setoption",
		Run: func(interaction cmd.Interaction) error {
			name, value, err := parseSetOptionOptions(interaction.Values)
			if err != nil {
				return err
			}

			return engine.OptionSchema.SetOption(name, value)
		},
		Flags: flags,
	}
}

func parseSetOptionOptions(values flag.Values) (string, []string, error) {
	if !values["name"].Set {
		return "", nil, errors.New("setoption: name flag not found")
	}

	name := values["name"].Value.(string)

	value := []string{}
	if values["value"].Set {
		value = values["value"].Value.([]string)
	}

	return name, value, nil
}
