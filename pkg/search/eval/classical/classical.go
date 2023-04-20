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

package classical

import (
	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/board/bitboard"
	"laptudirm.com/x/mess/pkg/board/move/attacks"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
	"laptudirm.com/x/mess/pkg/search/eval"
)

//go:generate go run laptudirm.com/x/mess/internal/generator/classical

// EfficientlyUpdatable (back-acronym of Efficiently Updatable PeSTO) is an efficiently
// updatable PeSTO evaluation function.
type EfficientlyUpdatable struct {
	// the board to evaluate
	Board *board.Board

	// the game phase to lerp between middle and end game
	phase eval.Eval

	// occupancy bitboards
	occupied      bitboard.Board
	occupiedMinus [piece.ColorN][piece.TypeN]bitboard.Board

	// king attackers information
	kingAreas          [piece.ColorN]bitboard.Board // area near the king
	kingAttacksCount   [piece.ColorN]int            // attacks in the king area
	kingAttackersCount [piece.ColorN]int            // attackers to the king area

	// various pawn bitboards
	pawnAttacks    [piece.ColorN]bitboard.Board // squares attacked by pawns
	pawnAttacksBy2 [piece.ColorN]bitboard.Board // squares attacked by 2 pawns
	blockedPawns   [piece.ColorN]bitboard.Board // pawns blocked by other pieces

	// areas in which the mobility of the pieces matter
	mobilityAreas [piece.ColorN]bitboard.Board

	// various attack bitboards
	attacked    [piece.ColorN]bitboard.Board // squares attacked
	attackedBy2 [piece.ColorN]bitboard.Board // squares attacked twice
	attackedBy  [piece.ColorN][piece.TypeN]bitboard.Board
}

// compile time check that OTSePUE implements eval.EfficientlyUpdatable
var _ eval.EfficientlyUpdatable = (*EfficientlyUpdatable)(nil)

// FillSquare adds the given piece to the given square of a chessboard.
func (classical *EfficientlyUpdatable) FillSquare(s square.Square, p piece.Piece) {
}

// ClearSquare removes the given piece from the given square.
func (classical *EfficientlyUpdatable) ClearSquare(s square.Square, p piece.Piece) {
}

// Accumulate accumulates the efficiently updated variables into the
// evaluation of the position from the perspective of the given side.
func (classical *EfficientlyUpdatable) Accumulate(stm piece.Color) eval.Eval {
	xtm := stm.Other()

	// initialize various tables
	classical.initialize()

	// piece evaluation terms
	score := classical.evaluatePawns(stm) - classical.evaluatePawns(xtm)   // pawns and structure
	score += classical.evaluatePieces(stm) - classical.evaluatePieces(xtm) // major and minor pieces
	score += classical.evaluateKing(stm) - classical.evaluateKing(xtm)     // king and king-safety

	// other evaluation terms
	score += classical.evaluateThreats(stm) - classical.evaluateThreats(xtm) // threats

	// linearly interpolate between the end game and middle game
	// evaluations using phase/startposPhase as the contribution
	// of the middle game to the final evaluation
	phase := util.Min(classical.phase, startposPhase)
	return util.Lerp(score.EG(), score.MG(), phase, startposPhase)
}

// evaluatePawns returns the static evaluation of our pawns.
func (classical *EfficientlyUpdatable) evaluatePawns(us piece.Color) Score {
	pawnPiece := piece.New(piece.Pawn, us)   // piece representing one of our pawns
	tempPawns := classical.Board.PawnsBB(us) // bitboard to temporarily store our pawns

	score := Score(0)

	// update attack bitboards with pawn attacks
	classical.attackedBy2[us] |= classical.pawnAttacks[us] & classical.attacked[us]
	classical.attacked[us] |= classical.pawnAttacks[us]
	classical.attackedBy[us][piece.Pawn] = classical.pawnAttacks[us]

	// penalty for having stacked pawns
	for file := square.FileA; file <= square.FileH; file++ {
		score += stackedPawnPenalty[(tempPawns & bitboard.Files[file]).Count()]
	}

	// evaluate every pawn
	for tempPawns != bitboard.Empty {
		// get next pawn
		pawn := tempPawns.Pop()

		// add psqt evaluation
		score += table[pawnPiece][pawn]
		classical.phase += phaseInc[piece.Pawn]
	}

	return score
}

