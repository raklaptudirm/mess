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

package main

import (
	_ "embed"

	"laptudirm.com/x/mess/internal/generator"
	"laptudirm.com/x/mess/pkg/board/move/castling"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/zobrist"
)

type zobristStruct struct {
	PieceSquare [piece.N][square.N]zobrist.Key
	EnPassant   [square.FileN]zobrist.Key
	Castling    [castling.N]zobrist.Key
	SideToMove  zobrist.Key
}

//go:embed .gotemplate
var template string

func main() {
	var z zobristStruct

	var rng util.PRNG
	rng.Seed(1070372) // seed used from Stockfish

	// piece square numbers
	for p := 0; p < piece.N; p++ {
		for s := square.A8; s <= square.H1; s++ {
			z.PieceSquare[p][s] = zobrist.Key(rng.Uint64())
		}
	}

	// en passant file numbers
	for f := square.FileA; f <= square.FileH; f++ {
		z.EnPassant[f] = zobrist.Key(rng.Uint64())
	}

	// castling right numbers
	for r := castling.NoCasl; r <= castling.All; r++ {
		z.Castling[r] = zobrist.Key(rng.Uint64())
	}

	// black to move
	z.SideToMove = zobrist.Key(rng.Uint64())

	generator.Generate("keys", template, z)
}
