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
	"laptudirm.com/x/mess/pkg/board/move/attacks"
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

	// general information
	occupied bitboard.Board

	pawnAttacks   [piece.ColorN]bitboard.Board
	blockedPawns  [piece.ColorN]bitboard.Board
	mobilityAreas [piece.ColorN]bitboard.Board

	occupiedMinusBishops [piece.ColorN]bitboard.Board
	occupiedMinusRooks   [piece.ColorN]bitboard.Board
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

	pesto.initialize()

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

var knightMobility = [9]Score{
	S(-104, -139), S(-45, -114), S(-22, -37), S(-8, 3),
	S(6, 15), S(11, 34), S(19, 38), S(30, 37),
	S(43, 17),
}

func (pesto *EfficientlyUpdatable) evaluateKnights(color piece.Color) Score {
	knightPiece := piece.New(piece.Knight, color)
	tempKnights := pesto.Board.KnightsBB(color)

	score := Score(0)

	for tempKnights != bitboard.Empty {
		knight := tempKnights.Pop()
		score += table[knightPiece][knight]
		pesto.phase += phaseInc[piece.Knight]

		attacks := attacks.Knight[knight]
		count := (attacks & pesto.mobilityAreas[color]).Count()
		score += knightMobility[count]
	}

	return score
}

var bishopMobility = [14]Score{
	S(-99, -186), S(-46, -124), S(-16, -54), S(-4, -14),
	S(6, 1), S(14, 20), S(17, 35), S(19, 39),
	S(19, 49), S(27, 48), S(26, 48), S(52, 32),
	S(55, 47), S(83, 2),
}

func (pesto *EfficientlyUpdatable) evaluateBishops(color piece.Color) Score {
	bishopPiece := piece.New(piece.Bishop, color)
	tempBishops := pesto.Board.BishopsBB(color)

	score := Score(0)

	for tempBishops != bitboard.Empty {
		bishop := tempBishops.Pop()
		score += table[bishopPiece][bishop]
		pesto.phase += phaseInc[piece.Bishop]

		attacks := attacks.Bishop(bishop, pesto.occupiedMinusBishops[color])
		count := (attacks & pesto.mobilityAreas[color]).Count()
		score += bishopMobility[count]
	}

	return score
}

var rookMobility = [15]Score{
	S(-127, -148), S(-56, -127), S(-25, -85), S(-12, -28),
	S(-10, 2), S(-12, 27), S(-11, 42), S(-4, 46),
	S(4, 52), S(9, 55), S(11, 64), S(19, 68),
	S(19, 73), S(37, 60), S(97, 15),
}

var rookSemiOpenFile = S(10, 9)
var rookFullOpenFile = S(34, 8)

func (pesto *EfficientlyUpdatable) evaluateRooks(color piece.Color) Score {
	them := color.Other()

	rookPiece := piece.New(piece.Rook, color)
	tempRooks := pesto.Board.RooksBB(color)

	score := Score(0)

	for tempRooks != bitboard.Empty {
		rook := tempRooks.Pop()
		score += table[rookPiece][rook]
		pesto.phase += phaseInc[piece.Rook]

		file := bitboard.Files[rook.File()]

		if pesto.Board.PawnsBB(color)&file == bitboard.Empty {
			if pesto.Board.PawnsBB(them)&file == bitboard.Empty {
				score += rookFullOpenFile
			} else {
				score += rookSemiOpenFile
			}
		}

		attacks := attacks.Rook(rook, pesto.occupiedMinusRooks[color])
		count := (attacks & pesto.mobilityAreas[color]).Count()
		score += rookMobility[count]
	}

	return score
}

var queenMobility = [28]Score{
	S(-111, -273), S(-253, -401), S(-127, -228), S(-46, -236),
	S(-20, -173), S(-9, -86), S(-1, -35), S(2, -1),
	S(8, 8), S(10, 31), S(15, 37), S(17, 55),
	S(20, 46), S(23, 57), S(22, 58), S(21, 64),
	S(24, 62), S(16, 65), S(13, 63), S(18, 48),
	S(25, 30), S(38, 8), S(34, -12), S(28, -29),
	S(10, -44), S(7, -79), S(-42, -30), S(-23, -50),
}

func (pesto *EfficientlyUpdatable) evaluateQueens(color piece.Color) Score {
	queenPiece := piece.New(piece.Queen, color)
	tempQueens := pesto.Board.QueensBB(color)

	score := Score(0)

	for tempQueens != bitboard.Empty {
		queen := tempQueens.Pop()
		score += table[queenPiece][queen]
		pesto.phase += phaseInc[piece.Queen]

		attacks := attacks.Queen(queen, pesto.occupied)
		count := (attacks & pesto.mobilityAreas[color]).Count()
		score += queenMobility[count]
	}

	return score
}

func (pesto *EfficientlyUpdatable) evaluateKing(color piece.Color) Score {
	kingPiece := piece.New(piece.King, color)
	king := (pesto.Board.KingBB(color)).FirstOne()

	return table[kingPiece][king]
}

func (pesto *EfficientlyUpdatable) initialize() {
	pesto.phase = 0

	black := pesto.Board.ColorBBs[piece.Black]
	white := pesto.Board.ColorBBs[piece.White]

	pesto.occupied = white | black

	blackKing := pesto.Board.KingBB(piece.Black)
	whiteKing := pesto.Board.KingBB(piece.White)

	blackPawns := pesto.Board.PawnsBB(piece.Black)
	whitePawns := pesto.Board.PawnsBB(piece.White)

	blackPawnsAdvanced := blackPawns.South()
	whitePawnsAdvanced := whitePawns.North()

	pesto.pawnAttacks[piece.Black] = blackPawnsAdvanced.East() | blackPawnsAdvanced.West()
	pesto.pawnAttacks[piece.White] = whitePawnsAdvanced.East() | whitePawnsAdvanced.West()

	pesto.blockedPawns[piece.White] = pesto.occupied.South() & whitePawns
	pesto.blockedPawns[piece.Black] = pesto.occupied.North() & blackPawns

	pesto.mobilityAreas[piece.Black] = ^(pesto.pawnAttacks[piece.White] | blackKing | pesto.blockedPawns[piece.Black])
	pesto.mobilityAreas[piece.White] = ^(pesto.pawnAttacks[piece.Black] | whiteKing | pesto.blockedPawns[piece.White])

	blackBishops := pesto.Board.BishopsBB(piece.Black)
	whiteBishops := pesto.Board.BishopsBB(piece.White)

	pesto.occupiedMinusBishops[piece.White] = pesto.occupied ^ whiteBishops
	pesto.occupiedMinusBishops[piece.Black] = pesto.occupied ^ blackBishops

	blackRooks := pesto.Board.RooksBB(piece.Black)
	whiteRooks := pesto.Board.RooksBB(piece.White)

	pesto.occupiedMinusRooks[piece.White] = pesto.occupied ^ whiteRooks
	pesto.occupiedMinusRooks[piece.Black] = pesto.occupied ^ blackRooks
}
