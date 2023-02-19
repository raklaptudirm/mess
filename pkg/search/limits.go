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

package search

import (
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/search/time"
)

// Limits contains the various limits which decide how long a search can
// run for. It should be passed to the main search function when starting
// a new search.
type Limits struct {
	// search tree limits
	Nodes int
	Depth int

	// TODO: implement searching selected moves
	// Moves []move.Move

	// search time limits
	Infinite        bool
	MoveTime        int
	Time, Increment [piece.ColorN]int
	MovesToGo       int
}

// UpdateLimits updates the search limits while a search is in progress.
// The caller should make sure that a search is indeed in progress before
// calling UpdateLimits.
func (search *Context) UpdateLimits(limits Limits) {
	search.limits = limits // update limits

	switch {
	case limits.Infinite:
		return

	case limits.MoveTime != 0:
		search.time = &time.MoveManager{Duration: limits.MoveTime}

	default:
		search.time = &time.NormalManager{
			Time:      limits.Time,
			Increment: limits.Increment,
			MovesToGo: limits.MovesToGo,
			Us:        search.sideToMove,
		}
	}

	search.time.GetDeadline() // get search deadline
}

// shouldStop checks the various limits provided for the search and
// reports if the search should be stopped at that moment.
func (search *Context) shouldStop() bool {

	// the depth limit is kept up in the iterative deepening
	// loop so it's breaching isn't tested in this function

	switch {
	case search.stopped:
		// search already stopped
		// no checking necessary
		return true

	case search.stats.Nodes&2047 != 0, search.limits.Infinite:
		// only check once every 2048 nodes to prevent
		// spending too much time here

		// if search is infinite never stop

		return false

	case search.stats.Nodes > search.limits.Nodes, search.time.Expired():
		// node limit or time limit crossed
		search.Stop()
		return true

	default:
		// no search stopping clause reached
		return false
	}
}
