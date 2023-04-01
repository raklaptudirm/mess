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

package pesto

import (
	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
	"laptudirm.com/x/mess/pkg/search/eval"
)

//go:generate go run laptudirm.com/x/mess/internal/generator/pesto

// EfficientlyUpdatable (back-acronym of Efficiently Updatable PeSTO) is an efficiently
// updatable PeSTO evaluation function.
type EfficientlyUpdatable struct {
	Board *board.Board
	phase eval.Eval // the game phase to lerp between middle and end game
}

// compile time check that OTSePUE implements eval.EfficientlyUpdatable
var _ eval.EfficientlyUpdatable = (*EfficientlyUpdatable)(nil)

// FillSquare adds the given piece to the given square of a chessboard.
func (pesto *EfficientlyUpdatable) FillSquare(s square.Square, p piece.Piece) {
}

// ClearSquare removes the given piece from the given square.
func (pesto *EfficientlyUpdatable) ClearSquare(s square.Square, p piece.Piece) {
}

// Accumulate accumulates the efficiently updated variables into the
// evaluation of the position from the perspective of the given side.
func (pesto *EfficientlyUpdatable) Accumulate(stm piece.Color) eval.Eval {
	xtm := stm.Other()

	pesto.phase = 0

	// score from side to move's perspective
	score := pesto.evaluatePawns(stm) - pesto.evaluatePawns(xtm)
	score += pesto.evaluateKnights(stm) - pesto.evaluateKnights(xtm)
	score += pesto.evaluateBishops(stm) - pesto.evaluateBishops(xtm)
	score += pesto.evaluateRooks(stm) - pesto.evaluateRooks(xtm)
	score += pesto.evaluateQueens(stm) - pesto.evaluateQueens(xtm)
	score += pesto.evaluateKing(stm) - pesto.evaluateKing(xtm)

	// linearly interpolate between the end game and middle game
	// evaluations using phase/startposPhase as the contribution
	// of the middle game to the final evaluation
	phase := util.Min(pesto.phase, startposPhase)
	return util.Lerp(score.EG(), score.MG(), phase, startposPhase)
}

func (pesto *EfficientlyUpdatable) evaluatePawns(color piece.Color) Score {
	pawnPiece := piece.New(piece.Pawn, color)
	tempPawns := pesto.Board.PawnsBB(color)

	score := Score(0)

	for file := square.FileA; file <= square.FileH; file++ {
		score += stackedPawnPenalty[(tempPawns & bitboard.Files[file]).Count()]
	}

	for tempPawns != bitboard.Empty {
		pawn := tempPawns.Pop()
		score += table[pawnPiece][pawn]
		pesto.phase += phaseInc[piece.Pawn]
	}

	return score
}

func (pesto *EfficientlyUpdatable) evaluateKnights(color piece.Color) Score {
	knightPiece := piece.New(piece.Knight, color)
	tempKnights := pesto.Board.KnightsBB(color)

	score := Score(0)

	for tempKnights != bitboard.Empty {
		knight := tempKnights.Pop()
		score += table[knightPiece][knight]
		pesto.phase += phaseInc[piece.Knight]
	}

	return score
}

func (pesto *EfficientlyUpdatable) evaluateBishops(color piece.Color) Score {
	bishopPiece := piece.New(piece.Bishop, color)
	tempBishops := pesto.Board.BishopsBB(color)

	score := Score(0)

	for tempBishops != bitboard.Empty {
		bishop := tempBishops.Pop()
		score += table[bishopPiece][bishop]
		pesto.phase += phaseInc[piece.Bishop]
	}

	return score
}

func (pesto *EfficientlyUpdatable) evaluateRooks(color piece.Color) Score {
	rookPiece := piece.New(piece.Rook, color)
	tempRooks := pesto.Board.RooksBB(color)

	score := Score(0)

	for tempRooks != bitboard.Empty {
		rook := tempRooks.Pop()
		score += table[rookPiece][rook]
		pesto.phase += phaseInc[piece.Rook]
	}

	return score
}

func (pesto *EfficientlyUpdatable) evaluateQueens(color piece.Color) Score {
	queenPiece := piece.New(piece.Queen, color)
	tempQueens := pesto.Board.QueensBB(color)

	score := Score(0)

	for tempQueens != bitboard.Empty {
		queen := tempQueens.Pop()
		score += table[queenPiece][queen]
		pesto.phase += phaseInc[piece.Queen]
	}

	return score
}

func (pesto *EfficientlyUpdatable) evaluateKing(color piece.Color) Score {
	kingPiece := piece.New(piece.King, color)
	king := (pesto.Board.KingBB(color)).FirstOne()

	return table[kingPiece][king]
}
