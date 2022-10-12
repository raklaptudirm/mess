package board

import (
	"laptudirm.com/x/mess/pkg/attacks"
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/castling"
	"laptudirm.com/x/mess/pkg/move"
	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
	"laptudirm.com/x/mess/pkg/util"
	"laptudirm.com/x/mess/pkg/zobrist"
)

// MakeMove plays a legal move on the Board.
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
		if b.PieceBBs[piece.Pawn]&b.ColorBBs[b.SideToMove.Other()]&attacks.Pawn[b.SideToMove][target] != 0 {
			b.EnPassantTarget = target
			// and new square to zobrist hash
			b.Hash ^= zobrist.EnPassant[b.EnPassantTarget.File()]
		}

	case isCastling:
		rookInfo := castling.Rooks[targetSq]
		b.ClearSquare(rookInfo.From)
		b.FillSquare(rookInfo.To, rookInfo.RookType)

	case isEnPassant:
		if b.SideToMove == piece.White {
			captureSq += 8
		} else {
			captureSq -= 8
		}
		fallthrough

	case isCapture:
		b.DrawClock = 0
		b.History[b.Plys].CapturedPiece = b.Position[captureSq]
		b.ClearSquare(captureSq)
	}

	// move the piece in the records
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

	sourceSq := m.Source()
	targetSq := m.Target()
	captureSq := targetSq
	fromPiece := m.FromPiece()
	pieceType := fromPiece.Type()
	capturedPiece := b.History[b.Plys].CapturedPiece

	isCastling := pieceType == piece.King && util.Abs(targetSq-sourceSq) == 2
	isEnPassant := pieceType == piece.Pawn && targetSq == b.EnPassantTarget
	isCapture := m.IsCapture()

	b.ClearSquare(targetSq)
	b.FillSquare(sourceSq, fromPiece)

	switch {
	case isCastling:
		rookInfo := castling.Rooks[targetSq]
		b.ClearSquare(rookInfo.To)
		b.FillSquare(rookInfo.From, rookInfo.RookType)

	case isEnPassant:
		if b.SideToMove == piece.White {
			captureSq += 8
		} else {
			captureSq -= 8
		}
		fallthrough

	case isCapture:
		b.FillSquare(captureSq, capturedPiece)
	}

	b.Hash = b.History[b.Plys].Hash
}

