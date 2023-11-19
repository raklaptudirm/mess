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

// Package tt implements a transposition table which is used to cache
// results from previous searches of a position to make search more
// efficient. It stores things like the score and pv move.
package tt

import (
	"math/bits"
	"unsafe"

	"laptudirm.com/x/mess/pkg/board/move"
	"laptudirm.com/x/mess/pkg/board/zobrist"
	"laptudirm.com/x/mess/pkg/search/eval"
)

// EntrySize stores the size in bytes of a tt entry.
var EntrySize = int(unsafe.Sizeof(Entry{}))

// NewTable creates a new transposition table with a size equal to or
// less than the given number of megabytes.
func NewTable(mbs int) *Table {
	// compute table size (number of entries)
	size := (mbs * 1024 * 1024) / EntrySize

	return &Table{
		table: make([]Entry, size),
		size:  size,
	}
}

// Table represents a transposition table.
type Table struct {
	table []Entry // hash table
	size  int     // table size
	epoch uint8   // table epoch
}

func (tt *Table) Clear() {
	clear(tt.table)
}

// NextEpoch increases the epoch number of the given tt.
func (tt *Table) NextEpoch() {
	tt.epoch++
}

// Resize resizes the given transposition table to the new size. The
// entries are copied from the old table to the new one. If the new table
// is smaller, some entries are discarded.
func (tt *Table) Resize(mbs int) {
	// compute new table size (number of entries)
	size := (mbs * 1024 * 1024) / EntrySize

	// create table with new size
	newTable := make([]Entry, size)

	// copy old elements
	copy(newTable, tt.table)

	// replace old table
	*tt = Table{
		table: newTable,
		size:  size,
	}
}

// Store puts the given data into the transposition table.
func (tt *Table) Store(entry Entry) {
	target := tt.fetch(entry.Hash)
	entry.epoch = tt.epoch

	// replace only if the new data has an equal or higher quality.
	if entry.quality() >= target.quality() {
		*target = entry
	}
}

// Probe fetches the data associated with the given zobrist key from the
// transposition table. It returns the fetched data and whether it is
// usable or not. It guards against hash collisions and empty entries.
// If the bool value is false, the entry should not be use for anything.
func (tt *Table) Probe(hash zobrist.Key) (Entry, bool) {
	entry := *tt.fetch(hash)
	return entry, entry.Type != NoEntry && entry.Hash == hash
}

// fetch returns a pointer pointing to the tt entry of the given hash.
func (tt *Table) fetch(hash zobrist.Key) *Entry {
	return &tt.table[tt.indexOf(hash)]
}

// indexOf is the hash function used by the transposition table.
func (tt *Table) indexOf(hash zobrist.Key) uint {
	// fast indexing function from Daniel Lemire's blog post
	// https://lemire.me/blog/2016/06/27/a-fast-alternative-to-the-modulo-reduction/
	index, _ := bits.Mul(uint(hash), uint(tt.size))
	return index
}

// Entry represents a transposition table entry.
type Entry struct {
	// complete hash of the position; to guard against
	// transposition table key collisions
	Hash zobrist.Key

	// best move in the position
	// used during iterative deepening as pv move
	Move move.Move

	// evaluation info
	Value Eval      // value of this position
	Type  EntryType // bound type of the value

	// entry metadata
	Depth uint8 // depth the position was searched to
	epoch uint8 // epoch/age of the entry from creation
}

// quality returns a quality measure of the given tt entry which will be
// used to determine whether a tt entry should be overwritten or not.
func (entry *Entry) quality() uint8 {
	return entry.epoch + entry.Depth/3
}

// EntryType represents the type of a transposition table entry's
// value, whether it exists, it is upper bound, lower bound, or exact.
type EntryType uint8

// constants representing various transposition table entry types
const (
	NoEntry EntryType = iota // no entry exists

	ExactEntry // the value is an exact score
	LowerBound // the value is a lower bound on the exact score
	UpperBound // the value is an upper bound on the exact score
)

// EvalFrom converts a given mate score from "n plys till mate from root"
// to "n plys till mate from current position" so that it is reusable in
// other positions with greater or lesser depth.
func EvalFrom(score eval.Eval, plys int) Eval {
	switch {
	case score > eval.WinInMaxPly:
		score += eval.Eval(plys)
	case score < eval.LoseInMaxPly:
		score -= eval.Eval(plys)
	}

	return Eval(score)
}

// Eval represents the evaluation of a transposition table entry. For mate
// scores, it stores "n plys till mate from current position" instead of the
// standard "n plys till mate from root" used in search.
type Eval eval.Eval

// Eval converts transposition table entry scores from "n plys to mate
// from current position" to "n plys till mate from root" which is the
// format used during search.
func (e Eval) Eval(plys int) eval.Eval {
	score := eval.Eval(e)

	// checkmate scores need to be changed from
	switch {
	case score > eval.WinInMaxPly:
		score -= eval.Eval(plys)
	case score < eval.LoseInMaxPly:
		score += eval.Eval(plys)
	}

	return score
}
