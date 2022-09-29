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
	"laptudirm.com/x/mess/pkg/move"
	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
	"laptudirm.com/x/mess/pkg/zobrist"
)

// Board represents the state of a chessboard at a given position.
type Board struct {
	// position data
	hash      zobrist.Key
	position  mailbox.Board      // 8x8 for fast lookup
	bitboards [13]bitboard.Board // bitboards for eval

	// useful bitboards
	friends bitboard.Board
	enemies bitboard.Board

	sideToMove      piece.Color
	enPassantTarget square.Square
	castlingRights  move.CastlingRights

	// move counters
	halfMoves int
	fullMoves int
}

// String converts a Board into a human readable string.
func (b Board) String() string {
	return fmt.Sprintf("%s\nFen: %s\nKey: %X\n", b.position, b.FEN(), b.hash)
}

// MakeMove plays a legal move on the Board.
func (b *Board) MakeMove(m move.Move) {
	if attackSet := b.MovesOf(m.From); !attackSet.IsSet(m.To) {
		// move not in attack board, illegal move
		panic(fmt.Sprintf("invalid move: piece can't move to given square\n%s", attackSet))
	}

	// update the half-move clock
	// it records the number of plys since the last pawn push or capture
	// for positions which are drawn by the 50-move rule
	switch {
	case m.FromPiece.Type() == piece.Pawn, m.IsCapture():
		// pawn push or capture: reset clock
		b.halfMoves = 0
	default:
		b.halfMoves++
	}

	// update castling rights
	// movement of the rooks or the king, or the capture of the rooks
	// leads to losing the right to castle: update it according to the move

	// rooks or king moved
	switch m.From {
	// white rights
	case square.H1:
		// kingside rook moved
		b.castlingRights &^= move.CastleWhiteKingside
	case square.A1:
		// queenside rook moved
		b.castlingRights &^= move.CastleWhiteQueenside
	case square.E1:
		// king moved
		b.castlingRights &^= move.CastleWhiteKingside
		b.castlingRights &^= move.CastleWhiteQueenside

	// black rights
	case square.H8:
		// kingside rook moved
		b.castlingRights &^= move.CastleBlackKingside
	case square.A8:
		// queenside rook moved
		b.castlingRights &^= move.CastleBlackQueenside
	case square.E8:
		// king moved
		b.castlingRights &^= move.CastleBlackKingside
		b.castlingRights &^= move.CastleBlackQueenside
	}

	// rooks captured
	switch m.To {
	// white rooks
	case square.H1:
		b.castlingRights &^= move.CastleWhiteKingside
	case square.A1:
		b.castlingRights &^= move.CastleWhiteQueenside

	// black rooks
	case square.H8:
		b.castlingRights &^= move.CastleBlackKingside
	case square.A8:
		b.castlingRights &^= move.CastleBlackKingside
	}

	b.hash ^= zobrist.Castling[m.CastlingRights] // remove old rights
	b.hash ^= zobrist.Castling[b.castlingRights] // put new rights

	// move the piece in the records

	if m.IsCapture() {
		// remove captured piece from records
		b.hash ^= zobrist.PieceSquare[m.CapturedPiece][m.Capture] // zobrist hash
		b.enemies.Unset(m.Capture)                                // enemy bitboard
		b.bitboards[m.CapturedPiece].Unset(m.Capture)             // piece bitboard
		b.position[m.Capture] = piece.Empty                       // mailbox board
	}

	// remove moved piece from initial square
	b.hash ^= zobrist.PieceSquare[m.FromPiece][m.From] // zobrist hash
	b.friends.Unset(m.From)                            // friends bitboard
	b.bitboards[m.FromPiece].Unset(m.From)             // piece bitboard
	b.position[m.From] = piece.Empty                   // mailbox board

	// add moved piece to destination square
	b.hash ^= zobrist.PieceSquare[m.FromPiece][m.To] // zobrist hash
	b.friends.Set(m.To)                              // friends bitboard
	b.bitboards[m.ToPiece].Set(m.To)                 // piece bitboard
	b.position[m.To] = m.ToPiece                     // mailbox board

	// update en passant target square
	// clear the previous square, and if current move was double a pawn
	// push, add set the en passant target to the new square

	if b.enPassantTarget != square.None {
		// remove previous square from zobrist hash
		b.hash ^= zobrist.EnPassant[b.enPassantTarget.File()]
	}

	// reset en passant target
	b.enPassantTarget = square.None

	if m.IsDoublePawnPush() {
		// double pawn push; set new en passant target
		b.enPassantTarget = m.From
		if b.sideToMove == piece.WhiteColor {
			b.enPassantTarget += 8
		} else {
			b.enPassantTarget -= 8
		}

		// and new square to zobrist hash
		b.hash ^= zobrist.EnPassant[b.enPassantTarget.File()]
	}

	// switch turn

	// update side to move
	switch b.sideToMove {
	case piece.WhiteColor:
		b.sideToMove = piece.BlackColor
	case piece.BlackColor:
		b.sideToMove = piece.WhiteColor
		b.fullMoves++ // turn completed
	}

	// switch bitboards
	b.friends, b.enemies = b.enemies, b.friends

	// switch zobrist hash
	b.hash ^= zobrist.SideToMove
}

