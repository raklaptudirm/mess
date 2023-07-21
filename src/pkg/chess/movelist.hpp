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

#ifndef CHESS_MOVE_LIST
#define CHESS_MOVE_LIST

#include <array>

#include "move.hpp"

namespace Chess {
    class MoveList {
        std::array<Move, Move::MaxInPosition> moves;
        int length = 0;

        public:
            [[nodiscard]] constexpr inline int32 Length() const {
                return length;
            }

            constexpr inline void Clear() {
                length = 0;
            }

            constexpr inline void operator +=(Move move) {
                moves[length++] = move;
            }

            constexpr inline Move operator [](int index) const {
                return moves[index];
            }

            // Iterator implements an iterator structure so that BitBoards can
            // be used inside range-for loops. The Iterator structure also keeps
            // the underlying BitBoard intact.
            struct Iterator {
                private:
                    const Move* internal;

                public:
                    // Constructor to convert a pointer to a Move (which is usually
                    // at some index inside a MoveList) to a MoveList Iterator.
                    explicit Iterator(const Move* ptr) : internal(ptr) {}

                    // Definition of the ++ operator in a for loop. It increments the
                    // pointer so that it points to the next element in the list.
                    inline Iterator operator ++() {
                        internal++;
                        return *this;
                    }

                    // Definition of an equality check with another
                    // Iterator at the rhs, usually Iterator(length).
                    inline bool operator !=(const Iterator& rhs) const {
                        return internal != rhs.internal;
                    }

                    // Definition of the dereference operator to convert
                    // an Iterator to the yield type Move. This returns
                    // the Move that the current pointer points to.
                    const Move& operator *() const {
                        return *internal;
                    }
            };

            // Definition of begin and end functions for construction an
            // iterator for the MoveList. The begin function returns an
            // Iterator with a pointer to the first element of the MoveList,
            // while the end function returns an Iterator with a pointer to
            // the 'length' element of the MoveList
            [[nodiscard]] Iterator begin() const { return Iterator(&moves[0x0000]); }
            [[nodiscard]] Iterator end()   const { return Iterator(&moves[length]); }
    };
}

#endif