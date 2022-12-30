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
	"errors"
	"strings"

	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/search"
	"laptudirm.com/x/mess/pkg/uci/cmd"
	"laptudirm.com/x/mess/pkg/uci/flag"
)

func newCmdUciNewGame(engine Engine) cmd.Command {
	return cmd.Command{
		Name: "ucinewgame",
		Run: func(interaction cmd.Interaction) error {
			// new context for new game
			*engine.search = search.NewContext(board.NewBoard(startpos))
			return nil
		},
	}
}

// fen string for the starting position
var startpos = strings.Fields("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

// parsePositionFlags parses the position data from the given flags.
func parsePositionFlags(values flag.Values) (*board.Board, error) {
	// set up new position here
	var b *board.Board

	switch {
	case values["startpos"].Set:
		b = board.NewBoard(startpos)
	case values["fen"].Set:
		b = board.NewBoard(values["fen"].Value.([]string))
	default:
		// one of fen or startpos have to be there
		return nil, errors.New("position: no startpos or fen option")
	}

	if values["moves"].Set {
		// play the provided moves on the board
		moves := values["moves"].Value.([]string)
		for _, m := range moves {
			b.MakeMove(b.NewMoveFromString(m))
		}
	}

	return b, nil
}

func newCmdPosition(engine Engine) cmd.Command {
	schema := flag.NewSchema()

	// position
	schema.Array("fen", len(startpos))
	schema.Button("startpos")

	// moves on position
	schema.Variadic("moves")

	return cmd.Command{
		Name: "position",
		Run: func(interaction cmd.Interaction) error {
			board, err := parsePositionFlags(interaction.Values)
			if err != nil {
				return err
			}

			engine.search.Board = board
			return nil
		},
		Flags: schema,
	}
}
