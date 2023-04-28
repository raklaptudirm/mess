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
	"bufio"
	"errors"
	"math"
	"os"
	"strings"

	"laptudirm.com/x/mess/pkg/board"
	"laptudirm.com/x/mess/pkg/board/piece"
	"laptudirm.com/x/mess/pkg/formats/fen"
	"laptudirm.com/x/mess/pkg/search/eval"
	"laptudirm.com/x/mess/pkg/search/eval/classical"
)

func NewDataset(filename string) (Dataset, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	reader := bufio.NewReader(file)

	evaluator := &classical.EfficientlyUpdatable{
		ShouldTrace: true,
	}
	chessboard := board.New(board.EU(evaluator))
	evaluator.Board = chessboard

	dataset := make(Dataset, 10)

	for rawEntry, err := reader.ReadString('\n'); err == nil; rawEntry, err = reader.ReadString('\n') {
		var entry Entry

		result, fenString, found := strings.Cut(rawEntry, " ")
		if !found {
			return nil, errors.New("read dataset: invalid entry")
		}

		switch result {
		case "[1.0]":
			entry.result = 1.0
		case "[0.0]":
			entry.result = 0.0
		case "[0.5]":
			entry.result = 0.5
		default:
			return nil, errors.New("read dataset: invalid entry")
		}

		chessboard.UpdateWithFEN(fen.FromString(fenString))

		entry.static = evaluator.Accumulate(piece.White)
		entry.coeffs = GetCoefficients(&evaluator.Trace)
		entry.eval = evaluator.Trace.Evaluation

		phase := classical.MaxPhase - evaluator.Phase
		entry.phaseFactors[0] = 1 - float64(phase)/float64(classical.MaxPhase)
		entry.phaseFactors[1] = 0 + float64(phase)/float64(classical.MaxPhase)
		entry.phase = phase / classical.MaxPhase

		entry.safety[piece.White] = evaluator.Trace.Safety[piece.White]
		entry.safety[piece.Black] = evaluator.Trace.Safety[piece.Black]

		dataset = append(dataset, entry)
	}

	return dataset, nil
}

// Dataset is the training data using which the
// evaluation terms will be tuned.
type Dataset []Entry

// ComputeK computes the optimal value of K, which is the sigmoid scaling
// factor which gives least error with wdl over the tuning dataset.
func (dataset Dataset) ComputeK(precision int) float64 {
	start, end, step := 0.0, 10.0, 1.0
	var current, err float64

	best := dataset.ComputeE(start)

	for i := 0; i <= precision; i++ {
		current = start - step
		for current < end {
			current += step
			err = dataset.ComputeE(current)
			if err <= best {
				best, start = err, current
			}
		}

		end = start + step
		start = start - step
		step = step / 10.0
	}

	return start
}

// ComputeE computes the value of E which is the error of wdl prediction
// by the static evaluation versus the actual match result. The static
// evaluation is converted to a wdl score using a scaled sigmoid function.
func (dataset Dataset) ComputeE(K float64) float64 {
	var total float64

	for _, entry := range dataset {
		// calculate the squared error with the current entry
		total += math.Pow(entry.result-Sigmoid(K, float64(entry.static)), 2)
	}

	// mean squared error
	return total / float64(len(dataset))
}

// Sigmoid implements a sigmoid function scaled by the factor K.
func Sigmoid(K, eval float64) float64 {
	return 1.0 / (1.0 + math.Exp(-K*eval/400.0))
}

// Entry contains the data of a single position which will be
// used to tune the evaluation terms.
type Entry struct {
	coeffs       []Coefficient // coefficients fo linear evaluation terms
	phaseFactors [2]float64    // factors for evaluation interpolation

	// linear safety scores for the position
	safety [piece.ColorN]classical.Score

	// wdl result of the position
	result float64

	eval   classical.Score // non-interpolated evaluation of the position
	static eval.Eval       // interpolated evaluation of the position
	phase  eval.Eval       // phase for evaluation interpolation
}