// terms for the mobility of pieces
var Mobility = [piece.TypeN][]Score{
	piece.Knight: {
		S(-104, -139), S(-45, -114), S(-22, -37), S(-8, 3),
		S(6, 15), S(11, 34), S(19, 38), S(30, 37),
		S(43, 17),
	},
	piece.Bishop: {
		S(-99, -186), S(-46, -124), S(-16, -54), S(-4, -14),
		S(6, 1), S(14, 20), S(17, 35), S(19, 39),
		S(19, 49), S(27, 48), S(26, 48), S(52, 32),
		S(55, 47), S(83, 2),
	},
	piece.Rook: {
		S(-127, -148), S(-56, -127), S(-25, -85), S(-12, -28),
		S(-10, 2), S(-12, 27), S(-11, 42), S(-4, 46),
		S(4, 52), S(9, 55), S(11, 64), S(19, 68),
		S(19, 73), S(37, 60), S(97, 15),
	},
	piece.Queen: {
		S(-111, -273), S(-253, -401), S(-127, -228), S(-46, -236),
		S(-20, -173), S(-9, -86), S(-1, -35), S(2, -1),
		S(8, 8), S(10, 31), S(15, 37), S(17, 55),
		S(20, 46), S(23, 57), S(22, 58), S(21, 64),
		S(24, 62), S(16, 65), S(13, 63), S(18, 48),
		S(25, 30), S(38, 8), S(34, -12), S(28, -29),
		S(10, -44), S(7, -79), S(-42, -30), S(-23, -50),
	},
}

// Bonuses for a rook being on an open file.
var (
	RookSemiOpenFile = S(10, 9) // no friendly pawns in the file
	RookFullOpenFile = S(34, 8) // no pawns in the file
)

// evaluatePieces evaluates our major and minor pieces.
func (classical *EfficientlyUpdatable) evaluatePieces(us piece.Color) Score {
	them := us.Other() // color of the opponent

	// bitboard containing all pieces except the king and pawns
	pieces := classical.Board.ColorBBs[us] &^
		classical.Board.PieceBBs[piece.Pawn] &^
		classical.Board.PieceBBs[piece.King]

	score := Score(0)

	// evaluate every piece
	for pieces != bitboard.Empty {
		// get next piece
		sq := pieces.Pop()
		pc := classical.Board.Position[sq]
		pt := pc.Type()

		// add psqt evaluation
		score += table[pc][sq]
		classical.phase += phaseInc[pt]

		// specialized evaluation terms for various pieces
		switch pt {
		case piece.Rook:
			// open file bonus for rooks
			file := bitboard.Files[sq.File()]

			switch bitboard.Empty {
			// no pawns on rook file: open file
			case classical.Board.PieceBBs[piece.Pawn] & file:
				score += RookFullOpenFile
			// no friendly pawns on rook file: semi-open file
			case classical.Board.PawnsBB(us) & file:
				score += RookSemiOpenFile
			}
		}

		// calculate the attacks of the current piece
		attacks := attacks.Of(pc, sq, classical.occupiedMinus[us][pt])

		// update attack bitboards with piece attacks
		classical.attackedBy2[us] |= attacks & classical.attacked[us]
		classical.attacked[us] |= attacks
		classical.attackedBy[us][pt] |= attacks

		// add mobility evaluation
		count := (attacks & classical.mobilityAreas[us]).Count()
		score += Mobility[pt][count]

		// update data for king attackers
		kingAttacks := attacks & classical.kingAreas[them] & ^classical.pawnAttacksBy2[them]
		if kingAttacks != bitboard.Empty {
			classical.kingAttacksCount[them] += kingAttacks.Count()
			classical.kingAttackersCount[them]++
		}
	}

	return score
}

// Bonuses/Penalties for having many/few pieces defending the king.
var KingDefenders = [12]Score{
	S(-37, -3), S(-17, 2), S(0, 6), S(11, 8),
	S(21, 8), S(32, 0), S(38, -14), S(10, -5),
	S(12, 6), S(12, 6), S(12, 6), S(12, 6),
}

