// Copyright © 2022 Rak Laptudirm <rak@laptudirm.com>
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
	"laptudirm.com/x/mess/internal/build"
	"laptudirm.com/x/mess/internal/engine/context"
	"laptudirm.com/x/mess/pkg/uci/cmd"
)

// UCI command uci:
//
// Tells engine to use the uci (universal chess interface), this will be
// sent once as a first command after program boot to tell the engine to
// switch to uci mode.
//
// After receiving the uci command the engine must identify itself with the
// id command and send the option commands to tell the GUI which engine
// settings the engine supports if any.
//
// After that the engine should send uciok to acknowledge the uci mode. If
// no uciok is sent within a certain time period, the engine task will be
// killed by the GUI.
func NewUci(engine *context.Engine) cmd.Command {
	return cmd.Command{
		Name: "uci",
		Run: func(interaction cmd.Interaction) error {

			// identify engine
			interaction.Replyf("id name Mess %s", build.Version)
			interaction.Reply("id author Rak Laptudirm")

			// declare uci support
			interaction.Reply("uciok")

			return nil
		},
	}
}
