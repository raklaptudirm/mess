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
	if !b.IsPseudoLegal(m) {
		// move not in attack board, illegal move
		panic(fmt.Sprintf("invalid move: %s can't move to %s\n%s\n%v", m.FromPiece, m.To, b, b.MoveList))
	}

	// update the half-move clock
	// it records the number of plys since the last pawn push or capture
	// for positions which are drawn by the 50-move rule
	if b.HalfMoves++; !m.IsReversible() {
		b.HalfMoves = 0
	}

	b.CastlingRights &^= m.CastlingRightUpdates()

	b.Hash ^= zobrist.Castling[m.CastlingRights] // remove old rights
	b.Hash ^= zobrist.Castling[b.CastlingRights] // put new rights

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
	if b.EnPassantTarget != square.None {
		b.Hash ^= zobrist.EnPassant[b.EnPassantTarget.File()] // reset hash
		b.EnPassantTarget = square.None                       // reset square
	}

	if m.IsDoublePawnPush() {
		// double pawn push; set new en passant target
		target := m.From
		if b.SideToMove == piece.White {
			target -= 8
		} else {
			target += 8
		}

		// only set en passant square if an enemy pawn can capture it
		if b.PieceBBs[piece.Pawn]&b.ColorBBs[b.SideToMove.Other()]&attacks.Pawn[b.SideToMove][target] != 0 {
			b.EnPassantTarget = target
			// and new square to zobrist hash
			b.Hash ^= zobrist.EnPassant[b.EnPassantTarget.File()]
		}
	}

	// switch turn

	// update side to move
	if b.SideToMove = b.SideToMove.Other(); b.SideToMove == piece.White {
		b.FullMoves++
	}
	b.Hash ^= zobrist.SideToMove // switch in zobrist hash

	b.MoveList = append(b.MoveList, m)
}

func (b *Board) UnmakeMove(m move.Move) {
	if b.EnPassantTarget != square.None {
		b.Hash ^= zobrist.EnPassant[b.EnPassantTarget.File()]
		b.EnPassantTarget = square.None
	}

	if m.EnPassantSquare != square.None {
		b.EnPassantTarget = m.EnPassantSquare
		b.Hash ^= zobrist.EnPassant[b.EnPassantTarget.File()]
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

	b.Hash ^= zobrist.Castling[b.CastlingRights]
	b.Hash ^= zobrist.Castling[m.CastlingRights]
	b.CastlingRights = m.CastlingRights

	b.HalfMoves = m.HalfMoves

	// update side to move

	b.Hash ^= zobrist.SideToMove
	if b.SideToMove = b.SideToMove.Other(); b.SideToMove == piece.Black {
		b.FullMoves--
	}

	b.MoveList = b.MoveList[:len(b.MoveList)-1]
}

func (b *Board) GenerateMoves() []move.Move {
	var moves []move.Move

	for from := square.A8; from <= square.H1; from++ {
		moveSet := b.MovesOf(from)

		for to := square.A8; to <= square.H1 && moveSet != bitboard.Empty; to++ {
			if !moveSet.IsSet(to) {
				continue
			}

			moves = append(moves, b.Moves(from, to)...)
			moveSet.Unset(to)
		}
	}

	return moves
}

func (b *Board) IsPseudoLegal(m move.Move) bool {
	return b.MovesOf(m.From).IsSet(m.To)
}

func (b *Board) MovesOf(index square.Square) bitboard.Board {
	p := b.Position[index]
	if p == piece.NoPiece || p.Color() != b.SideToMove {
		// other side has no moves
		return bitboard.Empty
	}

	var a bitboard.Board
	occ := b.Occupied()
	friends := b.ColorBBs[b.SideToMove]
	enemies := b.ColorBBs[b.SideToMove.Other()]

	switch p.Type() {
	case piece.King:
		cr := b.CastlingRights

		// even if the king and rooks haven't moved, the king can still
		// be prevented from castling with checks on it's path
		switch them := b.SideToMove.Other(); b.SideToMove {
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
		a = attacks.PawnAll(index, b.EnPassantTarget, b.SideToMove, occ, enemies)
	default:
		a = bitboard.Empty
	}

	return a &^ friends
}

func (b *Board) Moves(from, to square.Square) []move.Move {
	var moves []move.Move

	m := move.Move{
		From:    from,
		To:      to,
		Capture: to,

		FromPiece:     b.Position[from],
		ToPiece:       b.Position[from],
		CapturedPiece: b.Position[to],

		HalfMoves:       b.HalfMoves,
		CastlingRights:  b.CastlingRights,
		EnPassantSquare: b.EnPassantTarget,
	}

	// handle pawns separately for en passant and promotions
	if b.Position[from].Type() == piece.Pawn {
		switch {
		// pawn will promote
		case b.SideToMove == piece.White && to.Rank() == square.Rank8,
			b.SideToMove == piece.Black && to.Rank() == square.Rank1:
			// evaluate all possible promotions
			for _, promotion := range piece.Promotions {
				m.ToPiece = piece.New(promotion, b.SideToMove)
				moves = append(moves, m)
			}

			return moves

		// en passant capture
		case to == b.EnPassantTarget:
			m.Capture = to
			if b.SideToMove == piece.White {
				m.Capture += 8
			} else {
				m.Capture -= 8
			}
			m.CapturedPiece = b.Position[m.Capture]
		}
	}

	return []move.Move{m}
}
