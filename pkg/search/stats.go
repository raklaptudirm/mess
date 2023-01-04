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
	"fmt"
	"time"

	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/board/move"
	"laptudirm.com/x/mess/pkg/search/eval"
)

// Stats stores the search's various statistics.
type Stats struct {
	// time when search started
	SearchStart time.Time

	TTHits int // transposition table hits
	Nodes  int // positions (nodes) searched

	Depth    int // current iterative depth
	SelDepth int // maximum depth reached
}

// GenerateReport generates a statistics report from the current search
// context. It contains all the relevant stats that anyone can need to
// know about a search.
func (search *Context) GenerateReport() Report {
	searchTime := time.Since(search.stats.SearchStart)

	return Report{
		Depth:    search.stats.Depth,
		SelDepth: search.stats.Depth, // TODO: implement seldepth calculation

		Nodes: search.stats.Nodes,
		Nps:   float64(search.stats.Nodes) / util.Max(0.001, searchTime.Seconds()),

		Hashfull: 0, // TODO: implement hashfull calculation

		Time: searchTime,

		Score: search.pvScore,
		PV:    search.pv,
	}
}

// Report represents a report of various statistics about a search.
type Report struct {
	// depth stats
	Depth    int // current id depth
	SelDepth int // max depth reached

	// node stats
	Nodes int
	Nps   float64

	// tt stats
	Hashfull float64

	// search time stats
	Time time.Duration

	// principal variation stats
	Score eval.Eval
	PV    move.Variation
}

// String converts a Report into an UCI compatible info string.
func (report Report) String() string {
	return fmt.Sprintf(
		"info depth %d seldepth %d score %s nodes %d nps %.f hashfull %.f tbhits 0 time %d pv %s",
		report.Depth, report.SelDepth, report.Score, report.Nodes, report.Nps,
		report.Hashfull*1000, // convert fraction to per-mille
		report.Time.Milliseconds(), report.PV,
	)
}
