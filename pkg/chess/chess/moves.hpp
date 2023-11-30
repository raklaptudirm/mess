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

#ifndef CHESS_MOVES
#define CHESS_MOVES

#include <cassert>
#include <array>
#include <cstdint>

#include "bitboard.hpp"
#include "color.hpp"
#include "piece.hpp"

namespace Chess::MoveTable {
    namespace {
        struct magic {
            uint64_t Relevant = 0;
            uint64_t   Number = 0;
            int32_t    Offset = 0;

            constexpr inline magic() = default;

            constexpr inline magic(uint64_t relevant, uint64_t magic, int offset) :
                    Relevant(relevant), Number(magic), Offset(offset) {}
        };

        namespace blackMagic {
            namespace {
                template<Piece piece>
                constexpr inline uint8_t pieceIndex() {
                    assert(piece == Piece::Bishop || piece == Piece::Rook);
                    return static_cast<uint8_t>(piece) - 2;
                }
                // magic hash function shifts for each piece type.
                constexpr std::array<uint8_t, 2> pieceShifts { 9, 12 };

                // Table containing magic constants for sliding pieces.
                // magic numbers and offsets are from Analog Hors's CozyChess
                // Library: https://github.com/analog-hors/cozy-chess
                constexpr std::array<std::array<magic, Square::N>, 2> magics = {{
                    // Bishop Magics.
                    {{
                             { 0xffbfdfeff7fbfdff, 0xa7020080601803d8, 60984 }, { 0xffffbfdfeff7fbff, 0x13802040400801f1, 66046 },
                             { 0xffffffbfdfeff5ff, 0x0a0080181001f60c, 32910 }, { 0xffffffffbfddebff, 0x1840802004238008, 16369 },
                             { 0xfffffffffdbbd7ff, 0xc03fe00100000000, 42115 }, { 0xfffffffdfbf7afff, 0x24c00bffff400000,   835 },
                             { 0xfffffdfbf7efdfff, 0x0808101f40007f04, 18910 }, { 0xfffdfbf7efdfbfff, 0x100808201ec00080, 25911 },
                             { 0xffdfeff7fbfdffff, 0xffa2feffbfefb7ff, 63301 }, { 0xffbfdfeff7fbffff, 0x083e3ee040080801, 16063 },
                             { 0xffffbfdfeff5ffff, 0xc0800080181001f8, 17481 }, { 0xffffffbfddebffff, 0x0440007fe0031000, 59361 },
                             { 0xfffffffdbbd7ffff, 0x2010007ffc000000, 18735 }, { 0xfffffdfbf7afffff, 0x1079ffe000ff8000, 61249 },
                             { 0xfffdfbf7efdfffff, 0x3c0708101f400080, 68938 }, { 0xfffbf7efdfbfffff, 0x080614080fa00040, 61791 },
                             { 0xffeff7fbfdfffdff, 0x7ffe7fff817fcff9, 21893 }, { 0xffdfeff7fbfffbff, 0x7ffebfffa01027fd, 62068 },
                             { 0xffbfdfeff5fff5ff, 0x53018080c00f4001, 19829 }, { 0xffffbfddebffebff, 0x407e0001000ffb8a, 26091 },
                             { 0xfffffdbbd7ffd7ff, 0x201fe000fff80010, 15815 }, { 0xfffdfbf7afffafff, 0xffdfefffde39ffef, 16419 },
                             { 0xfffbf7efdfffdfff, 0xcc8808000fbf8002, 59777 }, { 0xfff7efdfbfffbfff, 0x7ff7fbfff8203fff, 16288 },
                             { 0xfff7fbfdfffdfbff, 0x8800013e8300c030, 33235 }, { 0xffeff7fbfffbf7ff, 0x0420009701806018, 15459 },
                             { 0xffdfeff5fff5efff, 0x7ffeff7f7f01f7fd, 15863 }, { 0xffbfddebffebddff, 0x8700303010c0c006, 75555 },
                             { 0xfffdbbd7ffd7bbff, 0xc800181810606000, 79445 }, { 0xfffbf7afffaff7ff, 0x20002038001c8010, 15917 },
                             { 0xfff7efdfffdfefff, 0x087ff038000fc001,  8512 }, { 0xffefdfbfffbfdfff, 0x00080c0c00083007, 73069 },
                             { 0xfffbfdfffdfbf7ff, 0x00000080fc82c040, 16078 }, { 0xfff7fbfffbf7efff, 0x000000407e416020, 19168 },
                             { 0xffeff5fff5efdfff, 0x00600203f8008020, 11056 }, { 0xffddebffebddbfff, 0xd003fefe04404080, 62544 },
                             { 0xffbbd7ffd7bbfdff, 0xa00020c018003088, 80477 }, { 0xfff7afffaff7fbff, 0x7fbffe700bffe800, 75049 },
                             { 0xffefdfffdfeff7ff, 0x107ff00fe4000f90, 32947 }, { 0xffdfbfffbfdfefff, 0x7f8fffcff1d007f8, 59172 },
                             { 0xfffdfffdfbf7efff, 0x0000004100f88080, 55845 }, { 0xfffbfffbf7efdfff, 0x00000020807c4040, 61806 },
                             { 0xfff5fff5efdfbfff, 0x00000041018700c0, 73601 }, { 0xffebffebddbfffff, 0x0010000080fc4080, 15546 },
                             { 0xffd7ffd7bbfdffff, 0x1000003c80180030, 45243 }, { 0xffafffaff7fbfdff, 0xc10000df80280050, 20333 },
                             { 0xffdfffdfeff7fbff, 0xffffffbfeff80fdc, 33402 }, { 0xffbfffbfdfeff7ff, 0x000000101003f812, 25917 },
                             { 0xfffffdfbf7efdfff, 0x0800001f40808200, 32875 }, { 0xfffffbf7efdfbfff, 0x084000101f3fd208,  4639 },
                             { 0xfffff5efdfbfffff, 0x080000000f808081, 17077 }, { 0xffffebddbfffffff, 0x0004000008003f80, 62324 },
                             { 0xffffd7bbfdffffff, 0x08000001001fe040, 18159 }, { 0xffffaff7fbfdffff, 0x72dd000040900a00, 61436 },
                             { 0xffffdfeff7fbfdff, 0xfffffeffbfeff81d, 57073 }, { 0xffffbfdfeff7fbff, 0xcd8000200febf209, 61025 },
                             { 0xfffdfbf7efdfbfff, 0x100000101ec10082, 81259 }, { 0xfffbf7efdfbfffff, 0x7fbaffffefe0c02f, 64083 },
                             { 0xfff5efdfbfffffff, 0x7f83fffffff07f7f, 56114 }, { 0xffebddbfffffffff, 0xfff1fffffff7ffc1, 57058 },
                             { 0xffd7bbfdffffffff, 0x0878040000ffe01f, 58912 }, { 0xffaff7fbfdffffff, 0x945e388000801012, 22194 },
                             { 0xffdfeff7fbfdffff, 0x0840800080200fda, 70880 }, { 0xffbfdfeff7fbfdff, 0x100000c05f582008, 11140 },
                     }},

                    // Rook Magics.
                    {{
                             { 0xfffefefefefefe81, 0x80280013ff84ffff, 10890 }, { 0xfffdfdfdfdfdfd83, 0x5ffbfefdfef67fff, 50579 },
                             { 0xfffbfbfbfbfbfb85, 0xffeffaffeffdffff, 62020 }, { 0xfff7f7f7f7f7f789, 0x003000900300008a, 67322 },
                             { 0xffefefefefefef91, 0x0050028010500023, 80251 }, { 0xffdfdfdfdfdfdfa1, 0x0020012120a00020, 58503 },
                             { 0xffbfbfbfbfbfbfc1, 0x0030006000c00030, 51175 }, { 0xff7f7f7f7f7f7f81, 0x0058005806b00002, 83130 },
                             { 0xfffefefefefe81ff, 0x7fbff7fbfbeafffc, 50430 }, { 0xfffdfdfdfdfd83ff, 0x0000140081050002, 21613 },
                             { 0xfffbfbfbfbfb85ff, 0x0000180043800048, 72625 }, { 0xfff7f7f7f7f789ff, 0x7fffe800021fffb8, 80755 },
                             { 0xffefefefefef91ff, 0xffffcffe7fcfffaf, 69753 }, { 0xffdfdfdfdfdfa1ff, 0x00001800c0180060, 26973 },
                             { 0xffbfbfbfbfbfc1ff, 0x4f8018005fd00018, 84972 }, { 0xff7f7f7f7f7f81ff, 0x0000180030620018, 31958 },
                             { 0xfffefefefe81feff, 0x00300018010c0003, 69272 }, { 0xfffdfdfdfd83fdff, 0x0003000c0085ffff, 48372 },
                             { 0xfffbfbfbfb85fbff, 0xfffdfff7fbfefff7, 65477 }, { 0xfff7f7f7f789f7ff, 0x7fc1ffdffc001fff, 43972 },
                             { 0xffefefefef91efff, 0xfffeffdffdffdfff, 57154 }, { 0xffdfdfdfdfa1dfff, 0x7c108007befff81f, 53521 },
                             { 0xffbfbfbfbfc1bfff, 0x20408007bfe00810, 30534 }, { 0xff7f7f7f7f817fff, 0x0400800558604100, 16548 },
                             { 0xfffefefe81fefeff, 0x0040200010080008, 46407 }, { 0xfffdfdfd83fdfdff, 0x0010020008040004, 11841 },
                             { 0xfffbfbfb85fbfbff, 0xfffdfefff7fbfff7, 21112 }, { 0xfff7f7f789f7f7ff, 0xfebf7dfff8fefff9, 44214 },
                             { 0xffefefef91efefff, 0xc00000ffe001ffe0, 57925 }, { 0xffdfdfdfa1dfdfff, 0x4af01f00078007c3, 29574 },
                             { 0xffbfbfbfc1bfbfff, 0xbffbfafffb683f7f, 17309 }, { 0xff7f7f7f817f7fff, 0x0807f67ffa102040, 40143 },
                             { 0xfffefe81fefefeff, 0x200008e800300030, 64659 }, { 0xfffdfd83fdfdfdff, 0x0000008780180018, 70469 },
                             { 0xfffbfb85fbfbfbff, 0x0000010300180018, 62917 }, { 0xfff7f789f7f7f7ff, 0x4000008180180018, 60997 },
                             { 0xffefef91efefefff, 0x008080310005fffa, 18554 }, { 0xffdfdfa1dfdfdfff, 0x4000188100060006, 14385 },
                             { 0xffbfbfc1bfbfbfff, 0xffffff7fffbfbfff,     0 }, { 0xff7f7f817f7f7fff, 0x0000802000200040, 38091 },
                             { 0xfffe81fefefefeff, 0x20000202ec002800, 25122 }, { 0xfffd83fdfdfdfdff, 0xfffff9ff7cfff3ff, 60083 },
                             { 0xfffb85fbfbfbfbff, 0x000000404b801800, 72209 }, { 0xfff789f7f7f7f7ff, 0x2000002fe03fd000, 67875 },
                             { 0xffef91efefefefff, 0xffffff6ffe7fcffd, 56290 }, { 0xffdfa1dfdfdfdfff, 0xbff7efffbfc00fff, 43807 },
                             { 0xffbfc1bfbfbfbfff, 0x000000100800a804, 73365 }, { 0xff7f817f7f7f7fff, 0x6054000a58005805, 76398 },
                             { 0xff81fefefefefeff, 0x0829000101150028, 20024 }, { 0xff83fdfdfdfdfdff, 0x00000085008a0014,  9513 },
                             { 0xff85fbfbfbfbfbff, 0x8000002b00408028, 24324 }, { 0xff89f7f7f7f7f7ff, 0x4000002040790028, 22996 },
                             { 0xff91efefefefefff, 0x7800002010288028, 23213 }, { 0xffa1dfdfdfdfdfff, 0x0000001800e08018, 56002 },
                             { 0xffc1bfbfbfbfbfff, 0xa3a80003f3a40048, 22809 }, { 0xff817f7f7f7f7fff, 0x2003d80000500028, 44545 },
                             { 0x81fefefefefefeff, 0xfffff37eefefdfbe, 36072 }, { 0x83fdfdfdfdfdfdff, 0x40000280090013c1,  4750 },
                             { 0x85fbfbfbfbfbfbff, 0xbf7ffeffbffaf71f,  6014 }, { 0x89f7f7f7f7f7f7ff, 0xfffdffff777b7d6e, 36054 },
                             { 0x91efefefefefefff, 0x48300007e8080c02, 78538 }, { 0xa1dfdfdfdfdfdfff, 0xafe0000fff780402, 28745 },
                             { 0xc1bfbfbfbfbfbfff, 0xee73fffbffbb77fe,  8555 }, { 0x817f7f7f7f7f7fff, 0x0002000308482882,  1009 },
                     }},
                }};
            }

