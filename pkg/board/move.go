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
	"strings"

	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/board/move"
	"laptudirm.com/x/mess/pkg/board/move/attacks"
	"laptudirm.com/x/mess/pkg/board/move/castling"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
	"laptudirm.com/x/mess/pkg/board/zobrist"
)

// MakeMove plays the given legal move on the Board.
func (b *Board) MakeMove(m move.Move) {
	// add current state to history
	b.History[b.Plys].Move = m
	b.History[b.Plys].CastlingRights = b.CastlingRights
	b.History[b.Plys].CapturedPiece = piece.NoPiece
	b.History[b.Plys].EnPassantTarget = b.EnPassantTarget
	b.History[b.Plys].DrawClock = b.DrawClock
	b.History[b.Plys].Hash = b.Hash

	// update the half-move clock
	// it records the number of plys since the last pawn push or capture
	// for positions which are drawn by the 50-move rule
	b.DrawClock++

	// parse move

	if m == move.Null {
		b.makeNullMove()
		return
	}

	sourceSq := m.Source()
	targetSq := m.Target()
	captureSq := targetSq
	fromPiece := m.FromPiece()
	pieceType := fromPiece.Type()
	toPiece := m.ToPiece()

	isDoublePush := pieceType == piece.Pawn && util.Abs(targetSq-sourceSq) == 16
	isCastling := pieceType == piece.King && util.Abs(targetSq-sourceSq) == 2
	isEnPassant := pieceType == piece.Pawn && targetSq == b.EnPassantTarget
	isCapture := m.IsCapture()

	if pieceType == piece.Pawn {
		b.DrawClock = 0
	}

	// update en passant target square
	if b.EnPassantTarget != square.None {
		b.Hash ^= zobrist.EnPassant[b.EnPassantTarget.File()] // reset hash
	}
	b.EnPassantTarget = square.None // reset square

	switch {
	case isDoublePush:
		// double pawn push; set new en passant target
		target := sourceSq
		if b.SideToMove == piece.White {
			target -= 8
		} else {
			target += 8
		}

		// only set en passant square if an enemy pawn can capture it
		if b.PawnsBB(b.SideToMove.Other())&attacks.Pawn[b.SideToMove][target] != 0 {
			b.EnPassantTarget = target
			// and new square to zobrist hash
			b.Hash ^= zobrist.EnPassant[b.EnPassantTarget.File()]
		}

	case isCastling:
		// castle the rook
		rookInfo := castling.Rooks[targetSq]
		b.ClearSquare(rookInfo.From)
		b.FillSquare(rookInfo.To, rookInfo.RookType)

	case isEnPassant:
		// capture square is different from target square during en passant
		if b.SideToMove == piece.White {
			captureSq += 8
		} else {
			captureSq -= 8
		}
		fallthrough

	case isCapture:
		b.History[b.Plys].CapturedPiece = b.Position[captureSq] // put captured piece in history
		b.DrawClock = 0                                         // reset draw clock since capture
		b.ClearSquare(captureSq)                                // clear the captured square
	}

	// move the piece
	b.ClearSquare(sourceSq)
	b.FillSquare(targetSq, toPiece)

	b.Hash ^= zobrist.Castling[b.CastlingRights] // remove old rights
	b.CastlingRights &^= castling.RightUpdates[sourceSq]
	b.CastlingRights &^= castling.RightUpdates[targetSq]
	b.Hash ^= zobrist.Castling[b.CastlingRights] // put new rights

	// switch turn
	b.Plys++

	// update side to move
	if b.SideToMove = b.SideToMove.Other(); b.SideToMove == piece.White {
		b.FullMoves++
	}
	b.Hash ^= zobrist.SideToMove // switch in zobrist hash
}

func (b *Board) makeNullMove() {
	// update en passant target square
	if b.EnPassantTarget != square.None {
		b.Hash ^= zobrist.EnPassant[b.EnPassantTarget.File()] // reset hash
	}
	b.EnPassantTarget = square.None // reset square

	// switch turn
	b.Plys++

	// update side to move
	if b.SideToMove = b.SideToMove.Other(); b.SideToMove == piece.White {
		b.FullMoves++
	}
	b.Hash ^= zobrist.SideToMove // switch in zobrist hash
}

// UnmakeMove unmakes the last move played on the Board.
func (b *Board) UnmakeMove() {
	if b.SideToMove = b.SideToMove.Other(); b.SideToMove == piece.Black {
		b.FullMoves--
	}

	b.Plys--

	b.EnPassantTarget = b.History[b.Plys].EnPassantTarget
	b.DrawClock = b.History[b.Plys].DrawClock
	b.CastlingRights = b.History[b.Plys].CastlingRights

	m := b.History[b.Plys].Move

	// parse move

	if m == move.Null {
		// use the hash stored in history
		b.Hash = b.History[b.Plys].Hash
		return
	}

	sourceSq := m.Source()
	targetSq := m.Target()
	captureSq := targetSq
	fromPiece := m.FromPiece()
	pieceType := fromPiece.Type()
	capturedPiece := b.History[b.Plys].CapturedPiece

	isCastling := pieceType == piece.King && util.Abs(targetSq-sourceSq) == 2
	isEnPassant := pieceType == piece.Pawn && targetSq == b.EnPassantTarget
	isCapture := m.IsCapture()

	// un-move the piece
	b.ClearSquare(targetSq)
	b.FillSquare(sourceSq, fromPiece)

	switch {
	case isCastling:
		// un-castle the rook
		rookInfo := castling.Rooks[targetSq]
		b.ClearSquare(rookInfo.To)
		b.FillSquare(rookInfo.From, rookInfo.RookType)

	case isEnPassant:
		// capture square is different from target square during en passant
		if b.SideToMove == piece.White {
			captureSq += 8
		} else {
			captureSq -= 8
		}
		fallthrough

	case isCapture:
		// put the captured piece back
		b.FillSquare(captureSq, capturedPiece)
	}

	// use the hash stored in history
	b.Hash = b.History[b.Plys].Hash
}

// NewMove returns a new move.Move representing moving a piece from `from`
// to `to` by adding the necessary contextual information from the Board.
// If the move is a promotion, the promotion square can be set using the
// (move).SetPromotion(piece.Piece) method.
func (b *Board) NewMove(from, to square.Square) move.Move {
	p := b.Position[from]
	return move.New(from, to, p, b.Position[to] != piece.NoPiece)
}

func (b *Board) NewMoveFromString(m string) move.Move {
	from := square.NewFromString(m[:2])
	to := square.NewFromString(m[2:4])

	move := b.NewMove(from, to)
	if len(m) == 5 {
		pieceID := m[4:]
		if b.SideToMove == piece.White {
			pieceID = strings.ToUpper(pieceID)
		}

		move = move.SetPromotion(piece.NewFromString(pieceID))
	}

	return move
}
