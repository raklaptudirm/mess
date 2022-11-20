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

// Package package chess contains the main board representation used by the
// mess chess engine. It may be used as a library when developing other
// engines. It also contains various sub-packages related to the board
// representation.
package chess

import (
	"fmt"

	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/chess/bitboard"
	"laptudirm.com/x/mess/pkg/chess/mailbox"
	"laptudirm.com/x/mess/pkg/chess/move"
	"laptudirm.com/x/mess/pkg/chess/move/attacks"
	"laptudirm.com/x/mess/pkg/chess/move/castling"
	"laptudirm.com/x/mess/pkg/chess/piece"
	"laptudirm.com/x/mess/pkg/chess/square"
	"laptudirm.com/x/mess/pkg/zobrist"
)

// Board represents the state of a chessboard at a given position. It
// contains two representations of the chess board: a 8x8 mailbox which is
// used to easily look up what piece occupies a given square, and a
// bitboard representation, used for various bitwise calculations like
// calculating the attack sets of pieces.
//
// Board contains the additional state information of a chessboard like the
// half-move clock, en passant square, number of plys, etc.
//
// Various pre-calculated utility information like check masks and pin
// masks are stored in Board to prevent the need for expensive calculation.
type Board struct {
	// main position data:
	// these are the basic position data for a chessboard

	// zobrist hash of current position
	Hash zobrist.Key

	// 8x8 mailbox board representation
	Position mailbox.Board

	// bitboard board representation
	PieceBBs [piece.TypeN]bitboard.Board
	ColorBBs [piece.ColorN]bitboard.Board

	// other necessary information
	SideToMove      piece.Color
	EnPassantTarget square.Square
	CastlingRights  castling.Rights

	// move counters
	Plys      int
	FullMoves int
	DrawClock int

	// game history
	History [move.MaxN]BoardState

	// utility information:
	// these data, used by movegen can be calculated from the main board
	// state but are time expensive, so are instead stored in Board

	// bitboards classified by side to move
	Friends bitboard.Board
	Enemies bitboard.Board

	// precalculated Friends | Enemies
	Occupied bitboard.Board

	// king square lookup table
	Kings [piece.ColorN]square.Square

	// places where pieces can move to
	// calculated as ^Friends & CheckMask
	Target bitboard.Board

	// check information
	CheckN    int
	CheckMask bitboard.Board

	// pinned piece information
	PinnedD  bitboard.Board
	PinnedHV bitboard.Board

	// squares attacked by enemy pieces
	SeenByEnemy bitboard.Board
}

// BoardState contains the irreversible position data of a given board
// state. This is used to rollback to a previous position in UnmakeMove.
type BoardState struct {
	// move information
	Move          move.Move   // move made on this BoardState
	CapturedPiece piece.Piece // piece captured by playing Move

	// irreversible information
	CastlingRights  castling.Rights
	EnPassantTarget square.Square
	DrawClock       int

	// zobrist key is reversible but is stored for repetition detection
	Hash zobrist.Key
}

// String converts a Board into a human readable string.
func (b Board) String() string {
	return fmt.Sprintf("%s\nFen: %s\nKey: %X\n", b.Position, b.FEN(), b.Hash)
}

// IsDraw checks if the given position is a draw either by the 50 move rule
// or by a repetition. Threefold repetition is not calculated as it is just
// simpler to evaluate any repetition as a draw.
func (b *Board) IsDraw() bool {
	return b.DrawClock >= 100 || b.IsRepetition()
}

// IsRepetition checks if the current position has occurred in the game
// before. This is done by probing the game history till the last
// irreversible move, which is pawn pushes or a capture.
func (b *Board) IsRepetition() bool {
	// probe till game start or last irreversible move, whichever is closer
	depth := util.Max(0, b.Plys-b.DrawClock)

	for i := b.Plys - 2; i >= depth; i -= 2 {
		if b.History[i].Hash == b.Hash {
			return true
		}
	}

	return false
}

// ClearSquare removes the piece occupying the given square and updates the
// dependent position information accordingly.
func (b *Board) ClearSquare(s square.Square) {
	p := b.Position[s]

	b.ColorBBs[p.Color()].Unset(s)

	// remove piece from other records
	b.PieceBBs[p.Type()].Unset(s)       // piece bitboard
	b.Position[s] = piece.NoPiece       // mailbox board
	b.Hash ^= zobrist.PieceSquare[p][s] // zobrist hash
}

