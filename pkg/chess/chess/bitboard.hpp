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
#include <cstdint>

#include "square.hpp"

namespace Chess {

    // A BitBoard efficiently represents a set of squares from the chessboard.
    // It also provides functions which enable easy manipulation of the set.
    struct BitBoard {
    private:
        // Internal uint64_t representation of the BitBoard.
        uint64_t internal;

    public:
        /* *************************
         * Constructor Definitions *
         ************************* */

        [[maybe_unused]] constexpr BitBoard() = default;

        // Constructor to convert uint64_t to a BitBoard.
        [[maybe_unused]] constexpr explicit BitBoard(const uint64_t bb)
            : internal(bb) {}

        // Constructor to convert Square to a BitBoard.
        [[maybe_unused]] constexpr explicit BitBoard(const Square square)
            : internal(1ull << static_cast<uint8_t>(square)) {}

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
            // The (internal - 1) gives a number with the lsb set to 0 and all
            // the lower bits set to 1. Doing a bitwise and of this number with
            // the original removes just the lsb (& 0) since all the lower bits
            // are already 0 by definition of the lsb. Therefore, the whole
            // operation is equivalent to a lsb-pop, thus making the number
            // 0 (false) if there are only 0-1 bits in the number.
            return internal & (internal - 1);
        }

        // Singular checks if the BitBoard has a single element.
        [[maybe_unused]] [[nodiscard]] constexpr inline bool Singular() const {
            return Some() && !Several();
        }

        // PopCount counts the number of elements in the BitBoard.
        [[maybe_unused]] [[nodiscard]] constexpr inline int32_t PopCount() const {
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
            internal ^= static_cast<uint64_t>(BitBoard(square));
        }

        // PopLSB removes the least significant set-bit from the BitBoard.
        [[maybe_unused]] constexpr inline Square PopLSB() {
            Square lsb = LSB();

            // Specifics of this operation are described in the documentation of
            // the Several function which uses the same lsb-popping mechanism.
            internal = internal & (internal - 1);
            return lsb;
        }

        // PopMSB removes the most significant set-bit from the BitBoard.
        [[maybe_unused]] constexpr inline Square PopMSB() {
            Square msb = MSB();
            internal ^= static_cast<uint64_t>(BitBoard(msb));
            return msb;
        }

        /* **********************
         * Conversion Functions *
         ********************** */

        // Conversion function to convert the BitBoard into an uint64_t.
        [[maybe_unused]] constexpr inline explicit operator uint64_t() const {
            return internal;
        }

        /* **********************
         * Operator Definitions *
         ********************** */

        // The == operator checks for equality between the two provided BitBoards.
        // The definition of the not-equal/ operator is also automatically generated
        // from this function.
        [[maybe_unused]] constexpr inline bool operator ==(const BitBoard &) const = default;

        // The ~ operator implements the set complement operation on the provided
        // BitBoard. A complement operation is defined as the operation which returns a
        // set containing all the elements missing from its operand.
        [[maybe_unused]] constexpr inline BitBoard operator ~() const {
            return BitBoard(~internal);
        }

        // The + operator implements the set union operation between the two
        // provided BitBoards. A union operation is defined as the operation
        // which returns a set containing all the elements present in either of
        // its two operands.
        [[maybe_unused]] constexpr inline BitBoard operator +(const BitBoard bb) const {
            return BitBoard(internal | static_cast<uint64_t>(bb));
        }

        // The += operator is the assigning version of the + operator, and it performs
        // a union operation with the target variable of the assignment.
        [[maybe_unused]] constexpr inline void operator +=(const BitBoard bb) {
            internal |= static_cast<uint64_t>(bb);
        }

        // The & operator implements the set intersection operation between the two
        // provided BitBoards. An intersection operation is defined as the operation
        // which returns a set containing the common elements of its operands.
        [[maybe_unused]] constexpr inline BitBoard operator & (const BitBoard bb) const {
            return BitBoard(internal & static_cast<uint64_t>(bb));
        }

