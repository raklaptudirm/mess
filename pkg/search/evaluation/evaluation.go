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
	"math"

	"laptudirm.com/x/mess/pkg/chess"
	"laptudirm.com/x/mess/pkg/chess/piece"
	"laptudirm.com/x/mess/pkg/chess/square"
)

// Type Func represents a board evaluation function.
type Func func(*chess.Board) Rel

// constants representing useful evaluations
const (
	Inf  = math.MaxInt32 / 2
	Mate = Inf - 1
	Draw = 0
)

// Of is a simple evaluation.Func which evaluates a position based on the
// material which each side has.
func Of(b *chess.Board) Rel {
	var eval Abs

	for s := square.A8; s <= square.H1; s++ {
		p := b.Position[s]
		if p == piece.NoPiece {
			continue
		}

		eval += material[p] + squareBonuses[p][s].Abs(p.Color())
	}

	return eval.Rel(b.SideToMove)
}