// FillSquare fills the given square with the given piece. Callers should
// make sure that the provided square is unoccupied, otherwise the
// incrementally updating the position will give wrong results.
func (b *Board) FillSquare(s square.Square, p piece.Piece) {
	c := p.Color()
	t := p.Type()

	b.ColorBBs[c].Set(s)

	if t == piece.King {
		b.Kings[c] = s
	}

	b.PieceBBs[t].Set(s)                // piece bitboard
	b.Position[s] = p                   // mailbox board
	b.Hash ^= zobrist.PieceSquare[p][s] // zobrist hash
}

// IsInCheck checks if the side with the given color is in check.
func (b *Board) IsInCheck(c piece.Color) bool {
	return b.IsAttacked(b.Kings[c], c.Other())
}

// IsAttacked checks if the given squares is attacked by pieces of the
// given color.
func (b *Board) IsAttacked(s square.Square, them piece.Color) bool {
	if attacks.Pawn[them.Other()][s]&b.Pawns(them) != bitboard.Empty {
		return true
	}

	if attacks.Knight[s]&b.Knights(them) != bitboard.Empty {
		return true
	}

	if attacks.King[s]&b.King(them) != bitboard.Empty {
		return true
	}

	queens := b.Queens(them)

	if attacks.Bishop(s, b.Occupied)&(b.Bishops(them)|queens) != bitboard.Empty {
		return true
	}

	return attacks.Rook(s, b.Occupied)&(b.Rooks(them)|queens) != bitboard.Empty
}

// Pawns returns a bitboard of all the pawns of the given color.
func (b *Board) Pawns(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Pawn] & b.ColorBBs[c]
}

// Knights returns a bitboard of all the knights of the given color.
func (b *Board) Knights(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Knight] & b.ColorBBs[c]
}

// Bishops returns a bitboard of all the bishops of the given color.
func (b *Board) Bishops(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Bishop] & b.ColorBBs[c]
}

// Rooks returns a bitboard of all the rooks of the given color.
func (b *Board) Rooks(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Rook] & b.ColorBBs[c]
}

// Queens returns a bitboard of all the queens of the given color.
func (b *Board) Queens(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.Queen] & b.ColorBBs[c]
}

// King returns a bitboard containing the king of the given color.
func (b *Board) King(c piece.Color) bitboard.Board {
	return b.PieceBBs[piece.King] & b.ColorBBs[c]
}

// InitBitboards initializes all the different utility bitboards which are
// calculated and necessary for move generation.
func (b *Board) InitBitboards() {
	b.Friends = b.ColorBBs[b.SideToMove]
	b.Enemies = b.ColorBBs[b.SideToMove.Other()]
	b.Occupied = b.Friends | b.Enemies
	b.CalculateCheckmask()
	b.CalculatePinmask()
	b.SeenByEnemy = b.SeenSquares(b.SideToMove.Other())
	b.Target = ^b.Friends & b.CheckMask
}

