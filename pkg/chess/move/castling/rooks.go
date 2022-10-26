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

package castling

import (
	"laptudirm.com/x/mess/pkg/chess/piece"
	"laptudirm.com/x/mess/pkg/chess/square"
)

const Sides = 4

var Rooks = [square.N]struct {
	From, To square.Square
	RookType piece.Piece
}{
	square.G1: {
		From:     square.H1,
		To:       square.F1,
		RookType: piece.WhiteRook,
	},
	square.C1: {
		From:     square.A1,
		To:       square.D1,
		RookType: piece.WhiteRook,
	},
	square.G8: {
		From:     square.H8,
		To:       square.F8,
		RookType: piece.BlackRook,
	},
	square.C8: {
		From:     square.A8,
		To:       square.D8,
		RookType: piece.BlackRook,
	},
}
