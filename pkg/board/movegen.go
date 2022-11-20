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
func (b *Board) GenerateMoves() []move.Move {
	// initialize the utility bitboards like check-mask and pin-masks
	b.InitBitboards()

	// 31 is the average number of chess moveList in a position
	// source: https://chess.stackexchange.com/a/24325/33336
	moveList := make([]move.Move, 0, 31)

	b.appendKingMoves(&moveList)

	if b.CheckN >= 2 {
		// only king moves are possible in double check
		return moveList
	}

	// moves of other pieces
	b.appendKnightMoves(&moveList)
	b.appendBishopMoves(&moveList)
	b.appendRookMoves(&moveList)
	b.appendQueenMoves(&moveList)
	b.appendPawnMoves(&moveList)

	return moveList
}

func (b *Board) appendKingMoves(moveList *[]move.Move) {
	king := piece.New(piece.King, b.SideToMove)
	kingSq := b.Kings[b.SideToMove]

	// king can't move to squares occupied by a friend or sen by an enemy
	kingMoves := attacks.King[kingSq] &^ (b.Friends | b.SeenByEnemy)
	b.serializeMoves(moveList, king, kingSq, kingMoves)

	if b.CheckN == 0 {
		// castling can only occur if king is not in check
		b.appendCastlingMoves(moveList)
	}
}

func (b *Board) appendKnightMoves(moveList *[]move.Move) {
	knight := piece.New(piece.Knight, b.SideToMove)
	// knights pinned in any direction can't move
	for knights := b.Knights(b.SideToMove) &^ (b.PinnedD | b.PinnedHV); knights != bitboard.Empty; {
		from := knights.Pop()
		knightMoves := attacks.Knight[from] & b.Target
		b.serializeMoves(moveList, knight, from, knightMoves)
	}
}

func (b *Board) appendBishopMoves(moveList *[]move.Move) {
	b.appendBishopTypeMoves(moveList, piece.New(piece.Bishop, b.SideToMove), b.Bishops(b.SideToMove))
}

func (b *Board) appendRookMoves(moveList *[]move.Move) {
	b.appendRookTypeMoves(moveList, piece.New(piece.Rook, b.SideToMove), b.Rooks(b.SideToMove))
}

func (b *Board) appendQueenMoves(moveList *[]move.Move) {
	queen := piece.New(piece.Queen, b.SideToMove)
	queens := b.Queens(b.SideToMove)

	b.appendBishopTypeMoves(moveList, queen, queens)
	b.appendRookTypeMoves(moveList, queen, queens)
}

// appendBishopTypeMoves appends the moves of any pieces which moves like a bishop.
func (b *Board) appendBishopTypeMoves(moveList *[]move.Move, bishop piece.Piece, bishops bitboard.Board) {
	bishops &^= b.PinnedHV

	pinned := bishops & b.PinnedD
	for pinned != bitboard.Empty {
		from := pinned.Pop()
		// pinned bishops can only move in their pin-mask
		bishopMoves := attacks.Bishop(from, b.Occupied) & b.Target & b.PinnedD
		b.serializeMoves(moveList, bishop, from, bishopMoves)
	}

	unpinned := bishops &^ b.PinnedD
	for unpinned != bitboard.Empty {
		from := unpinned.Pop()
		bishopMoves := attacks.Bishop(from, b.Occupied) & b.Target
		b.serializeMoves(moveList, bishop, from, bishopMoves)
	}
}

// appendRookTypeMoves appends the moves of any pieces which moves like a rook.
func (b *Board) appendRookTypeMoves(moveList *[]move.Move, rook piece.Piece, rooks bitboard.Board) {
	rooks &^= b.PinnedD

	pinned := rooks & b.PinnedHV
	for pinned != bitboard.Empty {
		from := pinned.Pop()
		// pinned rooks can only move in their pin-mask
		rookMoves := attacks.Rook(from, b.Occupied) & b.Target & b.PinnedHV
		b.serializeMoves(moveList, rook, from, rookMoves)
	}

	unpinned := rooks &^ b.PinnedHV
	for unpinned != bitboard.Empty {
		from := unpinned.Pop()
		rookMoves := attacks.Rook(from, b.Occupied) & b.Target
		b.serializeMoves(moveList, rook, from, rookMoves)
	}
}