            // Size of the hash table using this hash function.
            const int TableSize = 87988;

            using Table = std::array<BitBoard, blackMagic::TableSize>;

            template<Piece piece>
            constexpr inline magic GetMagic(const Square square) {
                return magics[pieceIndex<piece>()][static_cast<uint8_t>(square)];
            }

            template<Piece piece>
            constexpr inline uint32_t Index(const Square square, const BitBoard occupied) {
                constexpr uint8_t pieceShift = 64 - pieceShifts[pieceIndex<piece>()];  // Shift of given piece.

                // Get the relevant magic entry.
                const magic magic = GetMagic<piece>(square);

                // Mask off irrelevant blockers (outside or at end of rays).
                const uint64_t relevant = static_cast<uint64_t>(occupied) | magic.Relevant;

                // Compute hash and finally, the index.
                const uint64_t hash = relevant * magic.Number;
                return magic.Offset + static_cast<int32_t>(hash >> pieceShift);
            }
        }

        constexpr inline BitBoard bishopSlow(Square square, BitBoard blockers) {
            return BitBoard::Hyperbola(square, blockers, BitBoards::    Diagonal(square.    Diagonal())) |
                   BitBoard::Hyperbola(square, blockers, BitBoards::AntiDiagonal(square.AntiDiagonal()));
        }

