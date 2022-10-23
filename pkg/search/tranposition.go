package search

import (
	"laptudirm.com/x/mess/pkg/evaluation"
	"laptudirm.com/x/mess/pkg/zobrist"
)

type transpositionTable struct {
	table []ttEntry
	size  uint64
}

func NewTT(size int) *transpositionTable {
	return &transpositionTable{
		table: make([]ttEntry, size),
		size:  uint64(size),
	}
}

// Put puts the given data into the transposition table.
func (tt transpositionTable) Put(hash zobrist.Key, plys, depth int, score, alpha, beta evaluation.Rel) {
	value := score
	// checkmate values need to be changes from "n plys till mate from root" to
	// "n plys till mate from current position" to be reusable at a greater depth
	switch {
	case value > evaluation.WinInMaxPly:
		value += evaluation.Rel(plys)
	case value < evaluation.LoseInMaxPly:
		value -= evaluation.Rel(plys)
	}

	tt.table[uint64(hash)%tt.size] = ttEntry{
		value: value,
		depth: depth,
		eType: entryType(score, alpha, beta),
	}
}

// Get fetches the data associated with the given zobrist key from the
// transposition table and verifies if it is usable in the given context.
func (tt transpositionTable) Get(hash zobrist.Key, plys, depth int) (ttEntry, bool) {
	if entry := tt.table[uint64(hash)%tt.size]; entry.eType != none && entry.depth >= depth {
		// checkmate scores need to be changed from "n plys to mate from current position"
		// to "n plys till mate from root" which is the format used during comparison
		switch {
		case entry.value > evaluation.WinInMaxPly:
			entry.value -= evaluation.Rel(plys)
		case entry.value < evaluation.LoseInMaxPly:
			entry.value += evaluation.Rel(plys)
		}

		return entry, true
	}

	return ttEntry{}, false
}

type ttEntry struct {
	value evaluation.Rel
	depth int
	eType tteEntryType
}

func entryType(score, alpha, beta evaluation.Rel) tteEntryType {
	switch {
	case score <= alpha:
		return upperBound
	case score >= beta:
		return lowerBound
	default:
		return exact
	}
}

type tteEntryType int

// constants representing various transposition table entry types
const (
	none tteEntryType = iota
	exact
	lowerBound
	upperBound
)