func (b *Board) appendPawnMoves(moveList *[]move.Move) {
	// various properties which change depending on the side to move

	var down, left, right square.Square
	var promotionRank bitboard.Board
	var enPassantRank bitboard.Board
	var doublePushRank bitboard.Board
	var p piece.Piece

	left = -1
	right = 1

	switch b.SideToMove {
	case piece.White:
		down = 8

		promotionRank = bitboard.Rank8
		enPassantRank = bitboard.Rank5
		doublePushRank = bitboard.Rank3

		p = piece.WhitePawn

	case piece.Black:
		down = -8

		promotionRank = bitboard.Rank1
		enPassantRank = bitboard.Rank4
		doublePushRank = bitboard.Rank6

		p = piece.BlackPawn
	}

	pushTarget := b.CheckMask &^ b.Occupied
	captureTarget := b.Enemies & b.CheckMask

	pawns := b.Pawns(b.SideToMove)

	pawnsThatAttack := pawns &^ b.PinnedHV

	unpinnedPawnsThatAttack := pawnsThatAttack &^ b.PinnedD
	pinnedPawnsThatAttack := pawnsThatAttack & b.PinnedD

	pawnAttacksL := attacks.PawnsLeft(unpinnedPawnsThatAttack, b.SideToMove) & captureTarget
	pawnAttacksL |= attacks.PawnsLeft(pinnedPawnsThatAttack, b.SideToMove) & captureTarget & b.PinnedD

	pawnAttacksR := attacks.PawnsRight(unpinnedPawnsThatAttack, b.SideToMove) & captureTarget
	pawnAttacksR |= attacks.PawnsRight(pinnedPawnsThatAttack, b.SideToMove) & captureTarget & b.PinnedD

	simplePawnAttacksL := pawnAttacksL &^ promotionRank
	simplePawnAttacksR := pawnAttacksR &^ promotionRank

	for simplePawnAttacksL != bitboard.Empty {
		to := simplePawnAttacksL.Pop()
		from := to + down + right
		*moveList = append(*moveList, move.New(from, to, p, true))
	}

	for simplePawnAttacksR != bitboard.Empty {
		to := simplePawnAttacksR.Pop()
		from := to + down + left
		*moveList = append(*moveList, move.New(from, to, p, true))
	}

	promotionPawnAttacksL := pawnAttacksL & promotionRank
	promotionPawnAttacksR := pawnAttacksR & promotionRank

	for promotionPawnAttacksL != bitboard.Empty {
		to := promotionPawnAttacksL.Pop()
		from := to + down + right
		appendPromotions(moveList, move.New(from, to, p, true), b.SideToMove)
	}

	for promotionPawnAttacksR != bitboard.Empty {
		to := promotionPawnAttacksR.Pop()
		from := to + down + left
		appendPromotions(moveList, move.New(from, to, p, true), b.SideToMove)
	}

	pawnsThatPush := pawns &^ b.PinnedD

	unpinnedPawnsThatPush := pawnsThatPush &^ b.PinnedHV
	pinnedPawnsThatPush := pawnsThatPush & b.PinnedHV

	pawnPushesSingleUnpinned := attacks.PawnPush(unpinnedPawnsThatPush, b.SideToMove)
	pawnPushesSinglePinned := attacks.PawnPush(pinnedPawnsThatPush, b.SideToMove) & b.PinnedHV

	pawnPushesSingle := (pawnPushesSinglePinned | pawnPushesSingleUnpinned) &^ b.Occupied

	pawnPushesDouble := attacks.PawnPush(pawnPushesSingle&doublePushRank, b.SideToMove) & pushTarget

	pawnPushesSingle &= pushTarget

	simplePawnPushes := pawnPushesSingle &^ promotionRank

	for simplePawnPushes != bitboard.Empty {
		to := simplePawnPushes.Pop()
		from := to + down
		*moveList = append(*moveList, move.New(from, to, p, false))
	}

	for pawnPushesDouble != bitboard.Empty {
		to := pawnPushesDouble.Pop()
		from := to + down + down
		*moveList = append(*moveList, move.New(from, to, p, false))
	}

	promotionPawnPushes := pawnPushesSingle & promotionRank

	for promotionPawnPushes != bitboard.Empty {
		to := promotionPawnPushes.Pop()
		from := to + down
		appendPromotions(moveList, move.New(from, to, p, false), b.SideToMove)
	}

	if b.EnPassantTarget != square.None {
		epPawn := b.EnPassantTarget + down
		them := b.SideToMove.Other()

		epMask := bitboard.Squares[b.EnPassantTarget] | bitboard.Squares[epPawn]
		// check if en-passant leaves king in check
		// this does not account for the double rook pin
		if b.CheckMask&epMask == 0 {
			return
		}

		kingSq := b.Kings[b.SideToMove]
		kingMask := bitboard.Squares[kingSq] & enPassantRank

		enemyRooksQueens := (b.Rooks(them) | b.Queens(them)) & enPassantRank

		// if king and enemy horizontal sliding piece are on ep rank
		// a horizontal rook pin may be possible so more checks
		isPossiblePin := kingMask != bitboard.Empty && enemyRooksQueens != bitboard.Empty

		for fromBB := attacks.Pawn[them][b.EnPassantTarget] & pawnsThatAttack; fromBB != bitboard.Empty; {
			from := fromBB.Pop()

			// pawn is pinned in other direction
			if b.PinnedD.IsSet(from) && !b.PinnedD.IsSet(b.EnPassantTarget) {
				continue
			}

			// check for horizontal rook pin
			// remove the ep pawn and the enemy pawn from the blocker mask
			// and check if a rook ray from the king hits any rook or queen
			pawnsMask := bitboard.Squares[from] | bitboard.Squares[epPawn]
			if isPossiblePin && attacks.Rook(kingSq, b.Occupied&^pawnsMask)&enemyRooksQueens != 0 {
				break
			}

			*moveList = append(*moveList, move.New(from, b.EnPassantTarget, p, true))
		}
	}
}

