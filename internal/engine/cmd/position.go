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
	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/uci/cmd"
	"laptudirm.com/x/mess/pkg/uci/flag"
)

// UCI command position [ fen <fenstring> | startpos ] moves <move>...
//
// Set up the position described in fenstring on the internal board and
// play the moves on the internal chess board.
//
// If the game was played from the start position the string startpos will
// be sent
//
// Note: no "new" command is needed. However, if this position is from a
// different game than the last position sent to the engine, the GUI should
// have sent a ucinewgame in-between.
func NewPosition(engine *context.Engine) cmd.Command {
	schema := flag.NewSchema()

	// base position
	schema.Array("fen", len(board.StartFEN))
	schema.Button("startpos")

	// moves played on base position
	schema.Variadic("moves")

	return cmd.Command{
		Name: "position",
		Run: func(interaction cmd.Interaction) error {
			// parse flags into a board.Board
			board, err := parsePositionFlags(interaction.Values)
			if err != nil {
				return err
			}

			// update search board
			engine.Search.Board = board

			return nil
		},
		Flags: schema,
	}
}

// parsePositionFlags parses the position data from the given flags.
func parsePositionFlags(values flag.Values) (*board.Board, error) {
	// set up new position here
	var b *board.Board

	// parse base position
	switch {
	// only one of the base position descriptors should be set
	case values["startpos"].Set && values["fen"].Set:
		return nil, errors.New("position: both startpos and fen flags found")

	case values["startpos"].Set:
		// starting position
		b = board.NewBoard(board.StartFEN)

	case values["fen"].Set:
		// parse fen string for base position
		b = board.NewBoard(values["fen"].Value.([]string))

	default:
		// one of fen or startpos have to be there
		return nil, errors.New("position: no startpos or fen option")
	}

	// check if any moves are played on the base position
	if values["moves"].Set {
		// play the provided moves on the board
		moves := values["moves"].Value.([]string)
		for _, m := range moves {
			b.MakeMove(b.NewMoveFromString(m))
		}
	}

	return b, nil
}
