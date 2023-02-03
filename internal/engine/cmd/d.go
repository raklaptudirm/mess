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
	"laptudirm.com/x/mess/internal/engine/context"
	"laptudirm.com/x/mess/pkg/uci/cmd"
)

// Custom command d
//
// This command prints out the current position using ascii art, along with
// it's fen string, and zobrist key.
func NewD(engine *context.Engine) cmd.Command {
	return cmd.Command{
		Name: "d",
		Run: func(interaction cmd.Interaction) error {
			// print the current board with ascii art
			interaction.Reply(engine.Search.String())
			return nil
		},
	}
}
