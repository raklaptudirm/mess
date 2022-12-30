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
	"math"
	"strconv"

	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/search"
	"laptudirm.com/x/mess/pkg/search/time"
	"laptudirm.com/x/mess/pkg/uci/cmd"
	"laptudirm.com/x/mess/pkg/uci/flag"
)

func parseSearchLimits(engine Engine, values flag.Values) (search.Limits, error) {
	var limits search.Limits

	limits.Depth = search.MaxDepth
	if depth := values["depth"]; depth.Set {
		d, _ := strconv.Atoi(depth.Value.(string))
		limits.Depth = d
	}

	limits.Nodes = math.MaxInt32
	if nodes := values["nodes"]; nodes.Set {
		n, _ := strconv.Atoi(nodes.Value.(string))
		limits.Nodes = n
	}

	switch {
	case values["movetime"].Set:
		t, err := strconv.Atoi(values["movetime"].Value.(string))
		if err != nil {
			return limits, err
		}

		limits.Time = &time.MoveManager{Duration: t}

	case values["wtime"].Set:
		tc := &time.NormalManager{Us: engine.search.Board.SideToMove}

		var err error

		tc.Time[piece.White], err = strconv.Atoi(values["wtime"].Value.(string))
		if err != nil {
			return limits, err
		}

		tc.Time[piece.Black], err = strconv.Atoi(values["btime"].Value.(string))
		if err != nil {
			return limits, err
		}

		if values["winc"].Set {
			tc.Increment[piece.White], err = strconv.Atoi(values["winc"].Value.(string))
			if err != nil {
				return limits, err
			}

			tc.Increment[piece.Black], err = strconv.Atoi(values["binc"].Value.(string))
			if err != nil {
				return limits, err
			}
		}

		if values["movestogo"].Set {
			tc.MovesToGo, err = strconv.Atoi(values["movestogo"].Value.(string))
			if err != nil {
				return limits, err
			}
		}

		limits.Time = tc

	case values["infinite"].Set:
		limits.Infinite = true
		fallthrough

	default:
		limits.Time = &time.MoveManager{Duration: math.MaxInt32}
	}

	return limits, nil
}

func newCmdGo(engine Engine) cmd.Command {
	schema := flag.NewSchema()

	// schema.Variadic("searchmoves")
	// schema.Button("ponder")
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
			if engine.search.InProgress() {
				// search already ongoing
				return errors.New("error: search currently in progress")
			}

			limits, err := parseSearchLimits(engine, interaction.Values)
			if err != nil {
				return err
			}

			pv, _, err := engine.search.Search(limits)
			if err != nil {
				return err
			}

			interaction.Replyf("bestmove %s ponder %s", pv.Move(0), pv.Move(1))
			return nil
		},
		// execution of this function should not block the prompt loop
		Parallel: true,
		Flags:    schema,
	}
}

func newCmdStop(engine Engine) cmd.Command {
	return cmd.Command{
		Name: "stop",
		Run: func(interaction cmd.Interaction) error {
			// stop the search
			engine.search.Stop()
			return nil
		},
	}
}