func (b *Board) appendCastlingMoves(moveList *[]move.Move) {
	// for each castling move the following things are checked:
	// 1. if castling that side is legal (king and rook haven't moved)
	// 2. if pieces are occupying the space between the king and rook
	// 3. if the squares that the king moves through are seen by the enemy
	// if all the conditions are satisfied then castling that side is legal

	switch b.SideToMove {
	case piece.White:
		if b.CastlingRights&castling.WhiteK != 0 &&
			(b.Occupied|b.SeenByEnemy)&bitboard.F1G1 == bitboard.Empty {
			*moveList = append(*moveList, move.New(square.E1, square.G1, piece.WhiteKing, false))
		}

		if b.CastlingRights&castling.WhiteQ != 0 &&
			b.Occupied&bitboard.B1C1D1 == bitboard.Empty &&
			b.SeenByEnemy&bitboard.C1D1 == bitboard.Empty {
			*moveList = append(*moveList, move.New(square.E1, square.C1, piece.WhiteKing, false))
		}
	case piece.Black:
		if b.CastlingRights&castling.BlackK != 0 &&
			(b.Occupied|b.SeenByEnemy)&bitboard.F8G8 == bitboard.Empty {
			*moveList = append(*moveList, move.New(square.E8, square.G8, piece.BlackKing, false))
		}

		if b.CastlingRights&castling.BlackQ != 0 &&
			b.Occupied&bitboard.B8C8D8 == bitboard.Empty &&
			b.SeenByEnemy&bitboard.C8D8 == bitboard.Empty {
			*moveList = append(*moveList, move.New(square.E8, square.C8, piece.BlackKing, false))
		}
	}
}

// serializeMoves serialized the given move bitboard into the movelist.
func (b *Board) serializeMoves(moveList *[]move.Move, p piece.Piece, from square.Square, moves bitboard.Board) {
	for toBB := moves; toBB != bitboard.Empty; {
		to := toBB.Pop()
		*moveList = append(*moveList, move.New(from, to, p, b.Enemies.IsSet(to)))
	}
}

func appendPromotions(moveList *[]move.Move, m move.Move, c piece.Color) {
	*moveList = append(*moveList,
		m.SetPromotion(piece.New(piece.Queen, c)),
		m.SetPromotion(piece.New(piece.Rook, c)),
		m.SetPromotion(piece.New(piece.Bishop, c)),
		m.SetPromotion(piece.New(piece.Knight, c)),
	)
}
