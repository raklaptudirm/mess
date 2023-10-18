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
#include <iostream>

#include "../util/types.hpp"

#include "square.hpp"
#include "bitboard.hpp"

namespace Chess::Castling {
    struct Side {
        /* ******************************
         * Internal Enum Representation *
         ****************************** */

        // Number of Files, excluding None, on a Chessboard.
        static const int N = 2;

        // The internal enum representation of a Side.
        enum internal_type : uint8 { H, A };

        // The variable that stores the internal representation.
        internal_type internal;

        /***************************
         * Constructor Definitions *
         ***************************/

        // Constructor to convert an uint8 into a Side.
        // The Side with the given uint8 representation.
        constexpr explicit Side(uint8 side) : internal(static_cast<internal_type>(side)) {}

        // Constructor to convert an internal representation into a Side.
        // The Side with the given internal representation.
        constexpr Side(internal_type side) : internal(side) {}

        /************************
         * Conversion Functions *
         ************************/

        // Conversion function to convert a side into its uint8 representation.
        constexpr inline explicit operator uint8() const {
            return static_cast<uint8>(internal);
        }

        /************************
         * Operator Definitions *
         ************************/

        constexpr inline bool operator ==(const Side&) const = default;
    };

    struct Dimension {
        uint8 internal;

    public:
        constexpr inline Dimension(Color color, Side side)
                : internal(static_cast<uint8>(color) * Side::N + static_cast<uint8>(side)) {}

        [[nodiscard]] constexpr inline Color Color() const {
            return Chess::Color(internal / Chess::Color::N);
        }

        [[nodiscard]] constexpr inline Side Side() const {
            return Castling::Side(internal % Castling::Side::N);
        }

        constexpr inline explicit operator uint8() const {
            return internal;
        }
    };

    namespace Dimensions {
        constexpr static auto WhiteH = Dimension(Color::White, Side::H);
        constexpr static auto WhiteA = Dimension(Color::White, Side::A);
        constexpr static auto BlackH = Dimension(Color::Black, Side::H);
        constexpr static auto BlackA = Dimension(Color::Black, Side::A);
    }

    namespace Ends {
        constexpr File KingFileH = File::G;
        constexpr File RookFileH = File::F;
        constexpr File KingFileA = File::C;
        constexpr File RookFileA = File::D;

        constexpr Rank WhiteRank = Rank::First;
        constexpr Rank BlackRank = Rank::Eighth;

        constexpr Square WhiteKingH = Square(KingFileH, WhiteRank);
        constexpr Square WhiteRookH = Square(RookFileH, WhiteRank);
        constexpr Square WhiteKingA = Square(KingFileA, WhiteRank);
        constexpr Square WhiteRookA = Square(RookFileA, WhiteRank);

        constexpr Square BlackKingH = Square(KingFileH, BlackRank);
        constexpr Square BlackRookH = Square(RookFileH, BlackRank);
        constexpr Square BlackKingA = Square(KingFileA, BlackRank);
        constexpr Square BlackRookA = Square(RookFileA, BlackRank);
    }

    constexpr inline std::pair<Square, Square> EndSquares(const Dimension dimension) {
        switch (static_cast<uint8>(dimension)) {
            case static_cast<uint8>(Dimensions::WhiteH): return {Ends::WhiteKingH, Ends::WhiteRookH};
            case static_cast<uint8>(Dimensions::WhiteA): return {Ends::WhiteKingA, Ends::WhiteRookA};
            case static_cast<uint8>(Dimensions::BlackH): return {Ends::BlackKingH, Ends::BlackRookH};
            case static_cast<uint8>(Dimensions::BlackA): return {Ends::BlackKingA, Ends::BlackRookA};

            default: return {};
        }
    }

    class Rights {
        uint8 internal = 0;

    public:
        constexpr inline Rights() = default;
        constexpr inline explicit Rights(uint8 rights) : internal(rights) {}

        constexpr inline explicit Rights(Dimension dimension)
            : internal(1 << static_cast<uint8>(dimension)) {}

        [[nodiscard]] constexpr inline bool Has(const Rights subset) const {
            const auto ss = static_cast<uint8>(subset);
            return (internal & ss) == ss;
        }

        constexpr inline explicit operator uint8() const {
            return internal;
        }

        constexpr inline bool operator ==(const Rights&) const = default;

        constexpr inline Rights operator +(const Rights rhs) const {
            return Rights(internal | static_cast<uint8>(rhs));
        }

        constexpr inline void operator +=(const Rights rhs) {
            internal |= static_cast<uint8>(rhs);
        }

        constexpr inline Rights operator -(const Rights rhs) const {
            return Rights(internal &~ static_cast<uint8>(rhs));
        }

        constexpr inline void operator -=(const Rights rhs) {
            internal = internal &~ static_cast<uint8>(rhs);
        }

        constexpr inline Rights operator ~() const {
            return Rights(~internal);
        }

