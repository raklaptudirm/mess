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

package transposition

import (
	"reflect"

	"laptudirm.com/x/mess/pkg/search/evaluation"
	"laptudirm.com/x/mess/pkg/zobrist"
)

type Table struct {
	table []TableEntry
	size  int
}

func NewTable(mbs int) *Table {
	size := (mbs * 1024 * 1024) / int(reflect.TypeOf(TableEntry{}).Size())
	return &Table{
		table: make([]TableEntry, size),
		size:  size,
	}
}

// Put puts the given data into the transposition table.
func (tt *Table) Put(hash zobrist.Key, entry TableEntry) {
	entry.Hash = hash
	*tt.fetch(hash) = entry
}

// Get fetches the data associated with the given zobrist key from the
// transposition table.
func (tt *Table) Get(hash zobrist.Key) (TableEntry, bool) {
	entry := *tt.fetch(hash)
	return entry, entry.Type != NoEntry && entry.Hash == hash
}

func (tt *Table) fetch(hash zobrist.Key) *TableEntry {
	return &tt.table[tt.indexOf(hash)]
}

func (tt *Table) indexOf(hash zobrist.Key) int {
	return int(uint64(hash) % uint64(tt.size))
}

type TableEntry struct {
	Hash  zobrist.Key
	Value Eval
	Depth int
	Type  TableEntryType
}

type TableEntryType int

// constants representing various transposition table entry types
const (
	NoEntry TableEntryType = iota
	ExactEntry
	LowerBound
	UpperBound
)

func EvalFrom(eval evaluation.Rel, plys int) Eval {
	// checkmate values need to be changes from "n plys till mate from root" to
	// "n plys till mate from current position" to be reusable at a greater depth
	switch {
	case eval > evaluation.WinInMaxPly:
		eval += evaluation.Rel(plys)
	case eval < evaluation.LoseInMaxPly:
		eval -= evaluation.Rel(plys)
	}

	return Eval(eval)
}

type Eval evaluation.Rel

func (e Eval) Rel(plys int) evaluation.Rel {
	eval := evaluation.Rel(e)

	// checkmate scores need to be changed from "n plys to mate from current position"
	// to "n plys till mate from root" which is the format used during comparison
	switch {
	case eval > evaluation.WinInMaxPly:
		eval -= evaluation.Rel(plys)
	case eval < evaluation.LoseInMaxPly:
		eval += evaluation.Rel(plys)
	}

	return eval
}
