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
	"fmt"
	"math"
	"os"
	"strconv"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/schollz/progressbar/v3"
	"laptudirm.com/x/mess/pkg/search/eval/classical"
)

type Tuner struct {
	Config Config

	Dataset Dataset
	Delta   Vector

	K float64

	Gradient Vector

	Batch int
}

func (tuner *Tuner) Tune() {
	velocity := make(Vector, classical.TermsN)
	momentum := make(Vector, classical.TermsN)

	tuner.Gradient = make(Vector, classical.TermsN)
	tuner.Delta = make(Vector, classical.TermsN)

	rate := tuner.Config.LearningRate

	batchSize := float64(tuner.Config.BatchSize)

	errorName := make([]string, 0)
	errorData := make([]opts.LineData, 0)

	fmt.Println("tuner: computing optimal value of K")
	tuner.K = tuner.Dataset.ComputeK(tuner.Config.KPrecision)
	scale := (tuner.K * 2) / batchSize
	fmt.Printf("tuner: K = %v\n", tuner.K)

	E := tuner.ComputeE()
	fmt.Printf("tuner: E = %v\n", E)

	// plot the error data
	errorName = append(errorName, strconv.Itoa(0))
	errorData = append(errorData, opts.LineData{Value: E})

	errorPlot := charts.NewLine()
	errorPlot.SetXAxis(errorName).AddSeries("Error", errorData)

	plotFile, _ := os.Create("error-plot.html")
	_ = errorPlot.Render(plotFile)

	batches := len(tuner.Dataset) / tuner.Config.BatchSize

	for epoch := 0; epoch < tuner.Config.MaxEpochs; epoch++ {
		fmt.Printf("tuner: started new epoch (%d/%d)\n", epoch+1, tuner.Config.MaxEpochs)

		// create new progress bar for epoch
		progressBar := progressbar.NewOptions(
			batches,
			progressbar.OptionSetElapsedTime(true),
			progressbar.OptionSetItsString("batch"),
			progressbar.OptionSetPredictTime(true),
			progressbar.OptionSetRenderBlankState(true),
			progressbar.OptionShowCount(),
			progressbar.OptionShowIts(),
		)

		for tuner.Batch = 0; tuner.Batch < batches; tuner.Batch++ {
			tuner.Gradient = make(Vector, classical.TermsN)
			tuner.ComputeGradient()

			for i := 0; i < classical.TermsN; i++ {
				// calculate scaled gradients for term
				mgGradient := tuner.Gradient[i][MG] * scale
				egGradient := tuner.Gradient[i][EG] * scale

				// update momentum gradient
				momentum[i][MG] = momentum[i][MG]*0.9 + mgGradient*0.1
				momentum[i][EG] = momentum[i][EG]*0.9 + egGradient*0.1

				// update velocity gradient
				velocity[i][MG] = velocity[i][MG]*0.999 + mgGradient*mgGradient*0.001
				velocity[i][EG] = velocity[i][EG]*0.999 + egGradient*egGradient*0.001

				// tune parameters
				tuner.Delta[i][MG] += momentum[i][MG] * rate / math.Sqrt(1e-8+velocity[i][MG])
				tuner.Delta[i][EG] += momentum[i][EG] * rate / math.Sqrt(1e-8+velocity[i][EG])
			}

			// update progress bar
			_ = progressBar.Add(1)
		}

		// close the progress bar
		_ = progressBar.Close()

		// calculate mean squared error of tuned terms
		E := tuner.ComputeE()
		fmt.Printf("tuner: E = %v\n", E)

		// plot the error data
		errorName = append(errorName, strconv.Itoa(epoch+1))
		errorData = append(errorData, opts.LineData{Value: E})

		errorPlot := charts.NewLine()
		errorPlot.SetXAxis(errorName).AddSeries("Error", errorData)

		plotFile, _ := os.Create("error-plot.html")
		_ = errorPlot.Render(plotFile)

		if epoch != 0 {
			if epoch%tuner.Config.LearningStepRate == 0 {
				rate /= tuner.Config.LearningDropRate
			}

			if epoch%tuner.Config.ReportRate == 0 {
				fmt.Printf("%#v", tuner.Delta)
			}
		}
	}
}

func (tuner *Tuner) ComputeE() float64 {
	var total float64

	for _, entry := range tuner.Dataset {
		// calculate the squared error with the current entry
		static := tuner.LinerEvaluation(&entry, &GradientData{})
		total += math.Pow(entry.result-Sigmoid(tuner.K, static), 2)
	}

	// mean squared error
	return total / float64(len(tuner.Dataset))
}

type Config struct {
	KPrecision int

	ReportRate int

	LearningRate     float64
	LearningDropRate float64
	LearningStepRate int

	MaxEpochs int
	BatchSize int
}
