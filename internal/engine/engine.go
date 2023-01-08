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
	"laptudirm.com/x/mess/pkg/search"
	"laptudirm.com/x/mess/pkg/uci"
)

// NewClient returns a new uci.Client containing all of the engine's
// supported commands. The commands share a context.Engine among them.
func NewClient() uci.Client {

	client := uci.NewClient()

	// initialize search context
	search := search.NewContext(func(r search.Report) {
		client.Println(r)
	})
	// initialize engine context
	engine := &context.Engine{
		Search: &search,
	}

	// add the engine's commands to the client
	client.AddCommand(cmd.NewD(engine))
	client.AddCommand(cmd.NewGo(engine))
	client.AddCommand(cmd.NewUci(engine))
	client.AddCommand(cmd.NewStop(engine))
	client.AddCommand(cmd.NewBench(engine))
	client.AddCommand(cmd.NewPosition(engine))
	client.AddCommand(cmd.NewUciNewGame(engine))

	return client
}
