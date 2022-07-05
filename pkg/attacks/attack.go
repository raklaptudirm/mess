// Copyright © 2022 Rak Laptudirm <raklaptudirm@gmail.com>
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

package attacks

import (
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/square"
)

// lookup tables for precalculated attack boards of non-sliding pieces
var (
	kingAttacks      [64]bitboard.Board
	knightAttacks    [64]bitboard.Board
	whitePawnAttacks [64]bitboard.Board
	blackPawnAttacks [64]bitboard.Board
)

// init initializes the attack bitboard lookup tables for non-sliding
// pieces by computing the bitboards for each square.
func init() {
	for s := square.A8; s <= square.H1; s++ {
		// compute attack bitboards for current square
		kingAttacks[s] = kingAttacksFrom(s)
		knightAttacks[s] = knightAttacksFrom(s)
		whitePawnAttacks[s] = whitePawnAttacksFrom(s)
		blackPawnAttacks[s] = blackPawnAttacksFrom(s)
	}
}

type board struct {
	origin square.Square
	board  bitboard.Board
}

// addAttack adds the given square to the provided attack bitboard, but
// only if the square lies on the board, i.e, within A8 to H1.
func (b *board) addAttack(fileOffset square.File, rankOffset square.Rank) {
	attackFile := b.origin.File() + fileOffset
	attackRank := b.origin.Rank() + rankOffset

	attackSquare := square.From(attackFile, attackRank)

	switch {
	case attackFile < 0, attackFile > square.FileH, attackRank < 0, attackRank > square.Rank1:
		return
	}

	b.board.Set(attackSquare)
}
