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

import (
	"laptudirm.com/x/mess/pkg/square"
)

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

// useful bitboard definitions
var (
	Empty    Board = 0
	Universe Board = 0xffffffffffffffff
)

// IsSet checks whether the given Square is set in the bitboard.
func (b Board) IsSet(index square.Square) bool {
	return (b>>index)&1 == 1
}

// Set sets the given Square in the bitboard.
func (b *Board) Set(index square.Square) {
	*b |= buffer(index)
}

// Unset clears the given Square in the bitboard.
func (b *Board) Unset(index square.Square) {
	*b &^= buffer(index)
}

// buffer creates a utility bitboard in which only the given square is set.
func buffer(index square.Square) Board {
	return 1 << Board(index)
}
