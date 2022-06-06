// Copyright © 2022 Rak Laptudirm <raklaptudirm@gmail.com>
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
type Rank int

// Constants representing every rank.
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

// String converts a Rank into it's string representation.
func (r Rank) String() string {
	ranks := [...]string{
		Rank8: "8",
		Rank7: "7",
		Rank6: "6",
		Rank5: "5",
		Rank4: "4",
		Rank3: "3",
		Rank2: "2",
		Rank1: "1",
	}

	return ranks[r]
}

// rankFrom creates an instance of Rank from the given id.
func rankFrom(id string) Rank {
	switch id {
	case "8":
		return Rank8
	case "7":
		return Rank7
	case "6":
		return Rank6
	case "5":
		return Rank5
	case "4":
		return Rank4
	case "3":
		return Rank3
	case "2":
		return Rank2
	case "1":
		return Rank1
	default:
		panic("new rank: invalid rank id")
	}
}
