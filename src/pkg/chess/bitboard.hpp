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

#ifndef CHESS_BITBOARD
#define CHESS_BITBOARD

#include <bit>
#include <array>

#include "../util/types.hpp"
#include "../util/reverse.hpp"

#include "square.hpp"

namespace Chess {

    // A BitBoard efficiently represents a set of squares from the chessboard.
    // It also provides functions which enable easy manipulation of the set.
    struct BitBoard {
    private:
        // Internal uint64 representation of the BitBoard.
        uint64 internal;

    public:
        /* *************************
         * Constructor Definitions *
         ************************* */

        [[maybe_unused]] constexpr BitBoard() = default;

        // Constructor to convert uint64 to a BitBoard.
        [[maybe_unused]] constexpr explicit BitBoard(const uint64 bb)
            : internal(bb) {}

        // Constructor to convert Square to a BitBoard.
        [[maybe_unused]] constexpr explicit BitBoard(const Square square)
            : internal(1ull << static_cast<uint8>(square)) {}

        /* *********************
         * Methods Definitions *
         ********************* */

        // Some checks if the target BitBoard is populated.
        [[maybe_unused]] [[nodiscard]] constexpr inline bool Some() const {
            return !Empty();
        }

        // IsEmpty checks if the target BitBoard is empty.
        [[maybe_unused]] [[nodiscard]] constexpr inline bool Empty() const {
            return !internal;
        }

        // Several checks if the BitBoard has more than 1 element.
        [[maybe_unused]] [[nodiscard]] constexpr inline bool Several() const {
            return internal & (internal - 1);
        }

        // Singular checks if the BitBoard has a single element.
        [[maybe_unused]] [[nodiscard]] constexpr inline bool Singular() const {
            return Some() && !Several();
        }

        // PopCount counts the number of elements in the BitBoard.
        [[maybe_unused]] [[nodiscard]] constexpr inline int32 PopCount() const {
            return std::popcount(internal);
        }

        // IsDisjoint check if the target and the given BitBoard are disjoint,
        // i.e. don't have any elements(squares) in common between them.
        [[maybe_unused]] [[nodiscard]] constexpr inline bool IsDisjoint(BitBoard bb) const {
            return (*this & bb).Empty();
        }

        // Reverse reverses the given BitBoard.
        [[maybe_unused]] [[nodiscard]] constexpr inline BitBoard Reverse() const {
            return BitBoard(reverse(internal));
        }

        // LSB finds the least significant set-bit from the BitBoard.
        [[maybe_unused]] [[nodiscard]] constexpr inline Square LSB() const {
            return static_cast<Square>(std::countr_zero(internal));
        }

        // MSB finds the most significant set-bit from the BitBoard.
        [[maybe_unused]] [[nodiscard]] constexpr inline Square MSB() const {
            return static_cast<Square>(std::countl_zero(internal) ^ 63);
        }

        // Flip flips the given square in the BitBoard, i.e. removes
        // it if it is present in the set and vice versa.
        [[maybe_unused]] constexpr inline void Flip(Square square) {
            internal ^= static_cast<uint64>(BitBoard(square));
        }

        // PopLSB removes the least significant set-bit from the BitBoard.
        [[maybe_unused]] constexpr inline Square PopLSB() {
            Square lsb = LSB();
            internal = internal & (internal - 1);
            return lsb;
        }

        // PopMSB removes the most significant set-bit from the BitBoard.
        [[maybe_unused]] constexpr inline Square PopMSB() {
            Square msb = MSB();
            internal ^= static_cast<uint64>(BitBoard(msb));
            return msb;
        }

        /* **********************
         * Conversion Functions *
         *********************** */

        // Conversion function to convert the BitBoard into an uint64.
        [[maybe_unused]] constexpr inline explicit operator uint64() const {
            return internal;
        }

        /* **********************
         * Operator Definitions *
         ********************** */

        [[maybe_unused]] constexpr inline bool operator==(const BitBoard &) const = default;

        [[maybe_unused]] constexpr inline BitBoard operator|(const BitBoard bb) const {
            return BitBoard(internal | static_cast<uint64>(bb));
        }

