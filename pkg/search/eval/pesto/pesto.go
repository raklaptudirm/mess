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

	attacked    [piece.ColorN]bitboard.Board
	attackedBy2 [piece.ColorN]bitboard.Board
	attackedBy  [piece.ColorN][piece.TypeN]bitboard.Board
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

	// evaluate pieces
	score := pesto.evaluatePawns(stm) - pesto.evaluatePawns(xtm)
	score += pesto.evaluateKnights(stm) - pesto.evaluateKnights(xtm)
	score += pesto.evaluateBishops(stm) - pesto.evaluateBishops(xtm)
	score += pesto.evaluateRooks(stm) - pesto.evaluateRooks(xtm)
	score += pesto.evaluateQueens(stm) - pesto.evaluateQueens(xtm)
	score += pesto.evaluateKing(stm) - pesto.evaluateKing(xtm)

	// other evaluation terms
	score += pesto.evaluateThreats(stm) - pesto.evaluateThreats(xtm)

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

	pesto.attackedBy2[color] = pesto.pawnAttacks[color] & pesto.attacked[color]
	pesto.attacked[color] |= pesto.pawnAttacks[color]
	pesto.attackedBy[color][piece.Pawn] = pesto.pawnAttacks[color]

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

	pesto.attackedBy[color][piece.Knight] = bitboard.Empty

	for tempKnights != bitboard.Empty {
		knight := tempKnights.Pop()
		score += table[knightPiece][knight]
		pesto.phase += phaseInc[piece.Knight]

		attacks := attacks.Knight[knight]
		pesto.attackedBy2[color] = attacks & pesto.attacked[color]
		pesto.attacked[color] |= attacks
		pesto.attackedBy[color][piece.Knight] |= attacks

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

	pesto.attackedBy[color][piece.Bishop] = bitboard.Empty

	for tempBishops != bitboard.Empty {
		bishop := tempBishops.Pop()
		score += table[bishopPiece][bishop]
		pesto.phase += phaseInc[piece.Bishop]

		attacks := attacks.Bishop(bishop, pesto.occupiedMinusBishops[color])
		pesto.attackedBy2[color] = attacks & pesto.attacked[color]
		pesto.attacked[color] |= attacks
		pesto.attackedBy[color][piece.Bishop] |= attacks

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

	pesto.attackedBy[color][piece.Rook] = bitboard.Empty

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
		pesto.attackedBy2[color] = attacks & pesto.attacked[color]
		pesto.attacked[color] |= attacks
		pesto.attackedBy[color][piece.Rook] |= attacks

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

	pesto.attackedBy[color][piece.Queen] = bitboard.Empty

	for tempQueens != bitboard.Empty {
		queen := tempQueens.Pop()
		score += table[queenPiece][queen]
		pesto.phase += phaseInc[piece.Queen]

		attacks := attacks.Queen(queen, pesto.occupied)
		pesto.attackedBy2[color] = attacks & pesto.attacked[color]
		pesto.attacked[color] |= attacks
		pesto.attackedBy[color][piece.Queen] |= attacks

		count := (attacks & pesto.mobilityAreas[color]).Count()
		score += queenMobility[count]
	}

	return score
}

func (pesto *EfficientlyUpdatable) evaluateKing(color piece.Color) Score {
	kingPiece := piece.New(piece.King, color)
	king := (pesto.Board.KingBB(color)).FirstOne()

	attacks := attacks.King[king]
	pesto.attackedBy2[color] = attacks & pesto.attacked[color]
	pesto.attacked[color] |= attacks
	pesto.attackedBy[color][piece.King] |= attacks

	return table[kingPiece][king]
}

var ThreatWeakPawn = S(-11, -38)
var ThreatMinorAttackedByPawn = S(-55, -83)
var ThreatMinorAttackedByMinor = S(-25, -45)
var ThreatMinorAttackedByMajor = S(-30, -55)
var ThreatRookAttackedByLesser = S(-48, -28)
var ThreatMinorAttackedByKing = S(-43, -21)
var ThreatRookAttackedByKing = S(-33, -18)
var ThreatQueenAttackedByOne = S(-50, -7)
var ThreatOverloadedPieces = S(-7, -16)
var ThreatByPawnPush = S(15, 32)

func (pesto *EfficientlyUpdatable) evaluateThreats(us piece.Color) Score {
	score := Score(0)

	them := us.Other()
	enemies := pesto.Board.ColorBBs[them]

	pawns := pesto.Board.PawnsBB(us)
	knights := pesto.Board.KnightsBB(us)
	bishops := pesto.Board.BishopsBB(us)
	rooks := pesto.Board.RooksBB(us)
	queens := pesto.Board.QueensBB(us)

	pushRank := bitboard.Ranks[util.Ternary(us == piece.White, square.Rank3, square.Rank6)]

	attacksByKing := pesto.attackedBy[them][piece.King]
	attacksByPawns := pesto.attackedBy[them][piece.Pawn]
	attacksByMinors := pesto.attackedBy[them][piece.Knight] |
		pesto.attackedBy[them][piece.Bishop]
	attacksByMajors := pesto.attackedBy[them][piece.Rook] |
		pesto.attackedBy[them][piece.Queen]

	poorlyDefended := (pesto.attacked[them] & ^pesto.attacked[us]) |
		(pesto.attackedBy2[them] & ^pesto.attackedBy2[us] & ^pesto.attackedBy[us][piece.Pawn])

	weakMinors := (knights | bishops) & poorlyDefended

	poorlySupportedPawns := pawns & ^attacksByPawns & poorlyDefended
	score += Score(poorlySupportedPawns.Count()) * ThreatWeakPawn

	minorsAttackedByPawns := (knights | bishops) & attacksByPawns
	score += Score(minorsAttackedByPawns.Count()) * ThreatMinorAttackedByPawn

	minorsAttackedByMinors := (knights | bishops) & attacksByMinors
	score += Score(minorsAttackedByMinors.Count()) * ThreatMinorAttackedByMinor

	minorsAttackedByMajors := (knights | bishops) & attacksByMajors
	score += Score(minorsAttackedByMajors.Count()) * ThreatMinorAttackedByMajor

	rooksAttackedByLesser := rooks & (attacksByPawns | attacksByMinors)
	score += Score(rooksAttackedByLesser.Count()) * ThreatRookAttackedByLesser

	weakMinorsAttackedByKing := weakMinors & attacksByKing
	score += Score(weakMinorsAttackedByKing.Count()) * ThreatMinorAttackedByKing

	weakRooksAttackedByKing := rooks & poorlyDefended & attacksByKing
	score += Score(weakRooksAttackedByKing.Count()) * ThreatRookAttackedByKing

	attackedQueens := queens & pesto.attacked[them]
	score += Score(attackedQueens.Count()) * ThreatQueenAttackedByOne

	overloaded := (knights | bishops | rooks | queens) &
		pesto.attacked[us] & ^pesto.attackedBy2[us] &
		pesto.attacked[them] & ^pesto.attackedBy2[them]
	score += Score(overloaded.Count()) * ThreatOverloadedPieces

	pushThreat := attacks.PawnPush(pawns, us) &^ pesto.occupied
	pushThreat |= attacks.PawnPush(pushThreat & ^attacksByPawns & pushRank, us) &^ pesto.occupied
	pushThreat &= ^attacksByPawns & (pesto.attacked[us] | ^pesto.attacked[them])
	pushThreat = attacks.Pawns(pushThreat, us) & (enemies &^ pesto.attackedBy[us][piece.Pawn])
	score += Score(pushThreat.Count()) * ThreatByPawnPush

	return score
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

	pesto.attackedBy[piece.Black][piece.King] = attacks.King[blackKing.FirstOne()]
	pesto.attackedBy[piece.White][piece.King] = attacks.King[whiteKing.FirstOne()]

	pesto.attacked[piece.Black] = pesto.attackedBy[piece.Black][piece.King]
	pesto.attacked[piece.White] = pesto.attackedBy[piece.White][piece.King]

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
