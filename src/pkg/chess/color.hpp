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
#include <cassert>

#include "../util/types.hpp"

namespace Chess {
    // Color struct represents the chess colors. Usually this refers
    // to the colors of the chess pieces (White and Black), but this
    // type can also be used to refer to colors of squares on the
    // chessboard which are also one of White and Black.
    struct Color {
        // Constant representing number of valid Colors.
        static const int N = 2;

        // internal_type is the internal representation of a Color
        // used by this struct. Values of internal_type are implicitly
        // cast to values of type Color, and this type should
        // not be used or depended on outside of this definition.
        enum internal_type : uint8 {
            White, Black, None
        };

        internal_type internal = None;

        constexpr inline Color() = default;

        // Constructor to create a Color from its uint8 representation.
        constexpr explicit inline Color(uint8 color)
            : internal(static_cast<internal_type>(color)) {}

        // Constructor to create a Color from the given internal representation.
        // The constructor is not marked explicit cause the internal_type internal type
        // should implicitly act as a value of Color type outside this declaration.
        // NOLINTNEXTLINE(google-explicit-constructor)
        constexpr inline Color(internal_type color)
            : internal(color) {}

        // Constructor to create a new Color from its string representation.
        // White is represented by "w" and Black is represented by "b", and
        // the provided string is expected to be one of those two options.
        constexpr inline explicit Color(const std::string& color) {
            assert(color == "w" || color == "b");
            internal = color == "w" ? Color::White : Color::Black;
        }

        // The Bang (!) operator converts the current Color to the "other"
        // value, i.e. White -> Black or Black -> White.
        constexpr inline Color operator !() const {
            // ^ 1 flips the lsb which is used to represent the Color.
            return Color(static_cast<uint8>(internal) ^ 1);
        }

        // The EqualEqual (==) operator checks if two Colors are equal.
        // This definition in implicitly populated by the compiler. The
        // compiler also implicitly populates the definition of the
        // NotEqual (!=) operator which checks if two Colors are not equal.
        constexpr inline bool operator ==(const Color&) const = default;

        // uint8 is the conversion function to convert a Color to an uint8.
        constexpr inline explicit operator uint8() const {
            return static_cast<uint8>(internal);
        }

        // ToString converts the target Color to its string representation.
        // White -> "w" or Black -> "b" or finally None -> "-"
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