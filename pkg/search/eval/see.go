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
	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/board/move"
	"laptudirm.com/x/mess/pkg/board/move/attacks"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
)

var seeValue = [piece.TypeN]Eval{
	piece.Pawn:   100,
	piece.Knight: 400,
	piece.Bishop: 400,
	piece.Rook:   600,
	piece.Queen:  1000,
	piece.King:   30000,
}

// SEE performs a static exchange evaluation on the given board starting
// with the given move. It returns true if the capture sequence beats the
// provided threshold, and false otherwise.
func SEE(b *board.Board, m move.Move, threshold Eval) bool {
	// relevant squares
	source, target := m.Source(), m.Target()

	// relevant piece types
	attacker := m.ToPiece().Type()
	victim := util.Ternary(m.IsEnPassant(b.EnPassantTarget), piece.Pawn, b.Position[target].Type())

	balance := seeValue[victim] // win the victim
	if balance < threshold {
		// even if we win the captured piece for free, balance is still
		// less than the threshold, so we can't beat threshold
		return false
	}

	balance -= seeValue[attacker] // lose the attacker
	if balance >= threshold {
		// even if we lose the capturing piece for nothing, balance is
		// still greater than or equal to threshold, so this capture
		// will definitely beat threshold
		return true
	}

	// calculate occupied squares
	occupied := b.ColorBBs[piece.White] | b.ColorBBs[piece.Black]

	// make the capture
	occupied.Unset(source)             // remove the capturing piece
	sideToMove := b.SideToMove.Other() // switch sides after capture

	// calculate attackers to target square
	attackers := attackersTo(b, target, occupied) & occupied

	// calculate ray attackers to reveal x-rays
	diagonal := b.PieceBBs[piece.Bishop] | b.PieceBBs[piece.Queen] // diagonal attackers
	straight := b.PieceBBs[piece.Rook] | b.PieceBBs[piece.Queen]   // straight attackers

	for {
		// calculate friendly attackers
		friends := attackers & b.ColorBBs[sideToMove]
		if friends == bitboard.Empty {
			// no more friendly attackers: end see
			break
		}

		// find least valuable piece to attack with
		for attacker = piece.Pawn; attacker < piece.King; attacker++ {
			if friends&b.PieceBBs[attacker] != bitboard.Empty {
				// piece of this type has been found
				break
			}
		}

		if attacker == piece.King && (attackers&^friends) != bitboard.Empty {
			// king can't capture if other side still has attackers
			break
		}

		// get source square of new attacker
		source = (friends & b.PieceBBs[attacker]).FirstOne()

		// make the capture
		occupied.Unset(source)          // remove the capturing piece
		sideToMove = sideToMove.Other() // switch sides after capture

		// lose the current capturer
		balance = -balance - seeValue[attacker]

		if balance >= threshold {
			// capture is winning even if the current capturer is lost
			// so we can end the exchange evaluation safely
			break
		}

		// add attackers which were hidden by the capturing piece (x rays)
		switch attacker {
		case piece.Pawn, piece.Bishop:
			attackers |= attacks.Bishop(target, occupied) & diagonal
		case piece.Rook:
			attackers |= attacks.Rook(target, occupied) & straight
		case piece.Queen:
			switch {
			case source.File() == target.File(), source.Rank() == target.Rank():
				attackers |= attacks.Rook(target, occupied) & straight
			default:
				attackers |= attacks.Bishop(target, occupied) & diagonal
			}
		}

		// remove attackers which have already captured
		attackers &= occupied
	}

	// at the end of see sideToMove is the side which failed to capture
	// back. The capture sequence is only winning/equal if we are able
	// to capture back.
	return sideToMove != b.SideToMove
}

func attackersTo(b *board.Board, s square.Square, blockers bitboard.Board) bitboard.Board {
	diagonal := b.PieceBBs[piece.Bishop] | b.PieceBBs[piece.Queen]
	straight := b.PieceBBs[piece.Rook] | b.PieceBBs[piece.Queen]

	return attacks.King[s]&b.PieceBBs[piece.King] | // kings
		attacks.Knight[s]&b.PieceBBs[piece.Knight] | // knights
		attacks.Pawn[piece.White][s]&b.PawnsBB(piece.Black) | // black pawns
		attacks.Pawn[piece.Black][s]&b.PawnsBB(piece.White) | // white pawns
		attacks.Bishop(s, blockers)&(diagonal) | // bishops and queens
		attacks.Rook(s, blockers)&(straight) // rooks and queens
}
