package board

import (
	"fmt"

	"laptudirm.com/x/mess/pkg/piece"
	"laptudirm.com/x/mess/pkg/square"
)

type Move struct {
	From      square.Square
	To        square.Square
	Promotion piece.Piece
}

func (m Move) String() string {
	str := fmt.Sprintf("%s%s", m.From, m.To)
	if m.Promotion != piece.Empty {
		str += (m.Promotion + piece.Pawn).String()
	}
	return str
}
