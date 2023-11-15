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
	"math"

	"laptudirm.com/x/mess/pkg/search/eval"
	"laptudirm.com/x/mess/pkg/search/eval/classical"
)

type Vector [][2]float64
type ScoreVector []classical.Score

const (
	MG = 0 // index of mg value in Vector
	EG = 1 // index of eg value in Vector
)

func (vector Vector) EvaluationTerms() *classical.EvaluationTerms[classical.Score] {
	terms := classical.Terms
	for i := 0; i < classical.TermsN; i++ {
		term := terms.FetchTerm(i)
		*term += classical.S(
			eval.Eval(math.Round(vector[i][MG])),
			eval.Eval(math.Round(vector[i][EG])),
		)
	}

	return &terms
}

func VectorizeParams(terms classical.EvaluationTerms[classical.Score]) ScoreVector {
	vector := make(ScoreVector, classical.TermsN)

	for i := 0; i < classical.TermsN; i++ {
		vector[i] = *terms.FetchTerm(i)
	}

	return vector
}