        [[maybe_unused]] constexpr inline void operator|=(const BitBoard bb) {
            internal |= static_cast<uint64>(bb);
        }

        [[maybe_unused]] constexpr inline BitBoard operator&(const BitBoard bb) const {
            return BitBoard(internal & static_cast<uint64>(bb));
        }

        [[maybe_unused]] constexpr inline void operator&=(const BitBoard bb) {
            internal &= static_cast<uint64>(bb);
        }

        [[maybe_unused]] constexpr inline BitBoard operator^(const BitBoard bb) const {
            return BitBoard(internal ^ static_cast<uint64>(bb));
        }

        [[maybe_unused]] constexpr inline void operator^=(const BitBoard bb) {
            internal ^= static_cast<uint64>(bb);
        }

        [[maybe_unused]] constexpr inline BitBoard operator~() const {
            return BitBoard(~internal);
        }

        [[maybe_unused]] constexpr inline BitBoard operator!() const {
            return BitBoard(~internal);
        }

        [[maybe_unused]] constexpr inline BitBoard operator+(const BitBoard bb) const {
            return *this | bb;
        }

        [[maybe_unused]] constexpr inline void operator+=(const BitBoard bb) {
            internal |= static_cast<uint64>(bb);
        }

        [[maybe_unused]] constexpr inline BitBoard operator-(const BitBoard bb) const {
            return *this & ~bb;
        }

        [[maybe_unused]] constexpr inline void operator-=(const BitBoard bb) {
            *this &= ~bb;
        }

        [[maybe_unused]] constexpr inline BitBoard operator+(const Square square) const {
            return *this | BitBoard(square);
        }

        [[maybe_unused]] constexpr inline void operator+=(const Square square) {
            *this |= BitBoard(square);
        }

        [[maybe_unused]] constexpr inline BitBoard operator-(const Square square) const {
            return *this - BitBoard(square);
        }

        [[maybe_unused]] constexpr inline void operator-=(const Square square) {
            *this -= BitBoard(square);
        }

        // Definition of the less-than-equal operator, which checks if the
        // target BitBoard is a subset of the rhs BitBoard.
        [[maybe_unused]] constexpr inline bool operator<=(const BitBoard bb) const {
            return (*this & bb) == *this;
        }

        // Definition of the greater-than-equal operator, which checks if the
        // target BitBoard is a superset of the rhs BitBoard.
        [[maybe_unused]] constexpr inline bool operator>=(const BitBoard bb) const {
            return (*this & bb) == bb;
        }

        // Definition of the less-than operator, which checks if the target
        // BitBoard is a proper subset of the rhs BitBoard.
        [[maybe_unused]] constexpr inline bool operator<(const BitBoard bb) const {
            return *this <= bb && *this != bb;
        }

        // Definition of the greater-than operator, which checks if the rhs
        // BitBoard is a proper subset of the target BitBoard.
        [[maybe_unused]] constexpr inline bool operator>(const BitBoard bb) const {
            return *this >= bb && *this != bb;
        }

        // Definition of the indexing with Squares on a BitBoard.
        [[maybe_unused]] constexpr inline bool operator[](const Square square) const {
            return (*this & BitBoard(square)).Some();
        }

        [[maybe_unused]] constexpr inline BitBoard operator>>(const Direction direction) const {
            const auto shift = static_cast<int8>(direction);

            if (direction == Directions::North || direction == Directions::North+Directions::North)
                return BitBoard(internal << shift);
            if (direction == Directions::South || direction == Directions::South+Directions::South)
                return BitBoard(internal >> -shift);

            constexpr uint64 NOT_FILE_A = ~0x0101010101010101ULL;
            constexpr uint64 NOT_FILE_H = ~0x8080808080808080ULL;

            if (direction == Directions::West || direction == Directions::SouthWest)
                return BitBoard((internal & NOT_FILE_A) >> -shift);
            if (direction == Directions::NorthWest)
                return BitBoard((internal & NOT_FILE_A) << shift);

            if (direction == Directions::East || direction == Directions::NorthEast)
                return BitBoard((internal & NOT_FILE_H) << shift);
            if (direction == Directions::SouthEast)
                return BitBoard((internal & NOT_FILE_H) >> -shift);

            return *this;
        }

