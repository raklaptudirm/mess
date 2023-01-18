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

package eval

import (
	"laptudirm.com/x/mess/pkg/board/mailbox"
	"laptudirm.com/x/mess/pkg/board/move"
	"laptudirm.com/x/mess/pkg/board/piece"
)

// MoveFunc represents a move evaluation function.
type MoveFunc func(move.Move) Move

// Move represents the evaluation of a move.
type Move uint16

// constants representing move evaluations
const (
	PVMove Move = 60000

	MvvLvaOffset Move = 50000

	KillerMove1 Move = 42000
	KillerMove2 Move = 41000

	DefaultMove Move = 0
)

// MvvLva table taken from Blunder
// TODO: get better scores; may be redundant after see
// score = MvvLvaOffset + MvvLva[victim][attacker]
var MvvLva = [piece.TypeN][piece.TypeN]Move{
	// No piece (-) column is used as promotion scores
	// Attackers:  -   P   N   B   R   Q   K
	piece.Pawn:   {16, 15, 14, 13, 12, 11, 10},
	piece.Knight: {26, 25, 24, 23, 22, 21, 20},
	piece.Bishop: {36, 35, 34, 33, 32, 31, 30},
	piece.Rook:   {46, 45, 44, 43, 42, 41, 40},
	piece.Queen:  {56, 55, 54, 53, 52, 51, 50},
}

// OfMove is a move evaluation function which returns a move func which can
// be used for ordering moves. It takes the position and pv move as input.
func OfMove(info ModeEvalInfo) MoveFunc {
	return func(m move.Move) Move {
		switch {
		case m == info.PVMove:
			// pv move from previous iteration is most likely
			// to be the best move in the position
			return PVMove

		// captures and promotions
		case m.IsCapture(), m.IsPromotion():
			victim := info.Board[m.Target()].Type()
			attacker := m.FromPiece().Type() // piece.NoType for promotions

			// a less valuable piece capturing a more valuable
			// piece is very likely to be a good move
			return MvvLvaOffset + MvvLva[victim][attacker]

		// killer moves
		case m == info.Killers[0]:
			return KillerMove1
		case m == info.Killers[1]:
			return KillerMove2

		default:
			// default move evaluation
			return DefaultMove
		}
	}
}

// MoveEvalInfo stores the various search specific information which are
// required by the move ordering function.
type ModeEvalInfo struct {
	Board   *mailbox.Board
	PVMove  move.Move
	Killers [2]move.Move
}
