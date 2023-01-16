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
	"laptudirm.com/x/mess/pkg/search"
	"laptudirm.com/x/mess/pkg/uci/cmd"
)

// UCI command ucinewgame
//
// This is sent to the engine when the next search (started with position
// and go) will be from a different game. This can be a new game the
// engine should play or a new game it should analyze but also the next
// position from a test suite with positions only.
//
// [this clause is ignored and Mess depends on this command]
// If the GUI hasn't sent a ucinewgame before the first position command,
// the engine shouldn't expect any further ucinewgame commands as the GUI
// is probably not supporting the ucinewgame command. So the engine should
// not rely on this command even though all new GUIs should support it.
//
// As the engine's reaction to ucinewgame can take some time the GUI should
// always send isready after ucinewgame to wait for the engine to finish its
// operation.
func NewUciNewGame(engine *context.Engine) cmd.Command {
	return cmd.Command{
		Name: "ucinewgame",
		Run: func(interaction cmd.Interaction) error {
			// new context for new game
			engine.Search = search.NewContext(func(r search.Report) {
				interaction.Reply(r)
			}, engine.Options.Hash)
			return nil
		},
	}
}
