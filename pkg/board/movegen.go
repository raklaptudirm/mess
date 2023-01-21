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

package board

import (
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/board/move"
	"laptudirm.com/x/mess/pkg/board/move/attacks"
	"laptudirm.com/x/mess/pkg/board/move/castling"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
)

// GenerateMoves generates a move list of all the possible legal moves in
// the current position.
func (b *Board) GenerateMoves(tacticalOnly bool) []move.Move {
	// initialize movegen state
	state := moveGenState{Board: b}
	state.Init(tacticalOnly)

	// append moves to movelist
	if state.CheckN < 2 {
		// moves of other pieces are only possible
		// if the king is not in double check
		state.appendPawnMoves()
		state.appendKnightMoves()
		state.appendBishopMoves()
		state.appendRookMoves()
		state.appendQueenMoves()
	}

	// king moves are always possible
	state.appendKingMoves()

	b.UtilityInfo = &state

	return state.MoveList
}

func (s *moveGenState) appendKingMoves() {
	kingSq := s.Kings[s.SideToMove]

	// king can't move to squares occupied by a friend or sen by an enemy
	kingMoves := attacks.King[kingSq] & s.KingTarget
	s.serializeMoves(s.King, kingSq, kingMoves)

	if !s.TacticalOnly && s.CheckN == 0 {
		// castling can only occur if king is not in check
		s.appendCastlingMoves()
	}
}

func (s *moveGenState) appendKnightMoves() {
	// knights pinned in any direction can't move
	for knights := s.KnightsBB(s.SideToMove) &^ (s.PinnedD | s.PinnedHV); knights != bitboard.Empty; {
		from := knights.Pop()
		knightMoves := attacks.Knight[from] & s.Target
		s.serializeMoves(s.Knight, from, knightMoves)
	}
}

func (s *moveGenState) appendBishopMoves() {
	s.appendBishopTypeMoves(s.Bishop, s.BishopsBB(s.SideToMove))
}

func (s *moveGenState) appendRookMoves() {
	s.appendRookTypeMoves(s.Rook, s.RooksBB(s.SideToMove))
}

func (s *moveGenState) appendQueenMoves() {
	queens := s.QueensBB(s.SideToMove)

	s.appendBishopTypeMoves(s.Queen, queens)
	s.appendRookTypeMoves(s.Queen, queens)
}

// appendBishopTypeMoves appends the moves of any pieces which moves like a bishop.
func (s *moveGenState) appendBishopTypeMoves(bishop piece.Piece, bishops bitboard.Board) {
	bishops &^= s.PinnedHV

	pinned := bishops & s.PinnedD
	for pinned != bitboard.Empty {
		from := pinned.Pop()
		// pinned bishops can only move in their pin-mask
		bishopMoves := attacks.Bishop(from, s.Occupied) & s.Target & s.PinnedD
		s.serializeMoves(bishop, from, bishopMoves)
	}

	unpinned := bishops &^ s.PinnedD
	for unpinned != bitboard.Empty {
		from := unpinned.Pop()
		bishopMoves := attacks.Bishop(from, s.Occupied) & s.Target
		s.serializeMoves(bishop, from, bishopMoves)
	}
}

// appendRookTypeMoves appends the moves of any pieces which moves like a rook.
func (s *moveGenState) appendRookTypeMoves(rook piece.Piece, rooks bitboard.Board) {
	rooks &^= s.PinnedD

	pinned := rooks & s.PinnedHV
	for pinned != bitboard.Empty {
		from := pinned.Pop()
		// pinned rooks can only move in their pin-mask
		rookMoves := attacks.Rook(from, s.Occupied) & s.Target & s.PinnedHV
		s.serializeMoves(rook, from, rookMoves)
	}

	unpinned := rooks &^ s.PinnedHV
	for unpinned != bitboard.Empty {
		from := unpinned.Pop()
		rookMoves := attacks.Rook(from, s.Occupied) & s.Target
		s.serializeMoves(rook, from, rookMoves)
	}
}

