package castling

type Rights byte

func NewRights(r string) Rights {
	var rights Rights

	if r == "-" {
		return None
	}

	if r != "" && r[0] == 'K' {
		r = r[1:]
		rights |= WhiteKingside
	}

	if r != "" && r[0] == 'Q' {
		r = r[1:]
		rights |= WhiteQueenside
	}

	if r != "" && r[0] == 'k' {
		r = r[1:]
		rights |= BlackKingside
	}

	if r != "" && r[0] == 'q' {
		rights |= BlackQueenside
	}

	return rights
}

const (
	WhiteKingside  Rights = 1 << 0
	WhiteQueenside Rights = 1 << 1
	BlackKingside  Rights = 1 << 2
	BlackQueenside Rights = 1 << 3

	None Rights = 0

	White Rights = WhiteKingside | WhiteQueenside
	Black Rights = BlackKingside | BlackQueenside

	Kingside  Rights = WhiteKingside | BlackKingside
	Queenside Rights = WhiteQueenside | BlackQueenside

	All Rights = White | Black

	N = 16
)

func (c Rights) String() string {
	var str string

	if c&WhiteKingside != 0 {
		str += "K"
	}

	if c&WhiteQueenside != 0 {
		str += "Q"
	}

	if c&BlackKingside != 0 {
		str += "k"
	}

	if c&BlackQueenside != 0 {
		str += "q"
	}

	if str == "" {
		str = "-"
	}

	return str
}
