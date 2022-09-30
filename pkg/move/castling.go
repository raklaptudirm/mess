package move

type CastlingRights byte

func CastlingRightsFrom(r string) CastlingRights {
	var rights CastlingRights

	if r == "-" {
		return CastleNone
	}

	if r != "" && r[0] == 'K' {
		r = r[1:]
		rights |= CastleWhiteKingside
	}

	if r != "" && r[0] == 'Q' {
		r = r[1:]
		rights |= CastleWhiteQueenside
	}

	if r != "" && r[0] == 'k' {
		r = r[1:]
		rights |= CastleBlackKingside
	}

	if r != "" && r[0] == 'q' {
		rights |= CastleBlackQueenside
	}

	return rights
}

const (
	CastleWhiteKingside  CastlingRights = 1 << 0
	CastleWhiteQueenside CastlingRights = 1 << 1
	CastleBlackKingside  CastlingRights = 1 << 2
	CastleBlackQueenside CastlingRights = 1 << 3

	CastleNone CastlingRights = 0

	CastleWhite CastlingRights = CastleWhiteKingside | CastleWhiteQueenside
	CastleBlack CastlingRights = CastleBlackKingside | CastleBlackQueenside

	CastleKingside  CastlingRights = CastleWhiteKingside | CastleBlackKingside
	CastleQueenside CastlingRights = CastleWhiteQueenside | CastleBlackQueenside

	CastleAll CastlingRights = CastleWhite | CastleBlack
)

func (c CastlingRights) String() string {
	var str string

	if c&CastleWhiteKingside != 0 {
		str += "K"
	}

	if c&CastleWhiteQueenside != 0 {
		str += "Q"
	}

	if c&CastleBlackKingside != 0 {
		str += "k"
	}

	if c&CastleBlackQueenside != 0 {
		str += "q"
	}

	if str == "" {
		str = "-"
	}

	return str
}
