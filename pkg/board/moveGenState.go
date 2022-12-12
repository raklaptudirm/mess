package board

import (
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/board/move"
	"laptudirm.com/x/mess/pkg/board/move/attacks"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
)

// moveGenState stores various utility and generated data used during move
// generation. It is separate from board since this data is not necessary
// in the board representation.
type moveGenState struct {
	// board from which the moves are generated
	*Board

	// movelist that stores the generated moves
	MoveList []move.Move

	// utility information:
	// these data, used by movegen can be calculated from the main board
	// state but are time expensive or simply tedious to write in full,
	// and so are stored in Board instead.

	Us, Them piece.Color

	// adding Down to a square gives the square "below"
	// "below" is towards the player's own side
	Down square.Square

	// rank where the pawns get promoted
	PromotionRankBB bitboard.Board

	// rank at which pawns can en passant
	EnPassantRankBB bitboard.Board

	// rank from which a single push is the same
	// as a double push from home
	DoublePushRankBB bitboard.Board

	// king position look-up table
	Kings [piece.ColorN]square.Square

	// movegen type (tactical moves only?)
	TacticalOnly bool

	// color bitboards classified by side to move
	Friends bitboard.Board
	Enemies bitboard.Board

	// precalculated Friends | Enemies
	Occupied bitboard.Board

	// places where pieces can move to
	// calculated as ^Friends & CheckMask
	Target bitboard.Board
	// king target is special because it is
	// Enemies &^ SeenByEnemy
	KingTarget bitboard.Board

	// check information
	CheckN    int            // number of checkers [0, 2]
	CheckMask bitboard.Board // see docs for CalculateCheckmask

	// pinned piece information
	// see docs for CalculatePinmask
	PinnedD  bitboard.Board
	PinnedHV bitboard.Board

	// squares attacked by enemy pieces
	SeenByEnemy bitboard.Board

	// piece variables containing pieces of current color
	Pawn, Knight, Bishop, Rook, Queen, King piece.Piece
}

// AppendMoves appends the given moves to the current state's movelist.
func (s *moveGenState) AppendMoves(m ...move.Move) {
	s.MoveList = append(s.MoveList, m...)
}

// Init initializes all the different utility bitboards which are
// calculated and necessary for move generation.
func (s *moveGenState) Init(captureOnly bool) {
	// king positions for each color
	s.Kings[piece.White] = s.KingBB(piece.White).FirstOne()
	s.Kings[piece.Black] = s.KingBB(piece.Black).FirstOne()

	// move generation type
	s.TacticalOnly = captureOnly

	// occupancy bitboards
	s.Friends = s.ColorBBs[s.SideToMove]
	s.Enemies = s.ColorBBs[s.SideToMove.Other()]
	s.Occupied = s.Friends | s.Enemies

	// our and their colors
	s.Us = s.SideToMove
	s.Them = s.Us.Other()

	// side to move dependent variables
	if s.Us == piece.White {
		s.PromotionRankBB = bitboard.Rank8
		s.EnPassantRankBB = bitboard.Rank5
		s.DoublePushRankBB = bitboard.Rank3

		s.Down = 8

		s.Pawn = piece.WhitePawn
		s.Knight = piece.WhiteKnight
		s.Bishop = piece.WhiteBishop
		s.Rook = piece.WhiteRook
		s.Queen = piece.WhiteQueen
		s.King = piece.WhiteKing
	} else {
		s.PromotionRankBB = bitboard.Rank1
		s.EnPassantRankBB = bitboard.Rank4
		s.DoublePushRankBB = bitboard.Rank6

		s.Down = -8

		s.Pawn = piece.BlackPawn
		s.Knight = piece.BlackKnight
		s.Bishop = piece.BlackBishop
		s.Rook = piece.BlackRook
		s.Queen = piece.BlackQueen
		s.King = piece.BlackKing
	}

	s.CalculateCheckmask()
	s.CalculatePinmask()

	s.SeenByEnemy = s.SeenSquares(s.SideToMove.Other())

	// move generation type dependent variables
	if captureOnly {
		s.Target = s.Enemies & s.CheckMask
		s.KingTarget = s.Enemies &^ s.SeenByEnemy
	} else {
		s.Target = ^s.Friends & s.CheckMask
		s.KingTarget = ^s.Friends &^ s.SeenByEnemy
	}

	// 31 is the average number of chess s.MoveList in a position
	// source: https://chess.stackexchange.com/a/24325/33336
	s.MoveList = make([]move.Move, 0, 31)
}

