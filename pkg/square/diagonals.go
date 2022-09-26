package square

type Diagonal int

const (
	DiagonalH1H1 Diagonal = iota
	DiagonalH2G1
	DiagonalH3F1
	DiagonalH4E1
	DiagonalH5D1
	DiagonalH6C1
	DiagonalH7B1

	DiagonalH8A1

	DiagonalG8A2
	DiagonalF8A3
	DiagonalE8A4
	DiagonalD8A5
	DiagonalC8A6
	DiagonalB8A7
	DiagonalA8A8
)

type AntiDiagonal int

const (
	DiagonalA1A1 AntiDiagonal = iota
	DiagonalA2B1
	DiagonalA3C1
	DiagonalA4D1
	DiagonalA5E1
	DiagonalA6F1
	DiagonalA7G1

	DiagonalA8H1

	DiagonalB8H2
	DiagonalC8H3
	DiagonalD8H4
	DiagonalE8H5
	DiagonalF8H6
	DiagonalG8H7
	DiagonalH8H8
)
