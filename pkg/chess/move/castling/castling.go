package castling

import "laptudirm.com/x/mess/pkg/chess/square"

type Rights byte

func NewRights(r string) Rights {
	var rights Rights

	if r == "-" {
		return NoCasl
	}

	if r != "" && r[0] == 'K' {
		r = r[1:]
		rights |= WhiteK
	}

	if r != "" && r[0] == 'Q' {
		r = r[1:]
		rights |= WhiteQ
	}

	if r != "" && r[0] == 'k' {
		r = r[1:]
		rights |= BlackK
	}

	if r != "" && r[0] == 'q' {
		rights |= BlackQ
	}

	return rights
}

const (
	WhiteK Rights = 1 << 0
	WhiteQ Rights = 1 << 1
	BlackK Rights = 1 << 2
	BlackQ Rights = 1 << 3

	NoCasl Rights = 0

	WhiteA Rights = WhiteK | WhiteQ
	BlackA Rights = BlackK | BlackQ

	Kingside  Rights = WhiteK | BlackK
	Queenside Rights = WhiteQ | BlackQ

	All Rights = WhiteA | BlackA

	N = 16
)

var RightUpdates = [square.N]Rights{
	BlackQ, NoCasl, NoCasl, NoCasl, BlackA, NoCasl, NoCasl, BlackK,
	NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl,
	NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl,
	NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl,
	NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl,
	NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl,
	NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl, NoCasl,
	WhiteQ, NoCasl, NoCasl, NoCasl, WhiteA, NoCasl, NoCasl, WhiteK,
}

func (c Rights) String() string {
	var str string

	if c&WhiteK != 0 {
		str += "K"
	}

	if c&WhiteQ != 0 {
		str += "Q"
	}

	if c&BlackK != 0 {
		str += "k"
	}

	if c&BlackQ != 0 {
		str += "q"
	}

	if str == "" {
		str = "-"
	}

	return str
}
