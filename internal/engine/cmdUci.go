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

package engine

import (
	"laptudirm.com/x/mess/internal/build"
	"laptudirm.com/x/mess/pkg/uci/cmd"
)

func newCmdUci() cmd.Command {
	return cmd.Command{
		Name: "uci",
		Run: func(interaction cmd.Interaction) error {
			interaction.Replyf("id name Mess %s", build.Version)
			interaction.Reply("id author Rak Laptudirm")
			interaction.Reply("uciok")
			return nil
		},
	}
}
