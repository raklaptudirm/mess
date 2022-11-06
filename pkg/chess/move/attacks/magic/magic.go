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
	"laptudirm.com/x/mess/pkg/chess/bitboard"
	"laptudirm.com/x/mess/pkg/chess/square"
	"laptudirm.com/x/mess/internal/util"
)

// magicSeeds are optimized prng seeds which generate valid magics fastest
// these values have been taken from the Stockfish chess engine
var magicSeeds = [8]uint64{255, 16645, 15100, 12281, 32803, 55013, 10316, 728}

// Table represents a magic hash table.
type Table struct {
	MaxMaskN int      // maximum number of blocker mask permutations of piece
	MoveFunc MoveFunc // piece movegen function; see MoveFunc docs for info

	Magics [square.N]Magic            // list of magics for each square
	Table  [square.N][]bitboard.Board // the underlying move-list table
}

// MoveFunc is a sliding piece's move generation function. It takes the
// piece square, blocker mask, and bool which reports if the function is
// being used to generate blocker masks, so that it can mask out the edge
// bits. It returns a bitboard with all the movable squares set.
type MoveFunc func(square.Square, bitboard.Board, bool) bitboard.Board

// Probe probes the magic hash table for the move bitboard given the
// piece square and blocker mask. It returns the move bitboard.
func (t *Table) Probe(s square.Square, blockerMask bitboard.Board) bitboard.Board {
	return t.Table[s][t.Magics[s].Index(blockerMask)]
}

// Populate populates the magic hash table with all the magics and move bitboards.
func (t *Table) Populate() {
	var rand util.PRNG

	for s := square.A8; s <= square.H1; s++ {
		magic := &t.Magics[s]

		magic.BlockerMask = t.MoveFunc(s, bitboard.Empty, true)
		bitCount := magic.BlockerMask.CountBits()
		magic.Shift = uint8(64 - bitCount)

		permutationsN := 1 << bitCount
		permutations := make([]bitboard.Board, permutationsN)
		blockers := bitboard.Empty

		for index := 0; blockers != bitboard.Empty || index == 0; index++ {
			permutations[index] = blockers
			blockers = (blockers - magic.BlockerMask) & magic.BlockerMask
		}

		rand.Seed(magicSeeds[s.Rank()])

	searchingMagic:
		for {
			magic.Number = rand.SparseUint64()

			t.Table[s] = make([]bitboard.Board, t.MaxMaskN)

			for i := 0; i < permutationsN; i++ {
				blockers := permutations[i]
				index := magic.Index(blockers)
				attacks := t.MoveFunc(s, blockers, false)

				if t.Table[s][index] != bitboard.Empty && t.Table[s][index] != attacks {
					continue searchingMagic
				}

				t.Table[s][index] = attacks
			}

			break
		}
	}
}

// Magic represents a single magic entry.
type Magic struct {
	Number      uint64         // magic multiplication number
	BlockerMask bitboard.Board // mask of important blockers
	Shift       byte           // 64 - no of blocker permutations
}

// Index calculates the index of the given blocker mask given it's magic.
func (magic Magic) Index(blockerMask bitboard.Board) uint64 {
	blockerMask &= magic.BlockerMask // keep important bits
	return (uint64(blockerMask) * magic.Number) >> magic.Shift
}
