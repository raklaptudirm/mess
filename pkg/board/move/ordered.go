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

package move

// the following core types may represent move evaluations
// uint64 is excluded to prevent overflows during storage
type eval interface {
	~int | ~int8 | ~int16 | ~int32 |
		~uint | ~uint8 | ~uint16 | ~uint32
}

// ScoreMoves scores each move in the provided move-list according to the
// provided scorer function and returns an OrderedMoveList containing them.
func ScoreMoves[T eval](moveList []Move, scorer func(Move) T) OrderedMoveList[T] {
	ordered := make([]OrderedMove[T], len(moveList))

	for i, move := range moveList {
		ordered[i] = NewOrdered(move, scorer(move))
	}

	return OrderedMoveList[T]{
		moves:  ordered,
		Length: len(moveList),
	}
}

// OrderedMoveList represents an ordered/ranked move list.
type OrderedMoveList[T eval] struct {
	moves  []OrderedMove[T] // moves will be sorted later
	Length int              // number of moves in move-list
}

// PickMove finds the best move (move with the highest eval) from the
// unsorted moves and puts it at the index position.
func (list *OrderedMoveList[T]) PickMove(index int) Move {
	// perform a single selection sort iteration
	// the full array is not sorted as most of the moves
	// will not be searched due to alpha-beta pruning

	bestIndex := index
	bestScore := list.moves[index].Eval()

	for i := index + 1; i < list.Length; i++ {
		if eval := list.moves[i].Eval(); eval > bestScore {
			bestIndex = i
			bestScore = eval
		}
	}

	list.swap(index, bestIndex)

	return list.moves[index].Move()
}

func (list *OrderedMoveList[T]) swap(i, j int) {
	list.moves[i], list.moves[j] = list.moves[j], list.moves[i]
}

// NewOrdered creates a new ordered move with the provided move and
// evaluation. The evaluation's type should belong to eval.
func NewOrdered[T eval](m Move, eval T) OrderedMove[T] {
	// [ evaluation 32 bits ] [ move 32 bits ]
	return OrderedMove[T](uint64(eval)<<32 | uint64(m))
}

// An OrderedMove represents a move that can be ranked in a move-list.
type OrderedMove[T eval] uint64

// Eval returns the OrderedMove's eval.
func (m OrderedMove[T]) Eval() T {
	return T(m >> 32)
}

// Move returns the OrderedMove's move.
func (m OrderedMove[T]) Move() Move {
	return Move(m & 0xFFFFFFFF)
}