func (b *Board) GenerateMoves() []move.Move {
	moves := make([]move.Move, 0, 30)
	occ := b.Occupied()
	friends := b.ColorBBs[b.SideToMove]
	enemies := b.ColorBBs[b.SideToMove.Other()]

	if b.EnPassantTarget != square.None {
		pawns := b.PieceBBs[piece.Pawn] & b.ColorBBs[b.SideToMove]
		pawn := piece.New(piece.Pawn, b.SideToMove)

		for fromBB := attacks.Pawn[b.SideToMove.Other()][b.EnPassantTarget] & pawns; fromBB != bitboard.Empty; {
			from := fromBB.Pop()
			moves = append(moves, move.New(from, b.EnPassantTarget, pawn, true))
		}
	}

	{
		pawns := b.PieceBBs[piece.Pawn] & friends

		switch b.SideToMove {
		case piece.White:
			for fromBB := pawns & bitboard.Rank7; fromBB != bitboard.Empty; {
				from := fromBB.Pop()

				if to := from - 8; bitboard.Squares[to]&occ == 0 {
					addPromotions(&moves, move.New(from, to, piece.WhitePawn, false), piece.White)
				}

				if to := from - 7; from.File() < square.FileH && bitboard.Squares[to]&enemies != 0 {
					addPromotions(&moves, move.New(from, to, piece.WhitePawn, true), piece.White)
				}

				if to := from - 9; from.File() > square.FileA && bitboard.Squares[to]&enemies != 0 {
					addPromotions(&moves, move.New(from, to, piece.WhitePawn, true), piece.White)
				}
			}

			for fromBB := pawns &^ bitboard.Rank7; fromBB != bitboard.Empty; {
				from := fromBB.Pop()

				if to := from - 8; bitboard.Squares[to]&occ == 0 {
					moves = append(moves, move.New(from, to, piece.WhitePawn, false))

					if to := from - 16; from.Rank() == square.Rank2 && bitboard.Squares[to]&occ == 0 {
						moves = append(moves, move.New(from, to, piece.WhitePawn, false))
					}
				}

				if to := from - 7; from.File() < square.FileH && bitboard.Squares[to]&enemies != 0 {
					moves = append(moves, move.New(from, to, piece.WhitePawn, true))
				}

				if to := from - 9; from.File() > square.FileA && bitboard.Squares[to]&enemies != 0 {
					moves = append(moves, move.New(from, to, piece.WhitePawn, true))
				}
			}
		case piece.Black:
			for fromBB := pawns & bitboard.Rank2; fromBB != bitboard.Empty; {
				from := fromBB.Pop()

				if to := from + 8; bitboard.Squares[to]&occ == 0 {
					addPromotions(&moves, move.New(from, to, piece.BlackPawn, false), piece.Black)
				}

				if to := from + 9; from.File() < square.FileH && bitboard.Squares[to]&enemies != 0 {
					addPromotions(&moves, move.New(from, to, piece.BlackPawn, true), piece.Black)
				}

				if to := from + 7; from.File() > square.FileA && bitboard.Squares[to]&enemies != 0 {
					addPromotions(&moves, move.New(from, to, piece.BlackPawn, true), piece.Black)
				}
			}

			for fromBB := pawns &^ bitboard.Rank2; fromBB != bitboard.Empty; {
				from := fromBB.Pop()

				if to := from + 8; bitboard.Squares[to]&occ == 0 {
					moves = append(moves, move.New(from, to, piece.BlackPawn, false))

					if to := from + 16; from.Rank() == square.Rank7 && bitboard.Squares[to]&occ == 0 {
						moves = append(moves, move.New(from, to, piece.BlackPawn, false))
					}
				}

				if to := from + 9; from.File() < square.FileH && bitboard.Squares[to]&enemies != 0 {
					moves = append(moves, move.New(from, to, piece.BlackPawn, true))
				}

				if to := from + 7; from.File() > square.FileA && bitboard.Squares[to]&enemies != 0 {
					moves = append(moves, move.New(from, to, piece.BlackPawn, true))
				}
			}
		}
	}

	p := piece.New(piece.Knight, b.SideToMove)
	for fromBB := b.PieceBBs[piece.Knight] & friends; fromBB != bitboard.Empty; {
		from := fromBB.Pop()

		for toBB := attacks.Knight[from] &^ friends; toBB != bitboard.Empty; {
			to := toBB.Pop()
			isCap := b.Position[to] != piece.NoPiece
			moves = append(moves, move.New(from, to, p, isCap))
		}
	}

	p = piece.New(piece.Bishop, b.SideToMove)
	for fromBB := b.PieceBBs[piece.Bishop] & friends; fromBB != bitboard.Empty; {
		from := fromBB.Pop()

		for toBB := attacks.Bishop(from, occ) &^ friends; toBB != bitboard.Empty; {
			to := toBB.Pop()
			isCap := b.Position[to] != piece.NoPiece
			moves = append(moves, move.New(from, to, p, isCap))
		}
	}

	p = piece.New(piece.Rook, b.SideToMove)
	for fromBB := b.PieceBBs[piece.Rook] & friends; fromBB != bitboard.Empty; {
		from := fromBB.Pop()

		for toBB := attacks.Rook(from, occ) &^ friends; toBB != bitboard.Empty; {
			to := toBB.Pop()
			isCap := b.Position[to] != piece.NoPiece
			moves = append(moves, move.New(from, to, p, isCap))
		}
	}

	p = piece.New(piece.Queen, b.SideToMove)
	for fromBB := b.PieceBBs[piece.Queen] & friends; fromBB != bitboard.Empty; {
		from := fromBB.Pop()

		for toBB := attacks.Queen(from, occ) &^ friends; toBB != bitboard.Empty; {
			to := toBB.Pop()
			isCap := b.Position[to] != piece.NoPiece
			moves = append(moves, move.New(from, to, p, isCap))
		}
	}

	{
		p = piece.New(piece.King, b.SideToMove)

		from := b.Kings[b.SideToMove]
		for toBB := attacks.King[from] &^ friends; toBB != bitboard.Empty; {
			to := toBB.Pop()
			isCap := b.Position[to] != piece.NoPiece
			moves = append(moves, move.New(from, to, p, isCap))
		}

		switch b.SideToMove {
		case piece.White:
			if b.CastlingRights&castling.WhiteA == castling.NoCasl ||
				b.IsAttacked(square.E1, piece.Black) {
				break
			}

			if b.CastlingRights&castling.WhiteK != 0 &&
				occ&0x6000000000000000 == bitboard.Empty &&
				!b.IsAttacked(square.F1, piece.Black) {
				moves = append(moves, move.New(from, square.G1, p, false))
			}

			if b.CastlingRights&castling.WhiteQ != 0 &&
				occ&0xe00000000000000 == bitboard.Empty &&
				!b.IsAttacked(square.D1, piece.Black) {
				moves = append(moves, move.New(from, square.C1, p, false))
			}
		case piece.Black:
			if b.CastlingRights&castling.BlackA == castling.NoCasl ||
				b.IsAttacked(square.E8, piece.White) {
				break
			}

			if b.CastlingRights&castling.BlackK != 0 &&
				occ&0x60 == bitboard.Empty &&
				!b.IsAttacked(square.F8, piece.White) {
				moves = append(moves, move.New(from, square.G8, p, false))
			}

			if b.CastlingRights&castling.BlackQ != 0 &&
				occ&0xe == bitboard.Empty &&
				!b.IsAttacked(square.D8, piece.White) {
				moves = append(moves, move.New(from, square.C8, p, false))
			}
		}
	}

	return moves
}

func addPromotions(moveList *[]move.Move, m move.Move, c piece.Color) {
	*moveList = append(*moveList,
		m.SetPromotion(piece.New(piece.Queen, c)),
		m.SetPromotion(piece.New(piece.Rook, c)),
		m.SetPromotion(piece.New(piece.Bishop, c)),
		m.SetPromotion(piece.New(piece.Knight, c)),
	)
}

func (b *Board) NewMove(from, to square.Square) move.Move {
	p := b.Position[from]
	return move.New(from, to, p, b.Position[to] != piece.NoPiece)
}
