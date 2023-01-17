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

package engine

import (
	"laptudirm.com/x/mess/internal/engine/cmd"
	"laptudirm.com/x/mess/internal/engine/context"
	"laptudirm.com/x/mess/internal/engine/options"
	"laptudirm.com/x/mess/pkg/uci"
	"laptudirm.com/x/mess/pkg/uci/option"
)

// NewClient returns a new uci.Client containing all of the engine's
// supported commands. The commands share a context.Engine among them.
func NewClient() (uci.Client, error) {

	// initialize engine context
	engine := &context.Engine{}

	// add uci commands to engine
	engine.Client = uci.NewClient()
	engine.Client.AddCommand(cmd.NewD(engine))
	engine.Client.AddCommand(cmd.NewGo(engine))
	engine.Client.AddCommand(cmd.NewUci(engine))
	engine.Client.AddCommand(cmd.NewStop(engine))
	engine.Client.AddCommand(cmd.NewBench(engine))
	engine.Client.AddCommand(cmd.NewPosition(engine))
	engine.Client.AddCommand(cmd.NewSetOption(engine))
	engine.Client.AddCommand(cmd.NewPonderHit(engine))
	engine.Client.AddCommand(cmd.NewUciNewGame(engine))

	// run ucinewgame to initialize position
	if err := engine.Client.Run("ucinewgame"); err != nil {
		return uci.Client{}, err
	}

	// add uci options to engine
	engine.OptionSchema = option.NewSchema()
	engine.OptionSchema.AddOption("Hash", options.NewHash(engine))
	engine.OptionSchema.AddOption("Ponder", options.NewPonder(engine))
	engine.OptionSchema.AddOption("Threads", options.NewThreads(engine))

	// initialize options
	if err := engine.OptionSchema.SetDefaults(); err != nil {
		return uci.Client{}, err
	}

	return engine.Client, nil
}
