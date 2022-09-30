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

	switch m.From.Rank() - m.To.Rank() {
	case 32, -32:
		return true
	default:
		return false
	}
}
