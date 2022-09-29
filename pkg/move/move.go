package move

import (
	"fmt"

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

	HalfMoves int
	CastlingRights CastlingRights
	EnPassantSquare square.Square
}

func (m Move) String() string {
	str := fmt.Sprintf("%s%s", m.From, m.To)
	if m.IsPromotion() {
		str += (m.ToPiece + piece.Pawn).String()
	}
	return str
}

func (m Move) IsCapture() bool {
	return m.CapturedPiece != piece.Empty
}

func (m Move) IsEnPassant() bool {
	return m.FromPiece.Type() == piece.Pawn && m.To == m.EnPassantSquare
}

func (m Move) IsPromotion() bool {
	return m.FromPiece != m.ToPiece
}

func (m Move) IsDoublePawnPush() bool {
	if m.FromPiece != piece.Pawn {
		return false
	}

	switch m.From.Rank() - m.To.Rank() {
	case 32, -32:
		return true
	default:
		return false
	}
}
