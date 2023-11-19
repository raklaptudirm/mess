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

#include <cassert>
#include <cstdint>
#include <string>

#include "square.hpp"
#include "piece.hpp"
#include "castling.hpp"

namespace Chess {
    struct MoveFlag {
        enum internal_flag : uint8_t {
            Normal,
            NPromotion, BPromotion, RPromotion, QPromotion,
            EnPassant, DoublePush,
            CastleHSide, CastleASide,
        };

        uint8_t internal;

        constexpr MoveFlag(internal_flag flag) {
            internal = flag;
        }

        constexpr explicit MoveFlag(uint8_t flag) {
            internal = static_cast<internal_flag>(flag);
        }

        constexpr explicit MoveFlag(Castling::Side side) {
            internal = side == Castling::Side::H ? CastleHSide : CastleASide;
        }

        [[nodiscard]] constexpr bool IsPromotion() const {
            return NPromotion <= internal && internal <= QPromotion;
        }

        [[nodiscard]] constexpr bool IsCastling() const {
            return internal == CastleHSide || internal == CastleASide;
        }

        [[nodiscard]] constexpr Piece PromotedPiece() const {
            assert(IsPromotion());
            return static_cast<Piece>(internal);
        }

        constexpr explicit operator uint8_t() const {
            return static_cast<uint8_t>(internal);
        }
    };

    class Move {
    private:
        uint16_t internal;

        static constexpr int SOURCE_WIDTH = 6;
        static constexpr int TARGET_WIDTH = 6;
        static constexpr int MVFLAG_WIDTH = 4;

        static constexpr uint16_t SOURCE_MASK = (1 << SOURCE_WIDTH) - 1;
        static constexpr uint16_t TARGET_MASK = (1 << TARGET_WIDTH) - 1;
        static constexpr uint16_t MVFLAG_MASK = (1 << MVFLAG_WIDTH) - 1;

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

        constexpr inline Move(Square source, Square target, MoveFlag flag) {
            internal =
                (static_cast<uint16_t>(static_cast<uint8_t>(flag  )) << MVFLAG_OFFSET) |
                (static_cast<uint16_t>(static_cast<uint8_t>(source)) << SOURCE_OFFSET) |
                (static_cast<uint16_t>(static_cast<uint8_t>(target)) << TARGET_OFFSET);
        }

        [[nodiscard]] constexpr inline Square Source() const {
            return Square((internal >> SOURCE_OFFSET) & SOURCE_MASK);
        }

        [[nodiscard]] constexpr inline Square Target() const {
            return Square((internal >> TARGET_OFFSET) & TARGET_MASK);
        }

        [[nodiscard]] constexpr inline MoveFlag Flag() const {
            return MoveFlag((internal >> MVFLAG_OFFSET) & MVFLAG_MASK);
        }

        [[nodiscard]] constexpr inline std::string ToString() const {
            if (internal == 0) return "0000";
            return this->Source().ToString() + this->Target().ToString()
             + (this->Flag().IsPromotion() ? this->Flag().PromotedPiece().ToString() : "");
        }
    };

    constexpr inline std::ostream& operator<<(std::ostream& os, const Move& move) {
        os << move.ToString();
        return os;
    }
}

#endif