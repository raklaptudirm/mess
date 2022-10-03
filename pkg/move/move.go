package move

import (
	"fmt"

	"laptudirm.com/x/mess/pkg/castling"
	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
)

type Move struct {
	From    square.Square
	To      square.Square
	Capture square.Square

	FromPiece     piece.Piece
	ToPiece       piece.Piece
	CapturedPiece piece.Piece

	HalfMoves       int
	CastlingRights  castling.Rights
	EnPassantSquare square.Square
}

func (m Move) String() string {
	str := fmt.Sprintf("%s%s", m.From, m.To)
	if m.IsPromotion() {
		str += m.ToPiece.Type().String()
	}
	return str
}

func (m Move) CastlingRightUpdates() castling.Rights {
	toRemove := castling.None // castling rights to remove

	// update castling rights
	// movement of the rooks or the king, or the capture of the rooks
	// leads to losing the right to castle: update it according to the move

	// rooks or king moved
	switch m.From {
	// white rights
	case square.H1:
		// kingside rook moved
		toRemove |= castling.WhiteKingside
	case square.A1:
		// queenside rook moved
		toRemove |= castling.WhiteQueenside
	case square.E1:
		// king moved
		toRemove |= castling.White

	// black rights
	case square.H8:
		// kingside rook moved
		toRemove |= castling.BlackKingside
	case square.A8:
		// queenside rook moved
		toRemove |= castling.BlackQueenside
	case square.E8:
		// king moved
		toRemove |= castling.Black
	}

	// rooks captured
	switch m.To {
	// white rooks
	case square.H1:
		toRemove |= castling.WhiteKingside
	case square.A1:
		toRemove |= castling.WhiteQueenside

	// black rooks
	case square.H8:
		toRemove |= castling.BlackKingside
	case square.A8:
		toRemove |= castling.BlackKingside
	}

	return toRemove
}

func (m Move) IsReversible() bool {
	return !m.IsCapture() && m.FromPiece.Type() != piece.Pawn
}

func (m Move) IsCastle() bool {
	switch m.FromPiece {
	case piece.WhiteKing:
		return m.From == square.E1 && (m.To == square.G1 || m.To == square.C1)
	case piece.BlackKing:
		return m.From == square.E8 && (m.To == square.G8 || m.To == square.C8)
	default:
		return false
	}
}

func (m Move) IsCapture() bool {
	return m.CapturedPiece != piece.NoPiece
}

func (m Move) IsEnPassant() bool {
	return m.FromPiece.Type() == piece.Pawn && m.To == m.EnPassantSquare
}

func (m Move) IsPromotion() bool {
	return m.FromPiece != m.ToPiece
}

func (m Move) IsDoublePawnPush() bool {
	if m.FromPiece.Type() != piece.Pawn {
		return false
	}

	fromRank := m.From.Rank()
	toRank := m.To.Rank()

	switch {
	case fromRank == square.Rank2 && toRank == square.Rank4,
		fromRank == square.Rank7 && toRank == square.Rank5:
		return true
	default:
		return false
	}
}