// king-safety terms
var (
	// safety term for attacks in the king area
	SafetyAttackValue = S(-45, -34)

	// safety term for weak squares in the king area
	SafetyWeakSquares = S(-42, -41)

	// safety term for the absence of enemy queens
	SafetyNoEnemyQueens = S(237, 259)

	// safety terms for safe checks from enemies
	SafetySafeQueenCheck  = S(-93, -83)
	SafetySafeRookCheck   = S(-90, -98)
	SafetySafeBishopCheck = S(-59, -59)
	SafetySafeKnightCheck = S(-112, -117)

	// constant term for safety adjustment
	SafetyAdjustment = S(74, 26)
)

// evaluateKing returns evaluates our king and king-safety.
func (classical *EfficientlyUpdatable) evaluateKing(us piece.Color) Score {
	them := us.Other() // color of the opponent

	enemyQueens := classical.Board.QueensBB(them)

	score := Score(0)

	kingPiece := piece.New(piece.King, us)
	king := (classical.Board.KingBB(us)).FirstOne()

	// psqt evaluation of the king
	score += table[kingPiece][king]

	// defenders of king including pawns and minor pieces
	defenders := classical.Board.PawnsBB(us) |
		classical.Board.KnightsBB(us) |
		classical.Board.BishopsBB(us)

	// king defenders evaluation
	defenders &= classical.kingAreas[us]
	score += KingDefenders[defenders.Count()]

	// do safety evaluation if we have two attackers, or one
	// attacker with the potential for an enemy queen to join
	if classical.kingAttackersCount[us] >= 2-enemyQueens.Count() {
		// weak squares are squares which are attacked by the enemy, defended
		// once or less, and only defended by our king or queens
		weak := classical.attacked[them] &
			^classical.attackedBy2[us] &
			(^classical.attacked[us] | classical.attackedBy[us][piece.Queen] | classical.attackedBy[us][piece.King])

		// scale attack counts when the king area has more than the usual nine squares
		scaledAttackCount := 9 * classical.kingAttacksCount[us] / classical.kingAreas[us].Count()

		// safe squares are squares safe for our enemy, determined by the squares
		// which are not defended, or are weak and attacked twice
		safe := ^classical.Board.ColorBBs[them] &
			(^classical.attacked[us] | (weak & classical.attackedBy2[them]))

		// possible square and piece combinations that would check our king
		knightThreats := attacks.Knight[king]
		bishopThreats := attacks.Bishop(king, classical.occupied)
		rookThreats := attacks.Rook(king, classical.occupied)
		queenThreats := bishopThreats | rookThreats

		// safe check threats from enemy pieces
		knightChecks := knightThreats & safe & classical.attackedBy[them][piece.Knight]
		bishopChecks := bishopThreats & safe & classical.attackedBy[them][piece.Bishop]
		rookChecks := rookThreats & safe & classical.attackedBy[them][piece.Rook]
		queenChecks := queenThreats & safe & classical.attackedBy[them][piece.Queen]

		// calculate safety score
		safety := Score(0)

		// safety penalty for attacks in the king area
		safety += SafetyAttackValue * Score(scaledAttackCount)

		// safety penalty for weak squares in the king area
		safety += SafetyWeakSquares * Score((weak & classical.kingAreas[us]).Count())

		// safety penalty for safe checks from enemies
		safety += SafetySafeKnightCheck * Score(knightChecks.Count())
		safety += SafetySafeBishopCheck * Score(bishopChecks.Count())
		safety += SafetySafeRookCheck * Score(rookChecks.Count())
		safety += SafetySafeQueenCheck * Score(queenChecks.Count())

		// safety bonus for no enemy queens
		if enemyQueens == bitboard.Empty {
			safety += SafetyNoEnemyQueens
		}

		// constant safety adjustment
		safety += SafetyAdjustment

		mg, eg := safety.MG(), safety.EG()

		// convert safety to score with non-linear function
		score += S(
			-mg*util.Min(0, mg)/720,
			util.Min(0, eg)/20,
		)
	}

	// calculate the attacks of the king
	attacks := attacks.King[king]

	// update attack bitboards with king attacks
	classical.attackedBy2[us] |= attacks & classical.attacked[us]
	classical.attacked[us] |= attacks
	classical.attackedBy[us][piece.King] |= attacks

	return score
}

