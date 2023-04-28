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
	"fmt"

	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/board/square"
)

type EvaluationTerms[T any] struct {
	// general terms

	// mobility terms for number of legal moves
	Mobility [piece.TypeN][]T

	// piece-square table
	PieceSquare [piece.N][square.N]T

	// pawn specific
	StackedPawns [7]T

	// rook specific
	RookFullOpenFile T // no pawns in the file
	RookSemiOpenFile T // no friendly pawns in the file

	// king safety terms

	// term for number of minors/pawns defending the king
	KingDefenders [12]T

	// safety term for attacks in the king area
	SafetyAttackValue T

	// safety term for weak squares in the king area
	SafetyWeakSquares T

	// safety term for the absence of enemy queens
	SafetyNoEnemyQueens T

	// safety terms for safe checks from enemies
	SafetySafeQueenCheck  T
	SafetySafeRookCheck   T
	SafetySafeBishopCheck T
	SafetySafeKnightCheck T

	// constant term for safety adjustment
	SafetyAdjustment T

	// threat terms

	// threat term for weak pawns
	ThreatWeakPawn T

	// threat terms for attacked minors
	ThreatMinorAttackedByPawn  T
	ThreatMinorAttackedByMinor T
	ThreatMinorAttackedByMajor T
	ThreatMinorAttackedByKing  T

	// threat terms for attacked majors
	ThreatRookAttackedByLesser T
	ThreatRookAttackedByKing   T
	ThreatQueenAttackedByOne   T

	// threat term for overloaded pieces
	ThreatOverloadedPieces T

	// threat term for pawn push threats
	ThreatByPawnPush T
}

// FetchTerm returns a pointer to the evaluation term with the given
// index from the provided EvaluationTerms structure.
func (terms *EvaluationTerms[T]) FetchTerm(index int) *T {
	// singular terms
	switch index {
	case IndexRookFullOpenFile:
		return &terms.RookFullOpenFile
	case IndexRookSemiOpenFile:
		return &terms.RookSemiOpenFile

	case IndexSafetyAttackValue:
		return &terms.SafetyAttackValue
	case IndexSafetyWeakSquares:
		return &terms.SafetyWeakSquares
	case IndexSafetyNoEnemyQueens:
		return &terms.SafetyNoEnemyQueens
	case IndexSafetySafeQueenCheck:
		return &terms.SafetySafeQueenCheck
	case IndexSafetySafeRookCheck:
		return &terms.SafetySafeRookCheck
	case IndexSafetySafeBishopCheck:
		return &terms.SafetySafeBishopCheck
	case IndexSafetySafeKnightCheck:
		return &terms.SafetySafeKnightCheck
	case IndexSafetyAdjustment:
		return &terms.SafetyAdjustment

	case IndexThreatWeakPawn:
		return &terms.ThreatWeakPawn
	case IndexThreatMinorAttackedByPawn:
		return &terms.ThreatMinorAttackedByPawn
	case IndexThreatMinorAttackedByMinor:
		return &terms.ThreatMinorAttackedByMinor
	case IndexThreatMinorAttackedByMajor:
		return &terms.ThreatMinorAttackedByMajor
	case IndexThreatMinorAttackedByKing:
		return &terms.ThreatMinorAttackedByKing
	case IndexThreatRookAttackedByLesser:
		return &terms.ThreatRookAttackedByLesser
	case IndexThreatRookAttackedByKing:
		return &terms.ThreatRookAttackedByKing
	case IndexThreatQueenAttackedByOne:
		return &terms.ThreatQueenAttackedByOne
	case IndexThreatOverloadedPieces:
		return &terms.ThreatOverloadedPieces
	case IndexThreatByPawnPush:
		return &terms.ThreatByPawnPush
	}

	// term tables
	switch {
	case index >= IndexMobility &&
		index < IndexMobility+MobilityN:
		mobilityIndex := index - IndexMobility

		switch {
		case mobilityIndex < 9:
			return &terms.Mobility[piece.Knight][mobilityIndex]
		case mobilityIndex < 9+14:
			return &terms.Mobility[piece.Bishop][mobilityIndex-9]
		case mobilityIndex < 9+14+15:
			return &terms.Mobility[piece.Rook][mobilityIndex-9-14]
		default:
			return &terms.Mobility[piece.Queen][mobilityIndex-9-14-15]
		}

	case index >= IndexPSQT && index < IndexPSQT+PSQTN:
		psqtIndex := index - IndexPSQT

		sq := psqtIndex % 64
		psqtIndex = (psqtIndex - sq) / 64

		pc := psqtIndex % 2
		psqtIndex = (psqtIndex - pc) / 2

		pt := psqtIndex

		p := piece.New(piece.Type(pt+1), piece.Color(pc))
		return &terms.PieceSquare[p][sq]

	case index >= IndexStackedPawns &&
		index < IndexStackedPawns+StackedPawnsN:
		return &terms.StackedPawns[index-IndexStackedPawns]

	case index >= IndexKingDefenders &&
		index < IndexKingDefenders+KingDefendersN:
		return &terms.KingDefenders[index-IndexKingDefenders]
	}

	panic(fmt.Errorf("fetch term: invalid index %d", index))
}