func (b *Board) GenerateMoves() []move.Move {
	var moves []move.Move

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

				m := move.Move{
					From:    from,
					To:      to,
					Capture: to,

					FromPiece:     b.position[from],
					ToPiece:       b.position[from],
					CapturedPiece: b.position[to],

					HalfMoves:       b.halfMoves,
					CastlingRights:  b.castlingRights,
					EnPassantSquare: b.enPassantTarget,
				}

				switch {
				// pawn will promote
				case b.sideToMove == piece.WhiteColor && to.Rank() == square.Rank8:
					fallthrough
				case b.sideToMove == piece.BlackColor && to.Rank() == square.Rank1:
					// evaluate all possible promotions
					for _, promotion := range piece.Promotions {
						m.ToPiece = promotion
						moves = append(moves, m)
					}

				// en passant capture
				case to == b.enPassantTarget:
					m.Capture = to
					if b.sideToMove == piece.WhiteColor {
						m.Capture += 8
					} else {
						m.Capture -= 8
					}
					m.CapturedPiece = b.position[m.Capture]

				// simple push or capture
				default:
					moves = append(moves, m)
				}

				moveSet.Unset(to)
			}

		// other pieces move simply
		default:
			for j := 0; j < 64 && moveSet != bitboard.Empty; j++ {
				to := square.Square(j)
				if moveSet.IsSet(to) {
					m := move.Move{
						From: from,
						To:   to,

						FromPiece:     b.position[from],
						ToPiece:       b.position[from],
						CapturedPiece: b.position[to],

						HalfMoves:       b.halfMoves,
						CastlingRights:  b.castlingRights,
						EnPassantSquare: b.enPassantTarget,
					}
					moves = append(moves, m)
				}
				moveSet.Unset(to)
			}
		}
	}

	return moves
}

func (b *Board) Unmove(m move.Move) {
	b.hash ^= zobrist.SideToMove

	b.friends, b.enemies = b.enemies, b.friends

	// update side to move
	switch b.sideToMove {
	case piece.WhiteColor:
		b.sideToMove = piece.BlackColor
		b.fullMoves--
	case piece.BlackColor:
		b.sideToMove = piece.WhiteColor
	}

	if b.enPassantTarget != square.None {
		b.hash ^= zobrist.EnPassant[b.enPassantTarget.File()]
		b.enPassantTarget = square.None
	}

	if m.EnPassantSquare != square.None {
		b.enPassantTarget = m.EnPassantSquare
		b.hash ^= zobrist.EnPassant[b.enPassantTarget.File()]
	}

	b.hash ^= zobrist.PieceSquare[m.FromPiece][m.To] // zobrist hash
	b.friends.Unset(m.To)                            // friends bitboard
	b.bitboards[m.ToPiece].Unset(m.To)               // piece bitboard
	b.position[m.To] = piece.Empty                   // mailbox board

	b.hash ^= zobrist.PieceSquare[m.FromPiece][m.From] // zobrist hash
	b.friends.Unset(m.From)                            // friends bitboard
	b.bitboards[m.FromPiece].Unset(m.From)             // piece bitboard
	b.position[m.From] = m.FromPiece                   // mailbox board

	if m.IsCapture() {
		b.hash ^= zobrist.PieceSquare[m.CapturedPiece][m.Capture] // zobrist hash
		b.enemies.Set(m.Capture)                                  // enemy bitboard
		b.bitboards[m.CapturedPiece].Set(m.Capture)               // piece bitboard
		b.position[m.Capture] = m.CapturedPiece                   // mailbox board
	}

	b.hash ^= zobrist.Castling[b.castlingRights]
	b.hash ^= zobrist.Castling[m.CastlingRights]
	b.castlingRights = m.CastlingRights

	b.halfMoves = m.HalfMoves
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
