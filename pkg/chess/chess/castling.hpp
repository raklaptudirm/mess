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

#ifndef CHESS_CASTLING
#define CHESS_CASTLING

#include <string>
#include <cassert>
#include <cstdint>
#include <iostream>

#include "square.hpp"
#include "bitboard.hpp"

namespace Chess::Castling {
    struct Side {
        /* ******************************
         * Internal Enum Representation *
         ****************************** */

        // Number of castling sides on a chessboard.
        static const int N = 2;

        // The internal enum representation of a Side.
        enum internal_type : uint8_t { H, A };

        // The variable that stores the internal representation.
        internal_type internal;

        /***************************
         * Constructor Definitions *
         ***************************/

        // Constructor to convert an uint8_t into a Side.
        // The Side with the given uint8_t representation.
        constexpr explicit Side(uint8_t side) : internal(static_cast<internal_type>(side)) {}

        // Constructor to convert an internal representation into a Side.
        // The Side with the given internal representation.
        // NOLINTNEXTLINE
        [[maybe_unused]] constexpr Side(internal_type side) : internal(side) {}

        /************************
         * Conversion Functions *
         ************************/

        // Conversion function to convert a side into its uint8_t representation.
        constexpr inline explicit operator uint8_t() const {
            return static_cast<uint8_t>(internal);
        }

        /************************
         * Operator Definitions *
         ************************/

        // Equality operator to check if two Sides are equal.
        constexpr inline bool operator ==(const Side&) const = default;
    };

    // A Dimension represents a Color-Side pair, each of which uniquely
    // represents one "way" that castling is possible on a chessboard.
    struct Dimension {
    private:
        // Internal representation of a dimension.
        uint8_t internal;

    public:
        // Number of castling Dimension (2 Sides x 2 Colors).
        static constexpr int N = Side::N * Color::N;

        constexpr inline Dimension(Color color, Side side)
                : internal(static_cast<uint8_t>(color) * Side::N + static_cast<uint8_t>(side)) {}

        // Color returns the color of this Dimension.
        [[nodiscard]] constexpr inline Color Color() const {
            return Chess::Color(internal / Chess::Color::N);
        }

        // Side returns the Side of this Dimension.
        [[nodiscard]] constexpr inline Side Side() const {
            return Castling::Side(internal % Castling::Side::N);
        }

        constexpr inline explicit operator uint8_t() const {
            return internal;
        }
    };

    // List of all the possible castling Dimensions.
    namespace Dimensions {
        constexpr static auto WhiteH = Dimension(Color::White, Side::H);
        constexpr static auto WhiteA = Dimension(Color::White, Side::A);
        constexpr static auto BlackH = Dimension(Color::Black, Side::H);
        constexpr static auto BlackA = Dimension(Color::Black, Side::A);
    }

    // Ends name spaces all the Ranks, Files, and Squares relevant
    // to the end squares of a King and a Rook after castling.
    namespace Ends {
        // End Files of Kings and Rooks for each Side.
        constexpr File KingFileH = File::G;
        constexpr File RookFileH = File::F;
        constexpr File KingFileA = File::C;
        constexpr File RookFileA = File::D;

        // End Ranks of Kings and Rooks for each Color.
        constexpr Rank WhiteRank = Rank::First;
        constexpr Rank BlackRank = Rank::Eighth;

        // End squares of White Kings and Rooks.
        constexpr Square WhiteKingH = Square(KingFileH, WhiteRank);
        constexpr Square WhiteRookH = Square(RookFileH, WhiteRank);
        constexpr Square WhiteKingA = Square(KingFileA, WhiteRank);
        constexpr Square WhiteRookA = Square(RookFileA, WhiteRank);

        // End squares of Black Kings and Rooks.
        constexpr Square BlackKingH = Square(KingFileH, BlackRank);
        constexpr Square BlackRookH = Square(RookFileH, BlackRank);
        constexpr Square BlackKingA = Square(KingFileA, BlackRank);
        constexpr Square BlackRookA = Square(RookFileA, BlackRank);
    }

    // EndSquares returns a pair containing the end Squares of a King
    // and Rook respectively which are castling in the given dimension.
    constexpr inline std::pair<Square, Square> EndSquares(const Dimension dimension) {
        switch (static_cast<uint8_t>(dimension)) {
            case static_cast<uint8_t>(Dimensions::WhiteH): return {Ends::WhiteKingH, Ends::WhiteRookH};
            case static_cast<uint8_t>(Dimensions::WhiteA): return {Ends::WhiteKingA, Ends::WhiteRookA};
            case static_cast<uint8_t>(Dimensions::BlackH): return {Ends::BlackKingH, Ends::BlackRookH};
            case static_cast<uint8_t>(Dimensions::BlackA): return {Ends::BlackKingA, Ends::BlackRookA};

            default: return {};
        }
    }

    // Rights represents a set of the four different Dimensions of castling.
    class Rights {
        // Internal representation of the set.
        uint8_t internal = 0;

    public:
        constexpr inline Rights() = default;
        constexpr inline explicit Rights(uint8_t rights) : internal(rights) {}

