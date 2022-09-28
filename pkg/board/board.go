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

// Package board implements a complete chess board along with valid move
// generation and other related utilities.
package board

import (
	"fmt"

	"laptudirm.com/x/mess/pkg/attacks"
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/board/mailbox"
	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
)

// Board represents the state of a chessboard at a given position.
type Board struct {
	// position data
	position  mailbox.Board      // 8x8 for fast lookup
	bitboards [13]bitboard.Board // bitboards for eval

	// useful bitboards
	friends bitboard.Board
	enemies bitboard.Board

	sideToMove      piece.Color
	enPassantTarget square.Square
	castlingRights  CastlingRights

	// move counters
	halfMoves int
	fullMoves int
}

// String converts a Board into a human readable string.
func (b Board) String() string {
	return fmt.Sprintf("%s\nFEN: %s\n", b.position, b.FEN())
}

// MakeMove plays a legal move on the Board.
func (b *Board) MakeMove(move Move) {
	if attackSet := b.MovesOf(move.From); !attackSet.IsSet(move.To) {
		// move not in attack board, illegal move
		panic(fmt.Sprintf("invalid move: piece can't move to given square\n%s", attackSet))
	}

	isPawn := b.position[move.From].Type() == piece.Pawn
	isCapture := b.position[move.To] != piece.Empty || move.IsEnPassant

	// half-move clock stuff
	switch {
	case isPawn, isCapture:
		// reset clock
		b.halfMoves = 0
	default:
		b.halfMoves++
	}

	// update castling rights

	// rooks or king moved
	switch move.From {
	// white rights
	case square.H1:
		// kingside rook moved
		b.castlingRights &^= CastleWhiteKingside
	case square.A1:
		// queenside rook moved
		b.castlingRights &^= CastleWhiteQueenside
	case square.E1:
		// king moved
		b.castlingRights &^= CastleWhiteKingside
		b.castlingRights &^= CastleWhiteQueenside

	// black rights
	case square.H8:
		// kingside rook moved
		b.castlingRights &^= CastleBlackKingside
	case square.A8:
		// queenside rook moved
		b.castlingRights &^= CastleBlackQueenside
	case square.E8:
		// king moved
		b.castlingRights &^= CastleBlackKingside
		b.castlingRights &^= CastleBlackQueenside
	}

	// rooks captured
	switch move.To {
	// white rooks
	case square.H1:
		b.castlingRights &^= CastleWhiteKingside
	case square.A1:
		b.castlingRights &^= CastleWhiteQueenside

	// black rooks
	case square.H8:
		b.castlingRights &^= CastleBlackKingside
	case square.A8:
		b.castlingRights &^= CastleBlackKingside
	}

	captureSquare := move.To

	if move.IsEnPassant {
		// en-passant capture, captureSquare will be different to move.To
		captureSquare = b.enPassantTarget
		if b.sideToMove == piece.WhiteColor {
			captureSquare += 8
		} else {
			captureSquare -= 8
		}
	}

	if isCapture {
		b.enemies.Unset(captureSquare)
		b.bitboards[b.position[captureSquare]].Unset(captureSquare)
		b.position[captureSquare] = piece.Empty
	}

	b.friends.Unset(move.From)                          // friends bitboard
	b.bitboards[b.position[move.From]].Unset(move.From) // piece bitboard
	b.position[move.To] = b.position[move.From]         // 8x8 board

	b.friends.Set(move.To)                          // friends bitboard
	b.bitboards[b.position[move.From]].Set(move.To) // piece bitboard
	b.position[move.From] = piece.Empty             // 8x8 board

	// reset en passant target
	b.enPassantTarget = square.None

	if isPawn && move.IsDoublePawnPush() {
		b.enPassantTarget = move.From
		if b.sideToMove == piece.WhiteColor {
			b.enPassantTarget += 8
		} else {
			b.enPassantTarget -= 8
		}
	}

	b.switchTurn()
}

func (b *Board) switchTurn() {
	switch b.sideToMove {
	case piece.WhiteColor:
		b.sideToMove = piece.BlackColor
	case piece.BlackColor:
		b.sideToMove = piece.WhiteColor
		b.fullMoves++ // turn completed
	}

	// switch bitboards
	b.friends, b.enemies = b.enemies, b.friends

}

func (b *Board) GenerateMoves() []Move {
	var moves []Move

	for i := 0; i < 64; i++ {
		from := square.Square(i)
		moveSet := b.MovesOf(from)

		switch b.position[from] {
		// handle pawns separately for en passant and promotions
		case piece.Pawn:
			for j := 0; j < 64 && moveSet != bitboard.Empty; j++ {
				to := square.Square(j)
				if !moveSet.IsSet(to) {
					continue
				}

				move := Move{
					From: from,
					To:   to,
				}

				switch {
				// pawn will promote
				case b.sideToMove == piece.WhiteColor && to.Rank() == square.Rank8:
					fallthrough
				case b.sideToMove == piece.BlackColor && to.Rank() == square.Rank1:
					// evaluate all possible promotions
					for _, promotion := range piece.Promotions {
						move.Promotion = promotion
						moves = append(moves, move)
					}

				// en passant capture
				case to == b.enPassantTarget:
					// check for en passant
					move.IsEnPassant = true
					fallthrough

				// simple push or capture
				default:
					moves = append(moves, move)
				}

				moveSet.Unset(to)
			}

		// other pieces move simply
		default:
			for j := 0; j < 64 && moveSet != bitboard.Empty; j++ {
				to := square.Square(j)
				if moveSet.IsSet(to) {
					move := Move{
						From: from,
						To:   to,
					}
					moves = append(moves, move)
				}
				moveSet.Unset(to)
			}
		}
	}

	return moves
}

func (b *Board) MovesOf(index square.Square) bitboard.Board {
	p := b.position[index]
	if p.Color() != b.sideToMove {
		// other side has no moves
		return bitboard.Empty
	}

	switch p.Type() {
	case piece.King:
		return attacks.King(index, b.friends)
	case piece.Queen:
		return attacks.Queen(index, b.friends, b.friends|b.enemies)
	case piece.Rook:
		return attacks.Rook(index, b.friends, b.friends|b.enemies)
	case piece.Knight:
		return attacks.Knight(index, b.friends)
	case piece.Bishop:
		return attacks.Bishop(index, b.friends, b.friends|b.enemies)
	case piece.Pawn:
		return attacks.Pawn(index, b.enPassantTarget, b.sideToMove, b.friends, b.enemies)
	default:
		return bitboard.Empty
	}
}
