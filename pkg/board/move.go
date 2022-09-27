package board

import (
	"fmt"

	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
)

type Move struct {
	From square.Square
	To   square.Square

	Promotion piece.Piece

	IsEnPassant bool
}

func (m Move) String() string {
	str := fmt.Sprintf("%s%s", m.From, m.To)
	if m.Promotion != piece.Empty {
		str += (m.Promotion + piece.Pawn).String()
	}
	return str
}

func (m Move) IsPromotion() bool {
	return m.Promotion != piece.Empty
}

func (m Move) IsDoublePawnPush() bool {
	fromRank := m.From.Rank()
	toRank := m.To.Rank()

	switch {
	case fromRank != square.Rank2 && fromRank != square.Rank7:
		return false
	case toRank != square.Rank4 && toRank != square.Rank5:
		return false
	default:
		return true
	}
}