// constants for an universal indexing standard for all the
// classical evaluation terms, useful for tuning
const (
	IndexMobility = 0
	MobilityN     = 9 + 14 + 15 + 28

	// TODO: remove useless terms from psqt
	IndexPSQT = IndexMobility + MobilityN
	PSQTN     = piece.ColorN * (piece.TypeN - 1) * square.N

	IndexStackedPawns = IndexPSQT + PSQTN
	StackedPawnsN     = 7

	IndexRookFullOpenFile = IndexStackedPawns + StackedPawnsN
	IndexRookSemiOpenFile = IndexRookFullOpenFile + 1

	IndexKingDefenders = IndexRookSemiOpenFile + 1
	KingDefendersN     = 12

	IndexSafetyStart = IndexKingDefenders + KingDefendersN

	IndexSafetyAttackValue   = IndexSafetyStart
	IndexSafetyWeakSquares   = IndexSafetyAttackValue + 1
	IndexSafetyNoEnemyQueens = IndexSafetyWeakSquares + 1

	IndexSafetySafeQueenCheck  = IndexSafetyNoEnemyQueens + 1
	IndexSafetySafeRookCheck   = IndexSafetySafeQueenCheck + 1
	IndexSafetySafeBishopCheck = IndexSafetySafeRookCheck + 1
	IndexSafetySafeKnightCheck = IndexSafetySafeBishopCheck + 1

	IndexSafetyAdjustment = IndexSafetySafeKnightCheck + 1

	IndexSafetyEnd = IndexSafetyAdjustment

	IndexThreatWeakPawn = IndexSafetyAdjustment + 1

	IndexThreatMinorAttackedByPawn  = IndexThreatWeakPawn + 1
	IndexThreatMinorAttackedByMinor = IndexThreatMinorAttackedByPawn + 1
	IndexThreatMinorAttackedByMajor = IndexThreatMinorAttackedByMinor + 1
	IndexThreatMinorAttackedByKing  = IndexThreatMinorAttackedByMajor + 1

	IndexThreatRookAttackedByLesser = IndexThreatMinorAttackedByKing + 1
	IndexThreatRookAttackedByKing   = IndexThreatRookAttackedByLesser + 1
	IndexThreatQueenAttackedByOne   = IndexThreatRookAttackedByKing + 1

	IndexThreatOverloadedPieces = IndexThreatQueenAttackedByOne + 1

	IndexThreatByPawnPush = IndexThreatOverloadedPieces + 1

	TermsN = IndexThreatByPawnPush + 1
)

type EvaluationTrace struct {
	Evaluation Score               // non-interpolated evaluation
	Safety     [piece.ColorN]Score // king safety score

	// tracing data for each evaluation term
	EvaluationTerms[[piece.ColorN]int]
}
