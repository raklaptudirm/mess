package main

import (
	"fmt"
	"math"

	"laptudirm.com/x/mess/pkg/search/eval"
	"laptudirm.com/x/mess/pkg/search/eval/classical"
)

func main() {
	terms := classical.Terms

	for i := 0; i < classical.TermsN; i++ {
		term := terms.FetchTerm(i)
		original := *term
		delta := classical.S(
			eval.Eval(math.Round(data[i][0])), eval.Eval(math.Round(data[i][1])),
		)
		//fmt.Println(data[i])
		*term = original + delta
		if delta != 0 && *term != original+delta {
			panic("wtf")
		}
	}

	fmt.Printf("%#v\n", terms)
}
