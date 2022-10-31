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

package attacks

import (
	"laptudirm.com/x/mess/pkg/chess/bitboard"
	"laptudirm.com/x/mess/pkg/chess/move/attacks/magic"
	"laptudirm.com/x/mess/pkg/chess/piece"
	"laptudirm.com/x/mess/pkg/chess/square"
)

// lookup tables for precalculated attack boards of non-sliding pieces
var (
	King      [square.N]bitboard.Board
	Knight    [square.N]bitboard.Board
	PawnMoves [piece.NColor][square.N]bitboard.Board
	Pawn      [piece.NColor][square.N]bitboard.Board
)

// magic tables for precalculated attack boards of sliding pieces
var (
	RookTable   magic.Table
	BishopTable magic.Table
)

// init initializes the attack bitboard lookup tables for non-sliding
// pieces by computing the bitboards for each square.
func init() {
	// initialize lookup tables
	for s := square.A8; s <= square.H1; s++ {
		// compute attack bitboards for current square
		King[s] = kingAttacksFrom(s)
		Knight[s] = knightAttacksFrom(s)
		PawnMoves[piece.White][s] = whitePawnMovesFrom(s)
		PawnMoves[piece.Black][s] = blackPawnMovesFrom(s)
		Pawn[piece.White][s] = whitePawnAttacksFrom(s)
		Pawn[piece.Black][s] = blackPawnAttacksFrom(s)
	}

	// initialize magic tables

	RookTable = magic.Table{
		MaxMaskN: 4096, MoveFunc: rook,
	}

	BishopTable = magic.Table{
		MaxMaskN: 512, MoveFunc: bishop,
	}

	RookTable.Populate()
	BishopTable.Populate()
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

	switch {
	case attackFile < 0, attackFile > square.FileH, attackRank < 0, attackRank > square.Rank1:
		return
	}

	attackSquare := square.New(attackFile, attackRank)
	b.board.Set(attackSquare)
}
