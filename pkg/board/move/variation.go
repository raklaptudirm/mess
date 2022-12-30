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

package move

import "fmt"

// Variation represents a variation or a list of moves that can be played
// one after the other on a position.
type Variation struct {
	moves []Move
}

// Move return's the ith move of the variation. It returns move.Null if
// the ith move doesn't exist.
func (v *Variation) Move(i int) Move {
	if len(v.moves) <= i {
		return Null
	}

	return v.moves[i]
}

// Clear clears the variation.
func (v *Variation) Clear() {
	v.moves = v.moves[:0]
}

// Update updates the variation with the new move and it's child variation.
func (v *Variation) Update(pMove Move, line Variation) {
	v.Clear()
	v.moves = append(v.moves, pMove)
	v.moves = append(v.moves, line.moves...)
}

// String converts the variation into a human readable string.
func (v Variation) String() string {
	str := fmt.Sprintf("%v", v.moves)
	return str[1 : len(str)-1]
}
