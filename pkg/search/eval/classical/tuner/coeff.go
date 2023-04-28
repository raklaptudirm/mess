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

package tuner

import (
	"laptudirm.com/x/mess/internal/util"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/search/eval/classical"
)

// Coefficient represents the coefficient of a single evaluation term
// in some position with can be used to cheaply recalculate the new
// static evaluation when the terms are modified.
type Coefficient struct {
	Index        int16    // index of the evaluation term
	Type         TermType // type of the evaluation term
	White, Black int8     // coefficients for each color
}

// GetCoefficients gets the non-zero coefficients from an evaluation trace
// and returns those as a slice which can be used to recompute the static
// evaluation cheaply while tuning.
func GetCoefficients(trace *classical.EvaluationTrace) []Coefficient {
	// calculate number of non-zero coefficients so as to not
	// use an unnecessary amount of memory while allocating
	coefficientN := 0
	for i := 0; i < classical.TermsN; i++ {
		termTrace := *trace.FetchTerm(i)

		// check if coefficient is non-zero
		if termTrace[piece.White] != termTrace[piece.Black] ||
			(i >= classical.IndexSafetyStart && i <= classical.IndexSafetyEnd && termTrace[piece.White] != 0) {
			coefficientN++
		}
	}

	// refetch the non-zero coefficients and add them to the slice
	coefficients := make([]Coefficient, 0, coefficientN)
	for i := 0; i < classical.TermsN; i++ {
		termTrace := *trace.FetchTerm(i)

		if termTrace[piece.White] != termTrace[piece.Black] ||
			(i >= classical.IndexSafetyStart && i <= classical.IndexSafetyEnd && termTrace[piece.White] != 0) {

			// append to coefficients slice
			coefficients = append(coefficients, Coefficient{
				Index: int16(i),

				Type: util.Ternary(
					i >= classical.IndexSafetyStart &&
						i <= classical.IndexSafetyEnd,
					Safety, Normal,
				),

				White: int8(termTrace[piece.White]),
				Black: int8(termTrace[piece.Black]),
			})
		}
	}

	return coefficients
}

// TermType represents the type of a given linear term, which
// implies how it's manipulated before being added to the
// final evaluation.
type TermType int8

// constants representing various term types
const (
	Normal TermType = iota // linear evaluation
	Safety                 // king-safety evaluation

	TermTypeN // number of term types
)
