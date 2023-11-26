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

#ifndef CHESS_SQUARE
#define CHESS_SQUARE

#include <string>
#include <cassert>
#include <cstdint>

#include "direction.hpp"

namespace Chess {
    struct File {
        /* ******************************
         * Internal Enum Representation *
         ****************************** */

        // Number of Files, excluding None, on a Chessboard.
        static const int N = 8;

        // The internal enum representation of a File.
        enum internal_type : uint8_t {
            A, B, C, D, E, F, G, H, None
        };

        // The variable that stores the internal representation.
        internal_type internal = internal_type::None;

        /***************************
         * Constructor Definitions *
         ***************************/

        // Constructor to convert an uint8_t into a File.
        // The File with the given uint8_t representation.
        constexpr explicit File(uint8_t file) : internal(static_cast<internal_type>(file)) {}

        // Constructor to convert an internal representation into a File.
        // The File with the given internal representation.
        // NOLINTNEXTLINE(google-explicit-constructor)
        constexpr File(internal_type file) : internal(file) {}

        // Constructor to parse a string for a File.
        // Files are represented by the lowercase characters a-h.
        // File represented by the given string.
        constexpr inline explicit File(std::string file) {
            assert(file.length() == 1); assert('a' <= file.at(0) && file.at(0) <= 'h');
            internal = static_cast<internal_type>(static_cast<uint8_t>(file.at(0) - 'a'));
        }

        /************************
         * Conversion Functions *
         ************************/

        //   Conversion function to convert a file into its uint8_t representation.
        // The target file's uint8_t representation.
        constexpr inline explicit operator uint8_t() const {
            return static_cast<uint8_t>(internal);
        }

        //   Conversion function to convert a file into its string representation.
        // The target file's string representation.
        [[nodiscard]] constexpr inline std::string ToString() const {
            return {1, static_cast<char>(internal + 'a')};
        }

        /************************
         * Operator Definitions *
         ************************/

        constexpr inline bool operator ==(const File&) const = default;

        constexpr inline bool operator <(const File file) const {
            return internal < static_cast<uint8_t>(file);
        }

        constexpr inline bool operator >(const File file) const {
            return internal > static_cast<uint8_t>(file);
        }
    };

    struct Rank {
        /********************************
         * Internal Enum Representation *
         ********************************/

        static const int N = 8;

        enum internal_type : uint8_t {
            First, Second, Third, Fourth, Fifth, Sixth, Seventh, Eighth, None
        };

        internal_type internal = internal_type::None;

        /***************************
         * Constructor Definitions *
         ***************************/

        // NOLINTNEXTLINE(google-explicit-constructor)
        constexpr inline Rank(internal_type rank) {
            internal = rank;
        }

        constexpr inline explicit Rank(uint8_t rank) {
            internal = static_cast<internal_type>(rank);
        }

        constexpr inline explicit Rank(std::string rank) {
            assert(rank.length() == 1); assert('1' <= rank.at(0) && rank.at(0) <= '8');
            internal = static_cast<internal_type>(static_cast<uint8_t>(rank.at(0) - '1'));
        }

        /************************
         * Conversion Functions *
         ************************/

        constexpr inline explicit operator uint8_t() const {
            return static_cast<uint8_t>(internal);
        }

        [[nodiscard]] constexpr inline std::string ToString() const {
            return {1, static_cast<char>(internal + '1')};
        }

        /************************
         * Operator Definitions *
         ************************/

        constexpr inline bool operator ==(const Rank&) const = default;
    };

    struct Square {
        /********************************
         * Internal Enum Representation *
         ********************************/

        static constexpr int N = 64;

        enum internal_type : uint8_t {
            A1, B1, C1, D1, E1, F1, G1, H1,
            A2, B2, C2, D2, E2, F2, G2, H2,
            A3, B3, C3, D3, E3, F3, G3, H3,
            A4, B4, C4, D4, E4, F4, G4, H4,
            A5, B5, C5, D5, E5, F5, G5, H5,
            A6, B6, C6, D6, E6, F6, G6, H6,
            A7, B7, C7, D7, E7, F7, G7, H7,
            A8, B8, C8, D8, E8, F8, G8, H8, None
        };

        internal_type internal = internal_type::None;

        /***************************
         * Constructor Definitions *
         ***************************/

        constexpr Square() = default;

        // NOLINTNEXTLINE(google-explicit-constructor)
        constexpr inline Square(internal_type square) {
            internal = square;
        }

        constexpr inline explicit Square(uint8_t square) {
            internal = static_cast<internal_type>(square);
        }

        constexpr inline explicit Square(File file, Rank rank) {
            internal = static_cast<internal_type>(
                static_cast<uint8_t>(rank)*8 +
                static_cast<uint8_t>(file)
            );
        }

        constexpr inline explicit Square(const std::string& square) {
            if (square == "-") {
                internal = Square::None;
                return;
            }

            assert(square.length() == 2);

            const Chess::File file = Chess::File(square.substr(0, 1));
            const Chess::Rank rank = Chess::Rank(square.substr(1, 1));

            *this = Square(file, rank);
        }

        /*****************************
         * Property Getter Functions *
         *****************************/

        [[nodiscard]] constexpr inline File File() const {
            if (*this == None) return Chess::File::None;
            return Chess::File(static_cast<uint8_t>(internal) % Chess::File::N);
        }

        [[nodiscard]] constexpr inline Rank Rank() const {
            return Chess::Rank(static_cast<uint8_t>(internal) / Chess::File::N);
        }

        [[nodiscard]] constexpr inline uint8_t Diagonal() const {
            return 7 + static_cast<uint8_t>(this->Rank()) - static_cast<uint8_t>(this->File());
        }

        [[nodiscard]] constexpr inline uint8_t AntiDiagonal() const {
            return static_cast<uint8_t>(this->Rank()) + static_cast<uint8_t>(this->File());
        }

        /************************
         * Conversion Functions *
         ************************/

        constexpr inline explicit operator uint8_t() const {
            return static_cast<uint8_t>(internal);
        }

        [[nodiscard]] constexpr inline std::string ToString() const {
            if (internal == Square::None) return "-";
            return this->File().ToString() + this->Rank().ToString();
        }

        /************************
         * Operator Definitions *
         ************************/

        constexpr inline bool operator ==(const Square&) const = default;

        constexpr inline Square operator >>(const Direction direction) const {
            return Square(
                static_cast<int8_t>(static_cast<uint8_t>(internal)) +
                static_cast<int8_t>(direction)
            );
        }
    };

    constexpr inline std::ostream& operator<<(std::ostream& os, const File& file) {
        os << file.ToString();
        return os;
    }

    constexpr inline std::ostream& operator<<(std::ostream& os, const Rank& rank) {
        os << rank.ToString();
        return os;
    }

    constexpr inline std::ostream& operator<<(std::ostream& os, const Square& square) {
        os << square.ToString();
        return os;
    }
}


#endif