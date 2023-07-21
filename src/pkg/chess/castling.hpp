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
#include <assert.h>
#include <iostream>

#include "../util/types.hpp"

#include "square.hpp"
#include "bitboard.hpp"

namespace Chess::Castling {
        struct EndSquares {
            Square King, Rook;

            constexpr inline EndSquares(Square king, Square rook) : King(king), Rook(rook) {}
        };

        template<Color COLOR>
        constexpr inline EndSquares HEndSquares() {
            if (COLOR == Color::White)
                return {Square::G1, Square::F1};
            return {Square::G8, Square::F8};
        }

        constexpr inline EndSquares HEndSquares(const Color COLOR) {
            if (COLOR == Color::White)
                return {Square::G1, Square::F1};
            return {Square::G8, Square::F8};
        }

        template<Color COLOR>
        constexpr inline EndSquares AEndSquares() {
            if (COLOR == Color::White)
                return {Square::C1, Square::D1};
            return {Square::C8, Square::D8};
        }

        constexpr inline EndSquares AEndSquares(const Color COLOR) {
            if (COLOR == Color::White)
                return {Square::C1, Square::D1};
            return {Square::C8, Square::D8};
        }

        class Rights {
            uint8 internal = 0;

            public:
                constexpr inline Rights() = default;
                constexpr inline explicit Rights(uint8 rights) : internal(rights) {}

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

                    if (Has(Rights(1 << 0))) str += "K";
                    if (Has(Rights(1 << 1))) str += "Q";
                    if (Has(Rights(1 << 2))) str += "k";
                    if (Has(Rights(1 << 3))) str += "q";

                    return str;
                }
        };

        namespace {
            template<Color COLOR>
            constexpr inline uint8 offsetH() {
                return static_cast<uint8>(COLOR);
            }

            template<Color COLOR>
            constexpr inline uint8 offsetA() {
                return static_cast<uint8>(COLOR) + 2;
            }
        }

        constexpr static Rights WhiteH = Rights(1ull << offsetH<Color::White>());
        constexpr static Rights WhiteA = Rights(1ull << offsetA<Color::White>());
        constexpr static Rights BlackH = Rights(1ull << offsetH<Color::Black>());
        constexpr static Rights BlackA = Rights(1ull << offsetA<Color::Black>());

        constexpr static Rights White = WhiteH + WhiteA;
        constexpr static Rights Black = BlackH + BlackA;

        constexpr static Rights All  = White + Black;
        constexpr static Rights None = Rights(0ull);

        template<Color COLOR>
        constexpr inline Rights H() {
            if (COLOR == Color::White)
                return WhiteH;
            return BlackH;
        }

        template<Color COLOR>
        constexpr inline Rights A() {
            if (COLOR == Color::White)
                return WhiteA;
            return BlackA;
        }

        class Info {
            public:
                Castling::Rights Rights = None;

            private:
                bool chess960 = false;
                std::array<Square,   4> rooks = {};
                std::array<BitBoard, 4> paths = {};
                std::array<Castling::Rights, Square::N> masks = {};

            public:
                constexpr inline Info() = default;

                constexpr inline Info(std::string str, Square WhiteKing, Square BlackKing) {
                    if (str == "-") {
                        chess960 = false;
                        *this = Info(
                            None,
                            Square::E1, File::H, File::A,
                            Square::E8, File::H, File::A
                        );
                        return;
                    }

                    assert(0 < str.length() && str.length() <= 4);

                    const char id = str.at(0);
                    chess960 = id != 'K' && id != 'Q' && id != 'k' && id != 'q';

                    Castling::Rights rights = None;

                    File whiteH = File::H;
                    File whiteA = File::A;
                    File blackH = File::H;
                    File blackA = File::A;

                    for (const auto right : str) {
                        if (chess960) {
                            if ('a' <= right && right <= 'h') {
                                const auto file = File(std::string(1, right + 'a' - 'A'));
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

                    //std::cout << "black H file:" << (uint64)static_cast<uint8>(blackH) << std::endl;

                    *this = Info(
                            rights,
                            WhiteKing, whiteH, whiteA,
                            BlackKing, blackH, blackA
                    );
                }

                constexpr inline Info(
                    Castling::Rights rights,
                    Square whiteKing, File whiteRookHFile, File whiteRookAFile,
                    Square blackKing, File blackRookHFile, File blackRookAFile
                ) {
                    Rights = rights;

                    const Square whiteRooKH = Square(whiteRookHFile, Rank::First);
                    const Square whiteRookA = Square(whiteRookAFile, Rank::First);
                    const Square blackRookH = Square(blackRookHFile, Rank::Eighth);
                    const Square blackRookA = Square(blackRookAFile, Rank::Eighth);

                    rooks[offsetH<Color::White>()] = whiteRooKH;
                    rooks[offsetA<Color::White>()] = whiteRookA;
                    rooks[offsetH<Color::Black>()] = blackRookH;
                    rooks[offsetA<Color::Black>()] = blackRookA;

                    paths[offsetH<Color::White>()] = BitBoards::Between(whiteKing, Square::H1);
                    paths[offsetA<Color::White>()] = BitBoards::Between(whiteKing, Square::B1);
                    paths[offsetH<Color::Black>()] = BitBoards::Between(blackKing, Square::H8);
                    paths[offsetA<Color::Black>()] = BitBoards::Between(blackKing, Square::B8);

                    masks = {};
                    masks[static_cast<uint8>(whiteRooKH)] = WhiteH;
                    masks[static_cast<uint8>(whiteRookA)] = WhiteA;
                    masks[static_cast<uint8>(blackRookH)] = BlackH;
                    masks[static_cast<uint8>(blackRookA)] = BlackA;

                    masks[static_cast<uint8>(whiteKing)] = White;
                    masks[static_cast<uint8>(blackKing)] = Black;
                }

                [[nodiscard]] constexpr inline Castling::Rights Mask(const Square sq) const {
                    return masks[static_cast<uint8>(sq)];
                }

                template <Color STM>
                [[nodiscard]] constexpr inline Square RookH() const {
                    return rooks[offsetH<STM>()];
                }

                template <Color STM>
                [[nodiscard]] constexpr inline Square RookA() const {
                    return rooks[offsetA<STM>()];
                }

                template <Color STM>
                [[nodiscard]] constexpr inline BitBoard PathH() const {
                    return paths[offsetH<STM>()];
                }

                template <Color STM>
                [[nodiscard]] constexpr inline BitBoard PathA() const {
                    return paths[offsetA<STM>()];
                }

                [[nodiscard]] constexpr inline std::string ToString() const {
                    std::string str = "";

                    if (Rights.Has(Castling::WhiteH)) str += "K";
                    if (Rights.Has(Castling::WhiteA)) str += "Q";
                    if (Rights.Has(Castling::BlackH)) str += "k";
                    if (Rights.Has(Castling::BlackA)) str += "q";

                    return str;
                }
        };

        constexpr inline std::ostream& operator<<(std::ostream& os, const Info& rights) {
            os << rights.ToString();
            return os;
        }

        constexpr inline std::ostream& operator<<(std::ostream& os, const Rights& rights) {
            os << rights.ToString();
            return os;
        }
    }

#endif