        constexpr inline BitBoard rookSlow(Square square, BitBoard blockers) {
           return BitBoard::Hyperbola(square, blockers, BitBoards::File(square.File())) |
                  BitBoard::Hyperbola(square, blockers, BitBoards::Rank(square.Rank()));
        }

        constexpr std::array<std::array<uint64_t, Square::N>, Color::N> pawn {{
            {
                0x0000000000000200, 0x0000000000000500, 0x0000000000000a00, 0x0000000000001400,
                0x0000000000002800, 0x0000000000005000, 0x000000000000a000, 0x0000000000004000,
                0x0000000000020000, 0x0000000000050000, 0x00000000000a0000, 0x0000000000140000,
                0x0000000000280000, 0x0000000000500000, 0x0000000000a00000, 0x0000000000400000,
                0x0000000002000000, 0x0000000005000000, 0x000000000a000000, 0x0000000014000000,
                0x0000000028000000, 0x0000000050000000, 0x00000000a0000000, 0x0000000040000000,
                0x0000000200000000, 0x0000000500000000, 0x0000000a00000000, 0x0000001400000000,
                0x0000002800000000, 0x0000005000000000, 0x000000a000000000, 0x0000004000000000,
                0x0000020000000000, 0x0000050000000000, 0x00000a0000000000, 0x0000140000000000,
                0x0000280000000000, 0x0000500000000000, 0x0000a00000000000, 0x0000400000000000,
                0x0002000000000000, 0x0005000000000000, 0x000a000000000000, 0x0014000000000000,
                0x0028000000000000, 0x0050000000000000, 0x00a0000000000000, 0x0040000000000000,
                0x0200000000000000, 0x0500000000000000, 0x0a00000000000000, 0x1400000000000000,
                0x2800000000000000, 0x5000000000000000, 0xa000000000000000, 0x4000000000000000,
                0x0000000000000000, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000,
                0x0000000000000000, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000
            },
            {
                0x0000000000000000, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000,
                0x0000000000000000, 0x0000000000000000, 0x0000000000000000, 0x0000000000000000,
                0x0000000000000002, 0x0000000000000005, 0x000000000000000a, 0x0000000000000014,
                0x0000000000000028, 0x0000000000000050, 0x00000000000000a0, 0x0000000000000040,
                0x0000000000000200, 0x0000000000000500, 0x0000000000000a00, 0x0000000000001400,
                0x0000000000002800, 0x0000000000005000, 0x000000000000a000, 0x0000000000004000,
                0x0000000000020000, 0x0000000000050000, 0x00000000000a0000, 0x0000000000140000,
                0x0000000000280000, 0x0000000000500000, 0x0000000000a00000, 0x0000000000400000,
                0x0000000002000000, 0x0000000005000000, 0x000000000a000000, 0x0000000014000000,
                0x0000000028000000, 0x0000000050000000, 0x00000000a0000000, 0x0000000040000000,
                0x0000000200000000, 0x0000000500000000, 0x0000000a00000000, 0x0000001400000000,
                0x0000002800000000, 0x0000005000000000, 0x000000a000000000, 0x0000004000000000,
                0x0000020000000000, 0x0000050000000000, 0x00000a0000000000, 0x0000140000000000,
                0x0000280000000000, 0x0000500000000000, 0x0000a00000000000, 0x0000400000000000,
                0x0002000000000000, 0x0005000000000000, 0x000a000000000000, 0x0014000000000000,
                0x0028000000000000, 0x0050000000000000, 0x00a0000000000000, 0x0040000000000000
            }
        }};

