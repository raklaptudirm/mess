package board

import (
	"fmt"

	"laptudirm.com/x/mess/pkg/attacks"
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/castling"
	"laptudirm.com/x/mess/pkg/move"
	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
	"laptudirm.com/x/mess/pkg/zobrist"
)

// MakeMove plays a legal move on the Board.
func (b *Board) MakeMove(m move.Move) {
	if attackSet := b.MovesOf(m.From); !attackSet.IsSet(m.To) {
		// move not in attack board, illegal move
		panic(fmt.Sprintf("invalid move: piece can't move to given square\n%s", attackSet))
	}

	// update the half-move clock
	// it records the number of plys since the last pawn push or capture
	// for positions which are drawn by the 50-move rule
	if b.halfMoves++; !m.IsReversible() {
		b.halfMoves = 0
	}

	b.castlingRights &^= m.CastlingRightUpdates()

	b.hash ^= zobrist.Castling[m.CastlingRights] // remove old rights
	b.hash ^= zobrist.Castling[b.castlingRights] // put new rights

	// move the piece in the records

	if m.IsCapture() {
		b.ClearSquare(m.Capture)
	}

	b.ClearSquare(m.From)
	b.FillSquare(m.To, m.ToPiece)

	if m.IsCastle() {
		rookInfo := castling.Rooks[m.To]
		b.ClearSquare(rookInfo.From)
		b.FillSquare(rookInfo.To, rookInfo.RookType)
	}

	// update en passant target square
	// clear the previous square, and if current move was double a pawn
	// push, add set the en passant target to the new square

	// reset old en passant square
	if b.enPassantTarget != square.None {
		b.hash ^= zobrist.EnPassant[b.enPassantTarget.File()] // reset hash
		b.enPassantTarget = square.None                       // reset square
	}

	if m.IsDoublePawnPush() {
		// double pawn push; set new en passant target
		target := m.From
		if b.sideToMove == piece.White {
			target -= 8
		} else {
			target += 8
		}

		// only set en passant square if an enemy pawn can capture it
		if b.bitboards[piece.New(piece.Pawn, b.sideToMove.Other())]&attacks.Pawn[b.sideToMove][target] != 0 {
			b.enPassantTarget = target
			// and new square to zobrist hash
			b.hash ^= zobrist.EnPassant[b.enPassantTarget.File()]
		}
	}

	// switch turn

	// update side to move
	if b.sideToMove = b.sideToMove.Other(); b.sideToMove == piece.White {
		b.fullMoves++
	}

	// switch bitboards
	b.friends, b.enemies = b.enemies, b.friends

	// switch zobrist hash
	b.hash ^= zobrist.SideToMove
}

func (b *Board) Unmove(m move.Move) {
	if b.enPassantTarget != square.None {
		b.hash ^= zobrist.EnPassant[b.enPassantTarget.File()]
		b.enPassantTarget = square.None
	}

	if m.EnPassantSquare != square.None {
		b.enPassantTarget = m.EnPassantSquare
		b.hash ^= zobrist.EnPassant[b.enPassantTarget.File()]
	}

	b.ClearSquare(m.To)
	b.FillSquare(m.From, m.FromPiece)

	if m.IsCapture() {
		b.FillSquare(m.Capture, m.CapturedPiece)
	}

	if m.IsCastle() {
		rookInfo := castling.Rooks[m.To]
		b.ClearSquare(rookInfo.To)
		b.FillSquare(rookInfo.From, rookInfo.RookType)
	}

	b.hash ^= zobrist.Castling[b.castlingRights]
	b.hash ^= zobrist.Castling[m.CastlingRights]
	b.castlingRights = m.CastlingRights

	b.halfMoves = m.HalfMoves

	b.hash ^= zobrist.SideToMove

	b.friends, b.enemies = b.enemies, b.friends

	// update side to move
	if b.sideToMove = b.sideToMove.Other(); b.sideToMove == piece.Black {
		b.fullMoves--
	}
}

func (b *Board) GenerateMoves() []move.Move {
	var moves []move.Move

	for from := square.A8; from <= square.H1; from++ {
		moveSet := b.MovesOf(from)

		for to := square.A8; to <= square.H1 && moveSet != bitboard.Empty; to++ {
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

			// handle pawns separately for en passant and promotions
			if b.position[from].Type() == piece.Pawn {
				switch {
				// pawn will promote
				case b.sideToMove == piece.White && to.Rank() == square.Rank8,
					b.sideToMove == piece.Black && to.Rank() == square.Rank1:
					// evaluate all possible promotions
					for _, promotion := range piece.Promotions {
						m.ToPiece = piece.New(promotion, b.sideToMove)
						moves = append(moves, m)
					}

					moveSet.Unset(to)
					continue

				// en passant capture
				case to == b.enPassantTarget:
					m.Capture = to
					if b.sideToMove == piece.White {
						m.Capture += 8
					} else {
						m.Capture -= 8
					}
					m.CapturedPiece = b.position[m.Capture]
				}
			}

			moves = append(moves, m)
			moveSet.Unset(to)
		}
	}

	return moves
}

func (b *Board) MovesOf(index square.Square) bitboard.Board {
	p := b.position[index]
	if p == piece.NoPiece || p.Color() != b.sideToMove {
		// other side has no moves
		return bitboard.Empty
	}

	var a bitboard.Board
	occ := b.friends | b.enemies

	switch p.Type() {
	case piece.King:
		cr := b.castlingRights

		// even if the king and rooks haven't moved, the king can still
		// be prevented from castling with checks on it's path
		switch them := b.sideToMove.Other(); b.sideToMove {
		case piece.White:
			// return early
			if index != square.E1 {
				break
			}

			// king is in check
			if b.IsAttacked(square.E1, them) {
				cr = castling.None
				break
			}

			// can't castle through check
			if b.IsAttacked(square.F1, them) {
				cr &^= castling.WhiteKingside
			}

			if b.IsAttacked(square.D1, them) {
				cr &^= castling.WhiteQueenside
			}
		case piece.Black:
			// return early
			if index != square.E8 {
				break
			}

			// king is in check
			if b.IsAttacked(square.E8, them) {
				cr = castling.None
				break
			}

			// can't castle through check
			if b.IsAttacked(square.F8, them) {
				cr &^= castling.BlackKingside
			}

			if b.IsAttacked(square.D8, them) {
				cr &^= castling.BlackQueenside
			}
		}

		a = attacks.KingAll(index, occ, cr)
	case piece.Queen:
		a = attacks.Queen(index, occ)
	case piece.Rook:
		a = attacks.Rook(index, occ)
	case piece.Knight:
		a = attacks.Knight[index]
	case piece.Bishop:
		a = attacks.Bishop(index, occ)
	case piece.Pawn:
		a = attacks.PawnAll(index, b.enPassantTarget, b.sideToMove, occ, b.enemies)
	default:
		a = bitboard.Empty
	}

	return a &^ b.friends
}
