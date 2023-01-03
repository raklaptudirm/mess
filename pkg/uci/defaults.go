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

package uci

import (
	"errors"

	"laptudirm.com/x/mess/pkg/uci/cmd"
)

// default/preloaded commands
var cmdIsReady cmd.Command // isready
var cmdQuit cmd.Command    // quit

// errQuit is the error returned to quit the client
var errQuit = errors.New("client: quit")

func init() {
	// This is used to synchronize the engine with the GUI. When the GUI has
	// sent a command or multiple commands that can take some time to complete,
	// this command can be used to wait for the engine to be ready again or to
	// ping the engine to find out if it is still alive. E.g. this should be sent
	// after setting the path to the tablebases as this can take some time.
	//
	// This command is also required once before the engine is asked to do any
	// search to wait for the engine to finish initializing.
	//
	// This command must always be answered with readyok and can be sent also when
	// the engine is calculating in which case the engine should also immediately
	// answer with readyok without stopping the search.
	cmdIsReady = cmd.Command{
		Name: "isready",
		Run: func(interaction cmd.Interaction) error {
			interaction.Reply("readyok")
			return nil
		},
	}

	// quit the program as soon as possible
	cmdQuit = cmd.Command{
		Name: "quit",
		Run: func(cmd.Interaction) error {
			return errQuit
		},
	}
}