        constexpr std::array<uint64_t, Square::N> knight {
            0x0000000000020400, 0x0000000000050800, 0x00000000000A1100, 0x0000000000142200,
            0x0000000000284400, 0x0000000000508800, 0x0000000000A01000, 0x0000000000402000,
            0x0000000002040004, 0x0000000005080008, 0x000000000A110011, 0x0000000014220022,
            0x0000000028440044, 0x0000000050880088, 0x00000000A0100010, 0x0000000040200020,
            0x0000000204000402, 0x0000000508000805, 0x0000000A1100110A, 0x0000001422002214,
            0x0000002844004428, 0x0000005088008850, 0x000000A0100010A0, 0x0000004020002040,
            0x0000020400040200, 0x0000050800080500, 0x00000A1100110A00, 0x0000142200221400,
            0x0000284400442800, 0x0000508800885000, 0x0000A0100010A000, 0x0000402000204000,
            0x0002040004020000, 0x0005080008050000, 0x000A1100110A0000, 0x0014220022140000,
            0x0028440044280000, 0x0050880088500000, 0x00A0100010A00000, 0x0040200020400000,
            0x0204000402000000, 0x0508000805000000, 0x0A1100110A000000, 0x1422002214000000,
            0x2844004428000000, 0x5088008850000000, 0xA0100010A0000000, 0x4020002040000000,
            0x0400040200000000, 0x0800080500000000, 0x1100110A00000000, 0x2200221400000000,
            0x4400442800000000, 0x8800885000000000, 0x100010A000000000, 0x2000204000000000,
            0x0004020000000000, 0x0008050000000000, 0x00110A0000000000, 0x0022140000000000,
            0x0044280000000000, 0x0088500000000000, 0x0010A00000000000, 0x0020400000000000
        };

