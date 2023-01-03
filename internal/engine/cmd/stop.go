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
	"errors"

	"laptudirm.com/x/mess/internal/engine/context"
	"laptudirm.com/x/mess/pkg/uci/cmd"
)

// UCI command stop
//
// Stop calculating as soon as possible.
func NewStop(engine *context.Engine) cmd.Command {
	return cmd.Command{
		Name: "stop",
		Run: func(interaction cmd.Interaction) error {
			// check if any search is ongoing
			if !engine.Search.InProgress() {
				return errors.New("stop: no search ongoing")
			}

			// stop the search
			engine.Search.Stop()
			return nil
		},
	}
}
