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

package attacks

import (
	"laptudirm.com/x/mess/pkg/chess/bitboard"
	"laptudirm.com/x/mess/pkg/chess/square"
)

func rook(s square.Square, occ bitboard.Board, isMask bool) bitboard.Board {
	fileMask := bitboard.Files[s.File()]
	fileAttacks := hyperbola(s, occ, fileMask)

	rankMask := bitboard.Ranks[s.Rank()]
	rankAttacks := hyperbola(s, occ, rankMask)

	if isMask {
		fileAttacks &^= bitboard.Rank1 | bitboard.Rank8
		rankAttacks &^= bitboard.FileA | bitboard.FileH
	}

	return fileAttacks | rankAttacks
}

func Rook(s square.Square, blockers bitboard.Board) bitboard.Board {
	return RookTable.Probe(s, blockers)
}
