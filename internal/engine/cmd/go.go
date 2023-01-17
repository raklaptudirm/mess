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
	"errors"
	"math"
	"strconv"

	"laptudirm.com/x/mess/internal/engine/context"
	"laptudirm.com/x/mess/pkg/board/move"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/search"
	"laptudirm.com/x/mess/pkg/search/time"
	"laptudirm.com/x/mess/pkg/uci/cmd"
	"laptudirm.com/x/mess/pkg/uci/flag"
)

// UCI command go [flags]
//
// Start calculating on the current position set up with the position
// command.
//
// There are a number of commands that can follow this command, all will be
// sent in the same string. If one command is not sent its value should be
// interpreted as it would not influence the search.
//
// searchmoves move... [not implemented]
//
//	restrict search to this moves only
//	Example: After position startpos and go infinite searchmoves e2e4
//	d2d4 the engine should only search the two moves e2e4 and d2d4 in
//	the initial position.
//
// ponder [not implemented]
//
//	Start searching in pondering mode.
//	Do not exit the search in ponder mode, even if it's mate!
//	This means that the last move sent in in the position string is the
//	ponder move. The engine can do what it wants to do, but after a
//	ponderhit command it should execute the suggested move to ponder on.
//	This means that the ponder move sent by the GUI can be interpreted as
//	a recommendation about which move to ponder. However, if the engine
//	decides to ponder on a different move, it should not display any mainlines
//	as they are likely to be misinterpreted by the GUI because the GUI expects
//	the engine to ponder on the suggested move.
//
// wtime x
//
//	White has x msec left on the clock
//
// btime x
//
//	black has x msec left on the clock
//
// winc x
//
//	white increment per move in milliseconds if x > 0
//
// binc x
//
//	black increment per move in milliseconds if x > 0
//
// movestogo x
//
//	there are x moves to the next time control,
//	this will only be sent if x > 0,
//	if you don't get this and get the wtime and btime it's sudden death
//
// depth x
//
//	search x plies only.
//
// nodes x
//
//	search x nodes only,
//
// mate x [not implemented]
//
//	search for a mate in x moves
//
// movetime x
//
//	search exactly x milliseconds
//
// infinite
//
//	search until the stop command. Do not exit the search without being told so in this mode!
func NewGo(engine *context.Engine) cmd.Command {
	schema := flag.NewSchema()

	// schema.Variadic("searchmoves")
	schema.Button("ponder")
	schema.Single("wtime")
	schema.Single("btime")
	schema.Single("winc")
	schema.Single("binc")
	schema.Single("movestogo")
	schema.Single("depth")
	schema.Single("nodes")
	// schema.Single("mate")
	schema.Single("movetime")
	schema.Button("infinite")

	return cmd.Command{
		Name: "go",
		Run: func(interaction cmd.Interaction) error {
			if engine.Search.InProgress() {
				// search already ongoing
				return errors.New("error: search currently in progress")
			}

			// parse search limits from flags
			limits, err := parseSearchLimits(engine, interaction.Values)
			if err != nil {
				return err
			}

			// ponder search
			if interaction.Values["ponder"].Set {
				if !engine.Options.Ponder {
					return errors.New("go ponder: pondering is disabled")
				}

				engine.Pondering = true

				// store search limits for later
				engine.PonderLimits = limits

				// for now, start an infinite search
				limits = search.Limits{
					Depth:    search.MaxDepth,
					Nodes:    math.MaxInt,
					Infinite: true,
					Time:     &time.MoveManager{Duration: math.MaxInt32},
				}
			}

			// start searching
			engine.Searching = true
			// search in a separate thread that we don't block the repl
			go func() {
				defer func() {
					// set search booleans to false
					// since the search has ended
					engine.Searching = false
					engine.Pondering = false
				}()

				// start search
				pv, _, err := engine.Search.Search(limits)
				if err != nil {
					interaction.Reply(err)
					return
				}

				if bestMove, ponderMove := pv.Move(0), pv.Move(1); ponderMove == move.Null {
					// just print bestmove since pondermove is null
					interaction.Replyf("bestmove %s", bestMove)
				} else {
					// print bestmove and pondermove
					interaction.Replyf("bestmove %s ponder %s", bestMove, ponderMove)
				}
			}()

			return nil
		},

		Flags: schema,
	}
}

// parseSearchLimits parses the search flags and returns the limits.
func parseSearchLimits(engine *context.Engine, values flag.Values) (search.Limits, error) {
	var limits search.Limits

	// depth limit (default MaxDepth)
	limits.Depth = search.MaxDepth
	if depth := values["depth"]; depth.Set {
		d, _ := strconv.Atoi(depth.Value.(string))
		limits.Depth = d
	}

	// node limit (default MaxInt)
	limits.Nodes = math.MaxInt
	if nodes := values["nodes"]; nodes.Set {
		n, _ := strconv.Atoi(nodes.Value.(string))
		limits.Nodes = n
	}

	// check if wtime-btime controls are set
	timeSet := false
	if values["wtime"].Set || values["btime"].Set {
		if !values["wtime"].Set || !values["btime"].Set {
			return limits, errors.New("go: both wtime and btime should be set")
		}

		timeSet = true
	}

	switch {
	// only one of base time controls should be set
	case (values["movetime"].Set && values["infinite"].Set),
		(values["infinite"].Set && timeSet),
		(timeSet && values["movetime"].Set):

		return limits, errors.New("go: multiple time controls set")

	case values["movetime"].Set:
		// parse movetime
		t, err := strconv.Atoi(values["movetime"].Value.(string))
		if err != nil {
			return limits, err
		}

		limits.Time = &time.MoveManager{Duration: t}

	case timeSet:
		tc := &time.NormalManager{Us: engine.Search.Board.SideToMove}

		var err error

		// parse times

		tc.Time[piece.White], err = strconv.Atoi(values["wtime"].Value.(string))
		if err != nil {
			return limits, err
		}

		tc.Time[piece.Black], err = strconv.Atoi(values["btime"].Value.(string))
		if err != nil {
			return limits, err
		}

		if values["winc"].Set || values["binc"].Set {
			// if one is set, both should be set
			if !values["winc"].Set || !values["binc"].Set {
				return limits, errors.New("go: both winc and binc should be set")
			}

			// parse increments

			tc.Increment[piece.White], err = strconv.Atoi(values["winc"].Value.(string))
			if err != nil {
				return limits, err
			}

			tc.Increment[piece.Black], err = strconv.Atoi(values["binc"].Value.(string))
			if err != nil {
				return limits, err
			}
		}

		// parse moves to next time control
		if values["movestogo"].Set {
			tc.MovesToGo, err = strconv.Atoi(values["movestogo"].Value.(string))
			if err != nil {
				return limits, err
			}
		}

		limits.Time = tc

	case values["infinite"].Set:
		limits.Infinite = true

		// unnecessary, but keep as failsafe
		fallthrough

	default:
		// movetime manager with a very large value: effectively infinite
		limits.Time = &time.MoveManager{Duration: math.MaxInt32}
	}

	return limits, nil
}
