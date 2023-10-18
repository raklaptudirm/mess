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

#ifndef CHESS_MOVE
#define CHESS_MOVE

#include <cstdint>
#include <string>

#include "square.hpp"
#include "piece.hpp"
#include "castling.hpp"

namespace Chess {
    class Move {
        private:
            uint16 internal;

            static constexpr int SOURCE_WIDTH = 6;
            static constexpr int TARGET_WIDTH = 6;
            static constexpr int MVFLAG_WIDTH = 4;

            static constexpr uint16 SOURCE_MASK = (1 << SOURCE_WIDTH) - 1;
            static constexpr uint16 TARGET_MASK = (1 << TARGET_WIDTH) - 1;
            static constexpr uint16 MVFLAG_MASK = (1 << MVFLAG_WIDTH) - 1;

            static constexpr int SOURCE_OFFSET = 0;
            static constexpr int TARGET_OFFSET = SOURCE_OFFSET + SOURCE_WIDTH;
            static constexpr int MVFLAG_OFFSET = TARGET_OFFSET + TARGET_WIDTH;

        public:
            constexpr inline Move() = default;

            // MaxInGame represents the maximum number of moves that
            // can occur in a chess game. Games longer than 512 moves
            // are possible, but unlikely to occur in actual gameplay.
            constexpr static int MaxInGame = 512;

            // MaxInPosition represents the maximum number of moves
            // that can be legal in a chess position. Refer to
            // https://cutt.ly/ZwijiNYq for the source of the figure
            // of 218, which has been rounded to 220 here.
            constexpr static int MaxInPosition = 220;

            constexpr inline Move(Square source, Square target, uint16 flag) {
                internal = (flag << MVFLAG_OFFSET) |
                    (static_cast<uint16>(static_cast<uint8>(source) << SOURCE_OFFSET)) |
                    (static_cast<uint16>(static_cast<uint8>(target) << TARGET_OFFSET));
            }

            [[nodiscard]] constexpr inline Square Source() const {
                return Square((internal >> SOURCE_OFFSET) & SOURCE_MASK);
            }

            [[nodiscard]] constexpr inline Square Target() const {
                return Square((internal >> TARGET_OFFSET) & TARGET_MASK);
            }

            [[nodiscard]] constexpr inline uint16 Flag() const {
                return (internal >> MVFLAG_OFFSET) & MVFLAG_MASK;
            }

            struct Flag {
                // Flag for a Normal move i.e. none of the
                // other flags are applicable to the move.
                static constexpr uint16 Normal = 0;

                /***************************
                 * Pawn Special Move Flags *
                 ***************************/

                // Promotions Moves.

                static constexpr uint16 NPromotion = 1; // Flag for promotion to a Knight.
                static constexpr uint16 BPromotion = 2; // Flag for promotion to a Bishop.
                static constexpr uint16 RPromotion = 3; // Flag for promotion to a Rook.
                static constexpr uint16 QPromotion = 4; // Flag for promotion to a Queen.

                // Other Special Moves.

                static constexpr uint16 EnPassant  = 5; // Flag for En Passant capture.
                static constexpr uint16 DoublePush = 6; // Flag for Double Pawn Push.

                /******************
                 * Castling Flags *
                 ******************/

                static constexpr uint16 CastleHSide = 7; // Flag for Castling O-O.
                static constexpr uint16 CastleASide = 8; // Flag for Castling O-O-O.

                static constexpr uint16 FlagFrom(Castling::Side side) {
                    if (side == Castling::Side::H)
                        return CastleHSide;
                    return CastleASide;
                }

                static constexpr bool IsPromotion(uint16 flag) {
                    return NPromotion <= flag && flag <= QPromotion;
                }

                static constexpr bool IsCastling(uint16 flag) {
                    return flag == CastleHSide || flag == CastleASide;
                }
            };

            [[nodiscard]] constexpr inline std::string ToString() const {
                if (internal == 0) return "0000";
                return this->Source().ToString() + this->Target().ToString() 
                 + (Flag::IsPromotion(this->Flag()) ? Piece(this->Flag()).ToString() : "");
            }
    };

    constexpr inline std::ostream& operator<<(std::ostream& os, const Move& move) {
        os << move.ToString();
        return os;
    }
}

#endif