        [[maybe_unused]] constexpr inline void operator>>=(const Direction direction) {
            *this = *this >> direction;
        }

        /* **********************************
         * BitBoard Iterator Implementation *
         ********************************** */

        // Iterator implements an iterator structure so that BitBoards can
        // be used inside range-for loops. The Iterator structure also keeps
        // the underlying BitBoard intact.
        struct Iterator {
        private:
            // Internal representation of BitBoard we are iterating.
            uint64 internal;

        public:
            // Constructor to convert the given BitBoard uint64 into an iterable value.
            constexpr explicit Iterator(const uint64 bb) : internal(bb) {}

            // ++ takes the iterator forward by popping the LSB.
            // A reference to itself as required of an iterator.
            constexpr inline Iterator operator++() {
                internal = internal & (internal - 1);
                return *this;
            }

            // == implements an equality check between two Iterators.
            // Boolean describing whether the two are equal or not.
            constexpr inline bool operator ==(const Iterator&) const = default;

            // * operator finds the least significant set bit in the uint64.
            // Square representing the least significant set bit.
            constexpr Square operator*() const {
                return BitBoard(internal).LSB();
            }
        };

        /* ************************************************************
         * Definition of begin and end functions for construction an  *
         * iterator for the BitBoard. The begin function returns an   *
         * Iterator on the internal uint64, while the end function    *
         * returns an Iterator on 0, which is the end result for most *
         * BitBoard iterations.                                       *
         ************************************************************ */

        // begin functions returns an iterator for the BitBoard.
        [[nodiscard]]        constexpr inline Iterator begin() const { return Iterator(internal); }
        // end function returns an iterator for the empty BitBoard.
        [[nodiscard]] static constexpr inline Iterator end  ()       { return Iterator(0x000000); }

        [[nodiscard]] constexpr inline std::string ToString() const {
            std::string str;

            for (uint8 rank = 7; rank != 255; rank--) {
                for (uint8 file = 0; file < File::N; file++) {
                    str += (*this)[Square(rank * 8 + file)] ? "1 " : "0 ";
                }

                str += "\n";
            }

            return str;
        }

        /* ********************************
         * Static Functions for BitBoards *
         ******************************** */

        // Hyperbola implements the Hyperbola Quintessence algorithm for calculating ray
        // attacks. It uses the o - 2r trick to find the ray from the given blockers.
        //
        // square:   Square of the attacker.
        // blockers: Blockers blocking the attacks.
        // mask:     Mask of the ray attack.
        constexpr static inline BitBoard Hyperbola(Square square, BitBoard blockers, BitBoard mask) {
            const uint64 r = static_cast<uint64>(BitBoard(square)); // Piece's BitBoard as an uint64.
            const uint64 o = static_cast<uint64>(blockers & mask);  // Position's Masked Occupancy.

            // Calculate attack-set along the mask using the o - 2r trick.
            return BitBoard((o - 2 * r) ^ reverse(reverse(o) - 2 * reverse(r))) & mask;
        }
    };

    [[maybe_unused]] constexpr inline std::ostream &operator<<(std::ostream &os, const BitBoard &bb) {
        os << bb.ToString();
        return os;
    }

    namespace BitBoards {
        namespace {
            constexpr std::array<uint64, File::N> files = {
                    0x0101010101010101,
                    0x0202020202020202,
                    0x0404040404040404,
                    0x0808080808080808,
                    0x1010101010101010,
                    0x2020202020202020,
                    0x4040404040404040,
                    0x8080808080808080,
            };

            constexpr std::array<uint64, Rank::N> ranks = {
                    0x00000000000000FF,
                    0x000000000000FF00,
                    0x0000000000FF0000,
                    0x00000000FF000000,
                    0x000000FF00000000,
                    0x0000FF0000000000,
                    0x00FF000000000000,
                    0xFF00000000000000,
            };