// CalculateCheckmask calculates the check-mask of the current board state,
// along with the number of checkers.
//
// A checker is an enemy piece which is directly checking the king. The
// number of checkers can be a maximum of two (double check).
//
// The check-mask is defined as all the squares to which if a friendly
// piece is moved to will block all checks. This is defined as empty for
// double check, the checking piece and, if the checker is a sliding piece,
// the squares between the king and the checker. The bitboard is universe
// if the king is not in check.
func (s *moveGenState) CalculateCheckmask() {
	s.CheckN = 0
	s.CheckMask = bitboard.Empty

	kingSq := s.Kings[s.Us]

	pawns := s.PawnsBB(s.Them) & attacks.Pawn[s.Us][kingSq]
	knights := s.KnightsBB(s.Them) & attacks.Knight[kingSq]
	bishops := (s.BishopsBB(s.Them) | s.QueensBB(s.Them)) & attacks.Bishop(kingSq, s.Occupied)
	rooks := (s.RooksBB(s.Them) | s.QueensBB(s.Them)) & attacks.Rook(kingSq, s.Occupied)

	// a pawn and a knight cannot be checking the king at the same time as
	// they are not sliding pieces thus discovered attacks are impossible
	switch {
	case pawns != bitboard.Empty:
		s.CheckMask |= pawns
		s.CheckN++

	case knights != bitboard.Empty:
		s.CheckMask |= knights
		s.CheckN++
	}

	if bishops != bitboard.Empty {
		bishopSq := bishops.FirstOne()
		s.CheckMask |= bitboard.Between[kingSq][bishopSq] | bitboard.Squares[bishopSq]
		s.CheckN++
	}

	// 2 is the largest possible value for CheckN so short circuit if thats reached
	if s.CheckN < 2 && rooks != bitboard.Empty {
		if s.CheckN == 0 && rooks.Count() > 1 {
			// double check, don't set the check-mask
			s.CheckN++
		} else {
			rookSq := rooks.FirstOne()
			s.CheckMask |= bitboard.Between[kingSq][rookSq] | bitboard.Squares[rookSq]
			s.CheckN++
		}
	}

	if s.CheckN == 0 {
		// king is not in check so check-mask is universe
		s.CheckMask = bitboard.Universe
	}
}

// CalculatePinmask calculates the horizontal and vertical pin-masks.
// A pin-mask is defined as the mask containing all attack rays pieces
// pinning a piece in a given direction (horizontal or vertical).
func (s *moveGenState) CalculatePinmask() {
	kingSq := s.Kings[s.Us]

	friends := s.ColorBBs[s.Us]
	enemies := s.ColorBBs[s.Them]

	s.PinnedD = bitboard.Empty
	s.PinnedHV = bitboard.Empty

	// consider enemy rooks and queens which are attacking or would attack the king if not for intervening pieces
	// the king is considered as a rook for this and it's attack sets & with rooks and queens gives the bitboard
	for rooks := (s.RooksBB(s.Them) | s.QueensBB(s.Them)) & attacks.Rook(kingSq, enemies); rooks != bitboard.Empty; {
		rook := rooks.Pop()
		possiblePin := bitboard.Between[kingSq][rook] | bitboard.Squares[rook]

		// if there is only one friendly piece blocking the ray, it is pinned
		if (possiblePin & friends).Count() == 1 {
			s.PinnedHV |= possiblePin
		}
	}

	// consider enemy bishops and queens which are attacking or would attack the king if not for intervening pieces
	// the king is considered as a bishop for this and it's attack sets & with bishops and queens gives the bitboard
	for bishops := (s.BishopsBB(s.Them) | s.QueensBB(s.Them)) & attacks.Bishop(kingSq, enemies); bishops != bitboard.Empty; {
		bishop := bishops.Pop()
		possiblePin := bitboard.Between[kingSq][bishop] | bitboard.Squares[bishop]

		// if there is only one friendly piece blocking the ray, it is pinned
		if (possiblePin & friends).Count() == 1 {
			s.PinnedD |= possiblePin
		}
	}
}

// SeenSquares returns a bitboard containing all the squares that are
// seen(attacked) by pieces of the given color. The enemy king is not
// considered as a sliding ray blocker by SeenSquares since it has to
// move away from the attack exposing the blocked squares.
func (s *moveGenState) SeenSquares(by piece.Color) bitboard.Board {
	pawns := s.PawnsBB(by)
	knights := s.KnightsBB(by)
	bishops := s.BishopsBB(by)
	rooks := s.RooksBB(by)
	queens := s.QueensBB(by)
	kingSq := s.Kings[by]

	// don't consider the enemy king as a blocker
	blockers := s.Occupied &^ s.KingBB(by.Other())

	seen := attacks.PawnsLeft(pawns, by) | attacks.PawnsRight(pawns, by)

	for knights != bitboard.Empty {
		from := knights.Pop()
		seen |= attacks.Knight[from]
	}

	for bishops != bitboard.Empty {
		from := bishops.Pop()
		seen |= attacks.Bishop(from, blockers)
	}

	for rooks != bitboard.Empty {
		from := rooks.Pop()
		seen |= attacks.Rook(from, blockers)
	}

	for queens != bitboard.Empty {
		from := queens.Pop()
		seen |= attacks.Queen(from, blockers)
	}

	seen |= attacks.King[kingSq]

	return seen
}
