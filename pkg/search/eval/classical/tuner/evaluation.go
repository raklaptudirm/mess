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

func (tuner Tuner) LinerEvaluation(entry *Entry, gradientData *GradientData) float64 {
	var delta [TermTypeN][piece.ColorN][2]float64

	// calculate the change evaluation due to every tuned parameter
	for _, coeff := range entry.coeffs {
		delta[coeff.Type][piece.White][MG] += float64(coeff.White) * tuner.Delta[coeff.Index][MG]
		delta[coeff.Type][piece.White][EG] += float64(coeff.White) * tuner.Delta[coeff.Index][EG]
		delta[coeff.Type][piece.Black][MG] += float64(coeff.Black) * tuner.Delta[coeff.Index][MG]
		delta[coeff.Type][piece.Black][EG] += float64(coeff.Black) * tuner.Delta[coeff.Index][EG]
	}

	var normal [2]float64

	// calculate the new normal evaluation
	normal[MG] = float64(entry.eval.MG()) + delta[Normal][piece.White][MG] - delta[Normal][piece.Black][MG]
	normal[EG] = float64(entry.eval.EG()) + delta[Normal][piece.White][EG] - delta[Normal][piece.Black][EG]

	var wSafety, bSafety [2]float64

	// calculate the new linear king safety evaluation
	wSafety[MG] = float64(entry.safety[piece.White].MG()) + delta[Safety][piece.White][MG]
	wSafety[EG] = float64(entry.safety[piece.White].EG()) + delta[Safety][piece.White][EG]
	bSafety[MG] = float64(entry.safety[piece.Black].MG()) + delta[Safety][piece.Black][MG]
	bSafety[EG] = float64(entry.safety[piece.Black].EG()) + delta[Safety][piece.Black][EG]

	var safety [2]float64

	// put the linear king safety evaluations through the non-linear functions
	safety[MG] = (-wSafety[MG] * util.Min(0, wSafety[MG]) / 720) -
		(-bSafety[EG] * util.Min(0, bSafety[EG]) / 720)
	safety[EG] = (util.Min(0, wSafety[EG]) / 20) - (util.Min(0, bSafety[EG]) / 20)

	// store info for gradient calculation
	gradientData.egEval = normal[EG] + safety[EG]
	gradientData.wSafetyMG = wSafety[MG]
	gradientData.wSafetyEG = wSafety[EG]
	gradientData.bSafetyMG = bSafety[MG]
	gradientData.bSafetyEG = bSafety[EG]

	// interpolate mg and eg scores to get final static evaluation
	scoreMG := normal[MG] + safety[EG]
	scoreEG := normal[MG] + safety[EG]
	return util.Lerp(scoreEG, scoreMG, float64(entry.phase), float64(classical.MaxPhase))
}