            constexpr std::array<uint64, 15> diagonals = {
                    0x0000000000000080,
                    0x0000000000008040,
                    0x0000000000804020,
                    0x0000000080402010,
                    0x0000008040201008,
                    0x0000804020100804,
                    0x0080402010080402,
                    0x8040201008040201,
                    0x4020100804020100,
                    0x2010080402010000,
                    0x1008040201000000,
                    0x0804020100000000,
                    0x0402010000000000,
                    0x0201000000000000,
                    0x0100000000000000,
            };
            constexpr std::array<uint64, 15> antiDiagonals = {
                    0x0000000000000001,
                    0x0000000000000102,
                    0x0000000000010204,
                    0x0000000001020408,
                    0x0000000102040810,
                    0x0000010204081020,
                    0x0001020408102040,
                    0x0102040810204080,
                    0x0204081020408000,
                    0x0408102040800000,
                    0x0810204080000000,
                    0x1020408000000000,
                    0x2040800000000000,
                    0x4080000000000000,
                    0x8000000000000000,
            };
        }

        [[maybe_unused]] constexpr BitBoard Empty = BitBoard(0);
        [[maybe_unused]] constexpr BitBoard Full = BitBoard(~0);

        [[maybe_unused]] constexpr BitBoard White = BitBoard(0x55AA55AA55AA55AA);
        [[maybe_unused]] constexpr BitBoard Black = BitBoard(0xAA55AA55AA55AA55);

        [[maybe_unused]] constexpr BitBoard Edges = BitBoard(0xff818181818181ff);

        [[maybe_unused]] constexpr inline BitBoard File(Chess::File file) {
            if (file == File::None) return Empty;
            return BitBoard(files[static_cast<uint8>(file)]);
        }

        [[maybe_unused]] constexpr inline BitBoard Rank(Chess::Rank rank) {
            return BitBoard(ranks[static_cast<uint8>(rank)]);
        }

        [[maybe_unused]] constexpr inline BitBoard Diagonal(uint8 diagonal) {
            return BitBoard(diagonals[diagonal]);
        }

        [[maybe_unused]] constexpr inline BitBoard AntiDiagonal(uint8 antiDiagonal) {
            return BitBoard(antiDiagonals[antiDiagonal]);
        }

        namespace {
            constexpr std::array<std::array<BitBoard, Square::N>, Square::N> between = []() {
                std::array<std::array<BitBoard, Square::N>, Square::N> between = {};

                for (uint8 square1 = 0; square1 < Square::N; square1++) {
                    for (uint8 square2 = 0; square2 < Square::N; square2++) {
                        const Square sq1 = Square(square1);
                        const Square sq2 = Square(square2);

                        if (sq1 == sq2) continue;

                        BitBoard mask;

                        // Check for lateral blockerMask.
                        if (sq1.Rank() == sq2.Rank()) mask = BitBoards::Rank(sq1.Rank());
                        else if (sq1.File() == sq2.File()) mask = BitBoards::File(sq1.File());

                            // Check for diagonal blockerMask.
                        else if (sq1.Diagonal() == sq2.Diagonal()) mask = BitBoards::Diagonal(sq1.Diagonal());
                        else if (sq1.AntiDiagonal() == sq2.AntiDiagonal())
                            mask = BitBoards::AntiDiagonal(sq1.AntiDiagonal());

                            // No blockerMask between the two squares.
                        else continue;

                        const BitBoard blockers = BitBoard(sq1) + BitBoard(sq2);
                        between[square1][square2] = BitBoard::Hyperbola(sq1, blockers, mask) &
                                                    BitBoard::Hyperbola(sq2, blockers, mask);
                    }
                }

                return between;
            }();
        }

        [[maybe_unused]] constexpr inline BitBoard Between(Square square1, Square square2) {
            return between[static_cast<uint8>(square1)][static_cast<uint8>(square2)];
        }

        [[maybe_unused]] constexpr inline BitBoard Between1(Square square1, Square square2) {
            return Between(square1, square2) + BitBoard(square1);
        }

        [[maybe_unused]] constexpr inline BitBoard Between2(Square square1, Square square2) {
            return Between(square1, square2) + BitBoard(square2);
        }

        [[maybe_unused]] constexpr inline BitBoard Between12(Square square1, Square square2) {
            return Between(square1, square2) + BitBoard(square1) + BitBoard(square2);
        }
    }
}

#endif //CHESS_BITBOARD