        constexpr inline explicit Rights(Dimension dimension)
            : internal(1 << static_cast<uint8_t>(dimension)) {}

        // Has checks if the given Rights is a subset of the target.
        [[nodiscard]] constexpr inline bool Has(const Rights subset) const {
            const auto ss = static_cast<uint8_t>(subset);
            return (internal & ss) == ss;
        }

        [[nodiscard]] constexpr inline bool Has(const Dimension dim) const {
            return internal & static_cast<uint8_t>(Rights(dim));
        }

        constexpr inline explicit operator uint8_t() const {
            return internal;
        }

        constexpr inline bool operator ==(const Rights&) const = default;

        constexpr inline Rights operator +(const Rights rhs) const {
            return Rights(internal | static_cast<uint8_t>(rhs));
        }

        constexpr inline void operator +=(const Rights rhs) {
            internal |= static_cast<uint8_t>(rhs);
        }

        constexpr inline Rights operator -(const Rights rhs) const {
            return Rights(internal &~ static_cast<uint8_t>(rhs));
        }

        constexpr inline void operator -=(const Rights rhs) {
            internal = internal &~ static_cast<uint8_t>(rhs);
        }

        constexpr inline Rights operator ~() const {
            return Rights(~internal);
        }

        constexpr inline Rights operator &(const Rights rhs) const {
            return Rights(internal & static_cast<uint8_t>(rhs));
        }

        [[nodiscard]] constexpr inline std::string ToString() const {
            std::string str;

            if (Has(Dimensions::WhiteH)) str += "K";
            if (Has(Dimensions::WhiteA)) str += "Q";
            if (Has(Dimensions::BlackH)) str += "k";
            if (Has(Dimensions::BlackA)) str += "q";

            return str;
        }
    };

    constexpr static Rights WhiteH = Rights(Dimensions::WhiteH);
    constexpr static Rights WhiteA = Rights(Dimensions::WhiteA);
    constexpr static Rights BlackH = Rights(Dimensions::BlackH);
    constexpr static Rights BlackA = Rights(Dimensions::BlackA);

    constexpr static Rights White = WhiteH + WhiteA;
    constexpr static Rights Black = BlackH + BlackA;

    constexpr static Rights All  = White + Black;
    constexpr static Rights None = Rights(0ull);

    // Info contains all the castling metadata required to be able to determine
    // castling legality and the correct castling move in both Standard and FRC.
    class Info {
    private:
        bool chess960 = false;

        // Positions of the rooks.
        std::array<Square,   4> rooks = {};

        // Castling Legality checks.
        std::array<BitBoard, 4> blockerMask = {}; // Squares which need to be empty for target Dimension.
        std::array<BitBoard, 4> attacksMask = {}; // Squares which need to be safe  for target Dimension.

        // List of Rights to remove for Moves to and from each Square.
        // This ensures that the castling rights are updated when the
        // King moves or a Rook moves/is captured.
        std::array<Castling::Rights, Square::N> masks = {};

    public:
        constexpr inline Info() = default;

        // Parse parses the given castling rights string with the additional context
        // of the position of both the Kings, and returns a parsed Info and Rights.
        constexpr static inline std::pair<Info, Rights> Parse(std::string str, Square WhiteKing, Square BlackKing) {
            // - is the empty set of Rights.
            if (str == "-") {
                // Positions of rooks and whether we are playing FRC chess is
                // ambiguous/inconsequential and Standard chess is assumed.
                return {
                    Info(
                        Square::E1, File::H, File::A,
                        Square::E8, File::H, File::A,
                        false
                    ), None
                };
            }

            // Basic checks on the rights string.
            assert(0 < str.length() && str.length() <= 4);

            // FRC uses Shredder fen which uses a-h instead of k/q.
            const char id = str.at(0);
            const auto chess960 = id != 'K' && id != 'Q' && id != 'k' && id != 'q';

            Castling::Rights rights = None;

            // Default to Standard chess rook files.
            File whiteH = File::H;
            File whiteA = File::A;
            File blackH = File::H;
            File blackA = File::A;

            for (const auto right : str) {
                if (chess960) {
                    if ('a' <= right && right <= 'h') {
                        const auto file = File(std::string(1, right));
                        if (file > BlackKing.File()) {
                            blackH = file;
                            rights += Castling::BlackH;
                        } else {
                            blackA = file;
                            rights += Castling::BlackA;
                        }
                    } else {
                        const auto file = File(std::string(1, right + 'a' - 'A'));
                        if (file > WhiteKing.File()) {
                            whiteH = file;
                            rights += Castling::WhiteH;
                        } else {
                            whiteA = file;
                            rights += Castling::WhiteA;
                        }
                    }
                } else {
                    switch (right) {
                        case 'K': rights += Castling::WhiteH; break;
                        case 'Q': rights += Castling::WhiteA; break;
                        case 'k': rights += Castling::BlackH; break;
                        case 'q': rights += Castling::BlackA; break;

                        default: assert(false);
                    }
                }
            }

            return {
                Info(
                    WhiteKing, whiteH, whiteA,
                    BlackKing, blackH, blackA,
                    chess960
                ), rights
            };
        }

