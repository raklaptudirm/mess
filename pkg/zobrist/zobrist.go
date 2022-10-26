package zobrist

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

import (
	"laptudirm.com/x/mess/pkg/chess/move/castling"
	"laptudirm.com/x/mess/pkg/chess/piece"
	"laptudirm.com/x/mess/pkg/chess/square"
	"laptudirm.com/x/mess/pkg/util"
)

type Key uint64

var PieceSquare [piece.N][square.N]Key
var EnPassant [square.FileN]Key
var Castling [castling.N]Key
var SideToMove Key

func init() {
	var rng util.PRNG
	rng.Seed(1070372) // seed used from Stockfish

	// piece square numbers
	for p := 0; p < piece.N; p++ {
		for s := square.A8; s <= square.H1; s++ {
			PieceSquare[p][s] = Key(rng.Uint64())
		}
	}

	// en passant file numbers
	for f := square.FileA; f <= square.FileH; f++ {
		EnPassant[f] = Key(rng.Uint64())
	}

	// castling right numbers
	for r := castling.NoCasl; r <= castling.All; r++ {
		Castling[r] = Key(rng.Uint64())
	}

	// black to move number
	SideToMove = Key(rng.Uint64())
}
