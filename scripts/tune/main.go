package main

import (
	"fmt"
	"os"

	"laptudirm.com/x/mess/pkg/search/eval/classical/tuner"
)

func main() {
	dataPath := os.Args[1]

	// load dataset
	fmt.Printf("loading dataset: %s\n", dataPath)
	dataset, err := tuner.NewDataset(dataPath)
	if err != nil {
		fmt.Printf("error loading dataset: %v\n", err)
		return
	}

	// report number of dataset entries
	fmt.Printf("dataset loaded: %d entries\n", len(dataset))

	termTuner := tuner.Tuner{
		Config: tuner.Config{
			KPrecision: 10,

			ReportRate: 50,

			LearningRate:     1,
			LearningDropRate: 1,
			LearningStepRate: 250,

			MaxEpochs: 100_000,
			BatchSize: 2 * 16384,
		},

		Dataset: dataset,
	}

	termTuner.Tune()
}