        constexpr inline Info(
                Square whiteKing, File whiteRookHFile, File whiteRookAFile,
                Square blackKing, File blackRookHFile, File blackRookAFile,
                bool isChess960
        ) {
            chess960 = isChess960;

            // Convert the Rook Files to Squares.
            const Square whiteRookH = Square(whiteRookHFile, Rank::First);
            const Square whiteRookA = Square(whiteRookAFile, Rank::First);
            const Square blackRookH = Square(blackRookHFile, Rank::Eighth);
            const Square blackRookA = Square(blackRookAFile, Rank::Eighth);

            // Populate the rooks field of Info.
            rooks[static_cast<uint8_t>(Dimensions::WhiteH)] = whiteRookH;
            rooks[static_cast<uint8_t>(Dimensions::WhiteA)] = whiteRookA;
            rooks[static_cast<uint8_t>(Dimensions::BlackH)] = blackRookH;
            rooks[static_cast<uint8_t>(Dimensions::BlackA)] = blackRookA;

            // Blocker mask extracts the square which need to be empty in order for castling to be legal.
            // The king's path to its end square and the rooks path to its end square should be empty except for the
            // castling king and rook themselves. Therefore, the blocker mask is (kingPath + rookPath) - (king + rook).
            auto BLOCKER_MASK = [](Square king, Square rook, Square kingEnd, Square rookEnd) {
                return (BitBoards::Between2(king,  kingEnd) + BitBoards::Between2(rook, rookEnd)) - (BitBoard(king) + BitBoard(rook));
            };

            // Populate the blockerMask field of Info.
            blockerMask[static_cast<uint8_t>(Dimensions::WhiteH)] = BLOCKER_MASK(whiteKing, whiteRookH, Ends::WhiteKingH, Ends::WhiteRookH);
            blockerMask[static_cast<uint8_t>(Dimensions::WhiteA)] = BLOCKER_MASK(whiteKing, whiteRookA, Ends::WhiteKingA, Ends::WhiteRookA);
            blockerMask[static_cast<uint8_t>(Dimensions::BlackH)] = BLOCKER_MASK(blackKing, blackRookH, Ends::BlackKingH, Ends::BlackRookH);
            blockerMask[static_cast<uint8_t>(Dimensions::BlackA)] = BLOCKER_MASK(blackKing, blackRookA, Ends::BlackKingA, Ends::BlackRookA);

            // Populate the attacksMask field of Info.
            // Attack masks are the squares between the castling King and its end square both inclusive.
            // However, whether the king is in check is checks differently so only chess the end Square.
            attacksMask[static_cast<uint8_t>(Dimensions::WhiteH)] = BitBoards::Between2(whiteKing,  Ends::WhiteKingH);
            attacksMask[static_cast<uint8_t>(Dimensions::WhiteA)] = BitBoards::Between2(whiteKing,  Ends::WhiteKingA);
            attacksMask[static_cast<uint8_t>(Dimensions::BlackH)] = BitBoards::Between2(blackKing,  Ends::BlackKingH);
            attacksMask[static_cast<uint8_t>(Dimensions::BlackA)] = BitBoards::Between2(blackKing,  Ends::BlackKingA);

            // Zero out the entire masks array.
            masks = {};

            // Moves to and from the Rook's position imply the Rook
            // has moved or been captured so remove those Rights.
            masks[static_cast<uint8_t>(whiteRookH)] = WhiteH;
            masks[static_cast<uint8_t>(whiteRookA)] = WhiteA;
            masks[static_cast<uint8_t>(blackRookH)] = BlackH;
            masks[static_cast<uint8_t>(blackRookA)] = BlackA;

            // Moves from the Kings position imply that the King
            // has moved so remove all the Rights for that Color.
            masks[static_cast<uint8_t>(whiteKing)] = White;
            masks[static_cast<uint8_t>(blackKing)] = Black;
        }

        // Mask returns the relevant Rights mask for the given Square.
        [[nodiscard]] constexpr inline Castling::Rights Mask(const Square sq) const {
            return masks[static_cast<uint8_t>(sq)];
        }

        // Rook returns the position of the rook for the given Dimension.
        [[nodiscard]] constexpr inline Square Rook(Dimension dimension) const {
            return rooks[static_cast<uint8_t>(dimension)];
        }

        // BlockerMask returns the blocker mask for the given Dimension, which
        // are the Squares which need to be empty for castling to be legal.
        [[nodiscard]] constexpr inline BitBoard BlockerMask(Dimension dimension) const {
            return blockerMask[static_cast<uint8_t>(dimension)];
        }

        // AttackMask returns the attacks mask for the given Dimension, which
        // are the Squares which need to be safe for castling to be legal.
        [[nodiscard]] constexpr inline BitBoard AttackMask(Dimension dimension) const {
            return attacksMask[static_cast<uint8_t>(dimension)];
        }
    };

    constexpr inline std::ostream& operator<<(std::ostream& os, const Rights& rights) {
        os << rights.ToString();
        return os;
    }
}

#endif