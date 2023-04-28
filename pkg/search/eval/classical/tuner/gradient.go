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
)

type GradientData struct {
	egEval               float64
	wSafetyMG, bSafetyMG float64
	wSafetyEG, bSafetyEG float64
}

func (tuner *Tuner) ComputeGradient() {
	batchEnd := util.Min((tuner.Batch+1)*tuner.Config.BatchSize, len(tuner.Dataset))
	for i := tuner.Batch * tuner.Config.BatchSize; i < batchEnd; i++ {
		tuner.updateSingleGradient(&tuner.Dataset[i])
	}
}

func (tuner *Tuner) updateSingleGradient(entry *Entry) {
	var data GradientData

	E := tuner.LinerEvaluation(entry, &data) // updated static evaluation
	S := Sigmoid(tuner.K, E)                 // wdl prediction from static
	X := (entry.result - S) * S * (1 - S)

	// base values for each phase
	mgBase := X * entry.phaseFactors[MG]
	egBase := X * entry.phaseFactors[EG]

	// update the terms relevant to the position
	for _, coeff := range entry.coeffs {
		// calculate the difference in coefficients
		deltaCoeff := float64(coeff.White - coeff.Black)

		// update the current term
		switch coeff.Type {
		case Normal:
			// update linear term
			tuner.Gradient[coeff.Index][MG] += mgBase * deltaCoeff
			tuner.Gradient[coeff.Index][EG] += egBase * deltaCoeff

		case Safety:
			// update king safety term
			tuner.Gradient[coeff.Index][MG] += (mgBase / 360) *
				(util.Min(data.wSafetyMG, 0)*float64(coeff.White) - util.Min(data.bSafetyMG, 0)*float64(coeff.Black))
			tuner.Gradient[coeff.Index][EG] += (egBase / 20) *
				(util.Min(data.wSafetyEG, 0)*float64(coeff.White) - util.Min(data.bSafetyEG, 0)*float64(coeff.Black))
		}
	}
}
