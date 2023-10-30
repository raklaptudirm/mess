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

#ifndef CHESS_FEN
#define CHESS_FEN

#include <iostream>
#include <string>
#include <array>

#include "types/types.hpp"

#include "strutil/strutil.h"

#include "piece.hpp"
#include "square.hpp"
#include "castling.hpp"
#include "bitboard.hpp"

namespace Chess {
    class FEN {
        public:
            std::array<ColoredPiece, Square::N> Mailbox;

            Color SideToMove;
            Square EPTarget;

            uint16 PlysCount;
            uint8  DrawClock;

            Castling::Info CastlingInfo;
            Castling::Rights CastlingRights;

            bool FRC;

        private:
            static constexpr int MAILBOX_ID    = 0;
            static constexpr int SIDE_TM_ID    = 1;
            static constexpr int CASTLING_ID   = 2;
            static constexpr int EP_TARGET_ID  = 3;
            static constexpr int DRAW_CLOCK_ID = 4;
            static constexpr int MOVE_COUNT_ID = 5;

        public:
            static constexpr uint16 MoveToPlyCount(int mc, Color stm) {
                return static_cast<uint16>(mc * 2) - (stm == Color::White ? 2 : 1);
            }
            constexpr FEN(const std::string& fenString) {
                const std::vector<std::string> fields = strutil::split(fenString, " ");
                assert(fields.size() == 6);

                for (int sq = 0; sq < Square::N; sq++) {
                    Mailbox[sq] = ColoredPiece::None;
                }

                Square whiteKing = Square::None;
                Square blackKing = Square::None;

                const std::vector<std::string> ranks = strutil::split(fields[MAILBOX_ID], "/");
                for (uint8 rank = 0; rank < 8; rank++) {
                    uint8 file = 0;

                    for (const auto& info : ranks[rank]) {
                        if ('1' <= info && info <= '8') file += info - '1';
                        else Mailbox[(7-rank)*8 + file] = ColoredPiece(std::string(1, info));

                        if (info == 'K') whiteKing = Square((7-rank)*8 + file);
                        if (info == 'k') blackKing = Square((7-rank)*8 + file);

                        file++;
                    }
                }

                // rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1

                SideToMove = Color(fields[SIDE_TM_ID]);

                auto info = Castling::Info::Parse(fields[CASTLING_ID], whiteKing, blackKing);
                CastlingInfo = info.first;
                CastlingRights = info.second;

                EPTarget = Square(fields[EP_TARGET_ID]);

                DrawClock = std::stoi(fields[DRAW_CLOCK_ID]);
                PlysCount = MoveToPlyCount(std::stoi(fields[MOVE_COUNT_ID]), SideToMove);
            }
    };
}

#endif