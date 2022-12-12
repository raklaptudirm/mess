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

// Package bitboard implements a 64-bit bitboard and other related
// functions for manipulating them.
package bitboard

import (
	"math/bits"

	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
)

// Board is a 64-bit bitboard
type Board uint64

// String returns a string representation of the given BB.
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

// Up shifts the given BB up relative to the given color.
func (b Board) Up(c piece.Color) Board {
	switch c {
	case piece.White:
		return b.North()
	case piece.Black:
		return b.South()
	default:
		panic("bad color")
	}
}

// Down shifts the given BB down relative to the given color.
func (b Board) Down(color piece.Color) Board {
	switch color {
	case piece.White:
		return b.South()
	case piece.Black:
		return b.North()
	default:
		panic("bad color")
	}
}

// North shifts the given BB to the north.
func (b Board) North() Board {
	return b >> 8
}

// South shifts the given BB to the south.
func (b Board) South() Board {
	return b << 8
}

// East shifts the given BB to the east.
func (b Board) East() Board {
	return (b &^ FileH) << 1
}

// West shifts the given BB to the west.
func (b Board) West() Board {
	return (b &^ FileA) >> 1
}

// Pop returns the LSB of the given BB and removes it.
func (b *Board) Pop() square.Square {
	sq := b.FirstOne()
	*b &= *b - 1
	return sq
}

// Count returns the number of set bits in the given BB.
func (b Board) Count() int {
	return bits.OnesCount64(uint64(b))
}

// TODO: this is a duplicate of Count, remove it.
func (b Board) CountBits() int {
	return bits.OnesCount64(uint64(b))
}

// FirstOne returns the LSB of the given BB.
func (b Board) FirstOne() square.Square {
	return square.Square(bits.TrailingZeros64(uint64(b)))
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
