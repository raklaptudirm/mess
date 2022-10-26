package evaluation

import (
	"laptudirm.com/x/mess/pkg/chess"
	"laptudirm.com/x/mess/pkg/chess/move"
)

type Move int16

func OfMove(b *chess.Board, m move.Move) Move {
	switch {
	case m.IsPromotion():
		return Move(material[m.ToPiece().Type()])
	case m.IsCapture():
		return Move(material[b.Position[m.Target()].Type()] - material[m.FromPiece().Type()])
	default:
		return Move(squareBonuses[m.FromPiece()][m.Target()] - squareBonuses[m.FromPiece()][m.Source()])
	}
}