        // The &= operator is the assigning version of the & operator, and it performs
        // an intersection operation with the target variable of the assignment.
        [[maybe_unused]] constexpr inline void operator &= (const BitBoard bb) {
            internal &= static_cast<uint64_t>(bb);
        }

        // The - operator implements the set difference operation between the two
        // provided BitBoards. A set difference operation is defined as the operation
        // which returns a set containing all the elements present in the first set
        // but not present in the second one.
        [[maybe_unused]] constexpr inline BitBoard operator -(const BitBoard bb) const {
            return *this & ~bb;
        }

        // The -= operator is the assigning version for the - operator, and it performs
        // an exclusive or operation with the target variable of the assignment.
        [[maybe_unused]] constexpr inline void operator -=(const BitBoard bb) {
            *this &= ~bb;
        }

        // The ^ operator implements the set exclusive or operation between the two
        // provided BitBoards. An exclusive or operation is defined as the operation
        // which returns a set containing the elements missing or shared between its
        // operands.
        [[maybe_unused]] constexpr inline BitBoard operator ^ (const BitBoard bb) const {
            return BitBoard(internal ^ static_cast<uint64_t>(bb));
        }

        // The ^= operator is the assigning version for the ^ operator, and it performs
        // an exclusive or operation with the target variable of the assignment.
        [[maybe_unused]] constexpr inline void operator ^=(const BitBoard bb) {
            internal ^= static_cast<uint64_t>(bb);
        }

        // The + Square operator is an overload which converts the provided square to its
        // BitBoard representation and then performs a set union operation on its operands.
        [[maybe_unused]] constexpr inline BitBoard operator +(const Square square) const {
            return *this + BitBoard(square);
        }

        // The += Square operator is the assigning version of the + Square operator, which
        // performs the same operation with the target variable of the assignment.
        [[maybe_unused]] constexpr inline void operator +=(const Square square) {
            *this += BitBoard(square);
        }

        // The - Square operator is an overload which converts the provided square to its
        // BitBoard representation and then performs a set difference operation on its operands.
        [[maybe_unused]] constexpr inline BitBoard operator -(const Square square) const {
            return *this - BitBoard(square);
        }

        // The -= Square operator is the assigning version of the - Square operator, which
        // performs the same operation with the target variable of the assignment.
        [[maybe_unused]] constexpr inline void operator -=(const Square square) {
            *this -= BitBoard(square);
        }

        // The <= operator implements the subset checking operation between the two
        // provided BitBoards. A BitBoard is said to be the subset of another if the
        // other BitBoard contains all the elements present in this BitBoard.
        [[maybe_unused]] constexpr inline bool operator <=(const BitBoard bb) const {
            return (*this & bb) == *this;
        }

        // The >= operator implements the superset checking operation between the two
        // provided BitBoards. A BitBoard is said to be the superset of another if it
        // contains all the elements present in the other BitBoard.
        [[maybe_unused]] constexpr inline bool operator >=(const BitBoard bb) const {
            return (*this & bb) == bb;
        }

        // The < operator implements the proper subset checking operation between the two
        // provided BitBoards. A BitBoard is said to be the proper subset of another if
        // it is both a subset and not equal to the other BitBoard.
        [[maybe_unused]] constexpr inline bool operator <(const BitBoard bb) const {
            return *this <= bb && *this != bb;
        }

        // The > operator implements the proper superset checking operation between the two
        // provided BitBoards. A BitBoard is said to be the proper superset of another if
        // it is both a superset and not equal to the other BitBoard.
        [[maybe_unused]] constexpr inline bool operator >(const BitBoard bb) const {
            return *this >= bb && *this != bb;
        }

        // The [] operator implements indexing on the BitBoard with a Square. An index
        // operation returns a boolean representing if the BitBoard contains the square.
        [[maybe_unused]] constexpr inline bool operator [](const Square square) const {
            return (*this & BitBoard(square)).Some();
        }

