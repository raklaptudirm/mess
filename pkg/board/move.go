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
		if b.Pawns(b.SideToMove.Other())&attacks.Pawn[b.SideToMove][target] != 0 {
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
	b.CalculateCheckmask()
	b.CalculatePinmask()
	b.SeenByEnemy = b.SeenSquares(b.SideToMove.Other())

	moves := make([]move.Move, 0, 30)
	occ := b.Occupied()

	us := b.SideToMove
	friends := b.ColorBBs[us]

	target := ^friends & b.CheckMask

	{
		kingSq := b.Kings[us]
		king := piece.New(piece.King, us)
		for toBB := attacks.King[kingSq] &^ friends &^ b.SeenByEnemy; toBB != bitboard.Empty; {
			to := toBB.Pop()
			moves = append(moves, move.New(kingSq, to, king, occ.IsSet(to)))
		}
	}

	switch b.CheckN {
	case 0:
		b.genCastlingMoves(&moves)
	case 2:
		return moves
	}

	for pType := piece.Knight; pType <= piece.Queen; pType++ {
		p := piece.New(pType, us)
		for fromBB := b.PieceBBs[pType] & friends; fromBB != bitboard.Empty; {
			from := fromBB.Pop()

			for toBB := b.MovesOf(pType, from) & target; toBB != bitboard.Empty; {
				to := toBB.Pop()
				moves = append(moves, move.New(from, to, p, occ.IsSet(to)))
			}
		}
	}

	b.genPawnMoves(&moves)

	return moves
}

func (b *Board) genPawnMoves(moveList *[]move.Move) {
	us := b.SideToMove

	occ := b.Occupied()

	enemies := b.ColorBBs[us.Other()]
	enemies.Set(b.EnPassantTarget)

	var up, down, left, right square.Square
	var promotionRank bitboard.Board
	var p piece.Piece

	left = -1
	right = 1

	switch us {
	case piece.White:
		up = -8
		down = 8

		promotionRank = bitboard.Rank8

		p = piece.WhitePawn

	case piece.Black:
		up = 8
		down = -8

		promotionRank = bitboard.Rank1

		p = piece.BlackPawn
	}

	pushTarget := b.CheckMask &^ occ
	captureTarget := enemies & b.CheckMask

	pawns := b.Pawns(us)

	pawnsThatAttack := pawns &^ b.PinnedHV

	unpinnedPawnsThatAttack := pawnsThatAttack &^ b.PinnedD
	pinnedPawnsThatAttack := pawnsThatAttack & b.PinnedD

	pawnAttacksL := attacks.PawnsLeft(unpinnedPawnsThatAttack, us) & captureTarget
	pawnAttacksL |= attacks.PawnsLeft(pinnedPawnsThatAttack, us) & captureTarget & b.PinnedD

	pawnAttacksR := attacks.PawnsRight(unpinnedPawnsThatAttack, us) & captureTarget
	pawnAttacksR |= attacks.PawnsRight(pinnedPawnsThatAttack, us) & captureTarget & b.PinnedD

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
		addPromotions(moveList, move.New(from, to, p, true), us)
	}

	for promotionPawnAttacksR != bitboard.Empty {
		to := promotionPawnAttacksR.Pop()
		from := to + down + left
		addPromotions(moveList, move.New(from, to, p, true), us)
	}

	pawnsThatPush := pawns &^ b.PinnedD

	unpinnedPawnsThatPush := pawnsThatPush &^ b.PinnedHV
	pinnedPawnsThatPush := pawnsThatPush & b.PinnedHV

	pawnPushesSingle := attacks.PawnPushSingle(unpinnedPawnsThatPush, us) & pushTarget
	pawnPushesSingle |= attacks.PawnPushSingle(pinnedPawnsThatPush, us) & pushTarget & b.PinnedHV

	pawnPushesDouble := attacks.PawnPushDouble(unpinnedPawnsThatPush, us) & pushTarget
	pawnPushesDouble |= attacks.PawnPushDouble(pinnedPawnsThatPush, us) & pushTarget & b.PinnedHV

	simplePawnPushes := pawnPushesSingle &^ promotionRank

	for simplePawnPushes != bitboard.Empty {
		to := simplePawnPushes.Pop()
		from := to + down
		*moveList = append(*moveList, move.New(from, to, p, false))

		if to += up; pawnPushesDouble.IsSet(to) {
			*moveList = append(*moveList, move.New(from, to, p, false))
		}
	}

	promotionPawnPushes := pawnPushesSingle & promotionRank

	for promotionPawnPushes != bitboard.Empty {
		to := promotionPawnPushes.Pop()
		from := to + down
		addPromotions(moveList, move.New(from, to, p, false), us)
	}
}