func (s *moveGenState) appendPawnMoves() {
	s.appendPawnCaptures()

	pushTarget := s.CheckMask &^ s.Occupied

	// pawns that are pinned diagonally or blocked can't push
	pawnsThatPush := s.PawnsBB(s.SideToMove) &^ s.PinnedD &^ s.Occupied.Down(s.Us)

	pinnedPawnsThatPush := pawnsThatPush & s.PinnedHV
	unpinnedPawnsThatPush := pawnsThatPush &^ s.PinnedHV

	pinnedPawnPushesSingle := attacks.PawnPush(pinnedPawnsThatPush, s.SideToMove) & s.PinnedHV
	unpinnedPawnPushesSingle := attacks.PawnPush(unpinnedPawnsThatPush, s.SideToMove)

	pawnPushesSingle := (pinnedPawnPushesSingle | unpinnedPawnPushesSingle) & pushTarget

	// pawn pushes which result in promotions
	for promotionPawnPushes := pawnPushesSingle & s.PromotionRankBB; promotionPawnPushes != bitboard.Empty; {
		to := promotionPawnPushes.Pop()
		from := to + s.Down
		s.appendPromotions(move.New(from, to, s.Pawn, false), s.SideToMove)
	}

	if s.TacticalOnly {
		// don't append quiet moves
		return
	}

	// pawn pushes that don't result in promotions
	for simplePawnPushes := pawnPushesSingle &^ s.PromotionRankBB; simplePawnPushes != bitboard.Empty; {
		to := simplePawnPushes.Pop()
		from := to + s.Down
		s.AppendMoves(move.New(from, to, s.Pawn, false))
	}

	// double push is the same as a single push on the single pushed pawns
	// pawnPushes single is not used since pawn pushes which don't block
	// checks but whose double pushes do block them are removed
	pawnPushesDouble := (pinnedPawnPushesSingle | unpinnedPawnPushesSingle) & s.DoublePushRankBB
	pawnPushesDouble = attacks.PawnPush(pawnPushesDouble, s.Us) & pushTarget

	// double pawn pushes
	for pawnPushesDouble != bitboard.Empty {
		to := pawnPushesDouble.Pop()
		from := to + s.Down + s.Down
		s.AppendMoves(move.New(from, to, s.Pawn, false))
	}
}

func (s *moveGenState) appendPawnCaptures() {
	const left = -1
	const right = 1

	captureTarget := s.Enemies & s.CheckMask

	// pawns that aren't pinned horizantally or vertically
	// can freely move in diagonal directions
	pawnsThatAttack := s.PawnsBB(s.SideToMove) &^ s.PinnedHV

	unpinnedPawnsThatAttack := pawnsThatAttack &^ s.PinnedD
	pinnedPawnsThatAttack := pawnsThatAttack & s.PinnedD

	pawnAttacksL := attacks.PawnsLeft(unpinnedPawnsThatAttack, s.SideToMove) & captureTarget
	pawnAttacksL |= attacks.PawnsLeft(pinnedPawnsThatAttack, s.SideToMove) & captureTarget & s.PinnedD

	pawnAttacksR := attacks.PawnsRight(unpinnedPawnsThatAttack, s.SideToMove) & captureTarget
	pawnAttacksR |= attacks.PawnsRight(pinnedPawnsThatAttack, s.SideToMove) & captureTarget & s.PinnedD

	simplePawnAttacksL := pawnAttacksL &^ s.PromotionRankBB
	simplePawnAttacksR := pawnAttacksR &^ s.PromotionRankBB

	for simplePawnAttacksL != bitboard.Empty {
		to := simplePawnAttacksL.Pop()
		from := to + s.Down + right
		s.AppendMoves(move.New(from, to, s.Pawn, true))
	}

	for simplePawnAttacksR != bitboard.Empty {
		to := simplePawnAttacksR.Pop()
		from := to + s.Down + left
		s.AppendMoves(move.New(from, to, s.Pawn, true))
	}

	promotionPawnAttacksL := pawnAttacksL & s.PromotionRankBB
	promotionPawnAttacksR := pawnAttacksR & s.PromotionRankBB

	for promotionPawnAttacksL != bitboard.Empty {
		to := promotionPawnAttacksL.Pop()
		from := to + s.Down + right
		s.appendPromotions(move.New(from, to, s.Pawn, true), s.SideToMove)
	}

	for promotionPawnAttacksR != bitboard.Empty {
		to := promotionPawnAttacksR.Pop()
		from := to + s.Down + left
		s.appendPromotions(move.New(from, to, s.Pawn, true), s.SideToMove)
	}

	// append en passant capture
	if s.EnPassantTarget != square.None {
		epPawn := s.EnPassantTarget + s.Down

		epMask := bitboard.Squares[s.EnPassantTarget] | bitboard.Squares[epPawn]
		// check if en-passant leaves king in check
		// this does not account for the double rook pin
		if s.CheckMask&epMask == 0 {
			return
		}

		kingSq := s.Kings[s.SideToMove]
		kingMask := bitboard.Squares[kingSq] & s.EnPassantRankBB

		enemyRooksQueens := (s.RooksBB(s.Them) | s.QueensBB(s.Them)) & s.EnPassantRankBB

		// if king and enemy horizontal sliding piece are on ep rank
		// a horizontal rook pin may be possible so more checks
		isPossiblePin := kingMask != bitboard.Empty && enemyRooksQueens != bitboard.Empty

		for fromBB := attacks.Pawn[s.Them][s.EnPassantTarget] & pawnsThatAttack; fromBB != bitboard.Empty; {
			from := fromBB.Pop()

			// pawn is pinned in other direction
			if s.PinnedD.IsSet(from) && !s.PinnedD.IsSet(s.EnPassantTarget) {
				continue
			}

			// check for horizontal rook pin
			// remove the ep pawn and the enemy pawn from the blocker mask
			// and check if a rook ray from the king hits any rook or queen
			pawnsMask := bitboard.Squares[from] | bitboard.Squares[epPawn]
			if isPossiblePin && attacks.Rook(kingSq, s.Occupied&^pawnsMask)&enemyRooksQueens != 0 {
				break
			}

			s.AppendMoves(move.New(from, s.EnPassantTarget, s.Pawn, true))
		}
	}
}

