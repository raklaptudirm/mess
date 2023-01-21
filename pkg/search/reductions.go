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

package search

import (
	"math/bits"
)

// LMR reductions indexed by depth and move number.
var reductions [MaxDepth + 1][128]int

func init() {
	reductions[0][0] = 0

	log := func(n int) int {
		// fast log2 approximation
		return 63 - bits.LeadingZeros64(uint64(n))
	}

	for depth := 1; depth <= MaxDepth; depth++ {
		for moves := 1; moves < 128; moves++ {
			reductions[depth][moves] = 1 + log(depth)*log(moves)/2
		}
	}
}
