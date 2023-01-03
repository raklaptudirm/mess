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

package square

// Rank represents a rank on the chessboard.
// Every horizontal line of squares is called a rank.
//
//	8 8 8 8 8 8 8 8
//	7 7 7 7 7 7 7 7
//	6 6 6 6 6 6 6 6
//	5 5 5 5 5 5 5 5
//	4 4 4 4 4 4 4 4
//	3 3 3 3 3 3 3 3
//	2 2 2 2 2 2 2 2
//	1 1 1 1 1 1 1 1
type Rank int8

// constants representing various ranks
const (
	Rank8 Rank = iota
	Rank7
	Rank6
	Rank5
	Rank4
	Rank3
	Rank2
	Rank1
)

// RankN is the number of ranks.
const RankN = 8

// String converts a Rank into it's string representation.
func (r Rank) String() string {
	const rankToStr = "87654321"
	return string(rankToStr[r])
}

// RankFrom creates an instance of Rank from the given id.
func RankFrom(id string) Rank {
	return Rank1 - Rank(id[0]-'1')
}
