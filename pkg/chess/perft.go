// Copyright © 2022 Rak Laptudirm <rak@laptudirm.com>
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

package chess

import "fmt"

func Perft(b *Board, depth int) int {
	if depth == 0 {
		return 1
	}

	var nodes int
	moves := b.GenerateMoves()

	for _, move := range moves {
		b.MakeMove(move)
		newNodes := perft(b, depth-1)
		fmt.Printf("%s: %d\n", move, newNodes)
		nodes += newNodes
		b.UnmakeMove()
	}

	return nodes
}

func perft(b *Board, depth int) int {

	switch depth {
	case 0:
		return 1
	case 1:
		return len(b.GenerateMoves())
	default:
		var nodes int
		moves := b.GenerateMoves()

		for _, move := range moves {
			b.MakeMove(move)
			nodes += perft(b, depth-1)
			b.UnmakeMove()
		}

		return nodes
	}
}