        // The >> operator implements the BitBoard shift operation which shifts the
        // BitBoard in a given Direction. Shifting a BitBoard is equivalent to replacing
        // each of its element squares with another square where the difference between
        // the old and the new squares are the same and determined by the Direction.
        [[maybe_unused]] constexpr inline BitBoard operator >>(const Direction direction) const {
            const auto shift = static_cast<int8_t>(direction);

            // Straight up and down (and double that) shifts can be done without any masking
            // because of the internal representation used by the BitBoard.
            if (direction == Directions::North || direction == Directions::North+Directions::North)
                return BitBoard(internal << shift);
            if (direction == Directions::South || direction == Directions::South+Directions::South)
                return BitBoard(internal >> -shift);

            // Shifts to the east and west however need masking to prevent spills.

            constexpr uint64_t NOT_FILE_A = ~0x0101010101010101ULL;
            constexpr uint64_t NOT_FILE_H = ~0x8080808080808080ULL;

            if (direction == Directions::West || direction == Directions::SouthWest)
                return BitBoard((internal & NOT_FILE_A) >> -shift);
            if (direction == Directions::NorthWest)
                return BitBoard((internal & NOT_FILE_A) << shift);

            if (direction == Directions::East || direction == Directions::NorthEast)
                return BitBoard((internal & NOT_FILE_H) << shift);
            if (direction == Directions::SouthEast)
                return BitBoard((internal & NOT_FILE_H) >> -shift);

            // Ignore shifts towards unknown directions.
            return *this;
        }

        // The >>= operator is the assigning version of the >> operator, which performs
        // the same operation with the target variable of the assignment.
        [[maybe_unused]] constexpr inline void operator >>=(const Direction direction) {
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
            uint64_t internal;

        public:
            // Constructor to convert the given BitBoard uint64_t into an iterable value.
            constexpr explicit Iterator(const uint64_t bb) : internal(bb) {}

            // ++ takes the iterator forward by popping the LSB.
            // A reference to itself as required of an iterator.
            constexpr inline Iterator operator ++() {
                internal = internal & (internal - 1); // Pop LSB.
                return *this;
            }

            // == implements an equality check between two Iterators.
            // Boolean describing whether the two are equal or not.
            constexpr inline bool operator ==(const Iterator&) const = default;

            // * operator finds the least significant set bit in the uint64_t.
            // Square representing the least significant set bit.
            constexpr Square operator *() const {
                return BitBoard(internal).LSB();
            }
        };

        /* ************************************************************
         * Definition of begin and end functions for construction an  *
         * iterator for the BitBoard. The begin function returns an   *
         * Iterator on the internal uint64_t, while the end function    *
         * returns an Iterator on 0, which is the end result for most *
         * BitBoard iterations.                                       *
         ************************************************************ */

        // begin functions returns an iterator for the BitBoard.
        [[nodiscard]]        constexpr inline Iterator begin() const { return Iterator(internal); }
        // end function returns an iterator for the empty BitBoard.
        [[nodiscard]] static constexpr inline Iterator end  ()       { return Iterator(0x000000); }

        [[nodiscard]] constexpr inline std::string ToString() const {
            std::string str;

            for (uint8_t rank = 7; rank != 255; rank--) {
                for (uint8_t file = 0; file < File::N; file++) {
                    str += (*this)[Square(rank * 8 + file)] ? "1 " : "0 ";
                }

                str += "\n";
            }

            return str;
        }

        /* ********************************
         * Static Functions for BitBoards *
         ******************************** */

    private:
        // reverse reverses the bits of the given uint64_t number.
        // https://graphics.stanford.edu/~seander/bithacks.html#BitReverseTable
        constexpr static inline uint64_t reverse(uint64_t n) {
            // Lookup table with precomputed reverses for each byte value.
            constexpr uint64_t BitReverseTable256[256] = {
                #define R2(n)   (n),   (n + 2*64),   (n + 1*64),   (n + 3*64)
                #define R4(n) R2(n), R2(n + 2*16), R2(n + 1*16), R2(n + 3*16)
                #define R6(n) R4(n), R4(n + 2*4 ), R4(n + 1*4 ), R4(n + 3*4 )
                R6(0), R6(2), R6(1), R6(3)
            };

            // Reverse each byte in the number and append them together in reverse.
            return (BitReverseTable256[(n >>  0) & 0xff] << 56) |
                   (BitReverseTable256[(n >>  8) & 0xff] << 48) |
                   (BitReverseTable256[(n >> 16) & 0xff] << 40) |
                   (BitReverseTable256[(n >> 24) & 0xff] << 32) |
                   (BitReverseTable256[(n >> 32) & 0xff] << 24) |
                   (BitReverseTable256[(n >> 40) & 0xff] << 16) |
                   (BitReverseTable256[(n >> 48) & 0xff] <<  8) |
                   (BitReverseTable256[(n >> 56) & 0xff] <<  0);
        }

