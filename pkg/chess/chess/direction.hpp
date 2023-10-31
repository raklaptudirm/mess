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

#ifndef CHESS_DIRECTION
#define CHESS_DIRECTION

#include <cstdint>

#include "color.hpp"

namespace Chess {
    struct Direction {
        int8_t internal = 0;

        constexpr inline explicit Direction(int8_t direction) : internal(direction) {}

        constexpr inline explicit operator int8_t() const {
            return static_cast<int8_t>(internal);
        }

        constexpr inline bool operator ==(const Direction&) const = default;

        constexpr inline Direction operator +(const Direction rhs) const {
            return Direction(internal + static_cast<int8_t>(rhs));
        }

        constexpr inline Direction operator -() const {
            return Direction(-internal);
        }
    };

    namespace Directions {
        // The Cardinal Directions.
        constexpr Direction North = Direction(+8), East = Direction(+1);
        constexpr Direction South = Direction(-8), West = Direction(-1);

        // Combinations of Cardinal Directions.
        constexpr Direction NorthEast = North + East;
        constexpr Direction NorthWest = North + West;
        constexpr Direction SouthEast = South + East;
        constexpr Direction SouthWest = South + West;

        // No Direction.
        constexpr Direction None = Direction(0);

        template<Color STM>
        constexpr inline Direction Up() {
            if (STM == Color::White) return Directions::North;
            else                     return Directions::South;
        }

        template<Color STM>
        constexpr inline Direction Down() {
            if (STM == Color::White) return Directions::South;
            else                     return Directions::North;
        }

        static constexpr inline Direction Up(Color stm) {
            if (stm == Color::White) return Directions::North;
            else                     return Directions::South;
        }

        constexpr inline Direction Down(Color stm) {
            if (stm == Color::White) return Directions::South;
            else                     return Directions::North;
        }
    }
}

#endif