        [[nodiscard]] constexpr inline std::string ToString() const {
            std::string str;

            if (Has(Rights(Dimensions::WhiteH))) str += "K";
            if (Has(Rights(Dimensions::WhiteA))) str += "Q";
            if (Has(Rights(Dimensions::BlackH))) str += "k";
            if (Has(Rights(Dimensions::BlackA))) str += "q";

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

    class Info {
    private:
        bool chess960 = false;
        std::array<Square,   4> rooks = {};
        std::array<BitBoard, 4> blockerMask = {};
        std::array<BitBoard, 4> attacksMask = {};
        std::array<Castling::Rights, Square::N> masks = {};

    public:
        constexpr inline Info() = default;

        constexpr static inline std::pair<Info, Rights> Parse(std::string str, Square WhiteKing, Square BlackKing) {
            if (str == "-") {
                return {
                    Info(
                        Square::E1, File::H, File::A,
                        Square::E8, File::H, File::A,
                        false
                    ), None
                };
            }

            assert(0 < str.length() && str.length() <= 4);

            const char id = str.at(0);
            const auto chess960 = id != 'K' && id != 'Q' && id != 'k' && id != 'q';

            Castling::Rights rights = None;

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

            const Square whiteRookH = Square(whiteRookHFile, Rank::First);
            const Square whiteRookA = Square(whiteRookAFile, Rank::First);
            const Square blackRookH = Square(blackRookHFile, Rank::Eighth);
            const Square blackRookA = Square(blackRookAFile, Rank::Eighth);

            rooks[static_cast<uint8>(Dimensions::WhiteH)] = whiteRookH;
            rooks[static_cast<uint8>(Dimensions::WhiteA)] = whiteRookA;
            rooks[static_cast<uint8>(Dimensions::BlackH)] = blackRookH;
            rooks[static_cast<uint8>(Dimensions::BlackA)] = blackRookA;

            // Blocker mask extracts the square which need to be empty in order for castling to be legal.
            // The king's path to its end square and the rooks path to its end square should be empty except for the
            // castling king and rook themselves. Therefore, the blocker mask is (kingPath + rookPath) - (king + rook).
            #define BLOCKER_MASK(king, rook, kingEnd, rookEnd) \
            ((BitBoards::Between2(king,  kingEnd) + BitBoards::Between2(rook, rookEnd)) - (BitBoard(king) + BitBoard(rook)))

            blockerMask[static_cast<uint8>(Dimensions::WhiteH)] = BLOCKER_MASK(whiteKing, whiteRookH, Ends::WhiteKingH, Ends::WhiteRookH);
            blockerMask[static_cast<uint8>(Dimensions::WhiteA)] = BLOCKER_MASK(whiteKing, whiteRookA, Ends::WhiteKingA, Ends::WhiteRookA);
            blockerMask[static_cast<uint8>(Dimensions::BlackH)] = BLOCKER_MASK(blackKing, blackRookH, Ends::BlackKingH, Ends::BlackRookH);
            blockerMask[static_cast<uint8>(Dimensions::BlackA)] = BLOCKER_MASK(blackKing, blackRookA, Ends::BlackKingA, Ends::BlackRookA);

            #undef BLOCKER_MASK

            attacksMask[static_cast<uint8>(Dimensions::WhiteH)] = BitBoards::Between2(whiteKing,  Ends::WhiteKingH);
            attacksMask[static_cast<uint8>(Dimensions::WhiteA)] = BitBoards::Between2(whiteKing,  Ends::WhiteKingA);
            attacksMask[static_cast<uint8>(Dimensions::BlackH)] = BitBoards::Between2(blackKing,  Ends::BlackKingH);
            attacksMask[static_cast<uint8>(Dimensions::BlackA)] = BitBoards::Between2(blackKing,  Ends::BlackKingA);

            masks = {};
            masks[static_cast<uint8>(whiteRookH)] = WhiteH;
            masks[static_cast<uint8>(whiteRookA)] = WhiteA;
            masks[static_cast<uint8>(blackRookH)] = BlackH;
            masks[static_cast<uint8>(blackRookA)] = BlackA;

            masks[static_cast<uint8>(whiteKing)] = White;
            masks[static_cast<uint8>(blackKing)] = Black;
        }

        [[nodiscard]] constexpr inline Castling::Rights Mask(const Square sq) const {
            return masks[static_cast<uint8>(sq)];
        }

        [[nodiscard]] constexpr inline Square Rook(Dimension dimension) const {
            return rooks[static_cast<uint8>(dimension)];
        }

        [[nodiscard]] constexpr inline BitBoard BlockerMask(Dimension dimension) const {
            return blockerMask[static_cast<uint8>(dimension)];
        }

        [[nodiscard]] constexpr inline BitBoard AttackMask(Dimension dimension) const {
            return attacksMask[static_cast<uint8>(dimension)];
        }
    };

    constexpr inline std::ostream& operator<<(std::ostream& os, const Rights& rights) {
        os << rights.ToString();
        return os;
    }
}

#endif