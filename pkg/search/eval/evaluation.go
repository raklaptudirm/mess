// Copyright © 2023 Rak Laptudirm <rak@laptudirm.com>
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

// Package eval contains various types and functions related to evaluating
// a chess position. It is used by search to determine how good a position
// is and whether the moves leading to it should be played or not.
package eval

import (
	"fmt"
	"math"

	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/board/piece"
)

// EfficientlyUpdatable is an extension of board.EfficientlyUpdatable
// which represents an efficiently updatable evaluation function.
type EfficientlyUpdatable interface {
	// eval.EfficientlyUpdatable implements board.EfficientlyUpdatable
	// so that it can be efficiently updatable by a board.Board
	board.EfficientlyUpdatable

	// Accumulate the efficiently updatable variables and return the
	// evaluation of the position from the perspective of the given
	// color.
	Accumulate(piece.Color) Eval
}

// MatedIn returns the evaluation for being mated in the given plys.
func MatedIn(plys int) Eval {
	// prefer the longer lines when getting mated
	// so longer lines have higher scores (+plys)
	return -Mate + Eval(plys)
}

// RandDraw returns a random draw score based on the provided seed.
func RandDraw(seed int) Eval {
	return Eval(8 - (seed & 7))
}

// Eval represents a relative centipawn evaluation where > 0 is better for
// the side to move, while < 0 is better for the other side.
type Eval int

// constants representing useful relative evaluations
const (
	// basic evaluations
	Inf  Eval = math.MaxInt32 / 2 // prevent any overflows
	Mate Eval = Inf - 1           // Inf is king capture
	Draw Eval = 0

	// limits to differentiate between regular and mate in n evaluations
	WinInMaxPly  Eval = Mate - 2*10000
	LoseInMaxPly Eval = -WinInMaxPly
)

// String returns an UCI compliant string representation of the Eval.
func (r Eval) String() string {
	switch {
	// mate x
	case r > WinInMaxPly:
		plys := Mate - r
		mateInN := (plys / 2) + (plys % 2)
		return fmt.Sprintf("mate %d", mateInN)

	// mate -x
	case r < LoseInMaxPly:
		plys := -Mate - r
		mateInN := (plys / 2) + (plys % 2)
		return fmt.Sprintf("mate %d", mateInN)

	// cp x
	default:
		return fmt.Sprintf("cp %d", r)
	}
}
