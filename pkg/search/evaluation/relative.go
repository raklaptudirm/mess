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

package evaluation

import (
	"laptudirm.com/x/mess/pkg/chess/piece"
)

// Rel represents a relative centipawn evaluation where > 0 is better for
// the side to move, while < 0 is better for the other side.
type Rel int

// constants representing useful relative evaluations
const (
	// limits to differenciate between regular and mate in n evaluations
	WinInMaxPly  Rel = Mate - 2*10000
	LoseInMaxPly Rel = -WinInMaxPly
)

// Abs converts a Rel from the perspective of s to an absolute evaluation.
func (r Rel) Abs(s piece.Color) Abs {
	switch s {
	case piece.White:
		return Abs(r)
	case piece.Black:
		return Abs(-r)
	default:
		panic("bad color")
	}
}
