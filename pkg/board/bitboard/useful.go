package bitboard

import "laptudirm.com/x/mess/pkg/square"

// useful bitboard definitions
const (
	Empty    Board = 0
	Universe Board = 0xffffffffffffffff
)

// file bitboards
const (
	FileA Board = 0x0101010101010101
	FileB Board = 0x0202020202020202
	FileC Board = 0x0404040404040404
	FileD Board = 0x0808080808080808
	FileE Board = 0x1010101010101010
	FileF Board = 0x2020202020202020
	FileG Board = 0x4040404040404040
	FileH Board = 0x8080808080808080
)

var Files = [...]Board{
	square.FileA: FileA,
	square.FileB: FileB,
	square.FileC: FileC,
	square.FileD: FileD,
	square.FileE: FileE,
	square.FileF: FileF,
	square.FileG: FileG,
	square.FileH: FileH,
}

// rank bitboards
const (
	Rank1 Board = 0xff00000000000000
	Rank2 Board = 0x00ff000000000000
	Rank3 Board = 0x0000ff0000000000
	Rank4 Board = 0x000000ff00000000
	Rank5 Board = 0x00000000ff000000
	Rank6 Board = 0x0000000000ff0000
	Rank7 Board = 0x000000000000ff00
	Rank8 Board = 0x00000000000000ff
)

var Ranks = [...]Board{
	square.Rank1: Rank1,
	square.Rank2: Rank2,
	square.Rank3: Rank3,
	square.Rank4: Rank4,
	square.Rank5: Rank5,
	square.Rank6: Rank6,
	square.Rank7: Rank7,
	square.Rank8: Rank8,
}

// diagonal bitboards
const (
	DiagonalH1H1 Board = 0x8000000000000000
	DiagonalH2G1 Board = 0x4080000000000000
	DiagonalH3F1 Board = 0x2040800000000000
	DiagonalH4E1 Board = 0x1020408000000000
	DiagonalH5D1 Board = 0x0810204080000000
	DiagonalH6C1 Board = 0x0408102040800000
	DiagonalH7B1 Board = 0x0204081020408000

	DiagonalH8A1 Board = 0x0102040810204080

	DiagonalG8A2 Board = 0x0001020408102040
	DiagonalF8A3 Board = 0x0000010204081020
	DiagonalE8A4 Board = 0x0000000102040810
	DiagonalD8A5 Board = 0x0000000001020408
	DiagonalC8A6 Board = 0x0000000000010204
	DiagonalB8A7 Board = 0x0000000000000102
	DiagonalA8A8 Board = 0x0000000000000001
)

var Diagonals = [...]Board{
	square.DiagonalH1H1: DiagonalH1H1,
	square.DiagonalH2G1: DiagonalH2G1,
	square.DiagonalH3F1: DiagonalH3F1,
	square.DiagonalH4E1: DiagonalH4E1,
	square.DiagonalH5D1: DiagonalH5D1,
	square.DiagonalH6C1: DiagonalH6C1,
	square.DiagonalH7B1: DiagonalH7B1,

	square.DiagonalH8A1: DiagonalH8A1,

	square.DiagonalG8A2: DiagonalG8A2,
	square.DiagonalF8A3: DiagonalF8A3,
	square.DiagonalE8A4: DiagonalE8A4,
	square.DiagonalD8A5: DiagonalD8A5,
	square.DiagonalC8A6: DiagonalC8A6,
	square.DiagonalB8A7: DiagonalB8A7,
	square.DiagonalA8A8: DiagonalA8A8,
}

// anti-diagonal bitboards
const (
	DiagonalA1A1 Board = 0x0100000000000000
	DiagonalA2B1 Board = 0x0201000000000000
	DiagonalA3C1 Board = 0x0402010000000000
	DiagonalA4D1 Board = 0x0804020100000000
	DiagonalA5E1 Board = 0x1008040201000000
	DiagonalA6F1 Board = 0x2010080402010000
	DiagonalA7G1 Board = 0x4020100804020100

	DiagonalA8H1 Board = 0x8040201008040201

	DiagonalB8H2 Board = 0x0080402010080402
	DiagonalC8H3 Board = 0x0000804020100804
	DiagonalD8H4 Board = 0x0000008040201008
	DiagonalE8H5 Board = 0x0000000080402010
	DiagonalF8H6 Board = 0x0000000000804020
	DiagonalG8H7 Board = 0x0000000000008040
	DiagonalH8H8 Board = 0x0000000000000080
)

const (
	F1G1   Board = 0x6000000000000000
	F8G8   Board = 0x0000000000000060
	C1D1   Board = 0x0c00000000000000
	C8D8   Board = 0x000000000000000c
	B1C1D1 Board = 0x0e00000000000000
	B8C8D8 Board = 0x000000000000000e
)

var AntiDiagonals = [...]Board{
	square.DiagonalA1A1: DiagonalA1A1,
	square.DiagonalA2B1: DiagonalA2B1,
	square.DiagonalA3C1: DiagonalA3C1,
	square.DiagonalA4D1: DiagonalA4D1,
	square.DiagonalA5E1: DiagonalA5E1,
	square.DiagonalA6F1: DiagonalA6F1,
	square.DiagonalA7G1: DiagonalA7G1,

	square.DiagonalA8H1: DiagonalA8H1,

	square.DiagonalB8H2: DiagonalB8H2,
	square.DiagonalC8H3: DiagonalC8H3,
	square.DiagonalD8H4: DiagonalD8H4,
	square.DiagonalE8H5: DiagonalE8H5,
	square.DiagonalF8H6: DiagonalF8H6,
	square.DiagonalG8H7: DiagonalG8H7,
	square.DiagonalH8H8: DiagonalH8H8,
}

var Squares [square.N]Board

func init() {
	mask := Board(1)
	for s := square.A8; s <= square.H1; s++ {
		Squares[s] = mask
		mask <<= 1
	}
}