func (b *Board) genCastlingMoves(moveList *[]move.Move) {
	occ := b.Occupied()

	switch b.SideToMove {
	case piece.White:
		if b.CastlingRights&castling.WhiteA == castling.NoCasl ||
			b.IsAttacked(square.E1, piece.Black) {
			break
		}

		if b.CastlingRights&castling.WhiteK != 0 &&
			occ&bitboard.F1G1 == bitboard.Empty &&
			!b.IsAttacked(square.F1, piece.Black) {
			*moveList = append(*moveList, move.New(square.E1, square.G1, piece.WhiteKing, false))
		}

		if b.CastlingRights&castling.WhiteQ != 0 &&
			occ&bitboard.B1C1D1 == bitboard.Empty &&
			!b.IsAttacked(square.D1, piece.Black) {
			*moveList = append(*moveList, move.New(square.E1, square.C1, piece.WhiteKing, false))
		}
	case piece.Black:
		if b.CastlingRights&castling.BlackA == castling.NoCasl ||
			b.IsAttacked(square.E8, piece.White) {
			break
		}

		if b.CastlingRights&castling.BlackK != 0 &&
			occ&bitboard.F8G8 == bitboard.Empty &&
			!b.IsAttacked(square.F8, piece.White) {
			*moveList = append(*moveList, move.New(square.E8, square.G8, piece.BlackKing, false))
		}

		if b.CastlingRights&castling.BlackQ != 0 &&
			occ&bitboard.B8C8D8 == bitboard.Empty &&
			!b.IsAttacked(square.D8, piece.White) {
			*moveList = append(*moveList, move.New(square.E8, square.C8, piece.BlackKing, false))
		}
	}
}

func (b *Board) MovesOf(p piece.Type, s square.Square) bitboard.Board {
	switch p {
	case piece.Knight:
		return b.knightMoves(s)
	case piece.Bishop:
		return b.bishopMoves(s)
	case piece.Rook:
		return b.rookMoves(s)
	case piece.Queen:
		return b.queenMoves(s)
	default:
		panic("bad piece type")
	}
}

func (b *Board) knightMoves(s square.Square) bitboard.Board {
	switch {
	case b.PinnedD.IsSet(s),
		b.PinnedHV.IsSet(s):
		return bitboard.Empty
	default:
		return attacks.Knight[s]
	}
}

func (b *Board) bishopMoves(s square.Square) bitboard.Board {
	blockers := b.Occupied()

	switch {
	case b.PinnedHV.IsSet(s):
		return bitboard.Empty
	case b.PinnedD.IsSet(s):
		return attacks.Bishop(s, blockers) & b.PinnedD
	default:
		return attacks.Bishop(s, blockers)
	}
}

func (b *Board) rookMoves(s square.Square) bitboard.Board {
	blockers := b.Occupied()

	switch {
	case b.PinnedD.IsSet(s):
		return bitboard.Empty
	case b.PinnedHV.IsSet(s):
		return attacks.Rook(s, blockers) & b.PinnedHV
	default:
		return attacks.Rook(s, blockers)
	}
}

func (b *Board) queenMoves(s square.Square) bitboard.Board {
	return b.bishopMoves(s) | b.rookMoves(s)
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
