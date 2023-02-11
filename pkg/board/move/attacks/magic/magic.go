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

// Package magic provides reusable utility types and functions that are
// used to generate magic hash tables for any sliding piece.
//
// Blocker masks are uint64 bitboards and therefore there are too many
// permutations to exhaustively calculate. However, the relevant blockers
// for a given square are much fewer in number and can be calculated
// exhaustively. Therefore, we can design a perfect hash function which
// can index every blocker mask relevant to a given square by calculating
// a magic number such that mask * magic >> bits is a perfect contagious
// hash function. It is simplest to calculate this by generating random
// magic numbers and checking if they work.
package magic

import (
	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/board/square"
)

// magicSeeds are optimized prng seeds which generate valid magics fastest
// these values have been taken from the Stockfish chess engine
var magicSeeds = [8]uint64{255, 16645, 15100, 12281, 32803, 55013, 10316, 728}

// NewTable generates a new Magic Hash Table from the given moveFunc. It
// automatically generates the magics and thus is a slow function.
func NewTable(maskN int, moveFunc MoveFunc) *Table {
	var t Table

	// populate table

	var rand util.PRNG

	for s := square.A8; s <= square.H1; s++ {
		magic := &t.Magics[s] // get magic entry for the current square

		// calculate known info
		magic.BlockerMask = moveFunc(s, bitboard.Empty, true) // relevant blocker mask
		bitCount := magic.BlockerMask.Count()
		magic.Shift = uint8(64 - bitCount) // index function shift amount

		// calculate number of permutations of the blocker mask
		permutationsN := 1 << bitCount // 2^bitCount
		permutations := make([]bitboard.Board, permutationsN)

		// initialize blocker mask
		blockers := bitboard.Empty

		// generate all blocker mask permutations and store them, i.e. generate
		// all subsets of the set of the current blocker mask. This is achieved
		// using the Carry-Rippler Trick (https://bit.ly/3XlXipd)
		for index := 0; blockers != bitboard.Empty || index == 0; index++ {
			permutations[index] = blockers
			blockers = (blockers - magic.BlockerMask) & magic.BlockerMask
		}

		// seed random number generator
		rand.Seed(magicSeeds[s.Rank()])

	searchingMagic:
		for { // loop until a valid magic is found

			// initialize table entry
			t.Table[s] = make([]bitboard.Board, maskN)

			// generate a magic candidate
			magic.Number = rand.SparseUint64()

			// try to index all permutations of the blocker
			// mask using the new magic candidate
			for i := 0; i < permutationsN; i++ {
				blockers := permutations[i]

				index := magic.Index(blockers)          // permutation index
				attacks := moveFunc(s, blockers, false) // permutation attack set

				if t.Table[s][index] != bitboard.Empty && t.Table[s][index] != attacks {
					// the calculated index is not empty and the attack sets are not
					// equal: we have a hash collision. Continue searching the magic
					continue searchingMagic
				}

				// no hash collision: store the entry
				t.Table[s][index] = attacks
			}

			// all permutations were successfully stored without hash collisions,
			// so we have found a valid magic and can stop searching for others
			break
		}
	}

	return &t
}

// Table represents a magic hash table.
type Table struct {
	Magics [square.N]Magic            // list of magics for each square
	Table  [square.N][]bitboard.Board // the underlying move-list table
}

// Probe probes the magic hash table for the move bitboard given the
// piece square and blocker mask. It returns the move bitboard.
func (t *Table) Probe(s square.Square, blockerMask bitboard.Board) bitboard.Board {
	return t.Table[s][t.Magics[s].Index(blockerMask)]
}

// Magic represents a single magic entry. Each magic entry is used to
// index all the attack sets for a given square.
type Magic struct {
	Number      uint64         // magic multiplication number
	BlockerMask bitboard.Board // mask of relevant blockers
	Shift       byte           // 64 - no of blocker permutations
}

// Index calculates the index of the given blocker mask given it's magic.
func (magic Magic) Index(blockerMask bitboard.Board) uint64 {
	blockerMask &= magic.BlockerMask // remove irrelevant blockers
	return (uint64(blockerMask) * magic.Number) >> magic.Shift
}

// MoveFunc is a sliding piece's move generation function. It takes the
// piece square, blocker mask, and bool which reports if the function is
// being used to generate blocker masks, so that it can mask out the edge
// bits. It returns a bitboard with all the movable squares set.
type MoveFunc func(square.Square, bitboard.Board, bool) bitboard.Board
