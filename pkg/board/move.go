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

	{
		enemies := b.ColorBBs[b.SideToMove.Other()]
		enemies.Set(b.EnPassantTarget)

		up := square.Square(-8)
		down := square.Square(8)
		left := square.Square(-1)
		right := square.Square(1)
		promotionRank := bitboard.Rank7
		if b.SideToMove == piece.Black {
			up, down = down, up
			promotionRank = bitboard.Rank2
		}

		p := piece.New(piece.Pawn, b.SideToMove)

		pawns := b.PieceBBs[piece.Pawn] & friends
		normalPawns := pawns &^ promotionRank
		promotionPawns := pawns & promotionRank

		singlePush, doublePush := attacks.PawnPush(normalPawns, occ, b.SideToMove)
		leftAttacks := attacks.PawnLeft(normalPawns, enemies, b.SideToMove)
		rightAttacks := attacks.PawnRight(normalPawns, enemies, b.SideToMove)

		for singlePush != bitboard.Empty {
			to := singlePush.Pop()
			from := to + down
			moves = append(moves, move.New(from, to, p, false))
		}

		for doublePush != bitboard.Empty {
			to := doublePush.Pop()
			from := to + down + down
			moves = append(moves, move.New(from, to, p, false))
		}

		for leftAttacks != bitboard.Empty {
			to := leftAttacks.Pop()
			from := to + right + down
			moves = append(moves, move.New(from, to, p, true))
		}

		for rightAttacks != bitboard.Empty {
			to := rightAttacks.Pop()
			from := to + left + down
			moves = append(moves, move.New(from, to, p, true))
		}

		for promotionPawns != bitboard.Empty {
			from := promotionPawns.Pop()

			if to := from + up; bitboard.Squares[to]&occ == 0 {
				addPromotions(&moves, move.New(from, to, p, false), b.SideToMove)
			}

			if to := from + up + right; from.File() < square.FileH && bitboard.Squares[to]&enemies != 0 {
				addPromotions(&moves, move.New(from, to, p, true), b.SideToMove)
			}

			if to := from + up + left; from.File() > square.FileA && bitboard.Squares[to]&enemies != 0 {
				addPromotions(&moves, move.New(from, to, p, true), b.SideToMove)
			}
		}
	}

	for pType := piece.Knight; pType <= piece.King; pType++ {
		p := piece.New(pType, b.SideToMove)
		for fromBB := b.PieceBBs[pType] & friends; fromBB != bitboard.Empty; {
			from := fromBB.Pop()

			for toBB := b.MovesOf(pType, from, occ) &^ friends; toBB != bitboard.Empty; {
				to := toBB.Pop()
				moves = append(moves, move.New(from, to, p, occ.IsSet(to)))
			}
		}
	}

	switch b.SideToMove {
	case piece.White:
		if b.CastlingRights&castling.WhiteA == castling.NoCasl ||
			b.IsAttacked(square.E1, piece.Black) {
			break
		}

		if b.CastlingRights&castling.WhiteK != 0 &&
			occ&bitboard.F1G1 == bitboard.Empty &&
			!b.IsAttacked(square.F1, piece.Black) {
			moves = append(moves, move.New(square.E1, square.G1, piece.WhiteKing, false))
		}

		if b.CastlingRights&castling.WhiteQ != 0 &&
			occ&bitboard.B1C1D1 == bitboard.Empty &&
			!b.IsAttacked(square.D1, piece.Black) {
			moves = append(moves, move.New(square.E1, square.C1, piece.WhiteKing, false))
		}
	case piece.Black:
		if b.CastlingRights&castling.BlackA == castling.NoCasl ||
			b.IsAttacked(square.E8, piece.White) {
			break
		}

		if b.CastlingRights&castling.BlackK != 0 &&
			occ&bitboard.F8G8 == bitboard.Empty &&
			!b.IsAttacked(square.F8, piece.White) {
			moves = append(moves, move.New(square.E8, square.G8, piece.BlackKing, false))
		}

		if b.CastlingRights&castling.BlackQ != 0 &&
			occ&bitboard.B8C8D8 == bitboard.Empty &&
			!b.IsAttacked(square.D8, piece.White) {
			moves = append(moves, move.New(square.E8, square.C8, piece.BlackKing, false))
		}
	}

	return moves
}

func (b *Board) MovesOf(p piece.Type, s square.Square, blockers bitboard.Board) bitboard.Board {
	switch p {
	case piece.Knight:
		return attacks.Knight[s]
	case piece.Bishop:
		return attacks.Bishop(s, blockers)
	case piece.Rook:
		return attacks.Rook(s, blockers)
	case piece.Queen:
		return attacks.Queen(s, blockers)
	case piece.King:
		return attacks.King[s]
	default:
		panic("bad piece type")
	}
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