func (s *moveGenState) appendCastlingMoves() {
	// for each castling move the following things are checked:
	// 1. if castling that side is legal (king and rook haven't moved)
	// 2. if pieces are occupying the space between the king and rook
	// 3. if the squares that the king moves through are seen by the enemy
	// if all the conditions are satisfied then castling that side is legal

	switch s.SideToMove {
	case piece.White:
		if s.CastlingRights&castling.WhiteK != 0 &&
			(s.Occupied|s.SeenByEnemy)&bitboard.F1G1 == bitboard.Empty {
			s.AppendMoves(move.New(square.E1, square.G1, piece.WhiteKing, false))
		}

		if s.CastlingRights&castling.WhiteQ != 0 &&
			s.Occupied&bitboard.B1C1D1 == bitboard.Empty &&
			s.SeenByEnemy&bitboard.C1D1 == bitboard.Empty {
			s.AppendMoves(move.New(square.E1, square.C1, piece.WhiteKing, false))
		}
	case piece.Black:
		if s.CastlingRights&castling.BlackK != 0 &&
			(s.Occupied|s.SeenByEnemy)&bitboard.F8G8 == bitboard.Empty {
			s.AppendMoves(move.New(square.E8, square.G8, piece.BlackKing, false))
		}

		if s.CastlingRights&castling.BlackQ != 0 &&
			s.Occupied&bitboard.B8C8D8 == bitboard.Empty &&
			s.SeenByEnemy&bitboard.C8D8 == bitboard.Empty {
			s.AppendMoves(move.New(square.E8, square.C8, piece.BlackKing, false))
		}
	}
}

// serializeMoves serialized the given move bitboard into the movelist.
func (s *moveGenState) serializeMoves(p piece.Piece, from square.Square, moves bitboard.Board) {
	// append captures
	for captures := moves & s.Enemies; captures != bitboard.Empty; {
		to := captures.Pop()
		s.AppendMoves(move.New(from, to, p, true))
	}

	if s.TacticalOnly {
		// don't serialize quiet moves
		return
	}

	// append quiet moves
	for quiets := moves &^ s.Enemies; quiets != bitboard.Empty; {
		to := quiets.Pop()
		s.AppendMoves(move.New(from, to, p, false))
	}
}

// appendPromotions appends all the different promotion variations of the
// given move to the movelist.
func (s *moveGenState) appendPromotions(m move.Move, c piece.Color) {
	s.AppendMoves(
		m.SetPromotion(piece.New(piece.Queen, c)),
		m.SetPromotion(piece.New(piece.Rook, c)),
		m.SetPromotion(piece.New(piece.Bishop, c)),
		m.SetPromotion(piece.New(piece.Knight, c)),
	)
}