        const blackMagic::Table sliding = []() {
            blackMagic::Table table = {};

            for (uint8_t sq = 0; sq < Square::N; sq++) {
                const Square square = Square(sq);

                const uint64_t bishopMask = ~blackMagic::GetMagic<Piece::Bishop>(square).Relevant;
                const uint64_t rookMask   = ~blackMagic::GetMagic<Piece::Rook  >(square).Relevant;

                // Initialize Bishop Moves in table.
                BitBoard blockers = BitBoards::Empty;
                do {
                    const auto index = blackMagic::Index<Piece::Bishop>(square, blockers); // Get magic Index
                    const auto moves = MoveTable::bishopSlow(square, blockers);

                    // Assert we are not overwriting anything: the table is valid.
                    assert(index < blackMagic::TableSize);
                    assert(table[index] == BitBoards::Empty || table[index] == moves);

                    // Add the moves BitBoard to the table.
                    table[index] = moves;

                    // Calculate the next subset of the Relevant Blocker Mask using
                    // the Carry-Rippler Trick to generate all subsets.
                    blockers = BitBoard((static_cast<uint64_t>(blockers) - bishopMask) & bishopMask);
                } while (blockers != BitBoards::Empty);

                // Initialize Rook Moves in table.
                blockers = BitBoards::Empty;
                do {
                    const auto index = blackMagic::Index<Piece::Rook>(square, blockers);
                    const auto moves = MoveTable::rookSlow(square, blockers);

                    // Assert we are not overwriting anything: the table is valid.
                    assert(index < blackMagic::TableSize);
                    assert(table[index] == BitBoards::Empty || table[index] == moves);

                    // Add the moves BitBoard to the table.
                    table[index] = moves;

                    // Calculate the next subset of the Relevant Blocker Mask using
                    // the Carry-Rippler Trick to generate all subsets.
                    blockers = BitBoard((static_cast<uint64_t>(blockers) - rookMask) & rookMask);
                } while (blockers != BitBoards::Empty);
            }

            return table;
        }();

