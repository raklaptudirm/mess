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

import "laptudirm.com/x/mess/pkg/search/eval"

// S creates a new Score encapsulating the given mg and eg evaluations.
func S(mg, eg eval.Eval) Score {
	return Score(uint64(eg)<<32) + Score(mg)
}

// Score encapsulates the PeSTO middle game and end game stores into a
// single value of a single type.
type Score int64

// MG returns the given score's middle game evaluation.
func (score Score) MG() eval.Eval {
	return eval.Eval(int32(uint32(uint64(score))))
}

// EG return the given score's end game evaluation.
func (score Score) EG() eval.Eval {
	return eval.Eval(int32(uint32(uint64(score+(1<<31)) >> 32)))
}
