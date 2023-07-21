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

#ifndef CHESS_COLOR
#define CHESS_COLOR

#include <string>
#include <assert.h>

#include "../util/types.hpp"

namespace Chess {
    struct Color {
        static const int N = 2;

        enum internal_type : uint8 {
            White, Black, None
        };

        internal_type internal = None;

        constexpr inline Color() = default;

        constexpr explicit inline Color(uint8 color) {
            internal = static_cast<internal_type>(color);
        }

        constexpr inline Color(internal_type color) {
            internal = color;
        }

        constexpr inline explicit Color(const std::string& color) {
            assert(color == "w" || color == "b");
            internal = color == "w" ? Color::White : Color::Black;
        }

        constexpr inline Color operator !() const {
            return Color(static_cast<uint8>(internal) ^ 1);
        }

        constexpr inline bool operator ==(const Color&) const = default;

        constexpr inline explicit operator uint8() const {
            return static_cast<uint8>(internal);
        }

        [[nodiscard]] constexpr inline std::string ToString() const {
            if (internal == Color::White) return "w";
            if (internal == Color::Black) return "b";
            return "-";
        }
    };

    constexpr inline std::ostream& operator<<(std::ostream& os, const Color& color) {
        os << color.ToString();
        return os;
    }
}

#endif