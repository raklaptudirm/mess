// Copyright © 2022 Rak Laptudirm <rak@laptudirm.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package square

// Diagonal represents a diagonal on the chessboard.
// Every NE-SW diagonal is called a diagonal.
//
//     e d c b a 9 8 7
//     d c b a 9 8 7 6
//     c b a 9 8 7 6 5
//     b a 9 8 7 6 5 4
//     a 9 8 7 6 5 4 3
//     9 8 7 6 5 4 3 2
//     8 7 6 5 4 3 2 1
//     7 6 5 4 3 2 1 0
//
type Diagonal int

// constants representing various diagonals
const (
	// bottom diagonals
	DiagonalH1H1 Diagonal = iota
	DiagonalH2G1
	DiagonalH3F1
	DiagonalH4E1
	DiagonalH5D1
	DiagonalH6C1
	DiagonalH7B1

	// main diagonal
	DiagonalH8A1

	// top diagonals
	DiagonalG8A2
	DiagonalF8A3
	DiagonalE8A4
	DiagonalD8A5
	DiagonalC8A6
	DiagonalB8A7
	DiagonalA8A8
)

// AntiDiagonal represents an anti-diagonal on the chessboard.
// Every NW-SE diagonal is called an anti-diagonal.
//
//     7 8 9 a b c d e
//     6 7 8 9 a b c d
//     5 6 7 8 9 a b c
//     4 5 6 7 8 9 a b
//     3 4 5 6 7 8 9 a
//     2 3 4 5 6 7 8 9
//     1 2 3 4 5 6 7 8
//     0 1 2 3 4 5 6 7
//
type AntiDiagonal int

// constants representing various anti-diagonals
const (
	// bottom anti-diagonals
	DiagonalA1A1 AntiDiagonal = iota
	DiagonalA2B1
	DiagonalA3C1
	DiagonalA4D1
	DiagonalA5E1
	DiagonalA6F1
	DiagonalA7G1

	// main anti-diagonal
	DiagonalA8H1

	// top anti-diagonals
	DiagonalB8H2
	DiagonalC8H3
	DiagonalD8H4
	DiagonalE8H5
	DiagonalF8H6
	DiagonalG8H7
	DiagonalH8H8
)