// CalculateCheckmask calculates the check-mask of the current board state,
// along with the number of checkers.
//
// A checker is an enemy piece which is directly checking the king. The
// number of checkers can be a maximum of two (double check).
//
// The check-mask is defined as all the squares to which if a friendly
// piece is moved to will block all checks. This is defined as empty for
// double check, the checking piece and, if the checker is a sliding piece,
// the squares between the king and the checker. The bitboard is universe
// if the king is not in check.
func (b *Board) CalculateCheckmask() {
	us := b.SideToMove
	them := us.Other()

	b.CheckN = 0
	b.CheckMask = bitboard.Empty

	kingSq := b.Kings[us]

	pawns := b.Pawns(them) & attacks.Pawn[us][kingSq]
	knights := b.Knights(them) & attacks.Knight[kingSq]
	bishops := (b.Bishops(them) | b.Queens(them)) & attacks.Bishop(kingSq, b.Occupied)
	rooks := (b.Rooks(them) | b.Queens(them)) & attacks.Rook(kingSq, b.Occupied)

	// a pawn and a knight cannot be checking the king at the same time as
	// they are not sliding pieces thus discovered attacks are impossible
	switch {
	case pawns != bitboard.Empty:
		b.CheckMask |= pawns
		b.CheckN++

	case knights != bitboard.Empty:
		b.CheckMask |= knights
		b.CheckN++
	}

	if bishops != bitboard.Empty {
		bishopSq := bishops.FirstOne()
		b.CheckMask |= bitboard.Between[kingSq][bishopSq] | bitboard.Squares[bishopSq]
		b.CheckN++
	}

	// 2 is the largest possible value for CheckN so short circuit if thats reached
	if b.CheckN < 2 && rooks != bitboard.Empty {
		if b.CheckN == 0 && rooks.Count() > 1 {
			// double check, don't set the check-mask
			b.CheckN++
		} else {
			rookSq := rooks.FirstOne()
			b.CheckMask |= bitboard.Between[kingSq][rookSq] | bitboard.Squares[rookSq]
			b.CheckN++
		}
	}

	if b.CheckN == 0 {
		// king is not in check so check-mask is universe
		b.CheckMask = bitboard.Universe
	}
}

// CalculatePinmask calculates the horizontal and vertical pin-masks.
// A pin-mask is defined as the mask containing all attack rays pieces
// pinning a piece in a given direction (horizontal or vertical).
func (b *Board) CalculatePinmask() {
	us := b.SideToMove
	them := us.Other()

	kingSq := b.Kings[us]

	friends := b.ColorBBs[us]
	enemies := b.ColorBBs[them]

	b.PinnedD = bitboard.Empty
	b.PinnedHV = bitboard.Empty

	// consider enemy rooks and queens which are attacking or would attack the king if not for intervening pieces
	// the king is considered as a rook for this and it's attack sets & with rooks and queens gives the bitboard
	for rooks := (b.Rooks(them) | b.Queens(them)) & attacks.Rook(kingSq, enemies); rooks != bitboard.Empty; {
		rook := rooks.Pop()
		possiblePin := bitboard.Between[kingSq][rook] | bitboard.Squares[rook]

		// if there is only one friendly piece blocking the ray, it is pinned
		if (possiblePin & friends).Count() == 1 {
			b.PinnedHV |= possiblePin
		}
	}

	// consider enemy bishops and queens which are attacking or would attack the king if not for intervening pieces
	// the king is considered as a bishop for this and it's attack sets & with bishops and queens gives the bitboard
	for bishops := (b.Bishops(them) | b.Queens(them)) & attacks.Bishop(kingSq, enemies); bishops != bitboard.Empty; {
		bishop := bishops.Pop()
		possiblePin := bitboard.Between[kingSq][bishop] | bitboard.Squares[bishop]

		// if there is only one friendly piece blocking the ray, it is pinned
		if (possiblePin & friends).Count() == 1 {
			b.PinnedD |= possiblePin
		}
	}
}

// SeenSquares returns a bitboard containing all the squares that are
// seen(attacked) by pieces of the given color. The enemy king is not
// considered as a sliding ray blocker by SeenSquares since it has to
// move away from the attack exposing the blocked squares.
func (b *Board) SeenSquares(by piece.Color) bitboard.Board {
	pawns := b.Pawns(by)
	knights := b.Knights(by)
	bishops := b.Bishops(by)
	rooks := b.Rooks(by)
	queens := b.Queens(by)
	kingSq := b.Kings[by]

	// don't consider the enemy king as a blocker
	blockers := b.Occupied &^ b.King(by.Other())

	seen := attacks.PawnsLeft(pawns, by) | attacks.PawnsRight(pawns, by)

	for knights != bitboard.Empty {
		from := knights.Pop()
		seen |= attacks.Knight[from]
	}

	for bishops != bitboard.Empty {
		from := bishops.Pop()
		seen |= attacks.Bishop(from, blockers)
	}

	for rooks != bitboard.Empty {
		from := rooks.Pop()
		seen |= attacks.Rook(from, blockers)
	}

	for queens != bitboard.Empty {
		from := queens.Pop()
		seen |= attacks.Queen(from, blockers)
	}

	seen |= attacks.King[kingSq]

	return seen
}