        constexpr std::array<uint64_t, Square::N> king {
            0x0000000000000302, 0x0000000000000705, 0x0000000000000E0A, 0x0000000000001C14,
            0x0000000000003828, 0x0000000000007050, 0x000000000000E0A0, 0x000000000000C040,
            0x0000000000030203, 0x0000000000070507, 0x00000000000E0A0E, 0x00000000001C141C,
            0x0000000000382838, 0x0000000000705070, 0x0000000000E0A0E0, 0x0000000000C040C0,
            0x0000000003020300, 0x0000000007050700, 0x000000000E0A0E00, 0x000000001C141C00,
            0x0000000038283800, 0x0000000070507000, 0x00000000E0A0E000, 0x00000000C040C000,
            0x0000000302030000, 0x0000000705070000, 0x0000000E0A0E0000, 0x0000001C141C0000,
            0x0000003828380000, 0x0000007050700000, 0x000000E0A0E00000, 0x000000C040C00000,
            0x0000030203000000, 0x0000070507000000, 0x00000E0A0E000000, 0x00001C141C000000,
            0x0000382838000000, 0x0000705070000000, 0x0000E0A0E0000000, 0x0000C040C0000000,
            0x0003020300000000, 0x0007050700000000, 0x000E0A0E00000000, 0x001C141C00000000,
            0x0038283800000000, 0x0070507000000000, 0x00E0A0E000000000, 0x00C040C000000000,
            0x0302030000000000, 0x0705070000000000, 0x0E0A0E0000000000, 0x1C141C0000000000,
            0x3828380000000000, 0x7050700000000000, 0xE0A0E00000000000, 0xC040C00000000000,
            0x0203000000000000, 0x0507000000000000, 0x0A0E000000000000, 0x141C000000000000,
            0x2838000000000000, 0x5070000000000000, 0xA0E0000000000000, 0x40C0000000000000
        };
    }

    template<Color color>
    [[maybe_unused]] constexpr inline static BitBoard Pawn(Square square) {
        return BitBoard(pawn[static_cast<uint8_t>(color)][static_cast<uint8_t>(square)]);
    }

    [[maybe_unused]] constexpr inline static BitBoard Pawn(Color color, Square square) {
        return BitBoard(pawn[static_cast<uint8_t>(color)][static_cast<uint8_t>(square)]);
    }

    [[maybe_unused]] constexpr inline static BitBoard Knight(Square square) {
        return BitBoard(knight[static_cast<uint8_t>(square)]);
    }

    [[maybe_unused]] constexpr inline static BitBoard Bishop(Square square, BitBoard blockers) {
        return sliding[blackMagic::Index<Piece::Bishop>(square, blockers)];
    }

    [[maybe_unused]] constexpr inline static BitBoard Rook(Square square, BitBoard blockers) {
        return sliding[blackMagic::Index<Piece::Rook>(square, blockers)];
    }

    [[maybe_unused]] constexpr inline static BitBoard King(Square square) {
        return BitBoard(king[static_cast<uint8_t>(square)]);
    }
}

#endif