// threat terms
var (
	// threat term for weak pawns
	ThreatWeakPawn = S(-11, -38)

	// threat terms for attacked minors
	ThreatMinorAttackedByPawn  = S(-55, -83)
	ThreatMinorAttackedByMinor = S(-25, -45)
	ThreatMinorAttackedByMajor = S(-30, -55)
	ThreatMinorAttackedByKing  = S(-43, -21)

	// threat terms for attacked majors
	ThreatRookAttackedByLesser = S(-48, -28)
	ThreatRookAttackedByKing   = S(-33, -18)
	ThreatQueenAttackedByOne   = S(-50, -7)

	// threat term for overloaded pieces
	ThreatOverloadedPieces = S(-7, -16)

	// threat term for pawn push threats
	ThreatByPawnPush = S(15, 32)
)

// evaluateThreats evaluates various threats against our pieces.
func (classical *EfficientlyUpdatable) evaluateThreats(us piece.Color) Score {
	score := Score(0)

	them := us.Other()
	enemies := classical.Board.ColorBBs[them]

	// friendly piece bitboards
	pawns := classical.Board.PawnsBB(us)
	knights := classical.Board.KnightsBB(us)
	bishops := classical.Board.BishopsBB(us)
	rooks := classical.Board.RooksBB(us)
	queens := classical.Board.QueensBB(us)

	// bitboards for attacks by enemy pieces
	attacksByKing := classical.attackedBy[them][piece.King]
	attacksByPawns := classical.attackedBy[them][piece.Pawn]
	attacksByMinors := classical.attackedBy[them][piece.Knight] |
		classical.attackedBy[them][piece.Bishop]
	attacksByMajors := classical.attackedBy[them][piece.Rook] |
		classical.attackedBy[them][piece.Queen]

	// the 3rd rank relative to our color
	pushRank := bitboard.Ranks[util.Ternary(us == piece.White, square.Rank3, square.Rank6)]

	// safe pawn pushes
	safePush := attacks.PawnPush(pawns, us) &^ classical.occupied                                 // single push
	safePush |= attacks.PawnPush(safePush & ^attacksByPawns & pushRank, us) &^ classical.occupied // double push
	safePush &= ^attacksByPawns & (classical.attacked[us] | ^classical.attacked[them])            // push safety

	// poorly defended squares are squares which are attacked more times than they
	// are defended, and are not defended by any pawns
	poorlyDefended := (classical.attacked[them] & ^classical.attacked[us]) |
		(classical.attackedBy2[them] & ^classical.attackedBy2[us] & ^classical.attackedBy[us][piece.Pawn])

	// penalty for minor pieces which are poorly defended
	weakMinors := (knights | bishops) & poorlyDefended

	// penalty for pawns which can't be traded off and are poorly defended
	poorlySupportedPawns := pawns & ^attacksByPawns & poorlyDefended
	score += Score(poorlySupportedPawns.Count()) * ThreatWeakPawn

	// penalty for minors attacked by pawns
	minorsAttackedByPawns := (knights | bishops) & attacksByPawns
	score += Score(minorsAttackedByPawns.Count()) * ThreatMinorAttackedByPawn

	// penalty for minors attacked by minors
	minorsAttackedByMinors := (knights | bishops) & attacksByMinors
	score += Score(minorsAttackedByMinors.Count()) * ThreatMinorAttackedByMinor

	// penalty for minors attacked by majors
	minorsAttackedByMajors := (knights | bishops) & attacksByMajors
	score += Score(minorsAttackedByMajors.Count()) * ThreatMinorAttackedByMajor

	// penalty for rooks attacked by lesser pieces
	rooksAttackedByLesser := rooks & (attacksByPawns | attacksByMinors)
	score += Score(rooksAttackedByLesser.Count()) * ThreatRookAttackedByLesser

	// penalty for weak minors attacked by the king
	weakMinorsAttackedByKing := weakMinors & attacksByKing
	score += Score(weakMinorsAttackedByKing.Count()) * ThreatMinorAttackedByKing

	// penalty for weak rooks attacked by the king
	weakRooksAttackedByKing := rooks & poorlyDefended & attacksByKing
	score += Score(weakRooksAttackedByKing.Count()) * ThreatRookAttackedByKing

	// penalty for attacked queens
	attackedQueens := queens & classical.attacked[them]
	score += Score(attackedQueens.Count()) * ThreatQueenAttackedByOne

	// overloaded pieces are attacked and defended by exactly one piece
	overloaded := (knights | bishops | rooks | queens) &
		classical.attacked[us] & ^classical.attackedBy2[us] &
		classical.attacked[them] & ^classical.attackedBy2[them]
	score += Score(overloaded.Count()) * ThreatOverloadedPieces

	// bonus for giving threats to non-pawn enemy with safe pawn pushes
	// squares that are already threatened by our pawns is not considered
	pushThreat := attacks.Pawns(safePush, us) & (enemies &^ classical.attackedBy[us][piece.Pawn])
	score += Score(pushThreat.Count()) * ThreatByPawnPush

	return score
}

