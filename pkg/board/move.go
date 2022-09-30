package board

import (
	"fmt"

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
	if attackSet := b.MovesOf(m.From); !attackSet.IsSet(m.To) {
		// move not in attack board, illegal move
		panic(fmt.Sprintf("invalid move: piece can't move to given square\n%s", attackSet))
	}

	// update the half-move clock
	// it records the number of plys since the last pawn push or capture
	// for positions which are drawn by the 50-move rule
	switch {
	case m.FromPiece.Type() == piece.Pawn, m.IsCapture():
		// pawn push or capture: reset clock
		b.halfMoves = 0
	default:
		b.halfMoves++
	}

	// update castling rights
	// movement of the rooks or the king, or the capture of the rooks
	// leads to losing the right to castle: update it according to the move

	// rooks or king moved
	switch m.From {
	// white rights
	case square.H1:
		// kingside rook moved
		b.castlingRights &^= castling.WhiteKingside
	case square.A1:
		// queenside rook moved
		b.castlingRights &^= castling.WhiteQueenside
	case square.E1:
		// king moved
		b.castlingRights &^= castling.WhiteKingside
		b.castlingRights &^= castling.WhiteQueenside

	// black rights
	case square.H8:
		// kingside rook moved
		b.castlingRights &^= castling.BlackKingside
	case square.A8:
		// queenside rook moved
		b.castlingRights &^= castling.BlackQueenside
	case square.E8:
		// king moved
		b.castlingRights &^= castling.BlackKingside
		b.castlingRights &^= castling.BlackQueenside
	}

	// rooks captured
	switch m.To {
	// white rooks
	case square.H1:
		b.castlingRights &^= castling.WhiteKingside
	case square.A1:
		b.castlingRights &^= castling.WhiteQueenside

	// black rooks
	case square.H8:
		b.castlingRights &^= castling.BlackKingside
	case square.A8:
		b.castlingRights &^= castling.BlackKingside
	}

	b.hash ^= zobrist.Castling[m.CastlingRights] // remove old rights
	b.hash ^= zobrist.Castling[b.castlingRights] // put new rights

	// move the piece in the records

	if m.IsCapture() {
		b.ClearSquare(m.Capture)
	}

	b.ClearSquare(m.From)
	b.FillSquare(m.To, m.ToPiece)

	// update en passant target square
	// clear the previous square, and if current move was double a pawn
	// push, add set the en passant target to the new square

	if b.enPassantTarget != square.None {
		// remove previous square from zobrist hash
		b.hash ^= zobrist.EnPassant[b.enPassantTarget.File()]
	}

	// reset en passant target
	b.enPassantTarget = square.None

	if m.IsDoublePawnPush() {
		// double pawn push; set new en passant target
		b.enPassantTarget = m.From
		if b.sideToMove == piece.White {
			b.enPassantTarget += 8
		} else {
			b.enPassantTarget -= 8
		}

		// and new square to zobrist hash
		b.hash ^= zobrist.EnPassant[b.enPassantTarget.File()]
	}

	// switch turn

	// update side to move
	if b.sideToMove = b.sideToMove.Other(); b.sideToMove == piece.White {
		b.fullMoves++
	}

	// switch bitboards
	b.friends, b.enemies = b.enemies, b.friends

	// switch zobrist hash
	b.hash ^= zobrist.SideToMove
}

func (b *Board) Unmove(m move.Move) {
	if b.enPassantTarget != square.None {
		b.hash ^= zobrist.EnPassant[b.enPassantTarget.File()]
		b.enPassantTarget = square.None
	}

	if m.EnPassantSquare != square.None {
		b.enPassantTarget = m.EnPassantSquare
		b.hash ^= zobrist.EnPassant[b.enPassantTarget.File()]
	}

	b.ClearSquare(m.To)
	b.FillSquare(m.From, m.FromPiece)

	if m.IsCapture() {
		b.FillSquare(m.Capture, m.CapturedPiece)
	}

	b.hash ^= zobrist.Castling[b.castlingRights]
	b.hash ^= zobrist.Castling[m.CastlingRights]
	b.castlingRights = m.CastlingRights

	b.halfMoves = m.HalfMoves

	b.hash ^= zobrist.SideToMove

	b.friends, b.enemies = b.enemies, b.friends

	// update side to move
	if b.sideToMove = b.sideToMove.Other(); b.sideToMove == piece.Black {
		b.fullMoves--
	}
}

func (b *Board) GenerateMoves() []move.Move {
	var moves []move.Move

	for i := 0; i < square.N; i++ {
		from := square.Square(i)
		moveSet := b.MovesOf(from)

		switch b.position[from].Type() {
		// handle pawns separately for en passant and promotions
		case piece.Pawn:
			for j := 0; j < square.N && moveSet != bitboard.Empty; j++ {
				to := square.Square(j)
				if !moveSet.IsSet(to) {
					continue
				}

				m := move.Move{
					From:    from,
					To:      to,
					Capture: to,

					FromPiece:     b.position[from],
					ToPiece:       b.position[from],
					CapturedPiece: b.position[to],

					HalfMoves:       b.halfMoves,
					CastlingRights:  b.castlingRights,
					EnPassantSquare: b.enPassantTarget,
				}

				switch {
				// pawn will promote
				case b.sideToMove == piece.White && to.Rank() == square.Rank8:
					fallthrough
				case b.sideToMove == piece.Black && to.Rank() == square.Rank1:
					// evaluate all possible promotions
					for _, promotion := range piece.Promotions {
						m.ToPiece = piece.New(promotion, b.sideToMove)
						moves = append(moves, m)
					}

				// en passant capture
				case to == b.enPassantTarget:
					m.Capture = to
					if b.sideToMove == piece.White {
						m.Capture += 8
					} else {
						m.Capture -= 8
					}
					m.CapturedPiece = b.position[m.Capture]

				// simple push or capture
				default:
					moves = append(moves, m)
				}

				moveSet.Unset(to)
			}

		// other pieces move simply
		default:
			for j := 0; j < square.N && moveSet != bitboard.Empty; j++ {
				to := square.Square(j)
				if moveSet.IsSet(to) {
					m := move.Move{
						From: from,
						To:   to,

						FromPiece:     b.position[from],
						ToPiece:       b.position[from],
						CapturedPiece: b.position[to],

						HalfMoves:       b.halfMoves,
						CastlingRights:  b.castlingRights,
						EnPassantSquare: b.enPassantTarget,
					}
					moves = append(moves, m)
				}
				moveSet.Unset(to)
			}
		}
	}

	return moves
}

func (b *Board) MovesOf(index square.Square) bitboard.Board {
	p := b.position[index]
	if p.Color() != b.sideToMove {
		// other side has no moves
		return bitboard.Empty
	}

	switch p.Type() {
	case piece.King:
		return attacks.King(index, b.friends, b.friends|b.enemies, b.castlingRights)
	case piece.Queen:
		return attacks.Queen(index, b.friends, b.friends|b.enemies)
	case piece.Rook:
		return attacks.Rook(index, b.friends, b.friends|b.enemies)
	case piece.Knight:
		return attacks.Knight(index, b.friends)
	case piece.Bishop:
		return attacks.Bishop(index, b.friends, b.friends|b.enemies)
	case piece.Pawn:
		return attacks.Pawn(index, b.enPassantTarget, b.sideToMove, b.friends, b.enemies)
	default:
		return bitboard.Empty
	}
}
