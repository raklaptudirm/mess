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
	"os"
	"text/template"

	"laptudirm.com/x/mess/pkg/chess/move/castling"
	"laptudirm.com/x/mess/pkg/chess/piece"
	"laptudirm.com/x/mess/pkg/chess/square"
	"laptudirm.com/x/mess/pkg/util"
	"laptudirm.com/x/mess/pkg/zobrist"
)

type zobristStruct struct {
	PieceSquare [piece.N][square.N]zobrist.Key
	EnPassant   [square.FileN]zobrist.Key
	Castling    [castling.N]zobrist.Key
	SideToMove  zobrist.Key
}

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

	f, err := os.Create("gen_keys.go")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	zobristTemplate.Execute(f, z)
}

var zobristTemplate = template.Must(template.New("zobrist").Parse(`
// Code generated by go generate; DO NOT EDIT OR COMMIT THIS FILE
// The source code for the generator can be found at generator/attack

package zobrist

var PieceSquare = [16][64]Key{
	{{ range .PieceSquare }}
		{ {{ range . }} {{ printf "%#x" . }}, {{ end }} },
	{{ end }}
}
var EnPassant = [8]Key{
	{{ range .EnPassant }}
		{{ printf "%#x" .}},
	{{ end }}
}
var Castling = [16]Key{
	{{ range .Castling }}
		{{ printf "%#x" .}},
	{{ end }}
}
var SideToMove Key = {{ printf "%#x" .SideToMove }}
`))