// initialize empties and initializes various variables related to evaluation.
func (classical *EfficientlyUpdatable) initialize() {
	classical.phase = 0

	black := classical.Board.ColorBBs[piece.Black]
	white := classical.Board.ColorBBs[piece.White]

	classical.occupied = white | black

	blackKing := classical.Board.KingBB(piece.Black)
	whiteKing := classical.Board.KingBB(piece.White)

	classical.kingAreas[piece.Black] = bitboard.KingAreas[piece.Black][blackKing.FirstOne()]
	classical.kingAreas[piece.White] = bitboard.KingAreas[piece.White][whiteKing.FirstOne()]

	classical.kingAttackersCount[piece.Black] = 0
	classical.kingAttackersCount[piece.White] = 0

	classical.kingAttacksCount[piece.Black] = 0
	classical.kingAttacksCount[piece.White] = 0

	blackPawns := classical.Board.PawnsBB(piece.Black)
	whitePawns := classical.Board.PawnsBB(piece.White)

	blackPawnsAdvanced := blackPawns.South()
	whitePawnsAdvanced := whitePawns.North()

	classical.pawnAttacks[piece.Black] = blackPawnsAdvanced.East() | blackPawnsAdvanced.West()
	classical.pawnAttacks[piece.White] = whitePawnsAdvanced.East() | whitePawnsAdvanced.West()
	classical.pawnAttacksBy2[piece.Black] = blackPawnsAdvanced.East() & blackPawnsAdvanced.West()
	classical.pawnAttacksBy2[piece.White] = whitePawnsAdvanced.East() & whitePawnsAdvanced.West()

	classical.blockedPawns[piece.White] = classical.occupied.South() & whitePawns
	classical.blockedPawns[piece.Black] = classical.occupied.North() & blackPawns

	classical.attackedBy[piece.Black][piece.King] = attacks.King[blackKing.FirstOne()]
	classical.attackedBy[piece.White][piece.King] = attacks.King[whiteKing.FirstOne()]

	classical.attacked[piece.Black] = classical.attackedBy[piece.Black][piece.King]
	classical.attacked[piece.White] = classical.attackedBy[piece.White][piece.King]

	classical.mobilityAreas[piece.Black] = ^(classical.pawnAttacks[piece.White] | blackKing | classical.blockedPawns[piece.Black])
	classical.mobilityAreas[piece.White] = ^(classical.pawnAttacks[piece.Black] | whiteKing | classical.blockedPawns[piece.White])

	classical.occupiedMinus[piece.Black][piece.Bishop] = classical.occupied ^ classical.Board.BishopsBB(piece.Black)
	classical.occupiedMinus[piece.White][piece.Bishop] = classical.occupied ^ classical.Board.BishopsBB(piece.White)

	classical.occupiedMinus[piece.Black][piece.Rook] = classical.occupied ^ classical.Board.RooksBB(piece.Black)
	classical.occupiedMinus[piece.White][piece.Rook] = classical.occupied ^ classical.Board.RooksBB(piece.White)

	classical.occupiedMinus[piece.Black][piece.Queen] = classical.occupied ^ classical.Board.QueensBB(piece.Black)
	classical.occupiedMinus[piece.White][piece.Queen] = classical.occupied ^ classical.Board.QueensBB(piece.White)

	classical.attackedBy[piece.Black][piece.Knight] = bitboard.Empty
	classical.attackedBy[piece.Black][piece.Bishop] = bitboard.Empty
	classical.attackedBy[piece.Black][piece.Rook] = bitboard.Empty
	classical.attackedBy[piece.Black][piece.Queen] = bitboard.Empty

	classical.attackedBy[piece.White][piece.Knight] = bitboard.Empty
	classical.attackedBy[piece.White][piece.Bishop] = bitboard.Empty
	classical.attackedBy[piece.White][piece.Rook] = bitboard.Empty
	classical.attackedBy[piece.White][piece.Queen] = bitboard.Empty
}
