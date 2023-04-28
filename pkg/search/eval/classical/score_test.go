package classical_test

import (
	"fmt"
	"testing"

	"laptudirm.com/x/mess/pkg/search/eval"
	"laptudirm.com/x/mess/pkg/search/eval/classical"
)

func FuzzRecovery(f *testing.F) {
	f.Add(int32(1000), int32(-1000))
	f.Add(int32(2648), int32(7346))
	f.Add(int32(-3683), int32(-8374))

	f.Fuzz(func(t *testing.T, a, b int32) {
		mg, eg := eval.Eval(a), eval.Eval(b)
		s := classical.S(mg, eg)

		if s.MG() != mg || s.EG() != eg {
			t.Errorf("S(%d, %d) != S(%d, %d)", mg, eg, s.MG(), s.EG())
		}
	})
}

func FuzzAddition(f *testing.F) {
	f.Add(int32(1000), int32(-1000), int32(-1000), int32(1000))
	f.Add(int32(2648), int32(7346), int32(3683), int32(8374))
	f.Add(int32(-2648), int32(-7346), int32(-3683), int32(-8374))

	f.Fuzz(func(t *testing.T, a, b, c, d int32) {
		mg1, eg1, mg2, eg2 := eval.Eval(a), eval.Eval(b), eval.Eval(c), eval.Eval(d)

		s1 := classical.S(mg1, eg1)
		s2 := classical.S(mg2, eg2)

		fmt.Println(a, b, c, d)

		if sum := s1 + s2; sum != classical.S(mg1+mg2, eg1+eg2) {
			t.Errorf("S(%d, %d) + S(%d, %d) -> S(%d, %d)", a, b, c, d, sum.MG(), sum.EG())
		}
	})
}

func FuzzMultiplication(f *testing.F) {
	f.Add(int32(1000), int32(-1000), int32(-1000))
	f.Add(int32(2648), int32(7346), int32(3683))
	f.Add(int32(-2648), int32(-7346), int32(-3683))

	f.Fuzz(func(t *testing.T, a, b, c int32) {
		mg1, eg1, coeff := eval.Eval(a), eval.Eval(b), eval.Eval(c)

		s := classical.S(mg1, eg1)

		actual := classical.S(mg1*coeff, eg1*coeff)

		if product := classical.Score(coeff) * s; product != actual {
			t.Errorf("%d x S(%d, %d) -> S(%d, %d)\nshould be S(%d, %d)", c, a, b, product.MG(), product.EG(), actual.MG(), actual.EG())
		}
	})
}