    public:
        // Hyperbola implements the Hyperbola Quintessence algorithm for calculating ray
        // attacks. Provided with the sliding piece square, the piece's ray mask, and
        // the BitBoard of blockers, it returns a BitBoard of all the Squares to which
        // the given piece can move, without masking out any friendly squares.
        constexpr static inline BitBoard Hyperbola(Square square, BitBoard blockers, BitBoard mask) {
            const uint64_t r = static_cast<uint64_t>(BitBoard(square)); // Piece's BitBoard as an uint64_t.
            const uint64_t o = static_cast<uint64_t>(blockers & mask);  // Position's Masked Occupancy.

            // Calculate attack-set along the mask using the o - 2r trick.
            return BitBoard((o - 2 * r) ^ reverse(reverse(o) - 2 * reverse(r))) & mask;
        }
    };

    [[maybe_unused]] constexpr inline std::ostream &operator<<(std::ostream &os, const BitBoard &bb) {
        os << bb.ToString();
        return os;
    }

    namespace BitBoards {
        [[maybe_unused]] constexpr BitBoard Empty = BitBoard(0);
        [[maybe_unused]] constexpr BitBoard Full = ~BitBoard(0);

        [[maybe_unused]] constexpr BitBoard White = BitBoard(0x55AA55AA55AA55AA);
        [[maybe_unused]] constexpr BitBoard Black = BitBoard(0xAA55AA55AA55AA55);

        [[maybe_unused]] constexpr BitBoard Edges = BitBoard(0xff818181818181ff);

        // File returns the BitBoard representing the given File.
        [[maybe_unused]] constexpr inline BitBoard File(Chess::File file) {
            constexpr std::array<uint64_t, File::N> files = {
                    0x0101010101010101, 0x0202020202020202,
                    0x0404040404040404, 0x0808080808080808,
                    0x1010101010101010, 0x2020202020202020,
                    0x4040404040404040, 0x8080808080808080,
            };

            return BitBoard(files[static_cast<uint8_t>(file)]);
        }

        // Rank returns the BitBoard representing the given Rank.
        [[maybe_unused]] constexpr inline BitBoard Rank(Chess::Rank rank) {
            constexpr std::array<uint64_t, Rank::N> ranks = {
                    0x00000000000000FF, 0x000000000000FF00,
                    0x0000000000FF0000, 0x00000000FF000000,
                    0x000000FF00000000, 0x0000FF0000000000,
                    0x00FF000000000000, 0xFF00000000000000,
            };

            return BitBoard(ranks[static_cast<uint8_t>(rank)]);
        }

        // Diagonal returns the BitBoard representing the given Diagonal.
        [[maybe_unused]] constexpr inline BitBoard Diagonal(uint8_t diagonal) {
            constexpr std::array<uint64_t, 15> diagonals = {
                    0x0000000000000080, 0x0000000000008040, 0x0000000000804020,
                    0x0000000080402010, 0x0000008040201008, 0x0000804020100804,
                    0x0080402010080402, 0x8040201008040201, 0x4020100804020100,
                    0x2010080402010000, 0x1008040201000000, 0x0804020100000000,
                    0x0402010000000000, 0x0201000000000000, 0x0100000000000000,
            };

            return BitBoard(diagonals[diagonal]);
        }

        // AntiDiagonal returns the BitBoard representing the given AntiDiagonal.
        [[maybe_unused]] constexpr inline BitBoard AntiDiagonal(uint8_t antiDiagonal) {
            constexpr std::array<uint64_t, 15> antiDiagonals = {
                    0x0000000000000001, 0x0000000000000102, 0x0000000000010204,
                    0x0000000001020408, 0x0000000102040810, 0x0000010204081020,
                    0x0001020408102040, 0x0102040810204080, 0x0204081020408000,
                    0x0408102040800000, 0x0810204080000000, 0x1020408000000000,
                    0x2040800000000000, 0x4080000000000000, 0x8000000000000000,
            };

            return BitBoard(antiDiagonals[antiDiagonal]);
        }

        [[maybe_unused]] constexpr inline BitBoard Between(Square square1, Square square2) {
            // A Table of between BitBoards should be indexed with two squares, which
            // should index the between BitBoard of the two squares, i.e. a BitBoard
            // containing all the squares between the given two exclusive of both.
            using Table = std::array<std::array<BitBoard, Square::N>, Square::N>;

            // between is the above-mentioned between table. It is generated
            // automatically during compile-time using the given lambda function.
            constexpr Table between = []() {
                Table between = {};

                for (uint8_t square1 = 0; square1 < Square::N; square1++) {
                    for (uint8_t square2 = 0; square2 < Square::N; square2++) {
                        const Square sq1 = Square(square1);
                        const Square sq2 = Square(square2);

                        // Try to find a common Rank, File, or Diagonal between the squares, which,
                        // if it exists, will be a superset of the between BitBoard for the pair.
                        const BitBoard mask = sq1.Diagonal() == sq2.Diagonal() ? Diagonal(sq1.Diagonal()) :
                                sq1.AntiDiagonal() == sq2.AntiDiagonal() ? AntiDiagonal(sq1.AntiDiagonal()) :
                                sq1.File() == sq2.File() ? File(sq1.File()) :
                                sq1.Rank() == sq2.Rank() ? Rank(sq1.Rank()) : Empty;

                        // Check if the between BitBoard will be Empty. This will be the case if
                        // there is no ray linking the squares together, or the squares are equal.
                        if (mask.Empty() || sq1 == sq2) {
                            // Set the between BitBoard to Empty and continue.
                            between[square1][square2] = Empty;
                            continue;
                        }

                        const BitBoard blockers = BitBoard(sq1) + BitBoard(sq2);

                        // This step generates the between BitBoard for the current pair of Squares.
                        // We use the mask generated in the previous step to apply the Hyperbola
                        // algorithm along the rays joining the Squares together, using the two
                        // Squares as the blocker set.
                        //
                        // Blockers: . x . . . x . .
                        // Square 1: x x x x x x . .
                        // Square 2: . x x x x x x x
                        //
                        // The intersection between the two blocked rays will be the between
                        // BitBoard + Squares 1 and 2. Therefore, to get the between BitBoard, a
                        // final intersection operator with the union of Squares 1 and 2 is done.
                        between[square1][square2] = BitBoard::Hyperbola(sq1, blockers, mask) &
                                                    BitBoard::Hyperbola(sq2, blockers, mask);
                        between[square1][square2] -= blockers;
                    }
                }

                return between;
            }();

            // Query the between table with the given Squares and return the result.
            return between[static_cast<uint8_t>(square1)][static_cast<uint8_t>(square2)];
        }

        // Between returns a BitBoard containing all the squares between the two provided
        // squares, inclusive of the first square only.
        [[maybe_unused]] constexpr inline BitBoard Between1(Square square1, Square square2) {
            return Between(square1, square2) + BitBoard(square1);
        }

        // Between returns a BitBoard containing all the squares between the two provided
        // squares, inclusive of the second square only.
        [[maybe_unused]] constexpr inline BitBoard Between2(Square square1, Square square2) {
            return Between(square1, square2) + BitBoard(square2);
        }

        // Between returns a BitBoard containing all the squares between the two provided
        // squares, inclusive of both the squares.
        [[maybe_unused]] constexpr inline BitBoard Between12(Square square1, Square square2) {
            return Between(square1, square2) + BitBoard(square1) + BitBoard(square2);
        }
    }
}

#endif // CHESS_BITBOARD