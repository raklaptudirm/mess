package board

import (
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
}

func (b *Board) GenerateMoves() []move.Move {
	moves := make([]move.Move, 0, 20)
	occ := b.Occupied()
	friends := b.ColorBBs[b.SideToMove]
	enemies := b.ColorBBs[b.SideToMove.Other()]

	move := move.Move{
		HalfMoves:       b.HalfMoves,
		CastlingRights:  b.CastlingRights,
		EnPassantSquare: b.EnPassantTarget,
	}

	{
		pawns := b.PieceBBs[piece.Pawn] & friends
		enemies.Set(b.EnPassantTarget)

		switch b.SideToMove {
		case piece.White:
			move.FromPiece = piece.WhitePawn

			for fromBB := pawns & bitboard.Rank7; fromBB != bitboard.Empty; {
				from := fromBB.Pop()
				move.From = from
				move.CapturedPiece = piece.NoPiece

				if to := from - 8; bitboard.Squares[to]&occ == 0 {
					move.To = to
					for _, promotion := range piece.Promotions {
						move.ToPiece = piece.New(promotion, piece.White)
						moves = append(moves, move)
					}
				}

				if to := from - 7; from.File() < square.FileH && bitboard.Squares[to]&enemies != 0 {
					move.To = to
					move.Capture = to
					move.CapturedPiece = b.Position[to]

					for _, promotion := range piece.Promotions {
						move.ToPiece = piece.New(promotion, piece.White)
						moves = append(moves, move)
					}
				}

				if to := from - 9; from.File() > square.FileA && bitboard.Squares[to]&enemies != 0 {
					move.To = to
					move.Capture = to
					move.CapturedPiece = b.Position[to]

					for _, promotion := range piece.Promotions {
						move.ToPiece = piece.New(promotion, piece.White)
						moves = append(moves, move)
					}
				}
			}

			move.ToPiece = piece.WhitePawn

			for fromBB := pawns &^ bitboard.Rank7; fromBB != bitboard.Empty; {
				from := fromBB.Pop()
				move.From = from
				move.CapturedPiece = piece.NoPiece

				if to := from - 8; bitboard.Squares[to]&occ == 0 {
					move.To = to
					moves = append(moves, move)

					if to := from - 16; from.Rank() == square.Rank2 && bitboard.Squares[to]&occ == 0 {
						move.To = to
						moves = append(moves, move)
					}
				}

				if to := from - 7; from.File() < square.FileH && bitboard.Squares[to]&enemies != 0 {
					move.To = to
					if to == b.EnPassantTarget {
						to += 8
					}
					move.Capture = to
					move.CapturedPiece = b.Position[to]
					moves = append(moves, move)
				}

				if to := from - 9; from.File() > square.FileA && bitboard.Squares[to]&enemies != 0 {
					move.To = to
					if to == b.EnPassantTarget {
						to += 8
					}
					move.Capture = to
					move.CapturedPiece = b.Position[to]
					move := b.NewMove(from, to)
					moves = append(moves, move)
				}
			}
		case piece.Black:
			move.FromPiece = piece.BlackPawn

			for fromBB := pawns & bitboard.Rank2; fromBB != bitboard.Empty; {
				from := fromBB.Pop()
				move.From = from
				move.CapturedPiece = piece.NoPiece

				if to := from + 8; bitboard.Squares[to]&occ == 0 {
					move.To = to
					for _, promotion := range piece.Promotions {
						move.ToPiece = piece.New(promotion, piece.Black)
						moves = append(moves, move)
					}
				}

				if to := from + 9; from.File() < square.FileH && bitboard.Squares[to]&enemies != 0 {
					move.To = to
					move.Capture = to
					move.CapturedPiece = b.Position[to]
					for _, promotion := range piece.Promotions {
						move.ToPiece = piece.New(promotion, piece.Black)
						moves = append(moves, move)
					}
				}

				if to := from + 7; from.File() > square.FileA && bitboard.Squares[to]&enemies != 0 {
					move.To = to
					move.Capture = to
					move.CapturedPiece = b.Position[to]
					for _, promotion := range piece.Promotions {
						move.ToPiece = piece.New(promotion, piece.Black)
						moves = append(moves, move)
					}
				}
			}

			move.ToPiece = piece.BlackPawn

			for fromBB := pawns &^ bitboard.Rank2; fromBB != bitboard.Empty; {
				from := fromBB.Pop()
				move.From = from
				move.CapturedPiece = piece.NoPiece

				if to := from + 8; bitboard.Squares[to]&occ == 0 {
					move.To = to
					moves = append(moves, move)

					if to := from + 16; from.Rank() == square.Rank7 && bitboard.Squares[to]&occ == 0 {
						move.To = to
						moves = append(moves, move)
					}
				}

				if to := from + 9; from.File() < square.FileH && bitboard.Squares[to]&enemies != 0 {
					move.To = to
					if to == b.EnPassantTarget {
						to -= 8
					}
					move.Capture = to
					move.CapturedPiece = b.Position[to]
					moves = append(moves, move)
				}

				if to := from + 7; from.File() > square.FileA && bitboard.Squares[to]&enemies != 0 {
					move.To = to
					if to == b.EnPassantTarget {
						to -= 8
					}
					move.Capture = to
					move.CapturedPiece = b.Position[to]
					moves = append(moves, move)
				}
			}
		}
	}

	p := piece.New(piece.Knight, b.SideToMove)
	move.FromPiece = p
	move.ToPiece = p
	for fromBB := b.PieceBBs[piece.Knight] & friends; fromBB != bitboard.Empty; {
		from := fromBB.Pop()
		move.From = from

		for toBB := attacks.Knight[from] &^ friends; toBB != bitboard.Empty; {
			to := toBB.Pop()
			move.To = to
			move.Capture = to
			move.CapturedPiece = b.Position[to]

			moves = append(moves, move)
		}
	}

	p = piece.New(piece.Bishop, b.SideToMove)
	move.FromPiece = p
	move.ToPiece = p
	for fromBB := b.PieceBBs[piece.Bishop] & friends; fromBB != bitboard.Empty; {
		from := fromBB.Pop()
		move.From = from

		for toBB := attacks.Bishop(from, occ) &^ friends; toBB != bitboard.Empty; {
			to := toBB.Pop()
			move.To = to
			move.Capture = to
			move.CapturedPiece = b.Position[to]

			moves = append(moves, move)
		}
	}

	p = piece.New(piece.Rook, b.SideToMove)
	move.FromPiece = p
	move.ToPiece = p
	for fromBB := b.PieceBBs[piece.Rook] & friends; fromBB != bitboard.Empty; {
		from := fromBB.Pop()
		move.From = from

		for toBB := attacks.Rook(from, occ) &^ friends; toBB != bitboard.Empty; {
			to := toBB.Pop()
			move.To = to
			move.Capture = to
			move.CapturedPiece = b.Position[to]

			moves = append(moves, move)
		}
	}

	p = piece.New(piece.Queen, b.SideToMove)
	move.FromPiece = p
	move.ToPiece = p
	for fromBB := b.PieceBBs[piece.Queen] & friends; fromBB != bitboard.Empty; {
		from := fromBB.Pop()
		move.From = from

		for toBB := attacks.Queen(from, occ) &^ friends; toBB != bitboard.Empty; {
			to := toBB.Pop()
			move.To = to
			move.Capture = to
			move.CapturedPiece = b.Position[to]

			moves = append(moves, move)
		}
	}

	{
		p = piece.New(piece.King, b.SideToMove)
		move.FromPiece = p
		move.ToPiece = p

		from := b.Kings[b.SideToMove]
		move.From = from
		for toBB := attacks.King[from] &^ friends; toBB != bitboard.Empty; {
			to := toBB.Pop()
			move.To = to
			move.Capture = to
			move.CapturedPiece = b.Position[to]

			moves = append(moves, move)
		}

		move.CapturedPiece = piece.NoPiece

		switch b.SideToMove {
		case piece.White:
			if b.CastlingRights&castling.White == castling.None ||
				b.IsAttacked(square.E1, piece.Black) {
				break
			}

			if b.CastlingRights&castling.WhiteKingside != 0 &&
				occ&0x6000000000000000 == bitboard.Empty &&
				!b.IsAttacked(square.F1, piece.Black) {
				move.To = square.G1
				moves = append(moves, move)
			}

			if b.CastlingRights&castling.WhiteQueenside != 0 &&
				occ&0xe00000000000000 == bitboard.Empty &&
				!b.IsAttacked(square.D1, piece.Black) {
				move.To = square.C1
				moves = append(moves, move)
			}
		case piece.Black:
			if b.CastlingRights&castling.Black == castling.None ||
				b.IsAttacked(square.E8, piece.White) {
				break
			}

			if b.CastlingRights&castling.BlackKingside != 0 &&
				occ&0x60 == bitboard.Empty &&
				!b.IsAttacked(square.F8, piece.White) {
				move.To = square.G8
				moves = append(moves, move)
			}

			if b.CastlingRights&castling.BlackQueenside != 0 &&
				occ&0xe == bitboard.Empty &&
				!b.IsAttacked(square.D8, piece.White) {
				move.To = square.C8
				moves = append(moves, move)
			}
		}
	}

	return moves
}

func (b *Board) NewMove(from, to square.Square) move.Move {
	p := b.Position[from]
	m := move.Move{
		From:    from,
		To:      to,
		Capture: to,

		FromPiece:     p,
		ToPiece:       p,
		CapturedPiece: b.Position[to],

		HalfMoves:       b.HalfMoves,
		CastlingRights:  b.CastlingRights,
		EnPassantSquare: b.EnPassantTarget,
	}

	// handle pawns separately for en passant
	if b.Position[from].Type() == piece.Pawn &&
		to == b.EnPassantTarget {
		m.Capture = to
		if b.SideToMove == piece.White {
			m.Capture += 8
		} else {
			m.Capture -= 8
		}
		m.CapturedPiece = b.Position[m.Capture]
	}

	return m
}
