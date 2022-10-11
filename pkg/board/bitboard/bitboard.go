// Copyright © 2022 Rak Laptudirm <raklaptudirm@gmail.com>
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

// Package bitboard implements a 64-bit bitboard and other related
// functions for manipulating them.
package bitboard

import "laptudirm.com/x/mess/pkg/square"

// Board is a 64-bit bitboard
type Board uint64

func (b Board) String() string {
	var str string
	for s := square.A8; s <= square.H1; s++ {
		if b.IsSet(s) {
			str += "1"
		} else {
			str += "0"
		}

		if s.File() == square.FileH {
			str += "\n"
		} else {
			str += " "
		}
	}

	return str
}

func (b Board) Up() Board {
	return b << 8
}

func (b Board) Down() Board {
	return b >> 8
}

func (b Board) Right() Board {
	return (b &^ FileH) << 1
}

func (b Board) Left() Board {
	return (b &^ FileA) >> 1
}

func (b *Board) Pop() square.Square {
	sq := b.FirstOne()
	*b &= *b - 1
	return sq
}

var index = [square.N]square.Square{
	0, 47, 1, 56, 48, 27, 2, 60,
	57, 49, 41, 37, 28, 16, 3, 61,
	54, 58, 35, 52, 50, 42, 21, 44,
	38, 32, 29, 23, 17, 11, 4, 62,
	46, 55, 26, 59, 40, 36, 15, 53,
	34, 51, 20, 43, 31, 22, 10, 45,
	25, 39, 14, 33, 19, 30, 9, 24,
	13, 18, 8, 12, 7, 6, 5, 63,
}

func (b Board) FirstOne() square.Square {
	const deBruijn = 0x03f79d71b4cb0a89
	return index[((b^(b-1))*deBruijn)>>58]
}

// IsSet checks whether the given Square is set in the bitboard.
func (b Board) IsSet(index square.Square) bool {
	return b&Squares[index] != 0
}

// Set sets the given Square in the bitboard.
func (b *Board) Set(index square.Square) {
	if index == square.None {
		return
	}

	*b |= Squares[index]
}

// Unset clears the given Square in the bitboard.
func (b *Board) Unset(index square.Square) {
	if index == square.None {
		return
	}

	*b &^= Squares[index]
}
